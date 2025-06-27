package voting

type Vote struct {
	BlockHash string
	Voter     string
	Signature []byte
}
