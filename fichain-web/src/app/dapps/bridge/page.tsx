'use client';

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { formatEther } from 'ethers';
import { Copy, Check } from 'lucide-react';

// Zustand Store Imports
import { useAuthStore } from '@/stores/authStore';
import { useWebSocketStore } from '@/stores/webSocketStore';

// UI and Icon Imports
import Header from '@/components/header';
import { 
  Loader2, ArrowLeft, ArrowRight, CheckCircle, XCircle, Clock, RefreshCw, 
  ChevronLeft, ChevronRight, SearchX, Link as LinkIcon, Wallet
} from 'lucide-react';

// Configuration
import config from '@/lib/config';

// --- TYPE DEFINITIONS ---
interface DepositLog {
  id: number;
  sourceChainTxHash: string;
  fichainAddress: string;
  tokenName: string;
  amount: string;
  destChainTxHash: string;
  status: 'pending' | 'processing' | 'success' | 'failed';
  errorMessage: string;
  createdAt: string;
}

interface DepositWallet {
    address: string;
    tokenName: string;
}

interface PaginatedResponse {
  data: DepositLog[];
  pagination: { total: number; page: string; limit: string; };
}

// --- HELPER FUNCTIONS ---
const truncateHash = (hash: string, startChars = 8, endChars = 8) => {
  if (!hash || hash.length < startChars + endChars + 2) return hash || '';
  return `${hash.substring(0, startChars + 2)}...${hash.substring(hash.length - endChars)}`;
};

// --- CUSTOM HOOK for Copy-to-Clipboard ---
const useCopyToClipboard = () => {
    const [copied, setCopied] = useState(false);
  
    const copy = useCallback((text: string) => {
      if (!text) return;
      navigator.clipboard.writeText(text).then(() => {
        setCopied(true);
        setTimeout(() => setCopied(false), 2000); // Reset after 2 seconds
      });
    }, []);
  
    return { copied, copy };
};

// --- REUSABLE UI COMPONENTS ---
const StatusBadge = ({ status }: { status: DepositLog['status'] }) => {
    switch (status) {
        case 'success': return <span className="flex items-center text-xs font-medium text-green-400 bg-green-900/50 px-2 py-1 rounded-md"><CheckCircle size={14} className="mr-1.5"/>Thành công</span>;
        case 'failed': return <span className="flex items-center text-xs font-medium text-red-400 bg-red-900/50 px-2 py-1 rounded-md"><XCircle size={14} className="mr-1.5"/>Thất bại</span>;
        case 'processing': return <span className="flex items-center text-xs font-medium text-blue-400 bg-blue-900/50 px-2 py-1 rounded-md"><RefreshCw size={14} className="mr-1.5 animate-spin"/>Đang xử lý</span>;
        default: return <span className="flex items-center text-xs font-medium text-yellow-400 bg-yellow-900/50 px-2 py-1 rounded-md"><Clock size={14} className="mr-1.5"/>Đang chờ</span>;
    }
};

const BridgeLogItem = ({ log }: { log: DepositLog }) => {
  const formattedAmount = useMemo(() => new Intl.NumberFormat('vi-VN',{
    minimumFractionDigits: 0,
    maximumFractionDigits: 8,
  }).format(Number(formatEther(BigInt(log.amount || '0')))), [log.amount]);
    const bscExplorerUrl = `https://testnet.bscscan.com/tx/${log.sourceChainTxHash}`;
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 items-center gap-4 p-4 bg-gray-800/50 rounded-lg">
        <div className="flex flex-col">
          <span className="text-lg font-bold text-white">{formattedAmount} {log.tokenName}</span>
          <span className="text-xs text-gray-500 mt-1">{new Date(log.createdAt).toLocaleString()}</span>
        </div>
        <div className="flex flex-col space-y-2 text-sm">
          <div className="flex items-center">
            <span className="font-semibold text-gray-400 w-24">Từ (BSC):</span>
            <a href={bscExplorerUrl} target="_blank" rel="noopener noreferrer" className="font-mono text-purple-400 hover:underline flex items-center">
              {truncateHash(log.sourceChainTxHash, 6, 6)} <LinkIcon size={12} className="ml-1.5"/>
            </a>
          </div>
          <div className="flex items-center">
            <span className="font-semibold text-gray-400 w-24">Đến (Fichain):</span>
            <span className="font-mono text-green-400">{truncateHash(log.destChainTxHash, 6, 6)}</span>
          </div>
        </div>
        <div className="flex justify-start md:justify-end"><StatusBadge status={log.status} /></div>
      </div>
    );
};

const PaginationControls = ({ currentPage, totalPages, onPageChange }: { currentPage: number, totalPages: number, onPageChange: (page: number) => void }) => {
    if (totalPages <= 1) return null;
    return (
      <div className="flex items-center justify-center space-x-4 mt-8">
        <button onClick={() => onPageChange(currentPage - 1)} disabled={currentPage === 1} className="flex items-center px-4 py-2 text-sm font-medium text-white bg-gray-700 rounded-md hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
          <ChevronLeft size={16} className="mr-1" /> Trước
        </button>
        <span className="text-gray-400"> Trang {currentPage} trên {totalPages} </span>
        <button onClick={() => onPageChange(currentPage + 1)} disabled={currentPage === totalPages} className="flex items-center px-4 py-2 text-sm font-medium text-white bg-gray-700 rounded-md hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
          Sau <ChevronRight size={16} className="ml-1" />
        </button>
      </div>
    );
  };


// --- MAIN BRIDGE PAGE COMPONENT ---
export default function BridgePage() {
  const router = useRouter();
  
  // State for Deposit section
  const [selectedToken, setSelectedToken] = useState('USDT');
  const [depositWallet, setDepositWallet] = useState<DepositWallet | null>(null);
  const [isLoadingWallet, setIsLoadingWallet] = useState(false);
  const [walletError, setWalletError] = useState<string | null>(null);
  const { copied, copy } = useCopyToClipboard();

  // State for History section
  const [logs, setLogs] = useState<DepositLog[]>([]);
  const [isLoadingHistory, setIsLoadingHistory] = useState(true);
  const [historyError, setHistoryError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);

  // Global state
  const { isAuthenticated, _hasHydrated } = useAuthStore();
  const { wallet } = useWebSocketStore();

  const BRIDGEABLE_TOKENS = ['USDT', 'BTC', 'ETH'];
  const HISTORY_PAGE_LIMIT = 5;

  // --- DATA FETCHING EFFECTS ---

  // Effect to fetch the deposit wallet address when token changes
  useEffect(() => {
    if (!_hasHydrated || !isAuthenticated || !wallet?.address || !selectedToken) return;

    const fetchDepositWallet = async () => {
        setIsLoadingWallet(true);
        setWalletError(null);
        setDepositWallet(null);
        try {
            const url = `${config.bridgeApiBaseUrl}/deposit-wallet/${selectedToken.toLowerCase()}/${wallet.address}`;
            const response = await fetch(url);
            if (!response.ok) throw new Error("Không thể lấy địa chỉ nạp tiền.");
            
            const data: DepositWallet = await response.json();
            setDepositWallet(data);
        } catch (err: any) {
            setWalletError(err.message);
        } finally {
            setIsLoadingWallet(false);
        }
    };

    fetchDepositWallet();
  }, [wallet?.address, selectedToken, _hasHydrated, isAuthenticated]);

  // Effect to fetch bridge history when page or wallet changes
  useEffect(() => {
    if (!_hasHydrated || !isAuthenticated || !wallet?.address) return;

    const fetchLogs = async () => {
        setIsLoadingHistory(true);
        setHistoryError(null);
        try {
            const params = new URLSearchParams({ page: currentPage.toString(), limit: HISTORY_PAGE_LIMIT.toString() });
            const url = `${config.bridgeApiBaseUrl}/deposit-logs/${wallet.address}?${params.toString()}`;
            const response = await fetch(url);
            if (!response.ok) throw new Error("Lỗi API khi tải lịch sử.");
            
            const data: PaginatedResponse = await response.json();
            setLogs(data.data || []);
            setTotalPages(Math.ceil((data.pagination.total || 0) / HISTORY_PAGE_LIMIT));
        } catch (err: any) {
            setHistoryError(err.message);
        } finally {
            setIsLoadingHistory(false);
        }
    };
    fetchLogs(); // Initial fetch
    // const interval = setInterval(fetchLogs, 15000); // Poll for updates every 15 seconds

    // return () => clearInterval(interval); // Cleanup interval on component unmount
  }, [currentPage, wallet?.address, _hasHydrated, isAuthenticated]);

  // Auth check effect
  useEffect(() => {
    if (_hasHydrated && !isAuthenticated) {
      router.replace('/');
    }
  }, [_hasHydrated, isAuthenticated, router]);

  // --- EVENT HANDLERS ---
  const handlePageChange = (newPage: number) => {
    if (newPage > 0 && newPage <= totalPages) setCurrentPage(newPage);
  };

  // Main render logic
  if (!_hasHydrated || !isAuthenticated) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-900">
        <Loader2 className="h-12 w-12 animate-spin text-cyan-400" />
      </div>
    );
  }

  return (
    <main className="min-h-screen text-white">
      <Header />
      <div className="container mx-auto px-4 py-12 sm:py-16">
        <div className="max-w-4xl mx-auto">
          <Link href="/dapps" className="inline-flex items-center text-cyan-400 hover:text-cyan-300 mb-8 group p-2 rounded-md hover:bg-cyan-900/50 transition-colors -ml-2 z-[10]">
            <ArrowLeft className="h-5 w-5 mr-2 transition-transform group-hover:-translate-x-1" />
            Quay lại dApps
          </Link>
          
          {/* Main Content Card */}
          <div className="bg-gray-900/70 p-6 md:p-8 rounded-xl shadow-2xl backdrop-blur-sm border border-gray-700/50">
            
            {/* --- DEPOSIT SECTION --- */}
            <div className="text-center">
              <div className="flex items-center justify-center mb-4">
                <Wallet size={28} className="text-cyan-400 mr-4"/>
                <h1 className="text-3xl font-bold">Bridge Token</h1>
              </div>
              <p className="text-gray-400 max-w-lg mx-auto">
                  Chọn một token, sau đó gửi token (mạng BEP-20) đến địa chỉ BSC được cung cấp để nhận token tương ứng trên Fichain.
              </p>
            
              <div className="mt-6">
                  <label htmlFor="token-select" className="block text-sm font-medium text-gray-300 mb-2">Chọn Token</label>
                  <select 
                    id="token-select"
                    value={selectedToken} 
                    onChange={(e) => setSelectedToken(e.target.value)}
                    className="bg-gray-800 border border-gray-600 text-white text-lg rounded-lg focus:ring-cyan-500 focus:border-cyan-500 block w-full max-w-xs mx-auto p-3"
                  >
                    {BRIDGEABLE_TOKENS.map(token => <option key={token} value={token}>{token}</option>)}
                  </select>
              </div>
            
              <div className="mt-8 p-6 bg-gray-800/70 rounded-lg w-full max-w-2xl mx-auto border border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-300">Gửi {selectedToken} đến địa chỉ BSC này:</h3>
                  {isLoadingWallet && <div className="flex justify-center p-4"><Loader2 className="animate-spin text-cyan-400" /></div>}
                  {walletError && <div className="p-4 text-red-400">{walletError}</div>}
                  {depositWallet && (
                      <div className="mt-4 p-4 bg-black/30 rounded-md flex items-center justify-between font-mono text-cyan-300 text-sm md:text-base break-all">
                          <span>{depositWallet.address}</span>
                          <button onClick={() => copy(depositWallet.address)} className="ml-4 p-2 rounded-md hover:bg-cyan-900/50 transition-colors">
                              {copied ? <Check size={18} className="text-green-400" /> : <Copy size={18} />}
                          </button>
                      </div>
                  )}
                   <p className="text-xs text-gray-500 mt-4">
                      <span className='font-bold text-yellow-400'>Cảnh báo:</span> Chỉ gửi {selectedToken} (BEP-20) đến địa chỉ này. Gửi bất kỳ token nào khác có thể dẫn đến mất tiền vĩnh viễn.
                  </p>
              </div>
            </div>

            {/* --- DIVIDER & HISTORY SECTION --- */}
            <div className="mt-12 pt-8 border-t border-gray-700">
              <h2 className="text-2xl font-bold text-center mb-6">Lịch sử Giao dịch Bridge gần đây</h2>
              
              {isLoadingHistory && <div className="flex justify-center py-12"><Loader2 className="h-10 w-10 animate-spin text-cyan-400"/></div>}
              {historyError && <div className="text-center py-12 text-red-400">{historyError}</div>}
              
              {!isLoadingHistory && !historyError && logs.length === 0 && (
                <div className="text-center py-12 text-gray-500">
                  <SearchX size={48} className="mx-auto mb-4" />
                  <h3 className="text-xl font-semibold text-gray-300">Không có lịch sử bridge</h3>
                </div>
              )}

              {!isLoadingHistory && !historyError && logs.length > 0 && (
                <div className="space-y-3">
                  {logs.map(log => <BridgeLogItem key={log.id} log={log} />)}
                  <PaginationControls currentPage={currentPage} totalPages={totalPages} onPageChange={handlePageChange} />
                </div>
              )}
            </div>

          </div>
        </div>
      </div>
    </main>
  );
}
