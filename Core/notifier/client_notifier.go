package notifier

import (
	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/block_chain"
	"FichainCore/event"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	"FichainCore/params"
	"FichainCore/transaction"
)

type ClientNotifier struct {
	sender                 *message_sender.MessageSender
	chainEvent             chan event.ChainEvent
	chainEventSubscription event.Subscription
}

func NewClientNotifier(
	sender *message_sender.MessageSender,
) *ClientNotifier {
	return &ClientNotifier{
		sender:     sender,
		chainEvent: make(chan event.ChainEvent),
	}
}

func (cn *ClientNotifier) SubscribeChanEvent(bc *block_chain.BlockChain) {
	cn.chainEventSubscription = bc.SubscribeChainEvent(cn.chainEvent)
	go cn.HandleChanEvent()
}

func (cn *ClientNotifier) HandleChanEvent() {
	for {
		event := <-cn.chainEvent
		for _, v := range event.Block.Transactions {
			cn.NotifyTxMined(*v)
		}
		logger.Warn("TODO broad cast to subscribe event address")
	}
}

func (cn *ClientNotifier) NotifyTxMined(tx transaction.Transaction) {
	from, _ := tx.From(params.TempChainId)
	go cn.sender.SendMessageToAddress(
		from,
		message.MessageTxMined,
		&tx,
	)
	go cn.sender.SendMessageToAddress(
		tx.To(),
		message.MessageTxMined,
		&tx,
	)
}
