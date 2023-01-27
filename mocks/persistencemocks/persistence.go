// Code generated by mockery v2.14.1. DO NOT EDIT.

package persistencemocks

import (
	context "context"

	apitypes "github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"

	fftypes "github.com/hyperledger/firefly-common/pkg/fftypes"

	mock "github.com/stretchr/testify/mock"

	persistence "github.com/hyperledger/firefly-transaction-manager/internal/persistence"
)

// Persistence is an autogenerated mock type for the Persistence type
type Persistence struct {
	mock.Mock
}

// Close provides a mock function with given fields: ctx
func (_m *Persistence) Close(ctx context.Context) {
	_m.Called(ctx)
}

// DeleteCheckpoint provides a mock function with given fields: ctx, streamID
func (_m *Persistence) DeleteCheckpoint(ctx context.Context, streamID *fftypes.UUID) error {
	ret := _m.Called(ctx, streamID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) error); ok {
		r0 = rf(ctx, streamID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteListener provides a mock function with given fields: ctx, listenerID
func (_m *Persistence) DeleteListener(ctx context.Context, listenerID *fftypes.UUID) error {
	ret := _m.Called(ctx, listenerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) error); ok {
		r0 = rf(ctx, listenerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteStream provides a mock function with given fields: ctx, streamID
func (_m *Persistence) DeleteStream(ctx context.Context, streamID *fftypes.UUID) error {
	ret := _m.Called(ctx, streamID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) error); ok {
		r0 = rf(ctx, streamID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTransaction provides a mock function with given fields: ctx, txID
func (_m *Persistence) DeleteTransaction(ctx context.Context, txID string) error {
	ret := _m.Called(ctx, txID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, txID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCheckpoint provides a mock function with given fields: ctx, streamID
func (_m *Persistence) GetCheckpoint(ctx context.Context, streamID *fftypes.UUID) (*apitypes.EventStreamCheckpoint, error) {
	ret := _m.Called(ctx, streamID)

	var r0 *apitypes.EventStreamCheckpoint
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) *apitypes.EventStreamCheckpoint); ok {
		r0 = rf(ctx, streamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.EventStreamCheckpoint)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID) error); ok {
		r1 = rf(ctx, streamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetListener provides a mock function with given fields: ctx, listenerID
func (_m *Persistence) GetListener(ctx context.Context, listenerID *fftypes.UUID) (*apitypes.Listener, error) {
	ret := _m.Called(ctx, listenerID)

	var r0 *apitypes.Listener
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) *apitypes.Listener); ok {
		r0 = rf(ctx, listenerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.Listener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID) error); ok {
		r1 = rf(ctx, listenerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStream provides a mock function with given fields: ctx, streamID
func (_m *Persistence) GetStream(ctx context.Context, streamID *fftypes.UUID) (*apitypes.EventStream, error) {
	ret := _m.Called(ctx, streamID)

	var r0 *apitypes.EventStream
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID) *apitypes.EventStream); ok {
		r0 = rf(ctx, streamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.EventStream)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID) error); ok {
		r1 = rf(ctx, streamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByID provides a mock function with given fields: ctx, txID
func (_m *Persistence) GetTransactionByID(ctx context.Context, txID string) (*apitypes.ManagedTX, error) {
	ret := _m.Called(ctx, txID)

	var r0 *apitypes.ManagedTX
	if rf, ok := ret.Get(0).(func(context.Context, string) *apitypes.ManagedTX); ok {
		r0 = rf(ctx, txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.ManagedTX)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByNonce provides a mock function with given fields: ctx, signer, nonce
func (_m *Persistence) GetTransactionByNonce(ctx context.Context, signer string, nonce *fftypes.FFBigInt) (*apitypes.ManagedTX, error) {
	ret := _m.Called(ctx, signer, nonce)

	var r0 *apitypes.ManagedTX
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.FFBigInt) *apitypes.ManagedTX); ok {
		r0 = rf(ctx, signer, nonce)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apitypes.ManagedTX)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.FFBigInt) error); ok {
		r1 = rf(ctx, signer, nonce)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListListeners provides a mock function with given fields: ctx, after, limit, dir
func (_m *Persistence) ListListeners(ctx context.Context, after *fftypes.UUID, limit int, dir persistence.SortDirection) ([]*apitypes.Listener, error) {
	ret := _m.Called(ctx, after, limit, dir)

	var r0 []*apitypes.Listener
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, int, persistence.SortDirection) []*apitypes.Listener); ok {
		r0 = rf(ctx, after, limit, dir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*apitypes.Listener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID, int, persistence.SortDirection) error); ok {
		r1 = rf(ctx, after, limit, dir)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListStreamListeners provides a mock function with given fields: ctx, after, limit, dir, streamID
func (_m *Persistence) ListStreamListeners(ctx context.Context, after *fftypes.UUID, limit int, dir persistence.SortDirection, streamID *fftypes.UUID) ([]*apitypes.Listener, error) {
	ret := _m.Called(ctx, after, limit, dir, streamID)

	var r0 []*apitypes.Listener
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, int, persistence.SortDirection, *fftypes.UUID) []*apitypes.Listener); ok {
		r0 = rf(ctx, after, limit, dir, streamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*apitypes.Listener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID, int, persistence.SortDirection, *fftypes.UUID) error); ok {
		r1 = rf(ctx, after, limit, dir, streamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListStreams provides a mock function with given fields: ctx, after, limit, dir
func (_m *Persistence) ListStreams(ctx context.Context, after *fftypes.UUID, limit int, dir persistence.SortDirection) ([]*apitypes.EventStream, error) {
	ret := _m.Called(ctx, after, limit, dir)

	var r0 []*apitypes.EventStream
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, int, persistence.SortDirection) []*apitypes.EventStream); ok {
		r0 = rf(ctx, after, limit, dir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*apitypes.EventStream)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID, int, persistence.SortDirection) error); ok {
		r1 = rf(ctx, after, limit, dir)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTransactionsByCreateTime provides a mock function with given fields: ctx, after, limit, dir
func (_m *Persistence) ListTransactionsByCreateTime(ctx context.Context, after *apitypes.ManagedTX, limit int, dir persistence.SortDirection) ([]*apitypes.ManagedTX, error) {
	ret := _m.Called(ctx, after, limit, dir)

	var r0 []*apitypes.ManagedTX
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.ManagedTX, int, persistence.SortDirection) []*apitypes.ManagedTX); ok {
		r0 = rf(ctx, after, limit, dir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*apitypes.ManagedTX)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *apitypes.ManagedTX, int, persistence.SortDirection) error); ok {
		r1 = rf(ctx, after, limit, dir)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTransactionsByNonce provides a mock function with given fields: ctx, signer, after, limit, dir
func (_m *Persistence) ListTransactionsByNonce(ctx context.Context, signer string, after *fftypes.FFBigInt, limit int, dir persistence.SortDirection) ([]*apitypes.ManagedTX, error) {
	ret := _m.Called(ctx, signer, after, limit, dir)

	var r0 []*apitypes.ManagedTX
	if rf, ok := ret.Get(0).(func(context.Context, string, *fftypes.FFBigInt, int, persistence.SortDirection) []*apitypes.ManagedTX); ok {
		r0 = rf(ctx, signer, after, limit, dir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*apitypes.ManagedTX)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *fftypes.FFBigInt, int, persistence.SortDirection) error); ok {
		r1 = rf(ctx, signer, after, limit, dir)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTransactionsPending provides a mock function with given fields: ctx, after, limit, dir
func (_m *Persistence) ListTransactionsPending(ctx context.Context, after *fftypes.UUID, limit int, dir persistence.SortDirection) ([]*apitypes.ManagedTX, error) {
	ret := _m.Called(ctx, after, limit, dir)

	var r0 []*apitypes.ManagedTX
	if rf, ok := ret.Get(0).(func(context.Context, *fftypes.UUID, int, persistence.SortDirection) []*apitypes.ManagedTX); ok {
		r0 = rf(ctx, after, limit, dir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*apitypes.ManagedTX)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *fftypes.UUID, int, persistence.SortDirection) error); ok {
		r1 = rf(ctx, after, limit, dir)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteCheckpoint provides a mock function with given fields: ctx, checkpoint
func (_m *Persistence) WriteCheckpoint(ctx context.Context, checkpoint *apitypes.EventStreamCheckpoint) error {
	ret := _m.Called(ctx, checkpoint)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.EventStreamCheckpoint) error); ok {
		r0 = rf(ctx, checkpoint)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteListener provides a mock function with given fields: ctx, spec
func (_m *Persistence) WriteListener(ctx context.Context, spec *apitypes.Listener) error {
	ret := _m.Called(ctx, spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.Listener) error); ok {
		r0 = rf(ctx, spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteStream provides a mock function with given fields: ctx, spec
func (_m *Persistence) WriteStream(ctx context.Context, spec *apitypes.EventStream) error {
	ret := _m.Called(ctx, spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.EventStream) error); ok {
		r0 = rf(ctx, spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteTransaction provides a mock function with given fields: ctx, tx, new
func (_m *Persistence) WriteTransaction(ctx context.Context, tx *apitypes.ManagedTX, new bool) error {
	ret := _m.Called(ctx, tx, new)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *apitypes.ManagedTX, bool) error); ok {
		r0 = rf(ctx, tx, new)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPersistence interface {
	mock.TestingT
	Cleanup(func())
}

// NewPersistence creates a new instance of Persistence. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPersistence(t mockConstructorTestingTNewPersistence) *Persistence {
	mock := &Persistence{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
