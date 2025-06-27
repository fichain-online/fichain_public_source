// app/dapps/invoicing/create/page.tsx
'use client';

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { ethers, getBytes, formatEther, parseEther, hexlify } from 'ethers';

// --- Zustand Store Imports ---
import { useAuthStore } from '@/stores/authStore';
import { useWebSocketStore } from '@/stores/webSocketStore';

// --- UI and Icon Imports ---
import Header from '@/components/header';
import { ArrowLeft, Loader2, CheckCircle, XCircle, Repeat, FileText, User, ShoppingCart, Trash2, PlusCircle } from 'lucide-react';
import Link from 'next/link';

// --- Protobuf and Config Imports ---
import { type Receipt } from '@/proto/receipt_pb';
import config from '@/lib/config';
import InvoicingABI from '@/lib/abis/invoice.json';

import { 
  generateTransactionWithSign
} from '@/lib/utils';

// --- VM ---
import { encodeFunctionData } from 'viem';

const GAS_USE = 500000; // Creating invoices with items can be gas-intensive

// --- Type Definitions ---
interface Product {
  name: string;
  description: string;
  unitPrice: bigint;
}

interface CartItem {
  description: string; // From Product name
  unitPrice: bigint;
  quantity: bigint;
}

// --- Example Products ---
const EXAMPLE_PRODUCTS: Product[] = [
  { name: "Web Development Services", description: "10 hours of backend work", unitPrice: parseEther("1000000") },
  { name: "UI/UX Design Package", description: "Includes wireframes and mockups", unitPrice: parseEther("2500000") },
  { name: "Smart Contract Audit", description: "Security audit for one contract", unitPrice: parseEther("5000000") },
  { name: "Consulting Call", description: "1 hour strategy session", unitPrice: parseEther("500000") },
];

export default function CreateInvoicePage() {
  const router = useRouter();

  // --- State ---
  const [payerAddress, setPayerAddress] = useState('');
  const [description, setDescription] = useState('');
  const [cart, setCart] = useState<CartItem[]>([]);
  const [totalAmount, setTotalAmount] = useState<bigint>(BigInt(0));

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [receipt, setReceipt] = useState<Receipt | null>(null);
  const [createdInvoiceId, setCreatedInvoiceId] = useState<string | null>(null);

  const { isAuthenticated, _hasHydrated, privateKey } = useAuthStore();
  const { sendTransaction, getNonce, connected } = useWebSocketStore();

  // --- Auth check ---
  useEffect(() => {
    if (_hasHydrated && !isAuthenticated) {
      router.replace('/');
    }
  }, [_hasHydrated, isAuthenticated, router]);

  // --- Calculate total amount whenever cart changes ---
  useEffect(() => {
    const total = cart.reduce((acc, item) => {
      return acc + (item.unitPrice * item.quantity);
    }, BigInt(0));
    setTotalAmount(total);
  }, [cart]);

  // --- Cart Management ---
  const handleAddToCart = (product: Product) => {
    setCart(prevCart => {
      const existingItem = prevCart.find(item => item.description === product.name);
      if (existingItem) {
        // Increment quantity if item already in cart
        return prevCart.map(item => 
          item.description === product.name 
            ? { ...item, quantity: item.quantity + BigInt(1) } 
            : item
        );
      } else {
        // Add new item to cart
        return [...prevCart, { description: product.name, unitPrice: product.unitPrice, quantity: BigInt(1) }];
      }
    });
  };

  const handleRemoveFromCart = (index: number) => {
    setCart(prevCart => prevCart.filter((_, i) => i !== index));
  };
  
  // --- Transaction Handler ---
  const handleCreateInvoice = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!ethers.isAddress(payerAddress)) {
      setError("Vui lòng nhập địa chỉ người trả hợp lệ.");
      return;
    }
    if (cart.length === 0) {
      setError("Hoá đơn phải có ít nhất một sản phẩm.");
      return;
    }

    setError(null);
    setReceipt(null);
    setIsLoading(true);

    try {
      const nonce = await getNonce();
      const calldata = encodeFunctionData({
        abi: InvoicingABI,
        functionName: 'createInvoice',
        args: [
          payerAddress,
          description || "Thanh toán dịch vụ",
          cart // The cart's structure matches the required InvoiceItem[] struct
        ]
      });

      const tx = generateTransactionWithSign(
        privateKey!, BigInt(nonce), getBytes(config.invoiceContractAddress),
        ethers.toBeArray(0), // No value sent when creating
        getBytes(calldata), "Tạo Hoá Đơn",
        BigInt(GAS_USE), ethers.toBeArray(ethers.parseUnits("0.1", 'finney'))
      );

      const returnedReceipt: Receipt = await sendTransaction(tx);
      
      // Find the InvoiceCreated event to get the ID
      console.log("xx", returnedReceipt.logs)
      setCreatedInvoiceId(ethers.toBigInt(returnedReceipt.logs[0].topics[1]).toString());
      setReceipt(returnedReceipt);
    } catch (err: any) {
      setError(err.message || 'Giao dịch tạo hoá đơn thất bại.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleReset = () => {
    setReceipt(null);
    setError(null);
    setIsLoading(false);
    setPayerAddress('');
    setDescription('');
    setCart([]);
    setCreatedInvoiceId(null);
  };
  
  const formatAndShorten = (balance: bigint, decimals = 4) => {
    return new Intl.NumberFormat('vi-VN', { maximumFractionDigits: decimals }).format(parseFloat(formatEther(balance)));
  };

  if (!_hasHydrated || !isAuthenticated) {
    return <div className="flex items-center justify-center min-h-screen"><Loader2 className="h-12 w-12 animate-spin text-cyan-400" /></div>;
  }

  return (
    <main className="min-h-screen text-white">
      <Header />
      <div className="container mx-auto px-4 py-12 sm:py-16">
        <div className="max-w-4xl mx-auto">
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
                <h1 className="text-2xl font-bold">Hoá đơn đã được tạo!</h1>
                {createdInvoiceId && <p className="text-gray-300 mt-2">Mã hoá đơn của bạn là: <strong className="text-yellow-400 text-lg">{createdInvoiceId}</strong></p>}
                <p className="text-gray-400 mt-1">Gửi mã này cho người trả để họ thanh toán.</p>
                <button onClick={handleReset} className="mt-8 w-full max-w-sm mx-auto flex justify-center items-center py-3 px-4 rounded-md text-white bg-cyan-600 hover:bg-cyan-700">
                  <Repeat className="h-5 w-5 mr-2" /> Tạo hoá đơn khác
                </button>
              </div>
            ) : (
              <div>
                <h1 className="text-3xl font-bold text-yellow-400 mb-2 flex items-center">
                  <FileText className="h-8 w-8 mr-3"/> Tạo Hoá Đơn Mới
                </h1>
                <p className="text-gray-400 mb-6">Thêm sản phẩm và chỉ định người trả để tạo một hoá đơn trên chuỗi.</p>
                
                <form onSubmit={handleCreateInvoice} className="grid grid-cols-1 md:grid-cols-2 gap-8">
                  {/* Left Column: Form and Cart */}
                  <div className="space-y-6">
                    <div>
                      <h3 className="text-xl font-semibold mb-3">Thông tin hoá đơn</h3>
                      <div className="space-y-4">
                        <div>
                          <label htmlFor="payerAddress" className="flex items-center text-sm font-medium text-gray-300 mb-1"><User className="h-4 w-4 mr-1.5"/>Địa chỉ người trả</label>
                          <input id="payerAddress" type="text" value={payerAddress} onChange={e => setPayerAddress(e.target.value)} placeholder="0x..." required className="w-full bg-gray-700 border-gray-600 rounded-md py-2 px-3 focus:ring-cyan-500 focus:border-cyan-500" />
                        </div>
                        <div>
                          <label htmlFor="description" className="block text-sm font-medium text-gray-300 mb-1">Mô tả chung (tuỳ chọn)</label>
                          <input id="description" type="text" value={description} onChange={e => setDescription(e.target.value)} placeholder="VD: Hoá đơn dịch vụ Q4/2024" className="w-full bg-gray-700 border-gray-600 rounded-md py-2 px-3 focus:ring-cyan-500 focus:border-cyan-500" />
                        </div>
                      </div>
                    </div>

                    <div>
                      <h3 className="text-xl font-semibold mb-3 flex items-center"><ShoppingCart className="h-5 w-5 mr-2"/>Sản phẩm trong hoá đơn</h3>
                      <div className="space-y-2 bg-gray-900/50 p-4 rounded-lg min-h-[100px]">
                        {cart.length === 0 ? (
                          <p className="text-gray-500 text-center py-4">Chưa có sản phẩm nào.</p>
                        ) : (
                          cart.map((item, index) => (
                            <div key={index} className="flex justify-between items-center bg-gray-700/50 p-2 rounded-md">
                              <div>
                                <p className="font-semibold">{item.description}</p>
                                <p className="text-xs text-gray-400">{formatAndShorten(item.unitPrice)} VNĐ x {item.quantity.toString()}</p>
                              </div>
                              <div className="flex items-center gap-2">
                                <p className="font-bold text-cyan-400">{formatAndShorten(item.unitPrice * item.quantity)}</p>
                                <button type="button" onClick={() => handleRemoveFromCart(index)} className="p-1 text-red-400 hover:text-red-300">
                                  <Trash2 className="h-4 w-4"/>
                                </button>
                              </div>
                            </div>
                          ))
                        )}
                      </div>
                      <div className="mt-4 text-right">
                          <p className="text-gray-300">Tổng cộng:</p>
                          <p className="text-2xl font-bold text-yellow-400">{formatAndShorten(totalAmount)} VNĐ</p>
                      </div>
                    </div>
                  </div>

                  {/* Right Column: Product List */}
                  <div className="space-y-4">
                    <h3 className="text-xl font-semibold mb-3">Thêm sản phẩm mẫu</h3>
                    <div className="space-y-3">
                      {EXAMPLE_PRODUCTS.map((product, index) => (
                        <div key={index} className="bg-gray-900/50 p-4 rounded-lg flex justify-between items-center">
                          <div>
                            <p className="font-semibold">{product.name}</p>
                            <p className="text-sm text-gray-400">{formatAndShorten(product.unitPrice)} VNĐ</p>
                          </div>
                          <button type="button" onClick={() => handleAddToCart(product)} className="flex items-center gap-2 px-3 py-1.5 bg-cyan-600 hover:bg-cyan-700 rounded-md text-sm">
                            <PlusCircle className="h-4 w-4"/> Thêm
                          </button>
                        </div>
                      ))}
                    </div>
                     <div className="pt-4">
                        <button type="submit" disabled={isLoading || cart.length === 0} className="w-full flex justify-center items-center py-3 rounded-md bg-green-600 hover:bg-green-700 disabled:bg-gray-500">
                            {isLoading ? <Loader2 className="h-5 w-5 animate-spin"/> : 'Tạo Hoá Đơn'}
                        </button>
                     </div>
                  </div>
                </form>

                {error && (
                  <div className="mt-6 flex items-center p-4 text-sm text-red-300 bg-red-900/30 rounded-lg">
                    <XCircle className="h-5 w-5 mr-3 flex-shrink-0" />
                    <span>{error}</span>
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
