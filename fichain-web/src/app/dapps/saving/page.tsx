'use client';

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { ethers, keccak256, SigningKey, hexlify, getBytes, formatEther, parseUnits, Interface, toBigInt } from 'ethers';
import { formatDistanceToNow, fromUnixTime } from 'date-fns';
import { vi } from 'date-fns/locale';

// --- Zustand Store Imports ---
import { useAuthStore } from '@/stores/authStore';
import { useWebSocketStore } from '@/stores/webSocketStore';

// --- UI and Icon Imports ---
import Header from '@/components/header';
import { ArrowLeft, PiggyBank, Coins, Loader2, CheckCircle, XCircle, Repeat, Fuel, Landmark, Calendar, Percent, ShieldCheck, ShieldAlert, Lock, Unlock } from 'lucide-react';
import Link from 'next/link';

// --- Protobuf and Config Imports ---
import { create, toBinary } from "@bufbuild/protobuf";
import { TransactionSchema, TransactionSignDataSchema } from '@/proto/transaction_pb';
import { type Receipt } from '@/proto/receipt_pb';

import {
  CallSmartContractDataSchema,
  CallSmartContractHashDataSchema,
  CallSmartContractResponseSchema,
  type CallSmartContractData,
  type CallSmartContractHashData,
  type CallSmartContractResponse,
} from '@/proto/call_data_pb';

import config from '@/lib/config';
import SavingsContractABI from '@/lib/abis/saving.json'; // Import the ABI
import { formatAndShorten } from '@/lib/utils';

import { 
  generateCallDataWithSign,
  generateTransactionWithSign
} from '@/lib/utils';

// --- VM ---
import { encodeFunctionData, decodeFunctionData, decodeFunctionResult } from 'viem';

// --- Contract Address ---

// ✅ ADDED: Constant for minimum deposit amount for easy configuration
const MIN_DEPOSIT_AMOUNT = 100000;
const GAS_USE = 200000;

const OPEN_MESSAGE = "Mở Tiết Kiệm";
const WITHDRAW_MESSAGE = "Tất toán";
const WITHDRAW_EARLY_MESSAGE = "Tất toán sớm";

// --- Type Definitions for Clarity ---
interface Tier {
  id: number;
  duration: bigint;
  interestRateBps: bigint;
  earlyWithdrawalRateBps: bigint;
  isAvailable: boolean;
}

interface Saving {
  id: bigint;
  owner: string;
  principal: bigint;
  creationTime: bigint;
  unlockTime: bigint;
  tierId: bigint;
  isActive: boolean;
}

// Helper to create an ethers.js Interface instance
const savingsInterface = new Interface(SavingsContractABI);

export default function SavingPage() {
  const router = useRouter();

  // --- State for Data from Blockchain ---
  const [tiers, setTiers] = useState<Tier[]>([]);
  const [userSavings, setUserSavings] = useState<Saving[]>([]);
  const [isDataLoading, setIsDataLoading] = useState(true);

  // --- State for the "Create Saving" Form ---
  const [selectedTierId, setSelectedTierId] = useState<number | null>(null);
  const [amount, setAmount] = useState('');
  const [balance, setBalance] = useState<bigint>(BigInt(0)); // ✅ Initialized with 0n
  const [gasPrice, setGasPrice] = useState(0.1); 

  // --- State for Transaction Submission ---
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [receipt, setReceipt] = useState<Receipt | null>(null);
  const [activeTx, setActiveTx] = useState<'create' | 'withdraw' | 'withdrawEarly' | null>(null);

  // --- Get state from stores ---
  const { isAuthenticated, _hasHydrated, privateKey } = useAuthStore();
  const { sendTransaction, getNonce, getBalance, callSmartContract, connected, wallet } = useWebSocketStore();

  // --- Auth check effect ---
  useEffect(() => {
    if (_hasHydrated && !isAuthenticated) {
      router.replace('/');
    }
  }, [_hasHydrated, isAuthenticated, router]);

  // --- Data Fetching Effect ---
  const fetchData = useCallback(async () => {
    if (connected && wallet) { // ✅ Added wallet check
      setIsDataLoading(true);
      try {
        // get balance
        const fetchedBalance = await getBalance();
        setBalance(fetchedBalance);

        // Fetch Tiers
        const fetchedTiers: Tier[] = [];
        for (let i = 0; ; i++) {
          try {
            const encodedCalldata = encodeFunctionData({
              abi: SavingsContractABI,
              functionName: 'tiers',
              args: [BigInt(i)],
            });
            const callRq = generateCallDataWithSign(privateKey!, getBytes(config.savingContractAddress), getBytes(encodedCalldata));
            const callRs = await callSmartContract(callRq);
            const tierData = decodeFunctionResult({
              abi: SavingsContractABI,
              functionName: 'tiers',
              data: hexlify(callRs.Data) as `0x${string}`,
            }) as any[];
            if (tierData[0] === BigInt(0)) break;
            fetchedTiers.push({
              id: i,
              duration: tierData[0],
              interestRateBps: tierData[1],
              earlyWithdrawalRateBps: tierData[2],
              isAvailable: tierData[3],
            });
          } catch (e) {
            break;
          }
        }
        setTiers(fetchedTiers.filter(t => t.isAvailable));
        console.log(fetchedTiers)

        // Fetch User Savings
        const encodedCalldata = encodeFunctionData({
          abi: SavingsContractABI,
          functionName: 'getUserSavingIds',
          args: [wallet.address],
        });
        const callRq = generateCallDataWithSign(privateKey!, getBytes(config.savingContractAddress), getBytes(encodedCalldata));
        const callRs = await callSmartContract(callRq);

        if (callRs.Data.length === 0) {
          setUserSavings([]);
          return; // No need to proceed further
        }

        const savingIds = decodeFunctionResult({
          abi: SavingsContractABI,
          functionName: 'getUserSavingIds',
          data: hexlify(callRs.Data) as `0x${string}`,
        }) as bigint[];

        const fetchedSavings = [];
        for (const id of savingIds) {
          try {
            console.log(`Fetching details for saving ID: ${id}...`); // Good for debugging

            const encodedDetailsCalldata = encodeFunctionData({
              abi: SavingsContractABI,
              functionName: 'getSavingDetails',
              args: [id],
            });

            const callDetailsRq = generateCallDataWithSign(privateKey!, getBytes(config.savingContractAddress), getBytes(encodedDetailsCalldata));

            // The 'await' here will pause the loop until this specific call is finished
            const callDetailsRs = await callSmartContract(callDetailsRq);

            const savingDetail = decodeFunctionResult({
              abi: SavingsContractABI,
              functionName: 'getSavingDetails',
              data: hexlify(callDetailsRs.Data) as `0x${string}`,
            }) as Saving; // Assuming the result is an array
            savingDetail.id = id


            // Push the processed result into our array
            // Note: I've re-added the object creation from your first snippet
            // as it's much safer and more readable than using a raw array.
            fetchedSavings.push(savingDetail);

          } catch (error) {
            console.error(`Failed to fetch details for saving ID: ${id}`, error);
            // You can choose to continue the loop or break, depending on your needs
          }
        }

        // Now 'fetchedSavings' contains all the results, fetched sequentially.
        // You can sort it just like before.
        setUserSavings(fetchedSavings.sort((a, b) => Number(b.creationTime) - Number(a.creationTime)));

      } catch (err) {
        console.error("Failed to fetch contract data:", err);
        setError("Không thể tải dữ liệu từ hợp đồng. Vui lòng thử lại.");
      } finally {
        setIsDataLoading(false);
      }
    }
  }, [connected, wallet, privateKey]); // ✅ Added dependencies


  useEffect(() => {
    if (connected) {
      fetchData();
    }
  }, [connected]);


  // --- Form and UI Logic ---
  const handleReset = () => {
    setAmount('');
    setSelectedTierId(null);
    setError(null);
    setReceipt(null);
    setIsLoading(false);
    setActiveTx(null);
    fetchData(); // Refresh data after transaction
  };

  const formattedAmount = useMemo(() => {
    if (amount === '') return '';
    try {
      return new Intl.NumberFormat('vi-VN').format(parseInt(amount, 10));
    } catch {
      return amount;
    }
  }, [amount]);


  // --- Calculate the transfer fee using useMemo for efficiency ---
  const transferFee = useMemo(() => {
    try {
      const gasPriceWei = ethers.parseUnits(gasPrice.toString(), 'finney');
      const feeInWei = gasPriceWei * BigInt(GAS_USE);
      const feeAsNumber = Number(ethers.formatEther(feeInWei));
      return new Intl.NumberFormat('vi-VN').format(feeAsNumber)
    } catch (e) {
      return '0';
    }
  }, [gasPrice]);

  // ✅ ADDED: Memoized and formatted balance for display
  const formattedBalance = useMemo(() => {
    try {
        // formatEther effectively divides by 10^18
        const etherString = formatEther(balance);
        const number = parseFloat(etherString);
        return new Intl.NumberFormat('vi-VN', { maximumFractionDigits: 4 }).format(number);
    } catch {
        return '0';
    }
  }, [balance]);

  const handleAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const rawValue = e.target.value.replace(/[^0-9]/g, '');
    setAmount(rawValue);
  };

  // --- Core Transaction Handlers ---
  const handleCreateSaving = async (e: React.FormEvent) => {
    e.preventDefault();
    if (selectedTierId === null) {
      setError("Vui lòng chọn một kỳ hạn gửi.");
      return;
    }
    await submitTransaction('create', { tierId: selectedTierId });
  };

  const handleWithdraw = (savingId: bigint) => {
    submitTransaction('withdraw', { savingId });
  }

  const handleWithdrawEarly = (savingId: bigint) => {
    console.log("WithdrawEarly", savingId)
    submitTransaction('withdrawEarly', { savingId });
  }

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


  // --- Universal Transaction Submission Logic ---
  const submitTransaction = async (
    type: 'create' | 'withdraw' | 'withdrawEarly',
    params: { tierId?: number, savingId?: bigint } // ✅ Changed savingId to bigint
  ) => {
    setError(null);
    setReceipt(null);
    setIsLoading(true);
    setActiveTx(type);
  
    // ✅ UPDATED: Validation logic to check for minimum amount.
    if (type === 'create') {
        const numericAmount = parseInt(amount, 10);
        if (isNaN(numericAmount) || numericAmount < MIN_DEPOSIT_AMOUNT) {
          setError(`Số tiền gửi tối thiểu là ${new Intl.NumberFormat('vi-VN').format(MIN_DEPOSIT_AMOUNT)} VNĐ.`);
          setIsLoading(false);
          return;
        }
    }

    if (!privateKey) {
      setError('Không tìm thấy khóa riêng tư. Vui lòng đăng nhập lại.');
      setIsLoading(false);
      return;
    }

    try {
      const nonce = await getNonce();
      let data: string;
      let value: bigint;
      let message:  string; 

      if (type === 'create') {
        data = encodeFunctionData({
          abi: SavingsContractABI,
          functionName: 'createSaving',
          args: [params.tierId!],
        });
        value = parseUnits(amount, 18);

        message = OPEN_MESSAGE;
      } else if (type === 'withdraw') {
        data = encodeFunctionData({
          abi: SavingsContractABI,
          functionName: 'withdraw',
          args: [params.savingId!],
        });
        value = BigInt(0);

        message = WITHDRAW_MESSAGE;
      } else { // withdrawEarly
        console.log("params.savingId", params.savingId)
        data = encodeFunctionData({
          abi: SavingsContractABI,
          functionName: 'withdrawEarly',
          args: [params.savingId!],
        });
        value = BigInt(0);
        console.log("withdraw early data", data)

        message = WITHDRAW_EARLY_MESSAGE;
      }

    // Convert BigNumber to bytes array (Uint8Array)
      const gasPrice = 0.1
      const amountByteArray = ethers.toBeArray(value);   
      const tx = generateTransactionWithSign(
        privateKey!,
        BigInt(nonce),
        getBytes(config.savingContractAddress),
        amountByteArray,
        getBytes(data),
        message,
        BigInt(GAS_USE), // tmp hard code
        ethers.toBeArray(ethers.parseUnits(gasPrice.toString(), 'finney')),
      )

      const returnedReceipt: Receipt = await sendTransaction(tx);
      setReceipt(returnedReceipt);

    } catch (err: any) {
      console.error(`${type} failed:`, err);
      setError(err.message || 'Đã xảy ra lỗi không mong muốn.');
    } finally {
      setIsLoading(false);
    }
  };


  // --- Helper Functions for Display ---
  const formatBps = (bps: bigint) => `${(Number(bps) / 100).toFixed(2)}%`;
  const formatDuration = (seconds: bigint) => {
    const days = Number(seconds) / 86400;
    if (days < 1) return `${Number(seconds)} giây`;
    if (days < 30) return `${days} ngày`;
    return `${Math.round(days / 30)} tháng`;
  };

  // --- Auth Loading Screen ---
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
        <div className="max-w-4xl mx-auto">
          {!receipt && (
            <Link
              href="/dapps"
              className="inline-flex items-center text-cyan-400 hover:text-cyan-300 mb-8 group p-2 rounded-md hover:bg-cyan-900/50 transition-colors -ml-2 z-[10]"
            >
              <ArrowLeft className="h-5 w-5 mr-2 transition-transform group-hover:-translate-x-1" />
              Quay lại danh sách dApps
            </Link>
          )}

          <div className="bg-gray-800 p-6 sm:p-8 rounded-xl shadow-2xl">
            {/* --- RECEIPT DISPLAY (Unchanged) --- */}

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
                      </div>
                    </div>
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
                <div className="grid grid-cols-1 md:grid-cols-2 gap-8 md:gap-12">
                  {/* --- LEFT SIDE: CREATE SAVING --- */}
                  <div className="flex flex-col">
                    <h1 className="text-3xl font-bold text-cyan-400 mb-2 flex items-center">
                      <PiggyBank className="h-8 w-8 mr-3"/> Gửi Tiết Kiệm
                    </h1>
                    <p className="text-gray-400 mb-6">Chọn một kỳ hạn và gửi tiền để nhận lãi suất hấp dẫn.</p>

                    {isDataLoading ? (
                      <div className="flex items-center justify-center h-40"><Loader2 className="h-8 w-8 animate-spin" /></div>
                    ) : (
                        <form onSubmit={handleCreateSaving} className="space-y-6">
                          {/* Tier Selection */}
                          <div>
                            <label className="block text-sm font-medium text-gray-300 mb-2">Chọn kỳ hạn</label>
                            <div className="space-y-3">
                              {tiers.map(tier => (
                                <div key={tier.id} onClick={() => setSelectedTierId(tier.id)} className={`p-4 rounded-lg cursor-pointer transition-all border-2 ${selectedTierId === tier.id ? 'bg-cyan-900/50 border-cyan-500' : 'bg-gray-700/50 border-gray-600 hover:border-gray-500'}`}>
                                  <div className="flex justify-between items-center font-bold">
                                    <span><Calendar className="h-4 w-4 inline mr-2" />{formatDuration(tier.duration)}</span>
                                    <span className="text-cyan-300"><Percent className="h-4 w-4 inline mr-1" />{formatBps(tier.interestRateBps)}/năm</span>
                                  </div>
                                  <p className="text-xs text-gray-400 mt-1">Lãi suất rút sớm: {formatBps(tier.earlyWithdrawalRateBps)}/năm</p>
                                </div>
                              ))}
                            </div>
                          </div>

                          {/* Amount Input */}
                          <div>
                            <label htmlFor="amount" className="block text-sm font-medium text-gray-300 mb-2">Số Tiền Gửi</label>
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
                                className="w-full bg-gray-700 border-gray-600 rounded-md py-3 pl-10 pr-4 text-white focus:ring-cyan-500 focus:border-cyan-500"
                                // ✅ UPDATED: Placeholder to show the minimum amount
                                placeholder={`Tối thiểu ${new Intl.NumberFormat('vi-VN').format(MIN_DEPOSIT_AMOUNT)}`}
                                required
                                disabled={isLoading}
                              />
                            </div>
                            {/* ✅ ADDED: Balance display below the input */}
                            <p className="mt-2 text-sm text-gray-400 text-right">
                              Số dư: <span className="font-medium text-cyan-300">{formattedBalance} VNĐ</span>
                            </p>
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

                          {/* Submit Button */}
                          <button type="submit" disabled={isLoading && activeTx === 'create'} className="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-cyan-600 hover:bg-cyan-700 disabled:bg-gray-600 disabled:cursor-not-allowed transition-colors">
                            {isLoading && activeTx === 'create' ? <><Loader2 className="h-5 w-5 mr-2 animate-spin" />Đang xử lý...</> : 'Gửi Tiết Kiệm'}
                          </button>
                        </form>
                      )}
                  </div>

                  {/* --- RIGHT SIDE: MY SAVINGS (Unchanged, but ensure handleWithdraw uses bigint) --- */}
                  <div className="flex flex-col">
                    <h2 className="text-2xl font-bold text-gray-200 mb-2 flex items-center"><Landmark className="h-7 w-7 mr-3"/> Sổ Tiết Kiệm Của Tôi</h2>
                    <p className="text-gray-400 mb-6">Quản lý các khoản tiết kiệm đang hoạt động và đã tất toán.</p>
                     {isDataLoading ? (
                      <div className="flex items-center justify-center h-40"><Loader2 className="h-8 w-8 animate-spin" /></div>
                    ) : (
                        <div className="space-y-4 max-h-[60vh] overflow-y-auto pr-2">
                           {/* ... my savings list ... */}
                           {userSavings.map(saving => {
                            const isLocked = new Date() < fromUnixTime(Number(saving.unlockTime));
                            const canWithdraw = saving.isActive && !isLocked;
                            const canWithdrawEarly = saving.isActive && isLocked;

                            return (
                              <div key={String(saving.id)} className={`p-4 rounded-lg ${saving.isActive ? 'bg-gray-700/80' : 'bg-gray-900/50 opacity-60'}`}>
                                <div className="flex justify-between items-start">
                                  <div>
                                    <p className="font-bold text-lg">{formatAndShorten(saving.principal)} VNĐ</p>
                                    <p className="text-xs text-gray-400">Gửi ngày: {new Date(Number(saving.creationTime) * 1000).toLocaleDateString('vi-VN')}</p>
                                  </div>
                                  {!saving.isActive ? (
                                    <span className="text-sm font-medium text-gray-400 bg-gray-800 px-2 py-1 rounded-md">Đã tất toán</span>
                                  ) : isLocked ? (
                                      <span className="flex items-center text-sm font-medium text-yellow-400 bg-yellow-900/50 px-2 py-1 rounded-md">
                                        <Lock className="h-3 w-3 mr-1.5"/>
                                        Còn {formatDistanceToNow(fromUnixTime(Number(saving.unlockTime)), { addSuffix: true, locale: vi })}
                                      </span>
                                    ) : (
                                        <span className="flex items-center text-sm font-medium text-green-400 bg-green-900/50 px-2 py-1 rounded-md">
                                          <Unlock className="h-3 w-3 mr-1.5"/> Sẵn sàng tất toán
                                        </span>
                                      )}
                                </div>
                                {saving.isActive && (
                                  <div className="mt-4 flex gap-3">
                                    <button onClick={() => handleWithdraw(saving.id)} disabled={!canWithdraw || isLoading} className="flex-1 text-sm flex items-center justify-center py-2 px-3 rounded-md bg-green-600 hover:bg-green-700 disabled:bg-gray-500 disabled:cursor-not-allowed transition-colors">
                                      <ShieldCheck className="h-4 w-4 mr-2"/> Tất toán
                                    </button>
                                    <button onClick={() => handleWithdrawEarly(saving.id)} disabled={!canWithdrawEarly || isLoading} className="flex-1 text-sm flex items-center justify-center py-2 px-3 rounded-md bg-red-600 hover:bg-red-700 disabled:bg-gray-500 disabled:cursor-not-allowed transition-colors">
                                      <ShieldAlert className="h-4 w-4 mr-2"/> Rút sớm
                                    </button>
                                  </div>
                                )}
                              </div>
                            )
                          })}
                        </div>
                    )}
                  </div>
                </div>
              )}
            {error && !receipt && (
              <div className="mt-6 flex items-center p-4 text-sm text-red-300 bg-red-900/30 rounded-lg">
                <XCircle className="h-5 w-5 mr-3 flex-shrink-0" />
                <span>{error}</span>
              </div>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
