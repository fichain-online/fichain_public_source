package types

type Signer interface {
	Sign(tx Transaction, privateKey PrivateKey) (Transaction, error)
}
