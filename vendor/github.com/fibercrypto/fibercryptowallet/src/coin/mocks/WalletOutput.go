// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import core "github.com/fibercrypto/fibercryptowallet/src/core"
import mock "github.com/stretchr/testify/mock"

// WalletOutput is an autogenerated mock type for the WalletOutput type
type WalletOutput struct {
	mock.Mock
}

// GetOutput provides a mock function with given fields:
func (_m *WalletOutput) GetOutput() core.TransactionOutput {
	ret := _m.Called()

	var r0 core.TransactionOutput
	if rf, ok := ret.Get(0).(func() core.TransactionOutput); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.TransactionOutput)
		}
	}

	return r0
}

// GetWallet provides a mock function with given fields:
func (_m *WalletOutput) GetWallet() core.Wallet {
	ret := _m.Called()

	var r0 core.Wallet
	if rf, ok := ret.Get(0).(func() core.Wallet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Wallet)
		}
	}

	return r0
}
