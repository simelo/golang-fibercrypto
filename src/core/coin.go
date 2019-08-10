package core

type TransactionTimestamp uint64
type TransactionStatus uint32

const (
	TXN_STATUS_CREATED TransactionStatus = iota
	TXN_STATUS_PENDING
	TXN_STATUS_CONFIRMED
)

type Transaction interface {
	SupportedAssets() []string
	GetTimestamp() TransactionTimestamp
	GetStatus() TransactionStatus
	GetInputs() []TransactionInput
	GetOutputs() []TransactionOutput
	GetId() string
	ComputeFee(ticker string) uint64
}

type TransactionIterator interface {
	Value() Transaction
	Next() bool
	HasNext() bool
}

type TransactionInput interface {
	GetId() string
	GetSpentOutput() TransactionOutput
	GetCoins(ticker string) uint64
}

type TransactionInputIterator interface {
	Value() TransactionInput
	Next() bool
	HasNext() bool
}

type TransactionOutput interface {
	GetId() string
	IsSpent() bool
	GetAddress() Address
	GetCoins(ticker string) uint64
}

type TransactionOutputIterator interface {
	Value() TransactionOutput
	Next() bool
	HasNext() bool
}
