// app/dapps/bills/page.tsx
'use client';

import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { ethers, hexlify, getBytes, formatEther, toBigInt, decodeBytes32String } from 'ethers';
import { format, formatDistanceToNow, fromUnixTime, isPast } from 'date-fns';
import { vi } from 'date-fns/locale';

// --- Zustand Store Imports ---
import { useAuthStore } from '@/stores/authStore';
import { useWebSocketStore } from '@/stores/webSocketStore';

// --- UI and Icon Imports ---
import Header from '@/components/header';
import { ArrowLeft, Loader2, CheckCircle, XCircle, Repeat, FileText, Calendar, Zap, Droplets, Wifi, Home, HelpCircle, Hash, Clock, CircleDollarSign } from 'lucide-react';
import Link from 'next/link';

// --- Protobuf and Config Imports ---
import { type Receipt } from '@/proto/receipt_pb';
import config from '@/lib/config';
import ServiceBillManagerABI from '@/lib/abis/serviceBills.json'; // IMPORTANT: Use the new ABI

import { 
  generateCallDataWithSign,
  generateTransactionWithSign
} from '@/lib/utils';

// --- VM ---
import { encodeFunctionData, decodeFunctionResult } from 'viem';

// --- Type Definitions for this page ---
enum BillType { Water, Electric, Internet, Rent, ServiceFee, Other }

interface Bill {
  id: bigint;
  customer: string;
  description: string;
  billType: BillType;
  amount: bigint;
  dueDate: bigint;
  paymentDate: bigint;
  isPaid: boolean;
  usageValue: bigint;
  usageUnit: string; // Converted from bytes32
}

const PAY_MESSAGE = "Thanh toán hoá đơn";
const GAS_USE = 200000; // Adjust as needed for payBill function

export default function PayServiceBillsPage() {
  const router = useRouter();

  // --- State for Data from Blockchain ---
  const [bills, setBills] = useState<Bill[]>([]);
  const [isDataLoading, setIsDataLoading] = useState(true);

  // --- State for Transaction Submission ---
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [receipt, setReceipt] = useState<Receipt | null>(null);
  const [activeBillId, setActiveBillId] = useState<bigint | null>(null); // Track which bill is being paid

  // --- Get state from stores ---
  const { isAuthenticated, _hasHydrated, privateKey } = useAuthStore();
  const { sendTransaction, getNonce, callSmartContract, connected, wallet } = useWebSocketStore();

  // --- Auth check effect ---
  useEffect(() => {
    if (_hasHydrated && !isAuthenticated) {
      router.replace('/');
    }
  }, [_hasHydrated, isAuthenticated, router]);

  // --- Data Fetching Logic ---
  const fetchData = useCallback(async () => {
    if (connected && wallet && privateKey) {
      setIsDataLoading(true);
      try {
        // Step 1: Get the user's bill IDs
        const getIdsCalldata = encodeFunctionData({
          abi: ServiceBillManagerABI,
          functionName: 'getBillsForCustomer',
          args: [wallet.address],
        });

        const callIdsRq = generateCallDataWithSign(privateKey, getBytes(config.serviceBillContractAddress), getBytes(getIdsCalldata));
        const callIdsRs = await callSmartContract(callIdsRq);
        
        if (callIdsRs.Data.length === 0) {
          setBills([]);
          setIsDataLoading(false);
          return;
        }

        const billIds = decodeFunctionResult({
          abi: ServiceBillManagerABI,
          functionName: 'getBillsForCustomer',
          data: hexlify(callIdsRs.Data) as `0x${string}`,
        }) as bigint[];
        
        // Step 2: Fetch details for each bill ID
        const fetchedBills: Bill[] = [];
        for (const id of billIds) {
          const getDetailsCalldata = encodeFunctionData({
            abi: ServiceBillManagerABI,
            functionName: 'bills',
            args: [id],
          });
          const callDetailsRq = generateCallDataWithSign(privateKey, getBytes(config.serviceBillContractAddress), getBytes(getDetailsCalldata));
          const callDetailsRs = await callSmartContract(callDetailsRq);

          const billData = decodeFunctionResult({
            abi: ServiceBillManagerABI,
            functionName: 'bills',
            data: hexlify(callDetailsRs.Data) as `0x${string}`,
          }) as any[]; // viem decodes struct to array
          
          fetchedBills.push({
            id: billData[0],
            customer: billData[1],
            description: billData[2],
            billType: Number(billData[3]),
            amount: billData[4],
            dueDate: billData[5],
            paymentDate: billData[6],
            isPaid: billData[7],
            usageValue: billData[8],
            // Convert bytes32 usage unit to a readable string
            usageUnit: decodeBytes32String(billData[9]),
          });
        }
        
        // Sort bills by due date, unpaid first
        setBills(fetchedBills.sort((a, b) => {
          if (a.isPaid !== b.isPaid) return a.isPaid ? 1 : -1;
          return Number(a.dueDate) - Number(b.dueDate);
        }));

      } catch (err) {
        console.error("Failed to fetch bills data:", err);
        setError("Không thể tải danh sách hoá đơn. Vui lòng thử lại.");
      } finally {
        setIsDataLoading(false);
      }
    }
  }, [connected, wallet, privateKey, callSmartContract]);

  useEffect(() => {
    if (connected) {
      fetchData();
    }
  }, [connected, fetchData]);

  // --- Transaction Handler ---
  const handlePayBill = async (bill: Bill) => {
    setError(null);
    setReceipt(null);
    setIsLoading(true);
    setActiveBillId(bill.id);

    if (!privateKey) {
      setError('Không tìm thấy khóa riêng tư. Vui lòng đăng nhập lại.');
      setIsLoading(false);
      return;
    }

    try {
      const nonce = await getNonce();
      
      const calldata = encodeFunctionData({
        abi: ServiceBillManagerABI,
        functionName: 'payBill',
        args: [bill.id],
      });

      // For payBill, the `value` of the transaction is the bill's amount.
      const value = bill.amount;

      // Note: We use a fixed gas price here for simplicity, like in your original code.
      const gasPrice = 0.1;

      const tx = generateTransactionWithSign(
        privateKey,
        BigInt(nonce),
        getBytes(config.serviceBillContractAddress),
        ethers.toBeArray(value),
        getBytes(calldata),
        PAY_MESSAGE,
        BigInt(GAS_USE),
        ethers.toBeArray(ethers.parseUnits(gasPrice.toString(), 'finney')),
      );

      const returnedReceipt: Receipt = await sendTransaction(tx);
      setReceipt(returnedReceipt);

    } catch (err: any) {
      console.error(`Payment for bill #${bill.id} failed:`, err);
      setError(err.message || 'Thanh toán hoá đơn thất bại.');
    } finally {
      setIsLoading(false);
      setActiveBillId(null);
    }
  };

  const handleReset = () => {
    setError(null);
    setReceipt(null);
    setIsLoading(false);
    setActiveBillId(null);
    fetchData(); // Refresh bill list
  };

  // --- Helper Functions for Display ---
  const getBillTypeIcon = (type: BillType) => {
    switch (type) {
      case BillType.Water: return <Droplets className="h-5 w-5 mr-3 text-blue-400" />;
      case BillType.Electric: return <Zap className="h-5 w-5 mr-3 text-yellow-400" />;
      case BillType.Internet: return <Wifi className="h-5 w-5 mr-3 text-green-400" />;
      case BillType.Rent: return <Home className="h-5 w-5 mr-3 text-purple-400" />;
      default: return <HelpCircle className="h-5 w-5 mr-3 text-gray-400" />;
    }
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

  // --- Main Render ---
  return (
    <main className="min-h-screen text-white">
      <Header />
      <div className="container mx-auto px-4 py-12 sm:py-16">
        <div className="max-w-2xl mx-auto">
          {!receipt && (
            <Link href="/dapps" className="inline-flex items-center text-cyan-400 hover:text-cyan-300 mb-8 group p-2 rounded-md hover:bg-cyan-900/50 transition-colors -ml-2">
              <ArrowLeft className="h-5 w-5 mr-2 transition-transform group-hover:-translate-x-1" />
              Quay lại danh sách dApps
            </Link>
          )}

          <div className="bg-gray-800 p-6 sm:p-8 rounded-xl shadow-2xl">
            {receipt ? (
              // --- RECEIPT DISPLAY (Largely reused from your Savings page) ---
              <div>
                <div className="flex flex-col items-center text-center">
                  <CheckCircle className="h-16 w-16 text-green-400 mb-4" />
                  <h1 className="text-2xl font-bold">Thanh toán thành công</h1>
                  <p className="text-gray-400 mt-1">Hoá đơn của bạn đã được thanh toán và ghi nhận.</p>
                </div>
                {/* ... other receipt details (tx hash, block, fee) can be added here just like your savings page ... */}
                <button
                  onClick={handleReset}
                  className="mt-8 w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-cyan-600 hover:bg-cyan-700 transition-colors"
                >
                  <Repeat className="h-5 w-5 mr-2" />
                  Xem danh sách hoá đơn
                </button>
              </div>
            ) : (
              // --- BILL LIST DISPLAY ---
              <div>
                <h1 className="text-3xl font-bold text-cyan-400 mb-2 flex items-center">
                  <FileText className="h-8 w-8 mr-3"/> Hoá Đơn Dịch Vụ
                </h1>
                <p className="text-gray-400 mb-8">Xem và thanh toán các hoá đơn chưa thanh toán của bạn.</p>
                
                {isDataLoading ? (
                  <div className="flex items-center justify-center h-40"><Loader2 className="h-8 w-8 animate-spin" /></div>
                ) : bills.length === 0 ? (
                  <div className="text-center py-10 px-6 bg-gray-900/50 rounded-lg">
                    <CheckCircle className="h-12 w-12 mx-auto text-green-500" />
                    <h3 className="mt-4 text-xl font-semibold">Tuyệt vời!</h3>
                    <p className="mt-1 text-gray-400">Bạn không có hoá đơn nào cần thanh toán.</p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {bills.map(bill => {
                      const overdue = !bill.isPaid && isPast(fromUnixTime(Number(bill.dueDate)));
                      return (
                        <div key={String(bill.id)} className={`p-4 rounded-lg transition-all border-2 ${bill.isPaid ? 'bg-gray-900/30 border-gray-700 opacity-60' : 'bg-gray-700/60 border-gray-600'} ${overdue ? '!border-red-500' : ''}`}>
                          <div className="flex items-start justify-between">
                            <div className="flex items-center">
                                {getBillTypeIcon(bill.billType)}
                                <div>
                                    <p className="font-semibold text-white">{bill.description}</p>
                                    <p className="text-sm text-gray-400">#{String(bill.id)}</p>
                                </div>
                            </div>
                            <div className="text-right">
                                    <p className="text-xl font-bold text-cyan-300">{new Intl.NumberFormat('vi-VN',{
                                      minimumFractionDigits: 0,
                                      maximumFractionDigits: 8,
                                    }).format(Number(formatEther(bill.amount)))} VNĐ</p>
                                {bill.isPaid ? (
                                    <span className="text-xs font-medium text-green-400 bg-green-900/50 px-2 py-0.5 rounded-md">
                                        Đã thanh toán
                                    </span>
                                ) : overdue ? (
                                    <span className="text-xs font-medium text-red-400 bg-red-900/50 px-2 py-0.5 rounded-md">
                                        Quá hạn
                                    </span>
                                ) : null}
                            </div>
                          </div>
                          
                          <div className="mt-4 pt-4 border-t border-gray-700 text-sm space-y-2">
                              <div className="flex justify-between items-center text-gray-300">
                                  <span className="flex items-center"><Calendar className="h-4 w-4 mr-2 text-gray-500"/>Hạn chót</span>
                                  <span>{format(fromUnixTime(Number(bill.dueDate)), 'dd/MM/yyyy')}</span>
                              </div>
                              {bill.usageValue > 0 && (
                                <div className="flex justify-between items-center text-gray-300">
                                    <span className="flex items-center"><Hash className="h-4 w-4 mr-2 text-gray-500"/>Mức sử dụng</span>
                                    <span>{String(bill.usageValue)} {bill.usageUnit}</span>
                                </div>
                              )}
                               {bill.isPaid && (
                                <div className="flex justify-between items-center text-green-400">
                                    <span className="flex items-center"><Clock className="h-4 w-4 mr-2 text-gray-500"/>Ngày thanh toán</span>
                                    <span>{format(fromUnixTime(Number(bill.paymentDate)), 'dd/MM/yyyy HH:mm')}</span>
                                </div>
                              )}
                          </div>

                          {!bill.isPaid && (
                            <div className="mt-4">
                              <button 
                                onClick={() => handlePayBill(bill)}
                                disabled={isLoading}
                                className="w-full flex justify-center items-center py-2 px-4 border border-transparent rounded-md shadow-sm text-base font-medium text-white bg-cyan-600 hover:bg-cyan-700 disabled:bg-gray-500 disabled:cursor-not-allowed transition-colors"
                              >
                                {isLoading && activeBillId === bill.id ? (
                                  <><Loader2 className="h-5 w-5 mr-2 animate-spin" />Đang xử lý...</>
                                ) : (
                                  <><CircleDollarSign className="h-5 w-5 mr-2"/>Thanh toán ngay</>
                                )}
                              </button>
                            </div>
                          )}
                        </div>
                      )
                    })}
                  </div>
                )}
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
