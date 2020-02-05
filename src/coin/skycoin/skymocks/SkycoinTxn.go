// Code generated by mockery v1.0.0. DO NOT EDIT.

package skymocks

import mock "github.com/stretchr/testify/mock"

// SkycoinTxn is an autogenerated mock type for the SkycoinTxn type
type SkycoinTxn struct {
	mock.Mock
}

// EncodeSkycoinTransaction provides a mock function with given fields:
func (_m *SkycoinTxn) EncodeSkycoinTransaction() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
