import { create } from 'zustand';
import config from '@/lib/config';
import {
  MessageSchema,
  MessageHeaderSchema,
  HandshakeInitSchema,
  HandshakeAckSchema,
  HandshakeConfirmSchema,
  BytesMessageSchema,
  type Message,
  type MessageHeader,
  type HandshakeInit,
  type HandshakeAck,
  type HandshakeConfirm,
  type BytesMessage,
} from '@/proto/message_pb';
import {
  TransactionSchema,
  type Transaction
} from '@/proto/transaction_pb';
import {
  ReceiptSchema,
  type Receipt, 
} from '@/proto/receipt_pb';
import {
  CallSmartContractDataSchema,
  CallSmartContractResponseSchema,
  type CallSmartContractData, 
  type CallSmartContractResponse, 
} from '@/proto/call_data_pb';

import { create as createProto, toBinary, fromBinary } from "@bufbuild/protobuf";
import { ethers, Wallet, SigningKey, HDNodeWallet, hexlify, zeroPadBytes, toBigInt } from 'ethers';
import { v4 as uuidv4 } from 'uuid';
import { MessageType, Messages } from '@/types/message';
import { bytesToBigInt } from 'viem';


function bytesToUint64(bytes: Uint8Array): bigint {
  if (bytes.length > 8) {
    throw new Error("Input exceeds 8 bytes for uint64");
  }
  const padded = zeroPadBytes(bytes, 8);
  const hex = hexlify(padded);
  return toBigInt(hex);
}

interface PendingBalanceRequest {
  resolve: (value: bigint | PromiseLike<bigint>) => void;
  reject: (reason?: any) => void;
}

interface PendingNonceRequest {
  resolve: (value: number | PromiseLike<number>) => void;
  reject: (reason?: any) => void;
}

interface PendingTransactionReceipt {
  resolve: (value: Receipt | PromiseLike<Receipt>) => void;
  reject: (reason?: any) => void;
}

interface PendingCallSmartContractResponse {
  resolve: (value: CallSmartContractResponse | PromiseLike<CallSmartContractResponse>) => void;
  reject: (reason?: any) => void;
}


interface WebSocketStore {
  wallet: Wallet | null;
  socket: WebSocket | null;
  logs: string[];
  connected: boolean;

  pendingBalanceRequest: PendingBalanceRequest | null;
  pendingNonceRequest: PendingNonceRequest | null;
  pendingTransactionReceipt: PendingTransactionReceipt | null;
  pendingCallSmartContractResponse: PendingCallSmartContractResponse | null; // may update to map

  connect: () => void;
  disconnect: () => void;
  addLog: (log: string) => void;
  setWallet: (wallet: Wallet) => void;
  initHandShake: () => void;
  confirmHandShake: (ackMsg: HandshakeAck) => void;
  sendMessage: (msgType: MessageType, payload: Uint8Array) => void;
  sendTransaction: (tx: Transaction) => Promise<Receipt>;
  callSmartContract: (callData: CallSmartContractData) => Promise<CallSmartContractResponse>;
  getNonce: () => Promise<number>;

  getBalance: () => Promise<bigint>;
}

export const useWebSocketStore = create<WebSocketStore>((set, get) => ({
  socket: null,
  wallet: null,
  logs: [],
  connected: false,
  pendingBalanceRequest: null,
  pendingNonceRequest: null,
  pendingTransactionReceipt: null, 
  pendingCallSmartContractResponse: null,

  connect: () => {
    if (get().socket) return;

    const ws = new WebSocket(config.wsBaseUrl);
    ws.binaryType = "arraybuffer";

    ws.onopen = () => {
      get().addLog('✅ WebSocket connected');
      get().initHandShake();
    };

    ws.onclose = () => {
      get().addLog('❌ WebSocket disconnected');
      set({ socket: null, connected: false });
    };

    ws.onerror = (err) => {
      console.error('⚠️ WebSocket error:', err);
    };

    ws.onmessage = (event) => {
      const message = fromBinary(MessageSchema, new Uint8Array(event.data));
      get().addLog(`📩 Received message: ${message.header?.messageType}`);
      if (message.header != null) {
        switch (message.header.messageType) {
          case Messages.MessageHandshakeAck:
            const ackMessage = fromBinary(HandshakeAckSchema, new Uint8Array(message.payload));
            get().confirmHandShake(ackMessage);
            break;
          case Messages.MessageNonce:
            const nonceMessage = fromBinary(BytesMessageSchema, message.payload);
            const pendingRequest = get().pendingNonceRequest;
            if (pendingRequest) {
              const nonce = bytesToUint64(nonceMessage.data);
              pendingRequest.resolve(Number(nonce));
              // ✅ CORRECT: Update state using the 'set' function
              set({ pendingNonceRequest: null });
            } else {
              get().addLog(`  Received nonce for unknown request`);
            }
            break;
          case Messages.MessageBalance:
            const balanceMessage = fromBinary(BytesMessageSchema, message.payload);
            const pendingBalanceRequest = get().pendingBalanceRequest;
            if (pendingBalanceRequest) {
              const balance = bytesToBigInt(balanceMessage.data);
              pendingBalanceRequest.resolve(balance);
              set({ pendingBalanceRequest: null });
            } else {
              get().addLog(`  Received balance for unknown request`);
            }
            break;
          case Messages.MessageTxMined:
            get().addLog(`📩 Received : ${message.header?.messageType}`);
            // let query receipt
            const getReceipt = createProto(BytesMessageSchema, {
              data:  message.payload,
            });
            const bData = toBinary(BytesMessageSchema, getReceipt);
            get().sendMessage(Messages.MessageGetReceipt, bData);
            break;
          case Messages.MessageReceipt:
            //
            const receiptMessage = fromBinary(ReceiptSchema, message.payload);
            const pendingTransactionReceipt = get().pendingTransactionReceipt;
            if (pendingTransactionReceipt) {
              pendingTransactionReceipt.resolve(receiptMessage);
              // ✅ CORRECT: Update state using the 'set' function
              set({ pendingTransactionReceipt: null });
            } else {
              get().addLog(`  Received receipt for unknown transaction`);
            }
          case Messages.MessageCallResult:
            const result = fromBinary(CallSmartContractResponseSchema, message.payload)
            const pending = get().pendingCallSmartContractResponse
            if (pending) {
              pending.resolve(result);
              set({ pendingCallSmartContractResponse: null });
            } else {
              get().addLog(`  Received call result for unknown call smart contract`);
            }
          default:
            get().addLog(`📩 Received unknown message: ${message.header?.messageType}`);
        }
      }
    };

    set({ socket: ws });
  },

  disconnect: () => {
    const ws = get().socket;
    if (ws) {
      ws.close();
      set({ socket: null });
    }
  },

  addLog: (log: string) => {
    set(state => ({
      logs: [`[${new Date().toLocaleTimeString()}] ${log}`, ...state.logs.slice(0, 99)]
    }));
    console.log(log);
  },

  initHandShake: () => {
    const wallet = get().wallet!;
    const initPayloadString = `{"time":${Math.floor(Date.now() / 1000)},"uuid":"${uuidv4()}"}`;
    const initPayloadBytes = ethers.toUtf8Bytes(initPayloadString);
    const handshakeInitData = createProto(HandshakeInitSchema, {
      walletAddress: ethers.getBytes(wallet.address),
      payload: initPayloadBytes,
    });
    const bHandshakeData = toBinary(HandshakeInitSchema, handshakeInitData);
    get().sendMessage(Messages.MessageHandshakeInit, bHandshakeData);
  },

  confirmHandShake: (ackMsg: HandshakeAck) => {
    const digestToSign = ethers.keccak256(ackMsg.payload);
    get().addLog(`Signing hash ${digestToSign}`);
    try {
      const signingKey = new SigningKey(get().wallet!.privateKey);
      const signature = signingKey.sign(digestToSign);
      const signatureBytes = ethers.getBytes(ethers.concat([signature.r, signature.s, signature.yParity === 0 ? '0x00' : '0x01']));
      const handShakeConfirm = createProto(HandshakeConfirmSchema, {
        signature: signatureBytes
      });
      const bData = toBinary(HandshakeConfirmSchema, handShakeConfirm);
      get().sendMessage(Messages.MessageHandshakeConfirm, bData);
      get().addLog(`Handshake complete`);
    } catch (e: any) {
      get().addLog(`Crit Err: ${e.message}.`);
    }
    set({connected: true})
  },

  sendMessage: async (msgType: MessageType, payload: Uint8Array) => {
    const ws = get().socket!;
    const wallet = get().wallet!;
    const headerData = createProto(MessageHeaderSchema, {
      version: config.blockchainVersion,
      senderId: await wallet.getAddress(),
      messageType: msgType,
      timestamp: BigInt(Math.floor(Date.now() / 1000)),
      signature: new Uint8Array(),
    });
    const msg = createProto(MessageSchema, {
      header: headerData,
      payload: payload,
    });
    const bMsg = toBinary(MessageSchema, msg);
    ws.send(bMsg);
  },

  setWallet: (wallet: Wallet) => {
    set({ wallet });
  },

  getNonce: () => {
    return new Promise<number>((resolve, reject) => {
      const { socket, connected, addLog, sendMessage, wallet } = get();

      if (!socket || !connected || !wallet) {
        return reject(new Error('WebSocket is not connected or wallet is not set.'));
      }

      // Prevent making a new request if one is already pending
      if (get().pendingNonceRequest) {
        return reject(new Error('A nonce request is already in progress.'));
      }

      const timeoutId = setTimeout(() => {
        // ✅ CORRECT: Update state using 'set' on timeout
        set({ pendingNonceRequest: null });
        reject(new Error(`Nonce request timed out after 10 seconds.`));
      }, 10000);

      // ✅ CORRECT: Use 'set' to store the pending request object
      set({
        pendingNonceRequest: {
          resolve: (value) => {
            clearTimeout(timeoutId);
            resolve(value);
          },
          reject: (reason) => {
            clearTimeout(timeoutId);
            reject(reason);
          }
        }
      });

      const bAddress = ethers.getBytes(wallet.address);
      const getNoncePayload = createProto(BytesMessageSchema, { data: bAddress });
      const binaryPayload = toBinary(BytesMessageSchema, getNoncePayload);

      sendMessage(Messages.MessageGetNonce, binaryPayload);
    });
  },

  getBalance: () => {
    return new Promise<bigint>((resolve, reject) => {
      const { socket, connected, sendMessage, wallet } = get();

      if (!socket || !connected || !wallet) {
        return reject(new Error('WebSocket is not connected or wallet is not set.'));
      }

      // Prevent making a new request if one is already pending
      if (get().pendingBalanceRequest) {
        get().pendingBalanceRequest?.resolve(BigInt(0))
      }

      const timeoutId = setTimeout(() => {
        // ✅ CORRECT: Update state using 'set' on timeout
        set({ pendingBalanceRequest: null });
        reject(new Error(`get balance request timed out after 10 seconds.`));
      }, 10000);

      // ✅ CORRECT: Use 'set' to store the pending request object
      set({
        pendingBalanceRequest: {
          resolve: (value) => {
            clearTimeout(timeoutId);
            resolve(value);
          },
          reject: (reason) => {
            clearTimeout(timeoutId);
            reject(reason);
          }
        }
      });

      sendMessage(Messages.MessageGetBalance, new(Uint8Array));
    });
  },

  sendTransaction: (tx: Transaction) => {
    return new Promise<Receipt>((resolve, reject) => {
      const bData = toBinary(TransactionSchema, tx);
      get().sendMessage(Messages.MessageSendTransaction, bData);

      const timeoutId = setTimeout(() => {
        set({ pendingTransactionReceipt: null });
        reject(new Error(`Transaction request timed out after 30 seconds.`));
      }, 30000);

      set({
        pendingTransactionReceipt: {
          resolve: (value) => {
            clearTimeout(timeoutId);
            resolve(value);
          },
          reject: (reason) => {
            clearTimeout(timeoutId);
            reject(reason);
          }
        }
      });
    });
  },

  callSmartContract: (callData: CallSmartContractData) => {
    return new Promise<CallSmartContractResponse>((resolve, reject) => {
      const bData = toBinary(CallSmartContractDataSchema, callData);
      get().sendMessage(Messages.MessageCallSmartContract, bData);

      const timeoutId = setTimeout(() => {
        set({ pendingCallSmartContractResponse: null });
        reject(new Error(`Transaction request timed out after 30 seconds.`));
      }, 30000);

      set({
        pendingCallSmartContractResponse: {
          resolve: (value) => {
            clearTimeout(timeoutId);
            resolve(value);
          },
          reject: (reason) => {
            clearTimeout(timeoutId);
            reject(reason);
          }
        }
      });
    });
  }
}));
