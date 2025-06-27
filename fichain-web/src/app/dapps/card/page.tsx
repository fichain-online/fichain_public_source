// app/dapps/card/page.tsx
'use client';

import React, { useState, useMemo, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { formatEther } from 'ethers';

// --- Zustand Store Imports ---
import { useAuthStore } from '@/stores/authStore';
import { useWebSocketStore } from '@/stores/webSocketStore';

// --- UI and Icon Imports ---
import Header from '@/components/header';
import Link from 'next/link';
import { 
  ArrowLeft, Construction, CreditCard, Lock, Unlock, Eye, EyeOff, ShieldAlert, 
  Target, Landmark, CheckCircle, ArrowDownCircle, ArrowUpCircle, Wallet 
} from 'lucide-react';

// --- MOCK DATA ---
const mockCardData = {
  last4: '1234',
  expiry: '12/27',
  cardholder: 'Your Name Here',
  type: 'DEBIT',
};

const mockTransactions = [
  { id: 1, description: 'Netflix Subscription', amount: '250000', type: 'debit', date: '2024-05-20T10:30:00Z' },
  { id: 2, description: 'Salary Deposit', amount: '25000000', type: 'credit', date: '2024-05-20T09:00:00Z' },
  { id: 3, description: 'Coffee Shop', amount: '75000', type: 'debit', date: '2024-05-19T15:00:00Z' },
  { id: 4, description: 'Online Shopping', amount: '1250000', type: 'debit', date: '2024-05-18T20:15:00Z' },
  { id: 5, description: 'Friend Transfer', amount: '500000', type: 'credit', date: '2024-05-17T11:45:00Z' },
];
// --- END MOCK DATA ---

export default function CardPage() {
  const router = useRouter();

  // --- State for Simulated Functions ---
  const [isCardLocked, setIsCardLocked] = useState(false);
  const [showPin, setShowPin] = useState(false);
  const [spendingLimit, setSpendingLimit] = useState('50000000'); // Mock limit
  const [balance, setBalance] = useState<bigint>(BigInt(0)); // ✅ Initialized with 0n

  // --- Get state from stores ---
  const { isAuthenticated, _hasHydrated } = useAuthStore();
  const { getBalance, wallet, connected } = useWebSocketStore(); // Get real balance and wallet

  // --- Auth check effect ---
  useEffect(() => {
    if (_hasHydrated && !isAuthenticated) {
      router.replace('/');
    }
  }, [_hasHydrated, isAuthenticated, router]);

  // --- Simulated Action Handlers ---
  const handleToggleLock = () => {
    setIsCardLocked(prev => !prev);
    alert(isCardLocked ? 'Đã mô phỏng mở khoá thẻ.' : 'Đã mô phỏng khoá thẻ thành công.');
  };

  const handleShowPin = () => {
    if (!showPin) {
      alert('PIN giả lập của bạn là: 1234');
    }
    setShowPin(prev => !prev);
  };

  const handleSetLimit = () => {
    const newLimit = prompt('Nhập hạn mức chi tiêu mới (chỉ giả lập):', spendingLimit);
    if (newLimit && !isNaN(Number(newLimit))) {
      setSpendingLimit(newLimit);
      alert(`Đã mô phỏng đặt hạn mức mới: ${new Intl.NumberFormat('vi-VN').format(Number(newLimit))} VNĐ`);
    } else if (newLimit) {
      alert('Vui lòng nhập một số hợp lệ.');
    }
  };

  const handleReportLost = () => {
    setIsCardLocked(true);
    alert('Đã ghi nhận yêu cầu khoá thẻ và cấp lại. Thẻ của bạn đã được khoá (giả lập).');
  };

  // --- Auth check effect ---
  useEffect(() => {
    if (_hasHydrated && !isAuthenticated) {
      router.replace('/');
    }
  }, [_hasHydrated, isAuthenticated, router]);

  const fetchData = useCallback(async() => {
    if (connected && wallet) { // ✅ Added wallet check
      const fetchedBalance = await getBalance();
      setBalance(fetchedBalance);
    }
  },[connected, wallet])

  useEffect(() => {
    fetchData()
  }, [connected]);

  // --- Display Logic ---
  const formattedBalance = useMemo(() => {
    try {
      const etherString = formatEther(balance);
      const number = parseFloat(etherString);
      return new Intl.NumberFormat('vi-VN', { maximumFractionDigits: 4 }).format(number);
    } catch {
      return '0';
    }
  }, [balance]);

  return (
    <main className="min-h-screen text-white">
      <Header />
      <div className="container mx-auto px-4 py-12 sm:py-16">
        <div className="max-w-2xl mx-auto">
          <Link href="/dapps" className="inline-flex items-center text-cyan-400 hover:text-cyan-300 mb-8 group p-2 rounded-md hover:bg-cyan-900/50 transition-colors -ml-2">
            <ArrowLeft className="h-5 w-5 mr-2 transition-transform group-hover:-translate-x-1" />
            Quay lại danh sách dApps
          </Link>

          {/* --- DEVELOPMENT NOTICE --- */}
          <div className="p-4 mb-8 text-yellow-300 bg-yellow-900/50 border border-yellow-600 rounded-lg flex items-center">
            <Construction className="h-6 w-6 mr-4 flex-shrink-0" />
            <div>
              <h3 className="font-bold">Tính năng đang được phát triển</h3>
              <p className="text-sm">Dữ liệu thẻ và các chức năng trên trang này chỉ mang tính chất minh hoạ và chưa được kết nối với hệ thống thật.</p>
            </div>
          </div>
          
          <div className="space-y-8">
            {/* --- VISUAL CARD COMPONENT --- */}
            <div className={`relative p-6 rounded-2xl shadow-2xl overflow-hidden transition-all duration-300 ${isCardLocked ? 'bg-gray-600' : 'bg-gradient-to-br from-cyan-500 to-blue-700'}`}>
              {isCardLocked && (
                <div className="absolute inset-0 bg-black/50 flex items-center justify-center z-10">
                  <Lock className="h-16 w-16 text-white/50" />
                </div>
              )}
              <div className="flex justify-between items-start">
                <span className="font-bold text-xl">My App Bank</span>
                <Landmark className="h-7 w-7" />
              </div>
              <div className="mt-8">
                <div className="w-12 h-8 bg-yellow-400/80 rounded-md"></div>
              </div>
              <p className="mt-4 text-2xl font-mono tracking-widest">
                **** **** **** {mockCardData.last4}
              </p>
              <div className="mt-4 flex justify-between items-end">
                <div>
                  <p className="text-xs opacity-70">Card Holder</p>
                  <p className="font-medium">{wallet?.address.slice(0, 6)}...{wallet?.address.slice(-4)}</p>
                </div>
                <div>
                  <p className="text-xs opacity-70">Expires</p>
                  <p className="font-medium">{mockCardData.expiry}</p>
                </div>
              </div>
            </div>

            {/* --- ACTION BUTTONS --- */}
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 text-center">
              <button onClick={handleToggleLock} className="flex flex-col items-center p-3 bg-gray-700/50 rounded-lg hover:bg-gray-700 transition-colors">
                {isCardLocked ? <Unlock className="h-6 w-6 mb-2 text-green-400" /> : <Lock className="h-6 w-6 mb-2 text-red-400" />}
                <span className="text-sm">{isCardLocked ? 'Mở khoá thẻ' : 'Khoá thẻ'}</span>
              </button>
              <button onClick={handleShowPin} className="flex flex-col items-center p-3 bg-gray-700/50 rounded-lg hover:bg-gray-700 transition-colors">
                {showPin ? <EyeOff className="h-6 w-6 mb-2 text-cyan-400" /> : <Eye className="h-6 w-6 mb-2 text-cyan-400" />}
                <span className="text-sm">{showPin ? 'Ẩn mã PIN' : 'Xem mã PIN'}</span>
              </button>
              <button onClick={handleSetLimit} className="flex flex-col items-center p-3 bg-gray-700/50 rounded-lg hover:bg-gray-700 transition-colors">
                <Target className="h-6 w-6 mb-2 text-cyan-400" />
                <span className="text-sm">Hạn mức</span>
              </button>
              <button onClick={handleReportLost} className="flex flex-col items-center p-3 bg-gray-700/50 rounded-lg hover:bg-gray-700 transition-colors">
                <ShieldAlert className="h-6 w-6 mb-2 text-orange-400" />
                <span className="text-sm">Báo mất thẻ</span>
              </button>
            </div>

            {/* --- LINKED ACCOUNT --- */}
            <div>
              <h2 className="text-xl font-bold mb-4">Tài khoản liên kết</h2>
              <div className="bg-gray-800 p-4 rounded-lg flex justify-between items-center">
                <div className="flex items-center">
                  <Wallet className="h-8 w-8 mr-4 text-cyan-400" />
                  <div>
                    <p className="text-gray-400 text-sm">Số dư ví Blockchain</p>
                    <p className="font-mono text-xs break-all">{wallet?.address}</p>
                  </div>
                </div>
                <p className="text-lg font-bold text-cyan-300">{formattedBalance} VNĐ</p>
              </div>
            </div>

            {/* --- RECENT TRANSACTIONS --- */}
            <div>
              <h2 className="text-xl font-bold mb-4">Hoạt động gần đây</h2>
              <div className="space-y-3">
                {mockTransactions.map(tx => (
                  <div key={tx.id} className="bg-gray-800 p-4 rounded-lg flex items-center justify-between">
                    <div className="flex items-center">
                      <div className={`p-2 rounded-full mr-4 ${tx.type === 'credit' ? 'bg-green-900/50' : 'bg-red-900/50'}`}>
                        {tx.type === 'credit' ? <ArrowUpCircle className="h-5 w-5 text-green-400" /> : <ArrowDownCircle className="h-5 w-5 text-red-400" />}
                      </div>
                      <div>
                        <p className="font-medium">{tx.description}</p>
                        <p className="text-xs text-gray-400">{new Date(tx.date).toLocaleString('vi-VN')}</p>
                      </div>
                    </div>
                    <p className={`font-semibold ${tx.type === 'credit' ? 'text-green-400' : 'text-red-400'}`}>
                      {tx.type === 'credit' ? '+' : '-'} {new Intl.NumberFormat('vi-VN').format(Number(tx.amount))} VNĐ
                    </p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
