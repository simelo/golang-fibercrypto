// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StringAddress is an autogenerated mock type for the StringAddress type
type StringAddress struct {
	mock.Mock
}

// GetCoinType provides a mock function with given fields:
func (_m *StringAddress) GetCoinType() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// GetValue provides a mock function with given fields:
func (_m *StringAddress) GetValue() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// IsValid provides a mock function with given fields:
func (_m *StringAddress) IsValid() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SetCoinType provides a mock function with given fields: val
func (_m *StringAddress) SetCoinType(val []byte) {
	_m.Called(val)
}

// SetValue provides a mock function with given fields: val
func (_m *StringAddress) SetValue(val []byte) {
	_m.Called(val)
}
