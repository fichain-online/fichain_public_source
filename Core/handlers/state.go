package handlers

import (
	"encoding/binary"
	"math/big"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/block_chain"
	"FichainCore/call_data"
	"FichainCore/common"
	"FichainCore/evm"
	"FichainCore/p2p"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	"FichainCore/params"
	"FichainCore/state"
)

type StateHandler struct {
	stateDB *state.StateDB

	sender *message_sender.MessageSender
	bc     *block_chain.BlockChain
}

func NewStateHandler(
	stateDB *state.StateDB,
	sender *message_sender.MessageSender,
	bc *block_chain.BlockChain,
) *StateHandler {
	return &StateHandler{
		stateDB: stateDB,
		sender:  sender,
		bc:      bc,
	}
}

func (h *StateHandler) Handlers() map[string]func(p2p.Peer, *message.Message) error {
	return map[string]func(p2p.Peer, *message.Message) error{
		message.MessageGetBalance:        h.GetBalance,
		message.MessageGetNonce:          h.GetNonce,
		message.MessageCallSmartContract: h.CallSmartContract,
	}
}

func (h *StateHandler) GetBalance(peer p2p.Peer, msg *message.Message) error {
	address := peer.WalletAddress()
	balance := h.stateDB.GetBalance(address)

	err := h.sender.SendMessageToPeer(
		peer,
		message.MessageBalance,
		&message.BytesMessage{
			Data: balance.Bytes(),
		},
	)
	if err != nil {
		logger.Warn("error when send balance to peer", err)
	}
	return err
}

func (h *StateHandler) GetNonce(peer p2p.Peer, msg *message.Message) error {
	//
	address := common.BytesToAddress(msg.Payload.(*message.BytesMessage).Data)
	nonce := h.stateDB.GetNonce(address)
	logger.Debug("Nonce for ", address.Hex(), "is", nonce)

	bytes := make([]byte, 8) // uint64 takes 8 bytes
	binary.BigEndian.PutUint64(bytes, nonce)

	err := h.sender.SendMessageToPeer(
		peer,
		message.MessageNonce,
		&message.BytesMessage{
			Data: bytes,
		},
	)
	if err != nil {
		logger.Warn("error when send nonce to peer", err)
	}
	return err
}

func (h *StateHandler) CallSmartContract(peer p2p.Peer, msg *message.Message) error {
	//
	callData := msg.Payload.(*call_data.CallSmartContractData)
	from, err := callData.From()
	if err != nil {
		return err
	}
	header := h.bc.CurrentHeader()
	context := evm.NewEVMContext(from, big.NewInt(0), header, h.bc, h.bc.CurrentHeader().Proposer)
	vmenv := evm.NewEVM(context, h.stateDB, h.bc.Config(), evm.Config{})
	logger.DebugP("To address", callData.To().String())
	ret, _, err := vmenv.StaticCall(
		evm.AccountRef(from),
		callData.To(),
		callData.Data(),
		params.TempGasLimit,
	)
	if err != nil {
		logger.Debug("Error when call sc", err)
	}
	hash, err := callData.HashSign()
	if err != nil {
		logger.Error("error when get hashsign from call data", err)
		return err
	}
	res := &call_data.CallSmartContractResponse{
		Hash: hash,
		Data: ret,
	}
	logger.DebugP("call response", res)

	err = h.sender.SendMessageToPeer(
		peer,
		message.MessageCallResult,
		res,
	)
	if err != nil {
		logger.Warn("error when send nonce to peer", err)
	}

	return err
}
