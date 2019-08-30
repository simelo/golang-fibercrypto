package skycoin 
 
import (
	"strconv" 
	"github.com/skycoin/skycoin/src/readable"
	"github.com/fibercrypto/FiberCryptoWallet/src/core" 
	"github.com/fibercrypto/FiberCryptoWallet/src/util"
) 
 
/* 
SkycoinPendingTransaction 
*/ 
type SkycoinPendingTransaction struct{ //Implements Transaction interface
	Transaction readable.UnconfirmedTransactionVerbose
} 
 
func (txn *SkycoinPendingTransaction) SupportedAssets() []string { 
  	return []string{Sky, CoinHour} 
} 
 
func (txn *SkycoinPendingTransaction) GetTimestamp() core.Timestamp { 
  	return core.Timestamp(txn.Transaction.Received.Unix())
} 
 
func (txn *SkycoinPendingTransaction) GetStatus() core.TransactionStatus { 
  	return core.TXN_STATUS_PENDING
} 
 
func (txn *SkycoinPendingTransaction) GetInputs() []core.TransactionInput { 
	inputs := make([]core.TransactionInput, 0)
	for _ , input := range txn.Transaction.Transaction.In {
		inputs = append(inputs, &SkycoinTransactionInput{Input: input})
	}
	return inputs
} 
 
func (txn *SkycoinPendingTransaction) GetOutputs() []core.TransactionOutput { 
	outputs := make([]core.TransactionOutput, 0)
	for _ , output := range txn.Transaction.Transaction.Out {
		outputs = append(outputs, &SkycoinTransactionOutput{Output: output})
	}
	return outputs
} 
 
func (txn *SkycoinPendingTransaction) GetId() string { 
  	return txn.Transaction.Transaction.Hash
} 
 
func (txn *SkycoinPendingTransaction) ComputeFee(ticker string) uint64 { 
	if ticker == Sky {
		return uint64(0);
	}
	return txn.Transaction.Transaction.Fee
} 
 
/** 
 * SkycoinTransactionIterator 
 */ 
type SkycoinTransactionIterator struct { //Implements TransactionIterator interface
	Current int
	Transactions []core.Transaction
} 
 
func (it *SkycoinTransactionIterator) Value() core.Transaction { 
	return it.Transactions[it.Current]
} 
 
func (it *SkycoinTransactionIterator) Next() bool { 
	if it.HasNext() {
		it.Current++
		return true
	}
	return false
} 
 
func (it *SkycoinTransactionIterator) HasNext() bool { 
	return (it.Current + 1) < len(it.Transactions)
} 
 
func NewSkycoinTransactionIterator(transactions []core.Transaction) *SkycoinTransactionIterator {
	return &SkycoinTransactionIterator{Transactions: transactions, Current: -1}
}

type SkycoinTransactionInput struct { //Implements TransactionInput interface
	Input readable.TransactionInput
} 
 
func (in *SkycoinTransactionInput) GetId() string { 
  	return "" 
} 
 
func (in *SkycoinTransactionInput) IsSpent() bool { 
  	return true 
} 
 
func (in *SkycoinTransactionInput) GetSpentOutput() core.TransactionOutput { 
  	return nil 
} 
 
/** 
 * SkycoinTransactionInputIterator 
 */ 
type SkycoinTransactionInputIterator struct { 
} 
 
func (iter *SkycoinTransactionInputIterator) Value() core.TransactionInput { 
	return nil 
} 
 
func (iter *SkycoinTransactionInputIterator) Next() bool { 
	return false 
} 
 
func (iter *SkycoinTransactionInputIterator) HasNext() bool { 
  	return false
} 
 
/** 
 * SkycoinTransactionOutput 
 */ 
type SkycoinTransactionOutput struct { //Implements TransactionOutput interface 
	Output readable.TransactionOutput
} 

func (sto *SkycoinTransactionOutput) GetId() string { 
	return sto.Output.Hash 
} 

func (sto *SkycoinTransactionOutput) IsSpent() bool { 
	//TODO:
	return false 
} 

func (sto *SkycoinTransactionOutput) GetAddress() core.Address { 
	return SkycoinAddress{address: sto.Output.Address}
} 

func (sto *SkycoinTransactionOutput) GetCoins(ticker string) (uint64, error) { 
	accuracy, err := util.AltcoinQuotient(ticker)
	if err != nil {
		return uint64(0), err
	}
	if ticker == Sky { 
		coin, err2 := strconv.ParseFloat(sto.Output.Coins, 64)
		if err2 != nil {
			return uint64(0), err2
		}
		return uint64(coin * float64(accuracy)), nil
	} 
	return sto.Output.Hours * accuracy, nil
} 

type SkycoinTransactionOutputIterator struct { //Implements TransactionOutputIterator interface 
	Current int 
	Outputs []core.TransactionOutput 
} 

func (it *SkycoinTransactionOutputIterator) Value() core.TransactionOutput { 
	return it.Outputs[it.Current] 
} 

func (it *SkycoinTransactionOutputIterator) Next() bool { 
	if it.HasNext() { 
		it.Current++ 
		return true 
	} 
	return false 
} 

func (it *SkycoinTransactionOutputIterator) HasNext() bool { 
	return (it.Current + 1) < len(it.Outputs) 
}

func NewSkycoinTransactionOutputIterator(outputs []core.TransactionOutput) *SkycoinTransactionOutputIterator {
	return &SkycoinTransactionOutputIterator{Outputs: outputs, Current: -1}
}

