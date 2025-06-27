package main

// TransactionData represents the structure of a transaction.
// It uses struct tags to map the JSON keys to the Go fields.
type TransactionData struct {
	Name           string `json:"name"`
	To             string `json:"to"`
	Data           string `json:"data"`
	Amount         string `json:"amount"`    // Using string for large numbers to avoid precision loss
	Gas            uint64 `json:"gas"`       // uint64 is suitable for non-negative gas values
	GasPrice       string `json:"gas_price"` // Also string for potentially large numbers (e.g., in Gwei/Wei)
	ReplaceAddress []int  `json:"replace_address"`
}
