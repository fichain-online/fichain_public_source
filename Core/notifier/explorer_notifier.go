package notifier

import (
	"FichainCore/block_chain"
	"FichainCore/config"
	"FichainCore/event"
	"FichainCore/p2p/lookup_table"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	pb "FichainCore/proto"
	"FichainCore/types"
)

type ExplorerNotifier struct {
	sender                 *message_sender.MessageSender
	chainEvent             chan event.ChainEvent
	chainEventSubscription event.Subscription
	lookupTable            *lookup_table.LookupTable
}

func NewExplorerNotifier(
	sender *message_sender.MessageSender,
	lookupTable *lookup_table.LookupTable,
) *ExplorerNotifier {
	return &ExplorerNotifier{
		sender:      sender,
		chainEvent:  make(chan event.ChainEvent),
		lookupTable: lookupTable,
	}
}

func (cn *ExplorerNotifier) SubscribeChanEvent(bc *block_chain.BlockChain) {
	cn.chainEventSubscription = bc.SubscribeChainEvent(cn.chainEvent)
	go cn.HandleChanEvent()
}

func (cn *ExplorerNotifier) HandleChanEvent() {
	for {
		event := <-cn.chainEvent
		haveExpl := false
		for _, v := range config.GetConfig().GetExplorerAddresses() {
			if cn.lookupTable.Has(v) {
				haveExpl = true
				break
			}
		}
		if !haveExpl {
			continue
		}
		cn.NotifyChanEvent(event)
	}
}

func (cn *ExplorerNotifier) NotifyChanEvent(event event.ChainEvent) {
	pbChainEvent := &pb.ChainEvent{
		Block: event.Block.Proto(),
	}
	for _, v := range event.Logs {
		pbChainEvent.Logs = append(pbChainEvent.Logs, v.Proto())
	}
	w := &types.PbChainEventWrap{PbEvent: pbChainEvent}
	for _, v := range config.GetConfig().GetExplorerAddresses() {
		if cn.lookupTable.Has(v) {
			go cn.sender.SendMessageToAddress(v, message.MessageChainEvent, w)
		}
	}
}
