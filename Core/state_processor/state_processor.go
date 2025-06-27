package state_processor

import (
	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/block"
	"FichainCore/bloom"
	"FichainCore/common"
	"FichainCore/consensus"
	"FichainCore/crypto"
	"FichainCore/evm"
	"FichainCore/gas_pool"
	"FichainCore/log"
	"FichainCore/params"
	"FichainCore/receipt"
	"FichainCore/state"
	"FichainCore/transaction"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     evm.ChainContext    // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(
	config *params.ChainConfig,
	bc evm.ChainContext,
	engine consensus.Engine,
) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		// engine: engine,
	}
}

// Process processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(
	block *block.Block,
	statedb *state.StateDB,
	cfg evm.Config,
) ([]*receipt.Receipt, []*log.Log, uint64, error) {
	var (
		receipts []*receipt.Receipt
		usedGas  = new(uint64)
		header   = block.Header
		allLogs  []*log.Log
		gp       = new(gas_pool.GasPool).AddGas(block.GasLimit())
	)
	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions {
		statedb.Prepare(tx.Hash(), block.Header.ParentHash, i)
		receipt, _, err := ApplyTransaction(
			p.config,
			p.bc,
			&block.Header.Proposer,
			gp,
			statedb,
			header,
			tx,
			usedGas,
			cfg,
		)
		if err != nil {
			return nil, nil, 0, err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}

	// TODO
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	// p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), receipts)

	return receipts, allLogs, *usedGas, nil
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(
	config *params.ChainConfig,
	bc evm.ChainContext,
	author *common.Address,
	gp *gas_pool.GasPool,
	statedb *state.StateDB,
	header *block.BlockHeader,
	tx *transaction.Transaction,
	usedGas *uint64,
	cfg evm.Config,
) (*receipt.Receipt, uint64, error) {
	// Create a new context to be used in the EVM environment
	logger.DebugP("usedGas 1", *usedGas)
	from, err := tx.From(params.TempChainId)
	if err != nil {
		return nil, 0, err
	}
	context := evm.NewEVMContext(from, tx.GasPrice(), header, bc, *author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := evm.NewEVM(context, statedb, config, cfg)
	// Apply the transaction to the current state (included in the env)
	_, gas, failed, err := ApplyMessage(vmenv, tx, gp)
	if err != nil {
		return nil, 0, err
	}
	logger.DebugP("usedGas 2", gas)
	// Update the state with pending changes
	var root []byte
	statedb.Finalise(true)
	*usedGas += gas
	logger.DebugP("usedGas 3", *usedGas)
	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing wether the root touch-delete accounts.
	rc := receipt.NewReceipt(root, failed, *usedGas)
	rc.TxHash = tx.Hash()
	rc.GasUsed = gas
	// if the transaction created a contract, store the creation address in the receipt.
	if (tx.To() == common.Address{}) {
		rc.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
		logger.DebugP("Deploy contract address", rc.ContractAddress.String())
	}
	// Set the receipt logs and create a bloom for filtering
	rc.Logs = statedb.GetLogs(tx.Hash())
	rc.LogsBloom = bloom.CreateBloom([]*receipt.Receipt{rc}).Bytes()

	return rc, gas, err
}
