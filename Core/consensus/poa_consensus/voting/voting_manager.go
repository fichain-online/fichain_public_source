package voting

type VotingManager interface {
	SubmitVote(vote Vote) error
	HasReachedQuorum(blockHash string) (bool, error)
	GetVotes(blockHash string) ([]Vote, error)
}
