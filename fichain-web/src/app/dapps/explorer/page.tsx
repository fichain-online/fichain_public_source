'use client';

import React, { useState, useEffect, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { ethers, keccak256, SigningKey, hexlify, getBytes, formatEther, parseUnits, Interface, toBigInt } from 'ethers';
import Link from 'next/link';


// Zustand Store Imports
import { useAuthStore } from '@/stores/authStore';

// UI and Icon Imports
import Header from '@/components/header';
import { 
  Loader2, 
  ArrowLeft,
  ArrowUpRight,
  ArrowDownLeft,
  CheckCircle,
  XCircle,
  FileText,
  ChevronLeft,
  ChevronRight,
  SearchX
} from 'lucide-react';

// Configuration
import config from '@/lib/config';
import { useWebSocketStore } from '@/stores/webSocketStore';

// --- TYPE DEFINITIONS based on your API response ---
interface Log {
  transactionHash: string;
  logIndex: number;
  address: string;
  data: string;
  removed: boolean;
  topics: string[];
}

interface Receipt {
  transactionHash: string;
  status: number; // 1 for success, 0 for failure
  cumulativeGasUsed: number;
  gasUsed: number;
  contractAddress: string;
  logsBloom: string;
}

interface Transaction {
  hash: string;
  blockHash: string;
  blockHeight: number;
  transactionIndex: number;
  fromAddress: string;
  toAddress: string;
  nonce: number;
  amount: number | string; // API returns a number, but it's safer to handle as string for big numbers
  gasLimit: number;
  gasPrice: number | string;
  data: string;
  message: string;
  signature: string;
  logs: Log[];
  receipt?: Receipt;
}

// --- HELPER FUNCTIONS ---

// Truncates an address or hash for cleaner display
const truncateHash = (hash: string, startChars = 6, endChars = 6) => {
  if (!hash) return '';
  return `${hash.substring(0, startChars + 2)}...${hash.substring(hash.length - endChars)}`;
};


// --- REUSABLE COMPONENTS ---

// Component for a single transaction item in the list
const TransactionItem = ({ tx, currentUserAddress }: { tx: Transaction, currentUserAddress: string }) => {
  const isSender = tx.fromAddress.toLowerCase() === currentUserAddress.toLowerCase();
  
  const formattedAmount = useMemo(() => {
    try {
      // The amount is already a number in your sample, but for safety, we handle it
      // as if it were a BigInt string from a more robust API.
      const amountBigInt = formatEther(BigInt(tx.amount));
      const number = parseFloat(amountBigInt);
      return new Intl.NumberFormat('vi-VN', { maximumFractionDigits: 4 }).format(number);
    } catch {
      return 'N/A';
    }
  }, [tx.amount]);

  const TransactionIcon = isSender ? ArrowUpRight : ArrowDownLeft;
  const iconColor = isSender ? 'text-orange-400' : 'text-green-400';
  const amountColor = isSender ? 'text-orange-400' : 'text-green-400';
  const otherPartyAddress = isSender ? tx.toAddress : tx.fromAddress;

  return (
    <div className="flex items-center space-x-4 p-4 bg-gray-800/50 rounded-lg hover:bg-gray-800 transition-colors">
      <div className={`p-2 rounded-full bg-gray-700 ${iconColor}`}>
        <TransactionIcon size={20} />
      </div>

      <div className="flex-1 grid grid-cols-2 md:grid-cols-4 gap-4 items-center">
        {/* Column 1: Hash & Message */}
        <div className="flex flex-col">
          <span className="font-mono text-cyan-400 text-sm">{truncateHash(tx.hash)}</span>
          <span className="text-gray-400 text-xs mt-1">{tx.message || 'Giao dịch'}</span>
        </div>

        {/* Column 2: From/To (Mobile & Desktop) */}
        <div className="flex flex-col text-right md:text-left">
           <span className="text-xs text-gray-500">{isSender ? "Đến" : "Từ"}</span>
           <span className="font-mono text-sm">{truncateHash(otherPartyAddress)}</span>
        </div>

        {/* Column 3: Amount (Desktop) */}
        <div className="hidden md:flex flex-col text-right">
          <span className={`${amountColor} font-semibold`}>
            {isSender ? '-' : '+'} {formattedAmount} VNĐ
          </span>
          <span className="text-xs text-gray-500">
            Khối #{tx.blockHeight}
          </span>
        </div>

        {/* Column 4: Status (Desktop) */}
        <div className="hidden md:flex justify-end items-center">
          {tx.receipt?.status === 1 ? (
            <span className="flex items-center text-xs font-medium text-green-400 bg-green-900/50 px-2 py-1 rounded-md">
              <CheckCircle size={14} className="mr-1.5"/>
              Thành công
            </span>
          ) : (
            <span className="flex items-center text-xs font-medium text-red-400 bg-red-900/50 px-2 py-1 rounded-md">
              <XCircle size={14} className="mr-1.5" />
              Thất bại
            </span>
          )}
        </div>
      </div>
    </div>
  );
};


// Component for Pagination Controls
const PaginationControls = ({ currentPage, totalPages, onPageChange }: { currentPage: number, totalPages: number, onPageChange: (page: number) => void }) => {
  if (totalPages <= 1) return null;

  return (
    <div className="flex items-center justify-center space-x-4 mt-8">
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="flex items-center px-4 py-2 text-sm font-medium text-white bg-gray-700 rounded-md hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        <ChevronLeft size={16} className="mr-1" />
        Trước
      </button>

      <span className="text-gray-400">
        Trang {currentPage} trên {totalPages}
      </span>

      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="flex items-center px-4 py-2 text-sm font-medium text-white bg-gray-700 rounded-md hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        Sau
        <ChevronRight size={16} className="ml-1" />
      </button>
    </div>
  );
};


// --- MAIN EXPLORER PAGE COMPONENT ---
export default function ExplorerPage() {
  const router = useRouter();
  
  // State for data, loading, errors, and pagination
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);

  // Get user's wallet info from Zustand store
  const {  isAuthenticated, _hasHydrated } = useAuthStore();
  const { wallet,  } = useWebSocketStore();
  const PAGE_SIZE = 10;

  // Effect to fetch transactions when page or wallet changes
  useEffect(() => {
    // Don't fetch until the store is hydrated and the user is authenticated
    if (!_hasHydrated || !isAuthenticated || !wallet?.address) {
      return;
    }

    const fetchTransactions = async () => {
      setIsLoading(true);
      setError(null);
      
      try {
        const url = `${config.explorerApiBaseUrl}/transaction/${wallet.address}?page=${currentPage}&pageSize=${PAGE_SIZE}`;
        const response = await fetch(url);

        if (!response.ok) {
          throw new Error(`Lỗi API: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();
        
        setTransactions(data.data || []);
        setTotalPages(Math.ceil(data.total / data.pageSize));

      } catch (err: any) {
        console.error("Failed to fetch transactions:", err);
        setError("Không thể tải lịch sử giao dịch. Vui lòng thử lại sau.");
      } finally {
        setIsLoading(false);
      }
    };

    fetchTransactions();
  }, [currentPage, wallet?.address, _hasHydrated, isAuthenticated]);

  // Auth check effect
  useEffect(() => {
    if (_hasHydrated && !isAuthenticated) {
      router.replace('/');
    }
  }, [_hasHydrated, isAuthenticated, router]);


  // Handle page changes from pagination controls
  const handlePageChange = (newPage: number) => {
    if (newPage > 0 && newPage <= totalPages) {
      setCurrentPage(newPage);
    }
  };


  // Main render logic
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
          <Link
            href="/dapps"
            className="inline-flex items-center text-cyan-400 hover:text-cyan-300 mb-8 group p-2 rounded-md hover:bg-cyan-900/50 transition-colors -ml-2 z-[10]"
          >
            <ArrowLeft className="h-5 w-5 mr-2 transition-transform group-hover:-translate-x-1" />
            Quay lại danh sách dApps
          </Link>
          
          <div className="bg-gray-900/70 p-6 md:p-8 rounded-xl shadow-2xl backdrop-blur-sm border border-gray-700/50">
            <div className="flex items-center mb-6">
                <FileText className="h-7 w-7 text-cyan-400 mr-4" />
                <h1 className="text-3xl font-bold">Lịch sử Giao dịch</h1>
            </div>

            {isLoading && (
              <div className="flex justify-center items-center py-20">
                <Loader2 className="h-10 w-10 animate-spin text-cyan-400" />
              </div>
            )}

            {!isLoading && error && (
              <div className="text-center py-20 text-red-400">
                <p>{error}</p>
              </div>
            )}

            {!isLoading && !error && transactions.length === 0 && (
              <div className="text-center py-20 text-gray-500">
                <SearchX size={48} className="mx-auto mb-4" />
                <h3 className="text-xl font-semibold text-gray-300">Không tìm thấy giao dịch nào</h3>
                <p className="mt-2">Lịch sử giao dịch của bạn sẽ xuất hiện ở đây.</p>
              </div>
            )}

            {!isLoading && !error && transactions.length > 0 && (
              <div className="space-y-3">
                {transactions.map(tx => (
                  <TransactionItem key={tx.hash} tx={tx} currentUserAddress={wallet?.address || ''} />
                ))}
              </div>
            )}
            
            <PaginationControls
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={handlePageChange}
            />
          </div>
        </div>
      </div>
    </main>
  );
}
