import { create } from 'zustand';

interface WalletStore {
  wallet: HDNodeWallet | null;
  socket: WebSocket | null;
  logs: string[];
  connected: boolean;

  handshakeInitData:  Uint8Array | null;

  connect: () => void;
  disconnect: () => void;

  addLog: (log: string) => void;
  setWallet: (wallet: HDNodeWallet) => void;
  initHandShake: () => void;
  confirmHandShake: (ackMsg: HandshakeAck) => void;
  sendMessage: (
    msgType: MessageType, 
    payload: Uint8Array,
  ) => void;
}
