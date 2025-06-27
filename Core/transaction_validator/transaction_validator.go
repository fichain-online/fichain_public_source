package transaction_validator

import (
	"FichainCore/params"
	"FichainCore/transaction"
)

type TransactionValidator struct {
	EkycApiUrl  string
	ChainConfig *params.ChainConfig
}

func NewTransactionValidator() TransactionValidator {
	return TransactionValidator{}
}

func (v TransactionValidator) QuickVerify(tx transaction.Transaction) error {
	// from, err := tx.From(v.ChainConfig.ChainId)
	// if err != nil {
	// 	return err
	// }

	// client := &http.Client{Timeout: 1 * time.Second}
	// req, err := http.NewRequest("GET", v.EkycApiUrl+from.Hex(), nil)
	// if err != nil {
	// 	return err
	// }

	// resp, err := client.Do(req)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return errors.New("transaction verification failed")
	// }

	return nil
}

func (v TransactionValidator) DeepVerify(tx *transaction.Transaction) error {
	return nil
}
