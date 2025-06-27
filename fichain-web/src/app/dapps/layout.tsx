'use client';

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Wallet } from 'ethers';

// --- Zustand Store Imports ---
import { useAuthStore } from '@/stores/authStore';
import { useWebSocketStore } from '@/stores/webSocketStore';

// --- UI and Icon Imports ---
import Header from '@/components/header';
import { Loader2 } from 'lucide-react';

export default function DappsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();

  // 1. Get state from stores
  const { isAuthenticated, privateKey, _hasHydrated: authHasHydrated } = useAuthStore();
  const { connected, connect, setWallet, logs } = useWebSocketStore();
  
  // Local state to track the connection attempt
  const [isConnecting, setIsConnecting] = useState(false);

  // Effect for Authentication, Redirection, and WebSocket Connection
  useEffect(() => {
    // Wait until the auth store has been rehydrated from storage
    if (!authHasHydrated) {
      return; // Do nothing until we know the auth status
    }

    // If not authenticated, redirect to the login page
    if (!isAuthenticated) {
      router.replace('/');
      return;
    }

    // If authenticated but not connected, and we are not already trying to connect
    if (isAuthenticated && privateKey && !connected && !isConnecting) {
      console.log("DApps Layout: User is authenticated. Initiating WebSocket connection...");
      setIsConnecting(true); // Mark that we are attempting to connect
      try {
        const userWallet = new Wallet(privateKey);
        setWallet(userWallet);
        connect();
      } catch (error) {
        console.error("DApps Layout: Failed to create wallet or connect:", error);
        // Optional: handle this error, maybe redirect with an error message
        setIsConnecting(false); // Reset on failure
      }
    }
    
    // If the store reports it's connected, we can turn off our local 'isConnecting' flag.
    if (connected && isConnecting) {
        setIsConnecting(false);
    }

  }, [
    authHasHydrated,
    isAuthenticated,
    privateKey,
    connected,
    connect,
    setWallet,
    router,
    isConnecting
  ]);

  // 2. Render a loading state until everything is ready
  // We show a loader if:
  // - The auth store is still hydrating.
  // - The user is authenticated but the WebSocket is not yet connected.
  if (!authHasHydrated || !isAuthenticated || !connected) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen text-white ">
        <Loader2 className="h-12 w-12 animate-spin text-cyan-400" />
        <p className="mt-4 text-lg">
          { !authHasHydrated ? 'Đang xác thực...' : 'Đang kết nối đến máy chủ...' }
        </p>
        {/* Optional: Show recent logs during connection for debugging */}
        <div className="mt-4 p-2 bg-gray-800 rounded text-xs text-gray-400 w-full max-w-md h-24 overflow-y-auto">
          {logs.slice(-5).map((log, i) => <div key={i}>{log}</div>)}
        </div>
      </div>
    );
  }

  // 3. If authenticated and connected, render the main layout and the child page
  return (
    <main className="min-h-screen text-white ">
      <Header />
      {/* The {children} prop will be the actual page being rendered */}
      {children}
    </main>
  );
}
