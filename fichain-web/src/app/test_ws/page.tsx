'use client'
import { useState, useEffect, useRef, FC, ChangeEvent, KeyboardEvent } from 'react';
import { ethers, Wallet, SigningKey, Signature, HDNodeWallet } from 'ethers';
import { v4 as uuidv4 } from 'uuid';

import { create, toBinary, toJson } from "@bufbuild/protobuf";
import { useWebSocketStore } from '@/stores/webSocketStore';
import { Button } from "@/components/ui/button"
import config from '@/lib/config';
import { Messages } from '@/types/message';


// Create a registry if your messages use Any or extensions, otherwise often not strictly needed for basic messages
// const registry = createRegistry(MessageSchema, MessageHeaderSchema, HandshakeInitSchema, /* ... other schemas */);


const WebSocketTestPage: FC = () => {
  const { connect, addLog, logs, setWallet, wallet, sendMessage, connected } = useWebSocketStore();

  useEffect(() => {
    const randomWallet = Wallet.createRandom();
    setWallet(randomWallet);
    addLog(`Client Wallet Generated (for demo): Address: ${randomWallet.address}`);
  }, []);

  const connectWebSocket = () => {
    connect()
  };
  const sendPing = () => {
    sendMessage(Messages.MessagePing, new Uint8Array([])) 
  }

  return ( /* ... JSX remains the same ... */
    <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif' }}>
      <h1>WebSocket Client (protoc-gen-es Schemas)</h1>
      <p>Server URL: {config.wsBaseUrl}</p>
      {wallet && <p>Client Wallet: {wallet.address}</p>}
      <div>
        {!connected? (
          <Button onClick={connectWebSocket} disabled={!wallet}>Connect</Button>
        ) : (
          <Button onClick={() => alert("TODO")}>Disconnect</Button>
        )}
      </div>
      <p>Status: {connected? 'Connected' : 'Disconnected'}</p>
      {connected && (
        <div style={{ marginTop: '20px' }}>
          <button onClick={sendPing} style={{ padding: '8px 15px' }}>Send Ping</button>
        </div>
      )}
      <div style={{ marginTop: '20px', border: '1px solid #ccc', padding: '10px', height: '300px', overflowY: 'scroll', }}>
        <h2>Logs:</h2>
        {logs.map((log, index) => (
          <div key={index} style={{ fontSize: '0.9em', borderBottom: '1px dashed #eee', padding: '2px 0' }}>{log}</div>
        ))}
      </div>
    </div>
  );
};

export default WebSocketTestPage;
