// Copyright © 2023 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly-transaction-manager/internal/tmconfig"
	"github.com/hyperledger/firefly-transaction-manager/internal/tmmsgs"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type leveldbPersistence struct {
	db         *leveldb.DB
	syncWrites bool
	txMux      sync.RWMutex // allows us to draw conclusions on the cleanup of indexes
}

func NewLevelDBPersistence(ctx context.Context) (Persistence, error) {
	dbPath := config.GetString(tmconfig.PersistenceLevelDBPath)
	if dbPath == "" {
		return nil, i18n.NewError(ctx, tmmsgs.MsgLevelDBPathMissing)
	}
	db, err := leveldb.OpenFile(dbPath, &opt.Options{
		OpenFilesCacheCapacity: config.GetInt(tmconfig.PersistenceLevelDBMaxHandles),
	})
	if err != nil {
		return nil, i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceInitFailed, dbPath)
	}
	return &leveldbPersistence{
		db:         db,
		syncWrites: config.GetBool(tmconfig.PersistenceLevelDBSyncWrites),
	}, nil
}

const checkpointsPrefix = "checkpoints_0/"
const eventstreamsPrefix = "eventstreams_0/"
const eventstreamsEnd = "eventstreams_1"
const listenersPrefix = "listeners_0/"
const listenersEnd = "listeners_1"
const transactionsPrefix = "tx_0/"
const nonceAllocationPrefix = "nonce_0/"
const txPendingIndexPrefix = "tx_inflight_0/"
const txPendingIndexEnd = "tx_inflight_1"
const txCreatedIndexPrefix = "tx_created_0/"
const txCreatedIndexEnd = "tx_created_1"

func signerNoncePrefix(signer string) string {
	return fmt.Sprintf("%s%s_0/", nonceAllocationPrefix, signer)
}

func signerNonceEnd(signer string) string {
	return fmt.Sprintf("%s%s_1", nonceAllocationPrefix, signer)
}

func txNonceAllocationKey(signer string, nonce *fftypes.FFBigInt) []byte {
	return []byte(fmt.Sprintf("%s%s_0/%.24d", nonceAllocationPrefix, signer, nonce.Int()))
}

func txPendingIndexKey(sequenceID string) []byte {
	return []byte(fmt.Sprintf("%s%s", txPendingIndexPrefix, sequenceID))
}

func txCreatedIndexKey(tx *apitypes.ManagedTX) []byte {
	return []byte(fmt.Sprintf("%s%.19d/%s", txCreatedIndexPrefix, tx.Created.UnixNano(), tx.SequenceID))
}

func txDataKey(k string) []byte {
	return []byte(fmt.Sprintf("%s%s", transactionsPrefix, k))
}

func prefixedKey(prefix string, id fmt.Stringer) []byte {
	return []byte(fmt.Sprintf("%s%s", prefix, id))
}

func (p *leveldbPersistence) writeKeyValue(ctx context.Context, key, value []byte) error {
	err := p.db.Put(key, value, &opt.WriteOptions{Sync: p.syncWrites})
	if err != nil {
		return i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceWriteFailed)
	}
	return nil
}

func (p *leveldbPersistence) writeJSON(ctx context.Context, key []byte, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceMarshalFailed)
	}
	log.L(ctx).Debugf("Wrote %s", key)
	return p.writeKeyValue(ctx, key, b)
}

func (p *leveldbPersistence) getKeyValue(ctx context.Context, key []byte) ([]byte, error) {
	b, err := p.db.Get(key, &opt.ReadOptions{})
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, nil
		}
		return nil, i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceReadFailed, key)
	}
	return b, err
}

func (p *leveldbPersistence) readJSONByIndex(ctx context.Context, idxKey []byte, target interface{}) error {
	valKey, err := p.getKeyValue(ctx, idxKey)
	if err != nil || valKey == nil {
		return err
	}
	return p.readJSON(ctx, valKey, target)
}

func (p *leveldbPersistence) readJSON(ctx context.Context, key []byte, target interface{}) error {
	b, err := p.getKeyValue(ctx, key)
	if err != nil || b == nil {
		return err
	}
	err = json.Unmarshal(b, target)
	if err != nil {
		return i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceUnmarshalFailed)
	}
	log.L(ctx).Debugf("Read %s", key)
	return nil
}

func (p *leveldbPersistence) listJSON(ctx context.Context, collectionPrefix, collectionEnd, after string, limit int,
	dir SortDirection,
	val func() interface{}, // return a pointer to a pointer variable, of the type to unmarshal
	add func(interface{}), // passes back the val() for adding to the list, if the filters match
	indexResolver func(ctx context.Context, k []byte) ([]byte, error), // if non-nil then the initial lookup will be passed to this, to lookup the target bytes. Nil skips item
	filters ...func(interface{}) bool, // filters to apply to the val() after unmarshalling
) ([][]byte, error) {
	collectionRange := &util.Range{
		Start: []byte(collectionPrefix),
		Limit: []byte(collectionEnd),
	}
	var it iterator.Iterator
	switch dir {
	case SortDirectionAscending:
		afterKey := collectionPrefix + after
		if after != "" {
			collectionRange.Start = []byte(afterKey)
		}
		it = p.db.NewIterator(collectionRange, &opt.ReadOptions{DontFillCache: true})
		if after != "" && it.Next() {
			if !strings.HasPrefix(string(it.Key()), afterKey) {
				it.Prev() // skip back, as the first key was already after the "after" key
			}
		}
	default:
		if after != "" {
			collectionRange.Limit = []byte(collectionPrefix + after) // exclusive for limit, so no need to fiddle here
		}
		it = p.db.NewIterator(collectionRange, &opt.ReadOptions{DontFillCache: true})
	}
	defer it.Release()
	return p.iterateJSON(ctx, it, limit, dir, val, add, indexResolver, filters...)
}

func (p *leveldbPersistence) iterateJSON(ctx context.Context, it iterator.Iterator, limit int,
	dir SortDirection, val func() interface{}, add func(interface{}), indexResolver func(ctx context.Context, k []byte) ([]byte, error), filters ...func(interface{}) bool,
) (orphanedIdxKeys [][]byte, err error) {
	count := 0
	next := it.Next // forwards we enter this function before the first key
	if dir == SortDirectionDescending {
		next = it.Last // reverse we enter this function
	}
itLoop:
	for next() {
		if dir == SortDirectionDescending {
			next = it.Prev
		} else {
			next = it.Next
		}
		v := val()
		b := it.Value()
		if indexResolver != nil {
			valKey := b
			b, err = indexResolver(ctx, valKey)
			if err != nil {
				return nil, err
			}
			if b == nil {
				log.L(ctx).Warnf("Skipping orphaned index key '%s' pointing to '%s'", it.Key(), valKey)
				orphanedIdxKeys = append(orphanedIdxKeys, it.Key())
				continue itLoop
			}
		}
		err := json.Unmarshal(b, v)
		if err != nil {
			return nil, i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceUnmarshalFailed)
		}
		for _, f := range filters {
			if !f(v) {
				continue itLoop
			}
		}
		add(v)
		count++
		if limit > 0 && count >= limit {
			break
		}
	}
	log.L(ctx).Debugf("Listed %d items", count)
	return orphanedIdxKeys, nil
}

func (p *leveldbPersistence) deleteKeys(ctx context.Context, keys ...[]byte) error {
	for _, key := range keys {
		err := p.db.Delete(key, &opt.WriteOptions{Sync: p.syncWrites})
		if err != nil && err != leveldb.ErrNotFound {
			return i18n.WrapError(ctx, err, tmmsgs.MsgPersistenceDeleteFailed)
		}
		log.L(ctx).Debugf("Deleted %s", key)
	}
	return nil
}

func (p *leveldbPersistence) WriteCheckpoint(ctx context.Context, checkpoint *apitypes.EventStreamCheckpoint) error {
	return p.writeJSON(ctx, prefixedKey(checkpointsPrefix, checkpoint.StreamID), checkpoint)
}

func (p *leveldbPersistence) GetCheckpoint(ctx context.Context, streamID *fftypes.UUID) (cp *apitypes.EventStreamCheckpoint, err error) {
	err = p.readJSON(ctx, prefixedKey(checkpointsPrefix, streamID), &cp)
	return cp, err
}

func (p *leveldbPersistence) DeleteCheckpoint(ctx context.Context, streamID *fftypes.UUID) error {
	return p.deleteKeys(ctx, prefixedKey(checkpointsPrefix, streamID))
}

func (p *leveldbPersistence) ListStreams(ctx context.Context, after *fftypes.UUID, limit int, dir SortDirection) ([]*apitypes.EventStream, error) {
	streams := make([]*apitypes.EventStream, 0)
	if _, err := p.listJSON(ctx, eventstreamsPrefix, eventstreamsEnd, after.String(), limit, dir,
		func() interface{} { var v *apitypes.EventStream; return &v },
		func(v interface{}) { streams = append(streams, *(v.(**apitypes.EventStream))) },
		nil,
	); err != nil {
		return nil, err
	}
	return streams, nil
}

func (p *leveldbPersistence) GetStream(ctx context.Context, streamID *fftypes.UUID) (es *apitypes.EventStream, err error) {
	err = p.readJSON(ctx, prefixedKey(eventstreamsPrefix, streamID), &es)
	return es, err
}

func (p *leveldbPersistence) WriteStream(ctx context.Context, spec *apitypes.EventStream) error {
	return p.writeJSON(ctx, prefixedKey(eventstreamsPrefix, spec.ID), spec)
}

func (p *leveldbPersistence) DeleteStream(ctx context.Context, streamID *fftypes.UUID) error {
	return p.deleteKeys(ctx, prefixedKey(eventstreamsPrefix, streamID))
}

func (p *leveldbPersistence) ListListeners(ctx context.Context, after *fftypes.UUID, limit int, dir SortDirection) ([]*apitypes.Listener, error) {
	listeners := make([]*apitypes.Listener, 0)
	if _, err := p.listJSON(ctx, listenersPrefix, listenersEnd, after.String(), limit, dir,
		func() interface{} { var v *apitypes.Listener; return &v },
		func(v interface{}) { listeners = append(listeners, *(v.(**apitypes.Listener))) },
		nil,
	); err != nil {
		return nil, err
	}
	return listeners, nil
}

func (p *leveldbPersistence) ListStreamListeners(ctx context.Context, after *fftypes.UUID, limit int, dir SortDirection, streamID *fftypes.UUID) ([]*apitypes.Listener, error) {
	listeners := make([]*apitypes.Listener, 0)
	if _, err := p.listJSON(ctx, listenersPrefix, listenersEnd, after.String(), limit, dir,
		func() interface{} { var v *apitypes.Listener; return &v },
		func(v interface{}) { listeners = append(listeners, *(v.(**apitypes.Listener))) },
		nil,
		func(v interface{}) bool { return (*(v.(**apitypes.Listener))).StreamID.Equals(streamID) },
	); err != nil {
		return nil, err
	}
	return listeners, nil
}

func (p *leveldbPersistence) GetListener(ctx context.Context, listenerID *fftypes.UUID) (l *apitypes.Listener, err error) {
	err = p.readJSON(ctx, prefixedKey(listenersPrefix, listenerID), &l)
	return l, err
}

func (p *leveldbPersistence) WriteListener(ctx context.Context, spec *apitypes.Listener) error {
	return p.writeJSON(ctx, prefixedKey(listenersPrefix, spec.ID), spec)
}

func (p *leveldbPersistence) DeleteListener(ctx context.Context, listenerID *fftypes.UUID) error {
	return p.deleteKeys(ctx, prefixedKey(listenersPrefix, listenerID))
}

func (p *leveldbPersistence) indexLookupCallback(ctx context.Context, key []byte) ([]byte, error) {
	b, err := p.getKeyValue(ctx, key)
	switch {
	case err != nil:
		return nil, err
	case b == nil:
		return nil, nil
	}
	return b, err
}

func (p *leveldbPersistence) cleanupOrphanedTXIdxKeys(ctx context.Context, orphanedIdxKeys [][]byte) {
	p.txMux.Lock()
	defer p.txMux.Unlock()
	err := p.deleteKeys(ctx, orphanedIdxKeys...)
	if err != nil {
		log.L(ctx).Warnf("Failed to clean up orphaned index keys: %s", err)
	}
}

func (p *leveldbPersistence) listTransactionsByIndex(ctx context.Context, collectionPrefix, collectionEnd, afterStr string, limit int, dir SortDirection) ([]*apitypes.ManagedTX, error) {

	p.txMux.RLock()
	transactions := make([]*apitypes.ManagedTX, 0)
	orphanedIdxKeys, err := p.listJSON(ctx, collectionPrefix, collectionEnd, afterStr, limit, dir,
		func() interface{} { var v *apitypes.ManagedTX; return &v },
		func(v interface{}) { transactions = append(transactions, *(v.(**apitypes.ManagedTX))) },
		p.indexLookupCallback,
	)
	p.txMux.RUnlock()
	if err != nil {
		return nil, err
	}
	// If we find orphaned index keys we clean them up - which requires the write lock (hence dropping read-lock first)
	if len(orphanedIdxKeys) > 0 {
		p.cleanupOrphanedTXIdxKeys(ctx, orphanedIdxKeys)
	}
	return transactions, nil
}

func (p *leveldbPersistence) ListTransactionsByCreateTime(ctx context.Context, after *apitypes.ManagedTX, limit int, dir SortDirection) ([]*apitypes.ManagedTX, error) {
	afterStr := ""
	if after != nil {
		afterStr = fmt.Sprintf("%.19d/%s", after.Created.UnixNano(), after.SequenceID)
	}
	return p.listTransactionsByIndex(ctx, txCreatedIndexPrefix, txCreatedIndexEnd, afterStr, limit, dir)
}

func (p *leveldbPersistence) ListTransactionsByNonce(ctx context.Context, signer string, after *fftypes.FFBigInt, limit int, dir SortDirection) ([]*apitypes.ManagedTX, error) {
	afterStr := ""
	if after != nil {
		afterStr = fmt.Sprintf("%.24d", after.Int())
	}
	return p.listTransactionsByIndex(ctx, signerNoncePrefix(signer), signerNonceEnd(signer), afterStr, limit, dir)
}

func (p *leveldbPersistence) ListTransactionsPending(ctx context.Context, afterSequenceID string, limit int, dir SortDirection) ([]*apitypes.ManagedTX, error) {
	return p.listTransactionsByIndex(ctx, txPendingIndexPrefix, txPendingIndexEnd, afterSequenceID, limit, dir)
}

func (p *leveldbPersistence) GetTransactionByID(ctx context.Context, txID string) (tx *apitypes.ManagedTX, err error) {
	p.txMux.RLock()
	defer p.txMux.RUnlock()
	err = p.readJSON(ctx, txDataKey(txID), &tx)
	return tx, err
}

func (p *leveldbPersistence) GetTransactionByNonce(ctx context.Context, signer string, nonce *fftypes.FFBigInt) (tx *apitypes.ManagedTX, err error) {
	p.txMux.RLock()
	defer p.txMux.RUnlock()
	err = p.readJSONByIndex(ctx, txNonceAllocationKey(signer, nonce), &tx)
	return tx, err
}

func (p *leveldbPersistence) WriteTransaction(ctx context.Context, tx *apitypes.ManagedTX, new bool) (err error) {
	// We take a write-lock here, because we are writing multiple values (the indexes), and anybody
	// attempting to read the critical nonce allocation index must know the difference between a partial write
	// (we crashed before we completed all the writes) and an incomplete write that's in process.
	// The reading code detects partial writes and cleans them up if it finds them.
	p.txMux.Lock()
	defer p.txMux.Unlock()

	if tx.TransactionHeaders.From == "" ||
		tx.Nonce == nil ||
		tx.Created == nil ||
		tx.ID == "" ||
		tx.Status == "" {
		return i18n.NewError(ctx, tmmsgs.MsgPersistenceTXIncomplete)
	}
	idKey := txDataKey(tx.ID)
	if new {
		if tx.SequenceID != "" {
			// for new transactions sequence ID should always be generated by persistence layer
			// as the format of its value is persistence service specific
			log.L(ctx).Errorf("Sequence ID is not allowed for new transaction %s", tx.ID)
			return i18n.NewError(ctx, tmmsgs.MsgPersistenceSequenceIDNotAllowed)
		}
		tx.SequenceID = apitypes.NewULID().String()
		// This must be a unique ID, otherwise we return a conflict.
		// Note we use the final record we write at the end for the conflict check, and also that we're write locked here
		if existing, err := p.getKeyValue(ctx, idKey); err != nil {
			return err
		} else if existing != nil {
			return i18n.NewError(ctx, tmmsgs.MsgDuplicateID, idKey)
		}

		// We write the index records first - because if we crash, we need to be able to know if the
		// index records are valid or not. When reading under the read lock, if there is an index key
		// that does not have a corresponding managed TX available, we will clean up the
		// orphaned index (after swapping the read lock for the write lock)
		// See listTransactionsByIndex() for the other half of this logic.
		err = p.writeKeyValue(ctx, txCreatedIndexKey(tx), idKey)
		if err == nil && tx.Status == apitypes.TxStatusPending {
			err = p.writeKeyValue(ctx, txPendingIndexKey(tx.SequenceID), idKey)
		}
		if err == nil {
			err = p.writeKeyValue(ctx, txNonceAllocationKey(tx.TransactionHeaders.From, tx.Nonce), idKey)
		}
	}
	// If we are creating/updating a record that is not pending, we need to ensure there is no pending index associated with it
	if err == nil && tx.Status != apitypes.TxStatusPending {
		err = p.deleteKeys(ctx, txPendingIndexKey(tx.SequenceID))
	}
	if err == nil {
		err = p.writeJSON(ctx, idKey, tx)
	}
	return err
}

func (p *leveldbPersistence) DeleteTransaction(ctx context.Context, txID string) error {
	var tx *apitypes.ManagedTX
	err := p.readJSON(ctx, txDataKey(txID), &tx)
	if err != nil || tx == nil {
		return err
	}
	return p.deleteKeys(ctx,
		txDataKey(txID),
		txCreatedIndexKey(tx),
		txPendingIndexKey(tx.SequenceID),
		txNonceAllocationKey(tx.TransactionHeaders.From, tx.Nonce),
	)
}

func (p *leveldbPersistence) Close(ctx context.Context) {
	err := p.db.Close()
	if err != nil {
		log.L(ctx).Warnf("Error closing leveldb: %s", err)
	}
}
