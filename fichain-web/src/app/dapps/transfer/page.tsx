'use client';

// --- MODIFIED: Added useMemo ---
import React, { useState, useEffect, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { ethers, keccak256, SigningKey, Signature, HDNodeWallet, hexlify, getBytes, formatEther, toBigInt, parseUnits } from 'ethers';

// --- Zustand Store Imports ---
import { useAuthStore } from '@/stores/authStore';
import { useWebSocketStore } from '@/stores/webSocketStore';

// --- UI and Icon Imports ---
import Header from '@/components/header';
// --- MODIFIED: Added Fuel and more icons for user lookup ---
import { ArrowLeft, Wallet, Coins, MessageSquare, Loader2, CheckCircle, XCircle, Fuel, UserCheck, UserX, Repeat } from 'lucide-react';
import Link from 'next/link';

import { create, toBinary, toJson, fromBinary } from "@bufbuild/protobuf";
import {
  TransactionSchema,
  TransactionHashDataSchema,
  TransactionSignDataSchema,
  type Transaction,
  type TransactionSignData,
  type TransactionHashData,

} from '@/proto/transaction_pb'
// --- NEW: Import the Receipt type from your generated protobuf file ---
// Note: Adjust the path if your receipt protobuf file is located elsewhere.
import { type Receipt } from '@/proto/receipt_pb';
import config from '@/lib/config';

import { generateTransactionWithSign } from '@/lib/utils';

interface TransferResponse {
  status: 'success' | 'error';
  message: string;
  txHash?: string;
}

// A standard ETH transfer uses 21,000 gas.
const GAS_LIMIT_FOR_TRANSFER = 21000;

// --- NEW: Fake API for User Lookup ---
// This simulates a database of known users. In a real app, this would be an API call.
const FAKE_USER_DB: Record<string, string> = {
  '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266': 'Satoshi Nakamoto',
  '0x70997970C51812dc3A010C7d01b50e0d17dc79C8': 'Vitalik Buterin',
  '0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC': 'Gavin Wood',
};

// This function simulates fetching a user's name from an address with a 1-second delay.
const fetchUserNameByAddress = (address: string): Promise<string | null> => {
  console.log(`Querying name for address: ${address}`);
  return new Promise((resolve) => {
    setTimeout(() => {
      const name = FAKE_USER_DB[address] || null;
      resolve(name);
    }, 1000); // 1-second delay to simulate network latency
  });
};


export default function TransferPage() {
  const router = useRouter();

  // --- State for the form ---
  const [recipient, setRecipient] = useState('');
  const [amount, setAmount] = useState('');
  const [memo, setMemo] = useState('');
  const [gasPrice, setGasPrice] = useState(0.1); 

  // --- NEW: State for recipient name lookup ---
  const [recipientName, setRecipientName] = useState<string | null>(null);
  const [isNameLoading, setIsNameLoading] = useState(false);

  // --- MODIFIED: Combined loading, error, and success/receipt state ---
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [receipt, setReceipt] = useState<Receipt | null>(null); // State to hold the transaction receipt

  // --- Get state from stores ---
  const { isAuthenticated, _hasHydrated, privateKey } = useAuthStore();
  const { wallet, sendMessage, sendTransaction, getNonce } = useWebSocketStore();

  // --- Auth check effect ---
  useEffect(() => {
    if (_hasHydrated && !isAuthenticated) {
      router.replace('/');
    }
  }, [_hasHydrated, isAuthenticated, router]);


  // --- NEW: Function to reset the form state for a new transaction ---
  const handleReset = () => {
    setRecipient('');
    setAmount('');
    setMemo('');
    setError(null);
    setReceipt(null); // Clear the receipt to show the form again
    setIsLoading(false);
  };

  // --- NEW: Effect for debounced user name lookup ---
  useEffect(() => {
    // Reset name and loading state immediately on change
    setRecipientName(null);

    // Basic validation: check for 'd-' prefix and if the rest is a valid address
    const isValidFormat = ethers.isAddress(recipient);

    if (!isValidFormat) {
      setIsNameLoading(false);
      return; // Exit if the address format is not valid
    }

    setIsNameLoading(true);

    // Debounce: set a timer to fetch the name after 500ms of inactivity
    const debounceTimer = setTimeout(() => {
      fetchUserNameByAddress(recipient).then((name) => {
        setRecipientName(name); // Set the name (or null if not found)
        setIsNameLoading(false); // Stop loading
      });
    }, 500); // 500ms delay

    // Cleanup: clear the timer if the user types again before it fires
    return () => clearTimeout(debounceTimer);
  }, [recipient]); // This effect runs whenever the 'recipient' state changes

  // --- Calculate the transfer fee using useMemo for efficiency ---
  const transferFee = useMemo(() => {
    try {
      const gasPriceWei = ethers.parseUnits(gasPrice.toString(), 'finney');
      const feeInWei = gasPriceWei * BigInt(GAS_LIMIT_FOR_TRANSFER);
      const feeAsNumber = Number(ethers.formatEther(feeInWei));
      return new Intl.NumberFormat('vi-VN').format(feeAsNumber)
    } catch (e) {
      return '0';
    }
  }, [gasPrice]);

  // --- A derived state for the formatted amount display ---
  const formattedAmount = useMemo(() => {
    if (amount === '') return '';
    try {
      return new Intl.NumberFormat('vi-VN').format(parseInt(amount, 10));
    } catch {
      return amount;
    }
  }, [amount]);

  // --- Handler to update the raw amount from the formatted input ---
  const handleAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const rawValue = e.target.value.replace(/[^0-9]/g, '');
    setAmount(rawValue);
  };

  const handleTransfer = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setReceipt(null);
    setIsLoading(true);

    // --- MODIFIED: Updated validation to match our custom 'd-' prefix format ---
    const isValidRecipient = ethers.isAddress(recipient.substring(2));
    if (!isValidRecipient) {
      setError('Địa chỉ người nhận không hợp lệ. Phải có định dạng "0x...".');
      setIsLoading(false);
      return;
    }

    if (isNaN(parseFloat(amount)) || parseFloat(amount) <= 0) {
      setError('Số tiền phải là một số lớn hơn 0.');
      setIsLoading(false);
      return;
    }
    if (!wallet || !sendMessage) {
      setError('Kết nối chưa sẵn sàng. Vui lòng thử lại sau.');
      setIsLoading(false);
      return;
    }

    // Multiply by 10^18 using ethers.utils.parseUnits
    const bigAmount = ethers.parseUnits(amount, 18); // returns a BigNumber
    // Convert BigNumber to bytes array (Uint8Array)
    const amountByteArray = ethers.toBeArray(bigAmount);   

    try {
      const nonce = await getNonce()
      console.log("nonce is ", nonce)
      const tx = generateTransactionWithSign(
        privateKey!,
        BigInt(nonce),
        getBytes(recipient),
        amountByteArray,
        new Uint8Array(),
        memo || 'Chuyển khoản',
        BigInt(GAS_LIMIT_FOR_TRANSFER),
        ethers.toBeArray(ethers.parseUnits(gasPrice.toString(), 'finney')),
      )

      const returnedReceipt: Receipt = await sendTransaction(tx);

      console.log("Transaction successful, receipt:", returnedReceipt);

      // On success, set the receipt object to state. This will trigger the UI to display it.
      setReceipt(returnedReceipt);
    } catch (err: any) {
      console.error('Transfer failed:', err);
      setError(err.message || 'Đã xảy ra lỗi không mong muốn.');
    } finally {
      setIsLoading(false);
    }
  };

  const calculatedFeeEther = useMemo(() => {
    // Nếu chưa có receipt, không có phí để hiển thị
    if (!receipt) return '0';

    try {
      // Lấy gasPrice mà người dùng đã chọn trên form (đơn vị 'finney')
      // và chuyển nó sang đơn vị Wei (dạng BigInt)
      const gasPriceInWei = parseUnits(gasPrice.toString(), 'finney');
      
      // Lấy gasUsed từ receipt (đã là BigInt)
      const gasUsed = receipt.gasUsed;

      // Tính phí bằng Wei
      const feeInWei = gasUsed * gasPriceInWei;

      // Chuyển phí từ Wei sang Ether (chia cho 10^18) để hiển thị
      return formatEther(feeInWei);
    } catch (e) {
      console.error("Error calculating fee:", e);
      return '0';
    }
  }, [receipt, gasPrice]); // Phụ thuộc vào receipt và gasPrice

  // Loading state
  if (!_hasHydrated || !isAuthenticated) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen text-white bg-gray-900">
        <Loader2 className="h-12 w-12 animate-spin text-cyan-400" />
        <p className="mt-4 text-lg">Đang xác thực, vui lòng chờ...</p>
      </div>
    );
  }

  return (
    <main className="min-h-screen text-white">
      <Header />
      <div className="container mx-auto px-4 py-12 sm:py-16">
        <div className="max-w-2xl mx-auto">

          {!receipt && (
            <Link
              href="/dapps"
              className="inline-flex items-center text-cyan-400 hover:text-cyan-300 mb-8 group p-2 rounded-md hover:bg-cyan-900/50 transition-colors -ml-2 z-[10]"
            >
              <ArrowLeft className="h-5 w-5 mr-2 transition-transform group-hover:-translate-x-1" />
              Quay lại danh sách dApps
            </Link>
          )}
          <div className="bg-gray-800 p-8 rounded-xl shadow-2xl">
            {/* --- NEW: Conditional Rendering for Receipt vs. Form --- */}
            {receipt ? (
              // --- NEW: ENHANCED RECEIPT DISPLAY ---
              <div>
                <div className="flex flex-col items-center text-center">
                  <CheckCircle className="h-16 w-16 text-green-400 mb-4" />
                  <h1 className="text-2xl font-bold">Giao dịch thành công</h1>
                  <p className="text-gray-400 mt-1">Giao dịch của bạn đã được xác nhận trên blockchain.</p>
                </div>

                {/* Main Transaction Details */}
                <div className="mt-8 space-y-4 text-sm">
                  <div className="flex justify-between items-center bg-gray-900/50 p-4 rounded-lg">
                    <span className="text-gray-400">Số tiền:</span>
                    <span className="font-bold text-lg text-cyan-400">
                      {/* ✅ Đã sửa: Chuyển amount từ bytes -> BigInt -> Ether và định dạng */}
                      {new Intl.NumberFormat('vi-VN').format(
                        parseFloat(formatEther(toBigInt(receipt.amount)))
                      )} VNĐ
                    </span>
                  </div>
                  
                  <div className="space-y-3 pt-2">
                     <div className="flex justify-between items-start">
                      <span className="text-gray-400 shrink-0 mr-4">Từ:</span>
                      <span className="font-mono break-all text-right">{wallet?.address}</span>
                    </div>
                    <div className="flex justify-between items-start">
                      <span className="text-gray-400 shrink-0 mr-4">Đến:</span>
                       <div className="text-right">
                        <span className="font-mono break-all">{hexlify(receipt.to)}</span>
                        {/* If we found a name for the recipient, show it! */}
                        {recipientName && (
                          <span className="block text-xs text-green-400">({recipientName})</span>
                        )}
                      </div>
                    </div>
                    {/* Only show the memo if one was provided */}
                    {memo && (
                      <div className="flex justify-between items-start">
                        <span className="text-gray-400 shrink-0 mr-4">Ghi chú:</span>
                        <span className="text-right italic">"{memo}"</span>
                      </div>
                    )}
                  </div>
                </div>

                {/* Divider */}
                <hr className="my-6 border-gray-700" />
                
                {/* Technical Details */}
                <div className="space-y-3 text-sm">
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">Trạng thái:</span>
                    {receipt.status === 1 ? (
                       <span className="font-medium text-green-400 bg-green-900/50 px-2 py-1 rounded-md">Thành Công</span>
                    ) : (
                       <span className="font-medium text-red-400 bg-red-900/50 px-2 py-1 rounded-md">Thất Bại</span>
                    )}
                  </div>
                  <div className="flex justify-between items-start">
                    <span className="text-gray-400 shrink-0 mr-4">Mã giao dịch:</span>
                    <span className="font-mono break-all text-right text-cyan-400">{hexlify(receipt.txHash)}</span>
                  </div>
                   <div className="flex justify-between items-center">
                    <span className="text-gray-400">Số khối:</span>
                    <span className="font-mono">{receipt.blockNumber.toString()}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">Phí giao dịch:</span>
                    <span className="font-mono">
                      {new Intl.NumberFormat('vi-VN', {
                        minimumFractionDigits: 2, // Hiển thị ít nhất 6 chữ số thập phân
                        maximumFractionDigits: 2, // Tối đa 8 chữ số thập phân
                      }).format(parseFloat(calculatedFeeEther))} VNĐ
                    </span>
                  </div>
                </div>

                <button
                  onClick={handleReset}
                  className="mt-8 w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-cyan-600 hover:bg-cyan-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-cyan-500 transition-colors"
                >
                  <Repeat className="h-5 w-5 mr-2" />
                  Thực hiện giao dịch khác
                </button>
              </div>
            ) : (
                //* --- TRANSFER FORM --- */}
                <>
                  <h1 className="text-3xl font-bold text-center mb-2">Chuyển Khoản</h1>
                  <p className="text-center text-gray-400 mb-8">Gửi tài sản đến một địa chỉ ví khác một cách an toàn.</p>

                  <form onSubmit={handleTransfer} className="space-y-6">

                    {/* Recipient Address Input */}
                    <div>
                      <label htmlFor="recipient" className="block text-sm font-medium text-gray-300 mb-2">Địa chỉ người nhận</label>
                      <div className="relative">
                        <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                          <Wallet className="h-5 w-5 text-gray-500" />
                        </div>
                        <input
                          type="text"
                          id="recipient"
                          value={recipient}
                          onChange={(e) => setRecipient(e.target.value)}
                          className="w-full bg-gray-700 border-gray-600 rounded-md py-3 pl-10 pr-4 text-white placeholder-gray-500 focus:ring-cyan-500 focus:border-cyan-500"
                          placeholder="0x..."
                          required
                        />
                      </div>

                      {/* --- NEW: Recipient Name Lookup UI --- */}
                      <div className="h-6 mt-2 px-1 text-sm flex items-center">
                        {isNameLoading && (
                          <div className="flex items-center text-gray-400">
                            <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                            <span>Đang tìm kiếm tên người nhận...</span>
                          </div>
                        )}
                        {!isNameLoading && recipientName && (
                          <div className="flex items-center text-green-400">
                            <UserCheck className="h-4 w-4 mr-2 flex-shrink-0" />
                            <span>Người nhận: <span className="font-semibold">{recipientName}</span></span>
                          </div>
                        )}
                        {/* Show 'not found' only when loading is finished, name is null, and the format is valid */}
                        {!isNameLoading && !recipientName && ethers.isAddress(recipient) && (
                          <div className="flex items-center text-yellow-500">
                            <UserX className="h-4 w-4 mr-2 flex-shrink-0" />
                            <span>Không tìm thấy tên cho địa chỉ này.</span>
                          </div>
                        )}
                      </div>
                    </div>

                    {/* Amount Input */}
                    <div>
                      <label htmlFor="amount" className="block text-sm font-medium text-gray-300 mb-2">Số Tiền</label>
                      <div className="relative">
                        <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                          <Coins className="h-5 w-5 text-gray-500" />
                        </div>
                        <input
                          type="text"
                          inputMode="decimal"
                          id="amount"
                          value={formattedAmount}
                          onChange={handleAmountChange}
                          className="w-full bg-gray-700 border-gray-600 rounded-md py-3 pl-10 pr-4 text-white placeholder-gray-500 focus:ring-cyan-500 focus:border-cyan-500"
                          placeholder="1.000.000"
                          required
                        />
                      </div>
                    </div>

                    {/* Gas Price Slider and Fee Display */}
                    <div>
                      <label htmlFor="gasPrice" className="flex items-center text-sm font-medium text-gray-300 mb-2">
                        <Fuel className="h-4 w-4 mr-2" />
                        Phí Gas (finney)
                      </label>
                      <div className="flex items-center space-x-4">
                        <input
                          type="range"
                          id="gasPrice"
                          min="0.1"
                          max="0.3"
                          step="0.1"
                          value={gasPrice}
                          onChange={(e) => setGasPrice(Number(e.target.value))}
                          className="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-cyan-500"
                        />
                        <input 
                          type="number"
                          value={gasPrice}
                          onChange={(e) => setGasPrice(Number(e.target.value))}
                          className="w-20 bg-gray-700 border-gray-600 rounded-md py-1 px-2 text-white text-center focus:ring-cyan-500 focus:border-cyan-500"
                          min="0.1"
                          step="0.1"
                          max="0.3"
                        />
                      </div>
                      <p className="text-xs text-gray-400 mt-2">
                        Phí giao dịch dự kiến: <span className="font-medium text-cyan-400">{transferFee} VNĐ</span>
                      </p>
                    </div>

                    {/* Memo Input */}
                    <div>
                      <label htmlFor="memo" className="block text-sm font-medium text-gray-300 mb-2">Ghi chú (Tùy chọn)</label>
                      <div className="relative">
                        <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                          <MessageSquare className="h-5 w-5 text-gray-500" />
                        </div>
                        <input
                          type="text"
                          id="memo"
                          value={memo}
                          onChange={(e) => setMemo(e.target.value)}
                          className="w-full bg-gray-700 border-gray-600 rounded-md py-3 pl-10 pr-4 text-white placeholder-gray-500 focus:ring-cyan-500 focus:border-cyan-500"
                          placeholder="VD: Trả tiền cà phê"
                        />
                      </div>
                    </div>

                    {/* Feedback Messages */}
                    {error && (
                      <div className="flex items-center p-4 text-sm text-red-300 bg-red-900/30 rounded-lg">
                        <XCircle className="h-5 w-5 mr-3 flex-shrink-0" />
                        <span>{error}</span>
                      </div>
                    )}

                    {/* Submit Button */}
                    <div>
                      <button
                        type="submit"
                        disabled={isLoading}
                        className="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-cyan-600 hover:bg-cyan-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-gray-800 focus:ring-cyan-500 disabled:bg-gray-600 disabled:cursor-not-allowed transition-colors"
                      >
                        {isLoading ? (
                          <>
                            <Loader2 className="h-5 w-5 mr-2 animate-spin" />
                            Đang xử lý...
                          </>
                        ) : (
                            'Gửi'
                          )}
                      </button>
                    </div>
                  </form>
                </>
              )}
          </div>
        </div>
      </div>
    </main>
  );
}
