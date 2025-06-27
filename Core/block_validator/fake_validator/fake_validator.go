package fake_validator

import (
	"FichainCore/block"
	"FichainCore/evm"
	"FichainCore/log"
	"FichainCore/receipt"
	"FichainCore/state"
)

type FakeValidator struct{}

func (*FakeValidator) ValidateBody(*block.Block) error { return nil }

func (*FakeValidator) ValidateState(
	block, parent *block.Block,
	state *state.StateDB,
	receipts []*receipt.Receipt,
	usedGas uint64,
) error {
	return nil
}

func (*FakeValidator) Process(
	block *block.Block,
	statedb *state.StateDB,
	cfg evm.Config,
) ([]*receipt.Receipt, []*log.Log, uint64, error) {
	return nil, nil, 0, nil
}
