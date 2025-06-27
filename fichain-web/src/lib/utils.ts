import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

import { ByteArray } from "viem";
import { create, toBinary } from "@bufbuild/protobuf";
import { ethers, keccak256, SigningKey, hexlify, getBytes, formatEther, parseUnits, Interface, toBigInt } from 'ethers';
import config from '@/lib/config';

import {
  CallSmartContractDataSchema,
  CallSmartContractHashDataSchema,
  CallSmartContractResponseSchema,
  type CallSmartContractData, 
  type CallSmartContractHashData, 
  type CallSmartContractResponse, 
} from '@/proto/call_data_pb';

import {
  TransactionSchema,
  TransactionHashDataSchema,
  TransactionSignDataSchema,
  type Transaction,
  type TransactionSignData,
  type TransactionHashData,

} from '@/proto/transaction_pb'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const generateCallDataWithSign = (
  privateKey: string,
  to: Uint8Array, 
  data: Uint8Array,
) => {
  const callSignData = create(CallSmartContractHashDataSchema, { 
    ToAddress: to,
    Data: data,
  });
  const bSignData = toBinary(CallSmartContractHashDataSchema, callSignData)
  const hashSignData = keccak256(bSignData)
  let fmtPrivateKey = privateKey!
  if (!fmtPrivateKey.startsWith('0x')) {
    fmtPrivateKey = '0x' + fmtPrivateKey;
  }
  const signingKey = new SigningKey(fmtPrivateKey);
  const signature = signingKey.sign(hashSignData);
  // console.log("Hash sign data", hashSignData)
  const signatureBytes = ethers.getBytes(ethers.concat([signature.r, signature.s, signature.yParity === 0 ? '0x00' : '0x01']));
  const callRq = create(CallSmartContractDataSchema, { 
    ToAddress: to,
    Data: data,
    Sign: signatureBytes,
  });

  return callRq
} 

export const generateTransactionWithSign = (
  privateKey: string,
  nonce: bigint, 
  to: Uint8Array, 
  amount: Uint8Array, 
  data: Uint8Array,
  message: string,
  gas: bigint,
  gasPrice: Uint8Array,
) => {

  const txSignData = create(TransactionSignDataSchema, {
    ToAddress: to,
    Nonce: nonce,
    Amount: amount,
    Data: data,
    Gas: gas,
    GasPrice: gasPrice,
    Message: message,
    ChainId: ethers.toBeArray(BigInt(config.blockchainId)),
  })

  // calculate hash
  const bTxSignData = toBinary(TransactionSignDataSchema, txSignData)
  const hashSignData = keccak256(bTxSignData)
  console.log("bTxSignData",hexlify(bTxSignData))
  // ✅ FIX: Ensure the private key has a "0x" prefix
  //
  let fmtPrivateKey = privateKey!
  if (!fmtPrivateKey.startsWith('0x')) {
    fmtPrivateKey = '0x' + fmtPrivateKey;
  }
  const signingKey = new SigningKey(fmtPrivateKey);
  const signature = signingKey.sign(hashSignData);
  console.log("Hash sign data", hashSignData)

  const signatureBytes = ethers.getBytes(ethers.concat([signature.r, signature.s, signature.yParity === 0 ? '0x00' : '0x01']));

  const tx = create(TransactionSchema, {
    ToAddress: to,
    Nonce: nonce,
    Amount: amount,
    Data: data,
    Gas: gas,
    GasPrice: gasPrice,
    Message: message,
    Sign: signatureBytes,
  })

  return tx; 
}


  export const formatAndShorten = (balance: bigint, decimals = 4) => {
    try {
      const formatted = formatEther(balance);
      const num = parseFloat(formatted);
      // if (num > 0 && num < 0.0001) return num.toExponential(2);
      return new Intl.NumberFormat('vi-VN', { maximumFractionDigits: decimals }).format(num);
    } catch {
      return '0';
    }
  };
