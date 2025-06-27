// app/dapps/invoicing/pay/page.tsx
'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { ethers, getBytes, formatEther, hexlify } from 'ethers';

// --- Zustand Store Imports ---
import { useAuthStore } from '@/stores/authStore';
import { useWebSocketStore } from '@/stores/webSocketStore';

// --- UI and Icon Imports ---
import Header from '@/components/header';
import { ArrowLeft, Loader2, CheckCircle, XCircle, Repeat, Hash, Search, User, Wallet } from 'lucide-react';
import Link from 'next/link';

// --- Protobuf and Config Imports ---
import { type Receipt } from '@/proto/receipt_pb';
import config from '@/lib/config';
import InvoicingABI from '@/lib/abis/invoice.json';

import { 
  generateCallDataWithSign,
  generateTransactionWithSign
} from '@/lib/utils';

// --- VM ---
import { encodeFunctionData, decodeFunctionResult } from 'viem';

const GAS_USE = 300000;

// --- Type Definitions ---
interface InvoiceDetails {
  id: string;
  payee: string;
  payer: string;
  totalAmount: bigint;
  description: string;
  status: number; // 0: Created, 1: Paid, 2: Cancelled
  items: { description: string; unitPrice: bigint; quantity: bigint; }[];
}

const statusMap = {
  0: { text: "Chờ thanh toán", color: "text-yellow-400" },
  1: { text: "Đã thanh toán", color: "text-green-400" },
  2: { text: "Đã huỷ", color: "text-red-400" },
};

export default function PayInvoicePage() {
  const router = useRouter();

  // --- State ---
  const [invoiceIdInput, setInvoiceIdInput] = useState('');
  const [invoiceDetails, setInvoiceDetails] = useState<InvoiceDetails | null>(null);
  const [isFinding, setIsFinding] = useState(false);
  const [findError, setFindError] = useState<string | null>(null);

  const [isPaying, setIsPaying] = useState(false);
  const [payError, setPayError] = useState<string | null>(null);
  const [receipt, setReceipt] = useState<Receipt | null>(null);

  const { isAuthenticated, _hasHydrated, privateKey } = useAuthStore();
  const { sendTransaction, getNonce, connected, wallet } = useWebSocketStore();
  const { callSmartContract } = useWebSocketStore(); // Get callSmartContract

  // --- Auth check ---
  useEffect(() => {
    if (_hasHydrated && !isAuthenticated) {
      router.replace('/');
    }
  }, [_hasHydrated, isAuthenticated, router]);

  // --- Find Invoice Logic ---
  const handleFindInvoice = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!invoiceIdInput || parseInt(invoiceIdInput) <= 0) {
      setFindError("Vui lòng nhập mã hoá đơn hợp lệ.");
      return;
    }
    
    setIsFinding(true);
    setFindError(null);
    setInvoiceDetails(null);
    setPayError(null);

    try {
      const id = BigInt(invoiceIdInput);

      // Call getInvoiceDetails
      const detailsCalldata = encodeFunctionData({ abi: InvoicingABI, functionName: 'getInvoiceDetails', args: [id] });
      const detailsRq = generateCallDataWithSign(privateKey!, getBytes(config.invoiceContractAddress), getBytes(detailsCalldata));
      const detailsRs = await callSmartContract(detailsRq);
      const [payee, payer, totalAmount, description, status] = decodeFunctionResult({
        abi: InvoicingABI, functionName: 'getInvoiceDetails', data: hexlify(detailsRs.Data) as `0x${string}`
      }) as [string, string, bigint, string, number];

      if (payee === ethers.ZeroAddress) {
        throw new Error("Không tìm thấy hoá đơn với mã này.");
      }

      // Call getInvoiceItems
      const itemsCalldata = encodeFunctionData({ abi: InvoicingABI, functionName: 'getInvoiceItems', args: [id] });
      const itemsRq = generateCallDataWithSign(privateKey!, getBytes(config.invoiceContractAddress), getBytes(itemsCalldata));
      const itemsRs = await callSmartContract(itemsRq);
      const items = decodeFunctionResult({
        abi: InvoicingABI, functionName: 'getInvoiceItems', data: hexlify(itemsRs.Data) as `0x${string}`
      }) as { description: string; unitPrice: bigint; quantity: bigint; }[];

      setInvoiceDetails({ id: invoiceIdInput, payee, payer, totalAmount, description, status, items });

    } catch (err: any) {
      setFindError(err.message || "Đã xảy ra lỗi khi tìm hoá đơn.");
    } finally {
      setIsFinding(false);
    }
  };

  // --- Pay Invoice Logic ---
  const handlePayInvoice = async () => {
    if (!invoiceDetails) return;

    setIsPaying(true);
    setPayError(null);
    setReceipt(null);
    
    try {
      const nonce = await getNonce();
      const calldata = encodeFunctionData({
        abi: InvoicingABI,
        functionName: 'payInvoice',
        args: [BigInt(invoiceDetails.id)]
      });

      const tx = generateTransactionWithSign(
        privateKey!, BigInt(nonce), getBytes(config.invoiceContractAddress),
        ethers.toBeArray(invoiceDetails.totalAmount), // Send the required amount with the transaction
        getBytes(calldata), `Thanh toán HĐ #${invoiceDetails.id}`,
        BigInt(GAS_USE), ethers.toBeArray(ethers.parseUnits("0.1", 'finney'))
      );
      
      const returnedReceipt = await sendTransaction(tx);
      setReceipt(returnedReceipt);

    } catch (err: any) {
      setPayError(err.message || "Giao dịch thanh toán thất bại.");
    } finally {
      setIsPaying(false);
    }
  };
  
  const handleReset = () => {
    setReceipt(null);
    setInvoiceDetails(null);
    setInvoiceIdInput('');
    setFindError(null);
    setPayError(null);
  };
  
  const formatAndShorten = (balance: bigint, decimals = 4) => {
    return new Intl.NumberFormat('vi-VN', { maximumFractionDigits: decimals }).format(parseFloat(formatEther(balance)));
  };

  const isUserThePayer = wallet?.address.toLowerCase() === invoiceDetails?.payer.toLowerCase();

  if (!_hasHydrated || !isAuthenticated) {
    return <div className="flex items-center justify-center min-h-screen"><Loader2 className="h-12 w-12 animate-spin text-cyan-400" /></div>;
  }

  return (
    <main className="min-h-screen text-white">
      <Header />
      <div className="container mx-auto px-4 py-12 sm:py-16">
        <div className="max-w-2xl mx-auto">
          {!receipt && (
            <Link href="/dapps" className="inline-flex items-center text-cyan-400 hover:text-cyan-300 mb-8 group p-2 -ml-2">
              <ArrowLeft className="h-5 w-5 mr-2 transition-transform group-hover:-translate-x-1" />
              Quay lại danh sách dApps
            </Link>
          )}
          
          <div className="bg-gray-800 p-6 sm:p-8 rounded-xl shadow-2xl">
            {receipt ? (
              <div className="text-center">
                <CheckCircle className="h-16 w-16 text-green-400 mb-4 mx-auto" />
                <h1 className="text-2xl font-bold">Thanh toán thành công!</h1>
                <p className="text-gray-400 mt-1">Hoá đơn đã được thanh toán trên blockchain.</p>
                <button onClick={handleReset} className="mt-8 w-full max-w-sm mx-auto flex justify-center items-center py-3 px-4 rounded-md text-white bg-cyan-600 hover:bg-cyan-700">
                  <Repeat className="h-5 w-5 mr-2" /> Tìm hoá đơn khác
                </button>
              </div>
            ) : !invoiceDetails ? (
              // --- Find Invoice View ---
              <div>
                <h1 className="text-3xl font-bold text-yellow-400 mb-2 flex items-center">
                  <Search className="h-8 w-8 mr-3"/> Thanh Toán Hoá Đơn
                </h1>
                <p className="text-gray-400 mb-6">Nhập mã hoá đơn bạn nhận được để xem chi tiết và thanh toán.</p>
                <form onSubmit={handleFindInvoice} className="space-y-4">
                  <div>
                    <label htmlFor="invoiceId" className="flex items-center text-sm font-medium text-gray-300 mb-1"><Hash className="h-4 w-4 mr-1.5"/>Mã hoá đơn</label>
                    <input id="invoiceId" type="number" value={invoiceIdInput} onChange={e => setInvoiceIdInput(e.target.value)} placeholder="123" required className="w-full bg-gray-700 border-gray-600 rounded-md py-2 px-3 focus:ring-cyan-500 focus:border-cyan-500" />
                  </div>
                  <button type="submit" disabled={isFinding} className="w-full flex justify-center items-center py-3 rounded-md bg-cyan-600 hover:bg-cyan-700 disabled:bg-gray-500">
                    {isFinding ? <Loader2 className="h-5 w-5 animate-spin"/> : 'Tìm Hoá Đơn'}
                  </button>
                </form>
                {findError && (
                  <div className="mt-6 flex items-center p-4 text-sm text-red-300 bg-red-900/30 rounded-lg">
                    <XCircle className="h-5 w-5 mr-3 flex-shrink-0" />
                    <span>{findError}</span>
                  </div>
                )}
              </div>
            ) : (
              // --- Invoice Details View ---
              <div>
                <button onClick={() => setInvoiceDetails(null)} className="inline-flex items-center text-cyan-400 hover:text-cyan-300 mb-6 group p-2 -ml-2">
                  <ArrowLeft className="h-5 w-5 mr-2 transition-transform group-hover:-translate-x-1" />
                  Tìm hoá đơn khác
                </button>
                <div className="flex justify-between items-start">
                    <h2 className="text-2xl font-bold">Chi tiết Hoá đơn #{invoiceDetails.id}</h2>
                    <span className={`px-3 py-1 text-sm font-semibold rounded-full ${statusMap[invoiceDetails.status].color} bg-gray-900/50`}>
                        {statusMap[invoiceDetails.status].text}
                    </span>
                </div>
                <p className="text-gray-400 mt-1">{invoiceDetails.description}</p>
                
                <div className="mt-6 grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                  <div className="bg-gray-900/50 p-3 rounded-lg">
                    <p className="text-gray-400 flex items-center"><User className="h-4 w-4 mr-1.5"/>Người nhận tiền (Payee)</p>
                    <p className="font-mono break-all">{invoiceDetails.payee}</p>
                  </div>
                  <div className="bg-gray-900/50 p-3 rounded-lg">
                    <p className="text-gray-400 flex items-center"><Wallet className="h-4 w-4 mr-1.5"/>Người trả (Payer)</p>
                    <p className="font-mono break-all">{invoiceDetails.payer}</p>
                  </div>
                </div>

                <div className="mt-6">
                  <h3 className="font-semibold mb-2">Các mục</h3>
                  <div className="space-y-2 border border-gray-700 rounded-lg p-3">
                    {invoiceDetails.items.map((item, index) => (
                      <div key={index} className="flex justify-between items-center">
                        <p>{item.quantity.toString()} x {item.description}</p>
                        <p className="font-semibold">{formatAndShorten(item.unitPrice * item.quantity)} VNĐ</p>
                      </div>
                    ))}
                    <div className="border-t border-gray-600 my-2"></div>
                    <div className="flex justify-between items-center text-lg font-bold">
                      <p>TỔNG CỘNG</p>
                      <p className="text-yellow-400">{formatAndShorten(invoiceDetails.totalAmount)} VNĐ</p>
                    </div>
                  </div>
                </div>

                {invoiceDetails.status === 0 && isUserThePayer && (
                   <button onClick={handlePayInvoice} disabled={isPaying} className="mt-8 w-full flex justify-center items-center py-3 rounded-md bg-green-600 hover:bg-green-700 disabled:bg-gray-500">
                      {isPaying ? <Loader2 className="h-5 w-5 animate-spin"/> : `Thanh Toán Ngay ${formatAndShorten(invoiceDetails.totalAmount)} VNĐ`}
                  </button>
                )}
                
                {invoiceDetails.status === 0 && !isUserThePayer && (
                    <div className="mt-6 text-center p-3 text-sm text-yellow-300 bg-yellow-900/30 rounded-lg">
                        Hoá đơn này dành cho một địa chỉ khác. Bạn không thể thanh toán nó.
                    </div>
                )}
                
                {payError && (
                  <div className="mt-6 flex items-center p-4 text-sm text-red-300 bg-red-900/30 rounded-lg">
                    <XCircle className="h-5 w-5 mr-3 flex-shrink-0" />
                    <span>{payError}</span>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
