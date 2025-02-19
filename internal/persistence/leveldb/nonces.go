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

package leveldb

import (
	"context"
	"time"

	"github.com/hyperledger/firefly-common/pkg/log"
	"github.com/hyperledger/firefly-transaction-manager/internal/persistence"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
)

type lockedNonce struct {
	th       *leveldbPersistence
	nsOpID   string
	signer   string
	unlocked chan struct{}
	nonce    uint64
	spent    bool
}

// complete must be called for any lockedNonce returned from a successful assignAndLockNonce call
func (ln *lockedNonce) complete(ctx context.Context) {
	if ln.spent {
		log.L(ctx).Debugf("Next nonce %d for signer %s spent", ln.nonce, ln.signer)
	} else {
		log.L(ctx).Debugf("Returning next nonce %d for signer %s unspent", ln.nonce, ln.signer)
	}
	ln.th.nonceMux.Lock()
	delete(ln.th.lockedNonces, ln.signer)
	close(ln.unlocked)
	ln.th.nonceMux.Unlock()
}

func (p *leveldbPersistence) assignAndLockNonce(ctx context.Context, nsOpID, signer string, nextNonceCB persistence.NextNonceCallback) (*lockedNonce, error) {

	for {
		// Take the lock to query our nonce cache, and check if we are already locked
		p.nonceMux.Lock()
		doLookup := false
		locked, isLocked := p.lockedNonces[signer]
		if !isLocked {
			locked = &lockedNonce{
				th:       p,
				nsOpID:   nsOpID,
				signer:   signer,
				unlocked: make(chan struct{}),
			}
			p.lockedNonces[signer] = locked
			doLookup = true
		}
		p.nonceMux.Unlock()

		// If we're locked, then wait
		if isLocked {
			log.L(ctx).Debugf("Contention for next nonce for signer %s", signer)
			<-locked.unlocked
		} else if doLookup {
			// We have to ensure we either successfully return a nonce,
			// or otherwise we unlock when we send the error
			nextNonce, err := p.calcNextNonce(ctx, signer, nextNonceCB)
			if err != nil {
				locked.complete(ctx)
				return nil, err
			}
			locked.nonce = nextNonce
			return locked, nil
		}
	}

}

func (p *leveldbPersistence) calcNextNonce(ctx context.Context, signer string, nextNonceCB persistence.NextNonceCallback) (uint64, error) {

	// First we check our DB to find the last nonce we used for this address.
	// Note we are within the nonce-lock in assignAndLockNonce for this signer, so we can be sure we're the
	// only routine attempting this right now.
	var lastTxn *apitypes.ManagedTX
	txns, err := p.ListTransactionsByNonce(ctx, signer, nil, 1, 1)
	if err != nil {
		return 0, err
	}
	if len(txns) > 0 {
		lastTxn = txns[0]
		if time.Since(*lastTxn.Created.Time()) < p.nonceStateTimeout {
			nextNonce := lastTxn.Nonce.Uint64() + 1
			log.L(ctx).Debugf("Allocating next nonce '%s' / '%d' after TX '%s' (status=%s)", signer, nextNonce, lastTxn.ID, lastTxn.Status)
			return nextNonce, nil
		}
	}

	// If we don't have a fresh answer in our state store, then ask the node.
	nextNonce, err := nextNonceCB(ctx, signer)
	if err != nil {
		return 0, err
	}

	// If we had a stale answer in our state store, make sure this isn't re-used.
	// This is important in case we have transactions that have expired from the TX pool of nodes, but we still have them
	// in our state store. So basically whichever is further forwards of our state store and the node answer wins.
	if lastTxn != nil && nextNonce <= lastTxn.Nonce.Uint64() {
		log.L(ctx).Debugf("Node TX pool next nonce '%s' / '%d' is not ahead of '%d' in TX '%s' (status=%s)", signer, nextNonce, lastTxn.Nonce.Uint64(), lastTxn.ID, lastTxn.Status)
		nextNonce = lastTxn.Nonce.Uint64() + 1
	}

	return nextNonce, nil

}
