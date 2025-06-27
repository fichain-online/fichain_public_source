'use client'

import Head from "next/head";
import Link from 'next/link';
import { useRouter } from 'next/navigation'; // Use from 'next/navigation' in App Router
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { LogIn, UserPlus, LockKeyhole, Loader2 } from "lucide-react";
import { useEffect, useState } from "react";
import Header from "@/components/header";

// Import our new Zustand store
import { useAuthStore } from "@/stores/authStore";

export default function LoginPage() {
  // Local state for the input field
  const [privateKeyInput, setPrivateKeyInput] = useState('');
  
  // Get state and actions from the Zustand store
  const { login, isLoading, error, isAuthenticated } = useAuthStore();
  
  const router = useRouter();

  // Redirect if user is already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      router.push('/dapps'); 
    }
  }, [isAuthenticated, router]);

  const handleLogin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    await login(privateKeyInput);
    // The useEffect above will handle the redirect on successful login
  };

  return (
    <>
      <Head>
        <title>Đăng Nhập - Fichain</title>
        <meta name="description" content="Đăng nhập vào tài khoản Fichain của bạn." />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <main className="flex items-center justify-center min-h-screen bg-transparent text-slate-200 relative z-0 py-20 px-4">
        <Card className="w-full max-w-md bg-slate-800/80 backdrop-blur-md shadow-xl border-slate-700">
          <CardHeader className="text-center">
            <LockKeyhole className="mx-auto h-12 w-12 text-primary mb-4" />
            <CardTitle className="text-3xl font-bold text-slate-100">Đăng Nhập</CardTitle>
            <CardDescription className="text-slate-300">
              Truy cập vào tài khoản Fichain của bạn.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleLogin} className="space-y-6">
              <div>
                <label htmlFor="privateKey" className="block text-sm font-medium text-slate-300 mb-1">
                  Khóa Riêng Tư (Private Key)
                </label>
                <Input
                  id="privateKey"
                  type="password"
                  value={privateKeyInput}
                  onChange={(e) => setPrivateKeyInput(e.target.value)}
                  placeholder="Nhập khóa riêng tư của bạn"
                  required
                  className="bg-slate-700 border-slate-600 text-slate-100 placeholder-slate-400 focus:ring-primary focus:border-primary"
                  disabled={isLoading}
                />
                <p className="mt-2 text-xs text-amber-400">
                  <strong>LƯU Ý (PHIÊN BẢN THỬ NGHIỆM):</strong> Vui lòng lưu trữ Khóa Riêng Tư của bạn một cách cẩn thận và an toàn.
                </p>
                {/* Display login error from the store */}
                {error && (
                  <p className="mt-2 text-sm text-red-400">{error}</p>
                )}
              </div>
              <Button 
                type="submit" 
                variant="outline"
                className="w-full text-white border-gray-300 hover:bg-white/20"
                disabled={isLoading}
              >
                {isLoading ? (
                  <>
                    <Loader2 className="mr-2 h-5 w-5 animate-spin" />
                    Đang xác thực...
                  </>
                ) : (
                  <>
                    <LogIn className="mr-2 h-5 w-5" />
                    Đăng Nhập
                  </>
                )}
              </Button>
            </form>
            <div className="mt-6 text-center">
              <p className="text-sm text-slate-400">
                Chưa có tài khoản?{' '}
                <Link href="/register" className="font-medium text-primary hover:text-primary/80 hover:underline">
                  <UserPlus className="inline-block mr-1 h-4 w-4" />
                  Đăng ký ngay
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>
      </main>
    </>
  );
}
