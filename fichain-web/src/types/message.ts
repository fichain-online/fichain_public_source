export const Messages = {
  MessagePing: "ping",
  MessagePong: "pong",
  MessagePeerList: "peer_list",

  MessageHandshakeInit: "handshake_init",
  MessageHandshakeAck: "handshake_ack",
  MessageHandshakeConfirm: "handshake_confirm",

  MessageSendTransaction: "send_transaction",

  MessageCallSmartContract: "call_smart_contract",
  MessageCallResult: "call_result",

  MessageGetBalance: "get_balance",
  MessageBalance: "balance",

  MessageGetNonce: "get_nonce",
  MessageNonce: "nonce",

  MessageGetReceipt: "get_receipt",
  MessageReceipt: "receipt",

  MessageGetValidators: "get_validators",
  MessageValidators: "validator",

  MessageGetHeadBlock: "get_head_block",
  MessageHeadBlock: "head_block",

  MessageGetBlock: "get_block",
  MessageBlock: "block",

  MessageTxMined: "tx_mined"
} as const;

// Optional: You can also define a type for all possible message values
export type MessageType = typeof Messages[keyof typeof Messages];
