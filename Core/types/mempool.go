package types

// Mempool định nghĩa các hành vi chính mà bộ nhớ tạm của giao dịch cần hỗ trợ
type Mempool interface {
	// AddTransaction thêm một giao dịch mới vào mempool
	AddTransaction(tx Transaction) error

	// GetPendingTransactions trả về tất cả giao dịch đang chờ xử lý
	GetPendingTransactions() []Transaction

	// RemoveTransactions loại bỏ các giao dịch đã được xử lý
	RemoveTransactions(txs []Transaction) error

	// Size trả về số lượng giao dịch hiện có trong mempool
	Size() int

	// Clear xóa toàn bộ mempool
	Clear()
}
