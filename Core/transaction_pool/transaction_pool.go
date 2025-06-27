package transaction_pool

import (
	"errors"
	"sort"
	"sync"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/common"
	"FichainCore/params"
	"FichainCore/transaction"
)

type TransactionPool struct {
	mu           sync.Mutex
	transactions []*transaction.Transaction
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		transactions: []*transaction.Transaction{},
	}
}

func (m *TransactionPool) AddTransaction(tx *transaction.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// (Tùy chọn) kiểm tra trùng lặp
	for _, existing := range m.transactions {
		if existing.Hash() == tx.Hash() {
			return errors.New("transaction already exists in mempool")
		}
	}

	m.transactions = append(m.transactions, tx)
	return nil
}

// GetPendingTransactions returns a map of transactions grouped by fromAddress,
// and sorted by nonce in ascending order.
func (m *TransactionPool) GetPendingTransactions() map[common.Address][]*transaction.Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	grouped := make(map[common.Address][]*transaction.Transaction)

	// Group transactions by from address
	for _, tx := range m.transactions {
		from, err := tx.From(params.TempChainId)
		if err != nil {
			logger.Error("error when get from from tx", err)
		}
		grouped[from] = append(grouped[from], tx)
	}

	// Sort each group by nonce
	for from, txs := range grouped {
		sorted := txs
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Nonce() < sorted[j].Nonce()
		})
		grouped[from] = sorted
	}

	return grouped
}

func (m *TransactionPool) RemoveTransactions(txs []*transaction.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	hashSet := make(map[common.Hash]bool)
	for _, tx := range txs {
		hashSet[tx.Hash()] = true
	}

	filtered := m.transactions[:0]
	for _, tx := range m.transactions {
		if !hashSet[tx.Hash()] {
			filtered = append(filtered, tx)
		}
	}
	m.transactions = filtered
	return nil
}

func (m *TransactionPool) Size() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.transactions)
}

func (m *TransactionPool) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transactions = []*transaction.Transaction{}
}
