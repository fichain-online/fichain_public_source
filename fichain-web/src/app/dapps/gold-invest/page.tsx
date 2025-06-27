// app/dapps/gold/page.tsx
'use client';

import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { ethers, hexlify, getBytes, formatEther, parseUnits, toBigInt } from 'ethers';

// --- Zustand Store Imports ---
import { useAuthStore } from '@/stores/authStore';
import { useWebSocketStore } from '@/stores/webSocketStore';

// --- UI and Icon Imports ---
import Header from '@/components/header';
import { ArrowLeft, Loader2, CheckCircle, XCircle, Repeat, Gem, Wallet, ArrowRightLeft } from 'lucide-react';
import Link from 'next/link';

// --- Protobuf and Config Imports ---
import { type Receipt } from '@/proto/receipt_pb';
import config from '@/lib/config';
import GoldInvestABI from '@/lib/abis/gold_invest.json';
import SimpleGoldABI from '@/lib/abis/gold_token.json'; // The ERC20 token ABI

import { 
  generateCallDataWithSign,
  generateTransactionWithSign
} from '@/lib/utils';

// --- VM ---
import { encodeFunctionData, decodeFunctionResult } from 'viem';

const GAS_USE = 300000; // Gas might be higher for these interactions

export default function GoldInvestPage() {
  const router = useRouter();

  // --- State for Data from Blockchain ---
  const [isDataLoading, setIsDataLoading] = useState(true);
  const [ethBalance, setEthBalance] = useState<bigint>(BigInt(0));
  const [goldBalance, setGoldBalance] = useState<bigint>(BigInt(0));
  const [buyPrice, setBuyPrice] = useState<bigint>(BigInt(0));
  const [sellPrice, setSellPrice] = useState<bigint>(BigInt(0));

  // --- State for Forms ---
  const [activeTab, setActiveTab] = useState<'buy' | 'sell'>('buy');
  const [buyAmountEth, setBuyAmountEth] = useState(''); // Amount of ETH to spend
  const [sellAmountGold, setSellAmountGold] = useState(''); // Amount of Gold to sell
  
  // --- State for Transaction Submission ---
  const [isLoading, setIsLoading] = useState(false);
  const [activeTx, setActiveTx] = useState<'buy' | 'sell' | 'approve' | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [receipt, setReceipt] = useState<Receipt | null>(null);
  const [approvalSuccess, setApprovalSuccess] = useState(false);

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
    if (connected && wallet && privateKey) {
      setIsDataLoading(true);
      setError(null);
      try {
        // 1. Fetch ETH Balance
        const fetchedEthBalance = await getBalance();
        setEthBalance(fetchedEthBalance);

        // 2. Fetch Gold Prices from GoldInvest contract
        const pricesCalldata = encodeFunctionData({ abi: GoldInvestABI, functionName: 'getPrices' });
        const pricesCallRq = generateCallDataWithSign(privateKey, getBytes(config.goldInvestContractAddress), getBytes(pricesCalldata));
        const pricesCallRs = await callSmartContract(pricesCallRq);
        const [fetchedBuyPrice, fetchedSellPrice] = decodeFunctionResult({
          abi: GoldInvestABI,
          functionName: 'getPrices',
          data: hexlify(pricesCallRs.Data) as `0x${string}`,
        }) as [bigint, bigint];
        setBuyPrice(fetchedBuyPrice);
        setSellPrice(fetchedSellPrice);

        // 3. Fetch user's Gold token balance from SimpleGold contract
        const goldBalanceCalldata = encodeFunctionData({
            abi: SimpleGoldABI,
            functionName: 'balanceOf',
            args: [wallet.address]
        });
        const goldBalanceCallRq = generateCallDataWithSign(privateKey, getBytes(config.goldTokenContractAddress), getBytes(goldBalanceCalldata));
        const goldBalanceCallRs = await callSmartContract(goldBalanceCallRq);
        const fetchedGoldBalance = decodeFunctionResult({
          abi: SimpleGoldABI,
          functionName: 'balanceOf',
          data: hexlify(goldBalanceCallRs.Data) as `0x${string}`,
        }) as bigint;
        setGoldBalance(fetchedGoldBalance);

      } catch (err) {
        console.error("Failed to fetch contract data:", err);
        setError("Không thể tải dữ liệu. Vui lòng thử lại.");
      } finally {
        setIsDataLoading(false);
      }
    }
  }, [connected, wallet, privateKey, getBalance, callSmartContract]);

  useEffect(() => {
    if (connected) {
      fetchData();
    }
  }, [connected, fetchData]);

  // --- Transaction Handlers ---
  const handleBuyGold = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!buyAmountEth || parseFloat(buyAmountEth) <= 0) {
        setError("Vui lòng nhập số lượng hợp lệ.");
        return;
    }
    setError(null);
    setReceipt(null);
    setIsLoading(true);
    setActiveTx('buy');

    try {
        const nonce = await getNonce();
        const calldata = encodeFunctionData({ abi: GoldInvestABI, functionName: 'buyGold' });
        const value = parseUnits(buyAmountEth, 'ether'); 

        const tx = generateTransactionWithSign(
            privateKey!, BigInt(nonce), getBytes(config.goldInvestContractAddress),
            ethers.toBeArray(value), getBytes(calldata), "Mua Vàng",
            BigInt(GAS_USE), ethers.toBeArray(ethers.parseUnits("0.1", 'finney'))
        );
        const returnedReceipt: Receipt = await sendTransaction(tx);
        setReceipt(returnedReceipt);
    } catch (err: any) {
        setError(err.message || 'Giao dịch mua vàng thất bại.');
    } finally {
        setIsLoading(false);
        setActiveTx(null);
        setBuyAmountEth('');
    }
  };

  const handleApprove = async () => {
    if (!sellAmountGold || parseFloat(sellAmountGold) <= 0) {
        setError("Vui lòng nhập số lượng hợp lệ.");
        return;
    }
    setError(null);
    setReceipt(null);
    setIsLoading(true);
    setActiveTx('approve');
    setApprovalSuccess(false);

    try {
        const nonce = await getNonce();
        // parseUnits correctly handles float strings like "1.5"
        const amountToApprove = parseUnits(sellAmountGold, 18);
        const calldata = encodeFunctionData({
            abi: SimpleGoldABI, 
            functionName: 'approve',
            args: [config.goldInvestContractAddress, amountToApprove]
        });

        const tx = generateTransactionWithSign(
            privateKey!, BigInt(nonce), getBytes(config.goldTokenContractAddress),
            ethers.toBeArray(0), getBytes(calldata), "Uỷ quyền bán Vàng",
            BigInt(GAS_USE), ethers.toBeArray(ethers.parseUnits("0.1", 'finney'))
        );
        const returnedReceipt: Receipt = await sendTransaction(tx);
        if (returnedReceipt.status === 1) {
            setApprovalSuccess(true);
            setError(null); 
        } else {
            setError("Giao dịch uỷ quyền thất bại.");
        }
        // setReceipt(returnedReceipt);
    } catch (err: any) {
        setError(err.message || 'Giao dịch uỷ quyền thất bại.');
    } finally {
        setIsLoading(false);
        setActiveTx(null);
    }
  };

  const handleSellGold = async () => {
     if (!sellAmountGold || parseFloat(sellAmountGold) <= 0) {
        setError("Vui lòng nhập số lượng hợp lệ.");
        return;
    }
    setError(null);
    setReceipt(null);
    setIsLoading(true);
    setActiveTx('sell');

    try {
        const nonce = await getNonce();
        const amountToSell = parseUnits(sellAmountGold, 18);
        const calldata = encodeFunctionData({
            abi: GoldInvestABI,
            functionName: 'sellGold',
            args: [amountToSell]
        });

        const tx = generateTransactionWithSign(
            privateKey!, BigInt(nonce), getBytes(config.goldInvestContractAddress),
            ethers.toBeArray(0), getBytes(calldata), "Bán Vàng",
            BigInt(GAS_USE), ethers.toBeArray(ethers.parseUnits("0.1", 'finney'))
        );
        const returnedReceipt: Receipt = await sendTransaction(tx);
        setReceipt(returnedReceipt);
    } catch (err: any) {
        setError(err.message || 'Giao dịch bán vàng thất bại.');
    } finally {
        setIsLoading(false);
        setActiveTx(null);
        setSellAmountGold('');
        setApprovalSuccess(false);
    }
  };

  // --- UI Logic and Helpers ---
  const handleReset = () => {
    setReceipt(null);
    setError(null);
    setIsLoading(false);
    setActiveTx(null);
    setApprovalSuccess(false);
    fetchData(); 
  };
  
  const formatAndShorten = (balance: bigint, decimals = 4) => {
    try {
      const formatted = formatEther(balance);
      const num = parseFloat(formatted);
      // if (num > 0 && num < 0.0001) return num.toExponential(2);
      return new Intl.NumberFormat('vi-VN', { maximumFractionDigits: decimals }).format(num);
    } catch {
      return '0';
    }
  };

  const formattedBuyAmountEth = useMemo(() => {
    if (buyAmountEth === '') return '';
    try { return new Intl.NumberFormat('vi-VN').format(parseInt(buyAmountEth, 10)); }
    catch { return buyAmountEth; }
  }, [buyAmountEth]);

  // UPDATED: This formatting logic now correctly handles floats and trailing dots.
  const formattedSellAmountGold = useMemo(() => {
    if (sellAmountGold === '') return '';
    const [integerPart, decimalPart] = sellAmountGold.split('.');

    const formattedInteger = new Intl.NumberFormat('vi-VN').format(
        BigInt(integerPart || '0')
    );
    
    if (decimalPart !== undefined) {
        return `${formattedInteger}.${decimalPart}`;
    }
    if (sellAmountGold.endsWith('.')) {
        return `${formattedInteger}.`;
    }
    return formattedInteger;
  }, [sellAmountGold]);

  const handleBuyAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const rawValue = e.target.value.replace(/[^0-9]/g, '');
    setBuyAmountEth(rawValue);
  };
  
  // UPDATED: This handler now allows a single decimal point for float input.
  const handleSellAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    // Allow numbers and only one dot
    if (/^[0-9]*\.?[0-9]*$/.test(value)) {
        setSellAmountGold(value);
    }
  };

  const goldToReceive = useMemo(() => {
    if (!buyAmountEth || !buyPrice || buyPrice === BigInt(0)) return "0";
    try {
      const amountInWei = parseUnits(buyAmountEth, 'ether');
      const goldAmount = (amountInWei * BigInt(10**18)) / buyPrice;
      return formatAndShorten(goldAmount);
    } catch { return "0"; }
  }, [buyAmountEth, buyPrice]);

  const ethToReceive = useMemo(() => {
    if (!sellAmountGold || !sellPrice || sellPrice === BigInt(0)) return "0";
    try {
        const amountInGoldWei = parseUnits(sellAmountGold, 18);
        const ethAmount = (amountInGoldWei * sellPrice) / BigInt(10**18);
        return formatAndShorten(ethAmount);
    } catch { return "0"; }
  }, [sellAmountGold, sellPrice]);

  // --- Auth Loading Screen ---
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
              <div>
                <div className="flex flex-col items-center text-center">
                  <CheckCircle className="h-16 w-16 text-green-400 mb-4" />
                  <h1 className="text-2xl font-bold">Giao dịch thành công</h1>
                  <p className="text-gray-400 mt-1">Giao dịch của bạn đã được xác nhận trên blockchain.</p>
                </div>
                <button onClick={handleReset} className="mt-8 w-full flex justify-center items-center py-3 px-4 rounded-md text-white bg-cyan-600 hover:bg-cyan-700">
                  <Repeat className="h-5 w-5 mr-2" /> Thực hiện giao dịch khác
                </button>
              </div>
            ) : (
              <div>
                <h1 className="text-3xl font-bold text-yellow-400 mb-2 flex items-center">
                  <Gem className="h-8 w-8 mr-3"/> Đầu tư Vàng
                </h1>
                <p className="text-gray-400 mb-6">Mua và bán vàng kỹ thuật số (Chỉ) một cách an toàn và minh bạch.</p>

                {isDataLoading ? (
                  <div className="flex items-center justify-center h-40"><Loader2 className="h-8 w-8 animate-spin" /></div>
                ) : (
                  <div className="space-y-6">
                    {/* Price and Balance Info */}
                    <div className="grid grid-cols-2 gap-4 text-center">
                      <div className="p-4 bg-gray-900/50 rounded-lg">
                          <p className="text-sm text-gray-400">Giá Mua</p>
                          <p className="text-lg font-bold text-green-400">{formatAndShorten(buyPrice)} VNĐ / Chỉ</p>
                      </div>
                      <div className="p-4 bg-gray-900/50 rounded-lg">
                          <p className="text-sm text-gray-400">Giá Bán</p>
                          <p className="text-lg font-bold text-red-400">{formatAndShorten(sellPrice)} VNĐ / Chỉ</p>
                      </div>
                      <div className="p-4 bg-gray-700/50 rounded-lg">
                          <p className="text-sm text-gray-400 flex items-center justify-center"><Wallet className="h-4 w-4 mr-1.5"/>Số dư VNĐ</p>
                          <p className="text-lg font-bold">{formatAndShorten(ethBalance)}</p>
                      </div>
                      <div className="p-4 bg-gray-700/50 rounded-lg">
                          <p className="text-sm text-gray-400 flex items-center justify-center"><Gem className="h-4 w-4 mr-1.5"/>Số dư Vàng (Chỉ)</p>
                          <p className="text-lg font-bold">{formatAndShorten(goldBalance)}</p>
                      </div>
                    </div>

                    {/* Tab Switcher */}
                    <div className="flex bg-gray-900/50 rounded-lg p-1">
                      <button onClick={() => setActiveTab('buy')} className={`w-1/2 p-2 rounded-md font-semibold transition-colors ${activeTab === 'buy' ? 'bg-cyan-600 text-white' : 'text-gray-300 hover:bg-gray-700'}`}>Mua Vàng</button>
                      <button onClick={() => setActiveTab('sell')} className={`w-1/2 p-2 rounded-md font-semibold transition-colors ${activeTab === 'sell' ? 'bg-cyan-600 text-white' : 'text-gray-300 hover:bg-gray-700'}`}>Bán Vàng</button>
                    </div>

                    {/* Buy Form */}
                    {activeTab === 'buy' && (
                        <form onSubmit={handleBuyGold} className="space-y-4 p-4 bg-gray-900/30 rounded-lg">
                            <h3 className="font-semibold text-lg">Bạn muốn mua bao nhiêu?</h3>
                            <div>
                                <label htmlFor="buyAmount" className="block text-sm font-medium text-gray-300 mb-1">Số tiền (VNĐ)</label>
                                <input 
                                    id="buyAmount" 
                                    type="text"
                                    inputMode="decimal" 
                                    value={formattedBuyAmountEth} 
                                    onChange={handleBuyAmountChange} 
                                    placeholder="0" 
                                    className="w-full bg-gray-700 border-gray-600 rounded-md py-2 px-3 focus:ring-cyan-500 focus:border-cyan-500"
                                />
                                <p className="text-xs text-gray-400 mt-1">Số dư khả dụng: {formatAndShorten(ethBalance)} VNĐ</p>
                            </div>
                            <div className="text-center p-3 bg-gray-800 rounded-md">
                                <p className="text-sm text-gray-300">Bạn sẽ nhận được (ước tính)</p>
                                <p className="text-xl font-bold text-yellow-400">{goldToReceive} Chỉ</p>
                            </div>
                            <button type="submit" disabled={isLoading} className="w-full flex justify-center items-center py-3 rounded-md bg-green-600 hover:bg-green-700 disabled:bg-gray-500">
                                {isLoading && activeTx === 'buy' ? <Loader2 className="h-5 w-5 animate-spin"/> : 'Mua Vàng'}
                            </button>
                        </form>
                    )}

                    {/* Sell Form */}
                    {activeTab === 'sell' && (
                        <div className="space-y-4 p-4 bg-gray-900/30 rounded-lg">
                           <h3 className="font-semibold text-lg">Bạn muốn bán bao nhiêu?</h3>
                            <div>
                                <label htmlFor="sellAmount" className="block text-sm font-medium text-gray-300 mb-1">Số lượng (Chỉ)</label>
                                {/* UPDATED: Input for selling gold now handles floats */}
                                <input 
                                    id="sellAmount" 
                                    type="text" 
                                    inputMode="decimal"
                                    value={formattedSellAmountGold} 
                                    onChange={handleSellAmountChange} 
                                    placeholder="0.0" 
                                    className="w-full bg-gray-700 border-gray-600 rounded-md py-2 px-3 focus:ring-cyan-500 focus:border-cyan-500"
                                />
                                <p className="text-xs text-gray-400 mt-1">Số dư khả dụng: {formatAndShorten(goldBalance)} Chỉ</p>
                            </div>
                           <div className="text-center p-3 bg-gray-800 rounded-md">
                                <p className="text-sm text-gray-300">Bạn sẽ nhận được (ước tính)</p>
                                <p className="text-xl font-bold text-cyan-400">{ethToReceive} VNĐ</p>
                           </div>
                           <div className="flex items-center p-3 text-sm text-yellow-300 bg-yellow-900/30 rounded-lg">
                                <ArrowRightLeft className="h-8 w-8 mr-3 flex-shrink-0" />
                                <span>Việc bán vàng yêu cầu 2 bước: <b>Uỷ quyền</b> cho hợp đồng, sau đó <b>Bán</b>.</span>
                           </div>
                           <div className="flex gap-4">
                                <button onClick={handleApprove} disabled={isLoading} className="w-1/2 flex justify-center items-center py-3 rounded-md bg-blue-600 hover:bg-blue-700 disabled:bg-gray-500">
                                  {isLoading && activeTx === 'approve' ? <Loader2 className="h-5 w-5 animate-spin"/> : '1. Uỷ quyền'}
                                </button>
                                <button onClick={handleSellGold} disabled={isLoading || !approvalSuccess} className="w-1/2 flex justify-center items-center py-3 rounded-md bg-red-600 hover:bg-red-700 disabled:bg-gray-500 disabled:cursor-not-allowed">
                                  {isLoading && activeTx === 'sell' ? <Loader2 className="h-5 w-5 animate-spin"/> : '2. Bán Vàng'}
                                </button>
                           </div>
                           {approvalSuccess && <div className="text-center text-green-400 text-sm">✓ Uỷ quyền thành công! Bây giờ bạn có thể bán.</div>}
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
            )}
          </div>
        </div>
      </div>
    </main>
  );
}
