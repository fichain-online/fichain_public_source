'use client'
import config from '@/lib/config';
import Head from "next/head";
import Link from 'next/link';
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
  CardFooter
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Progress } from "@/components/ui/progress";
import {
  UserPlus, LogIn, UploadCloud, ScanFace, KeyRound, ShieldCheck, Copy, ClipboardCheck, Eye, EyeOff, ArrowRight
} from "lucide-react";
import { ethers } from "ethers";

// Import for tsParticles
import { useEffect, useMemo, useState, useCallback } from "react";
import Particles, { initParticlesEngine } from "@tsparticles/react";
import { loadSlim } from "@tsparticles/slim";

// Import Header
import Header from "@/components/header";

// Import react-hot-toast
import toast, { Toaster } from 'react-hot-toast';

export default function RegisterPage() {
  const [currentStep, setCurrentStep] = useState(1);

  // EKYC State (Step 1)
  const [frontIdFile, setFrontIdFile] = useState<File | null>(null);
  const [backIdFile, setBackIdFile] = useState<File | null>(null);
  const [livenessVideoFile, setLivenessVideoFile] = useState<File | null>(null);

  // Wallet State (Step 2 & 3)
  const [seedPhrase, setSeedPhrase] = useState<string>('');
  const [generatedPrivateKey, setGeneratedPrivateKey] = useState<string>('');
  const [walletAddress, setWalletAddress] = useState<string>('');
  const [isSeedPhraseGenerated, setIsSeedPhraseGenerated] = useState(false);
  const [seedPhraseCopied, setSeedPhraseCopied] = useState(false);
  const [showSeedPhrase, setShowSeedPhrase] = useState(false);

  // Seed Phrase Verification State (Step 3)
  const [seedPhraseVerificationInput, setSeedPhraseVerificationInput] = useState('');
  const [isSeedPhraseVerified, setIsSeedPhraseVerified] = useState(false);

  // Email State (Step 4)
  const [email, setEmail] = useState('');

  // General State
  const [isLoading, setIsLoading] = useState(false);

  const handleFileChange = (setter: React.Dispatch<React.SetStateAction<File | null>>) => (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setter(e.target.files[0]);
    }
  };

  const handleGenerateWallet = () => {
    try {
      const wallet = ethers.Wallet.createRandom();
      setSeedPhrase(wallet.mnemonic?.phrase || 'Lỗi tạo seed phrase');
      setGeneratedPrivateKey(wallet.privateKey);
      setWalletAddress(wallet.address)
      setIsSeedPhraseGenerated(true);
      setSeedPhraseCopied(false);
      setIsSeedPhraseVerified(false);
      setSeedPhraseVerificationInput('');
      setShowSeedPhrase(false);
      toast.success("Ví và Seed Phrase đã được tạo!");
    } catch (error) {
      console.error("Lỗi tạo ví:", error);
      toast.error("Đã xảy ra lỗi khi tạo ví. Vui lòng thử lại.");
    }
  };

  const handleCopySeedPhrase = async () => {
    if (!seedPhrase) return;
    try {
      await navigator.clipboard.writeText(seedPhrase);
      setSeedPhraseCopied(true); // Still useful for UI feedback on button
      toast.success('Đã sao chép Seed Phrase vào clipboard!');
      setTimeout(() => setSeedPhraseCopied(false), 2000);
    } catch (err) {
      console.error('Không thể sao chép seed phrase: ', err);
      toast.error('Không thể sao chép. Vui lòng sao chép thủ công.');
    }
  };

  const handleVerifySeedPhrase = () => {
    if (seedPhraseVerificationInput.trim() === seedPhrase.trim()) {
      setIsSeedPhraseVerified(true);
      toast.success("Xác nhận Seed Phrase thành công! Bạn có thể tiếp tục.");
    } else {
      setIsSeedPhraseVerified(false);
      toast.error("Seed Phrase không khớp. Vui lòng kiểm tra lại.");
    }
  };

  const handleNextStep = () => {
    let canProceed = false;
    let errorMessage = "";

    switch (currentStep) {
      case 1:
        if (frontIdFile && backIdFile && livenessVideoFile) {
          canProceed = true;
        } else {
          errorMessage = "Vui lòng hoàn thành tất cả các mục EKYC để tiếp tục.";
        }
        break;
      case 2:
        if (isSeedPhraseGenerated) {
          canProceed = true;
        } else {
          errorMessage = "Vui lòng tạo ví và Seed Phrase để tiếp tục.";
        }
        break;
      case 3:
        if (isSeedPhraseVerified) {
          canProceed = true;
        } else {
          errorMessage = "Vui lòng xác nhận Seed Phrase để tiếp tục.";
        }
        break;
      case 4:
        if (email && email.includes('@') && email.includes('.')) {
          canProceed = true;
        } else {
          errorMessage = "Vui lòng nhập địa chỉ email hợp lệ để tiếp tục.";
        }
        break;
      default:
        break;
    }

    if (canProceed) {
      setCurrentStep(currentStep + 1);
    } else if (errorMessage) {
      toast.error(errorMessage);
    }
  };

  const handleRegister = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    if (!frontIdFile || !backIdFile || !livenessVideoFile || !isSeedPhraseGenerated || !isSeedPhraseVerified || !email || !email.includes('@')) {
      toast.error("Đã xảy ra lỗi. Vui lòng đảm bảo tất cả các bước đã được hoàn thành đúng cách.");
      setIsLoading(false);
      return;
    }

    console.log("Thông tin đăng ký (DEMO):");
    console.log("Mặt trước CMND/CCCD:", frontIdFile.name);
    console.log("Mặt sau CMND/CCCD:", backIdFile.name);
    console.log("Video Liveness:", livenessVideoFile.name);
    console.log("Seed Phrase (đã xác nhận):", seedPhrase);
    console.log("Private Key (đã tạo):", generatedPrivateKey);
    console.log("Email:", email);
    console.log("Address:", walletAddress);

    try {
      const formData = new FormData();
      formData.append("id_front", frontIdFile);
      formData.append("id_back", backIdFile);
      formData.append("selfie", livenessVideoFile);
      formData.append("email", email);
      formData.append("address", walletAddress);

      const response = await fetch(`${config.ekycApiBaseUrl}/user`, {
        method: "POST",
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Đăng ký thất bại");
      }

      toast.success("Đăng ký thành công!");

      toast.custom((t) => (
        <div
          className={`${t.visible ? 'animate-enter' : 'animate-leave'
            } max-w-md w-full bg-slate-800 shadow-lg rounded-lg pointer-events-auto flex ring-1 ring-black ring-opacity-5 border border-green-500`}
        >
          <div className="flex-1 w-0 p-4">
            <div className="flex items-start">
              <div className="flex-shrink-0 pt-0.5">
                <ShieldCheck className="h-10 w-10 text-green-500" />
              </div>
              <div className="ml-3 flex-1">
                <p className="text-sm font-medium text-green-400">
                  Đăng ký thành công (DEMO)!
                </p>
                <p className="mt-1 text-sm text-slate-300">
                  Email: {email} (Trạng thái duyệt EKYC sẽ được gửi qua email này).
                </p>
                <p className="mt-2 text-sm font-semibold text-red-400">
                  LƯU TRỮ CẨN THẬN KHÓA RIÊNG TƯ:
                </p>
                <p className="mt-1 text-xs text-red-300 bg-slate-900 p-2 rounded select-all font-mono break-all">
                  {generatedPrivateKey}
                </p>
                <p className="mt-1 text-xs text-amber-400">
                  Bạn sẽ cần Khóa Riêng Tư này để đăng nhập trong phiên bản DEMO.
                  Trong ứng dụng thực tế, Khóa Riêng Tư sẽ không được hiển thị như thế này.
                </p>
              </div>
            </div>
          </div>
          <div className="flex border-l border-slate-700">
            <button
              onClick={() => toast.dismiss(t.id)}
              className="w-full border border-transparent rounded-none rounded-r-lg p-4 flex items-center justify-center text-sm font-medium text-primary hover:text-primary/80 focus:outline-none focus:ring-2 focus:ring-primary"
            >
              Đóng
            </button>
          </div>
        </div>
      ), { duration: Infinity });

      setTimeout(() => {
        setIsLoading(false);
        window.location.href = '/login';
      }, 3000);

    } catch (error) {
      toast.error(String(error));
      setIsLoading(false);
    }
  };

  const progressValue = (currentStep / 5) * 100;

  return (
    <>
      <Head>
        <title>Đăng Ký - Fichain</title>
        <meta name="description" content="Tạo tài khoản Fichain mới của bạn." />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      {/* Toaster component needs to be rendered */}
      <Toaster
        position="top-center"
        reverseOrder={false}
        toastOptions={{
          // Define default options
          className: '',
          duration: 5000,
          style: {
            background: '#1e293b', // slate-800
            color: '#e2e8f0', // slate-200
            border: '1px solid #334155', // slate-700
          },
          // Default options for specific types
          success: {
            duration: 3000,
            iconTheme: {
              primary: '#10b981', // emerald-500
              secondary: '#1e293b',
            },
          },
          error: {
            iconTheme: {
              primary: '#f43f5e', // rose-500
              secondary: '#1e293b',
            },
          },
        }}
      />
      <Header />


      <main className="flex items-center justify-center min-h-screen bg-transparent text-slate-200 relative z-0 py-20 px-4">
        <Card className="w-full max-w-2xl bg-slate-800/80 backdrop-blur-md shadow-xl border-slate-700">
          <CardHeader className="text-center">
            <UserPlus className="mx-auto h-12 w-12 text-primary mb-4" />
            <CardTitle className="text-3xl font-bold text-slate-100">Tạo Tài Khoản Mới</CardTitle>
            <CardDescription className="text-slate-300">
              Bước {currentStep} / 5
            </CardDescription>
            <Progress value={progressValue} className="w-full mt-4" />
          </CardHeader>
          <CardContent>
            <form onSubmit={(e) => {
              e.preventDefault();
              if (currentStep === 5) {
                handleRegister(e);
              } else {
                handleNextStep();
              }
            }} className="space-y-6">

              {/* Step 1: EKYC */}
              {currentStep === 1 && (
                <div className="space-y-4">
                  <h3 className="text-lg font-semibold text-slate-100">1. Thông Tin E-KYC</h3>
                  <div>
                    <Label htmlFor="frontId" className="block text-sm font-medium text-slate-300 mb-1">
                      Ảnh Mặt Trước CMND/CCCD <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id="frontId"
                      type="file"
                      accept="image/*"
                      onChange={handleFileChange(setFrontIdFile)}
                      required
                      className="bg-slate-700 border-slate-600 text-slate-100 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-primary/80 file:text-primary-foreground hover:file:bg-primary"
                    />
                    {frontIdFile && <p className="text-xs text-slate-400 mt-1">Đã chọn: {frontIdFile.name}</p>}
                  </div>
                  <div>
                    <Label htmlFor="backId" className="block text-sm font-medium text-slate-300 mb-1">
                      Ảnh Mặt Sau CMND/CCCD <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id="backId"
                      type="file"
                      accept="image/*"
                      onChange={handleFileChange(setBackIdFile)}
                      required
                      className="bg-slate-700 border-slate-600 text-slate-100 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-primary/80 file:text-primary-foreground hover:file:bg-primary"
                    />
                    {backIdFile && <p className="text-xs text-slate-400 mt-1">Đã chọn: {backIdFile.name}</p>}
                  </div>
                  <div>
                    <Label htmlFor="liveness" className="block text-sm font-medium text-slate-300 mb-1">
                      Video Selfie Xác Thực (Liveness Check) <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id="liveness"
                      type="file"
                      accept="video/*"
                      onChange={handleFileChange(setLivenessVideoFile)}
                      required
                      className="bg-slate-700 border-slate-600 text-slate-100 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-primary/80 file:text-primary-foreground hover:file:bg-primary"
                    />
                    {livenessVideoFile && <p className="text-xs text-slate-400 mt-1">Đã chọn: {livenessVideoFile.name}</p>}
                    <p className="mt-1 text-xs text-slate-400">
                      (Demo: Tải lên video selfie ngắn. Trong thực tế, bạn sẽ được yêu cầu thực hiện các hành động như quay đầu, mỉm cười.)
                    </p>
                  </div>
                  <Button
                    variant="outline"
                    type="button"
                    onClick={handleNextStep}
                    className="w-full text-white border-gray-300 hover:bg-white/20 disabled:text-slate-500 disabled:border-slate-600 disabled:hover:bg-transparent"
                    disabled={!frontIdFile || !backIdFile || !livenessVideoFile} >
                    Tiếp Tục <ArrowRight className="ml-2 h-5 w-5" />
                  </Button>
                </div>
              )}

              {/* Step 2: Generate & Copy Seed Phrase */}
              {currentStep === 2 && (
                <div className="space-y-4">
                  <h3 className="text-lg font-semibold text-slate-100">2. Tạo Ví & Lưu Seed Phrase</h3>
                  {!isSeedPhraseGenerated ? (
                    <Button type="button" onClick={handleGenerateWallet} className="w-full bg-green-600 hover:bg-green-700 text-white">
                      <KeyRound className="mr-2 h-5 w-5" />
                      Tạo Ví Mới & Seed Phrase
                    </Button>
                  ) : (
                    <div className="space-y-4">
                      <div>
                        <Label className="block text-sm font-medium text-slate-300 mb-1">
                          Seed Phrase (Cụm Từ Khôi Phục) - LƯU TRỮ CẨN THẬN!
                        </Label>
                        <div className="relative">
                          <Textarea
                            readOnly
                            value={showSeedPhrase ? seedPhrase : '****************************'}
                            className="bg-slate-900 border-slate-700 text-amber-300 font-mono text-sm h-24 select-all resize-none pr-12"
                            rows={3}
                          />
                          <Button
                            type="button"
                            variant="outline"
                            className="absolute top-2 right-2 text-slate-400 hover:text-primary border-slate-600 hover:bg-slate-700"
                            onClick={() => setShowSeedPhrase(!showSeedPhrase)}
                            aria-label={showSeedPhrase ? 'Ẩn Seed Phrase' : 'Hiện Seed Phrase'}
                          >
                            {showSeedPhrase ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                          </Button>
                        </div>
                        <p className="mt-1 text-xs text-amber-400">
                          <strong>QUAN TRỌNG:</strong> Viết xuống và cất giữ Seed Phrase này ở một nơi AN TOÀN. Đây là cách duy nhất để khôi phục ví của bạn nếu mất quyền truy cập. <strong>KHÔNG chia sẻ với bất kỳ ai.</strong>
                        </p>
                        <Button type="button" onClick={handleCopySeedPhrase} variant="outline" className="mt-2 text-slate-300 border-slate-500 hover:bg-slate-700">
                          {seedPhraseCopied ? <ClipboardCheck className="mr-2 h-4 w-4 text-green-500" /> : <Copy className="mr-2 h-4 w-4" />}
                          {seedPhraseCopied ? 'Đã Sao Chép!' : 'Sao Chép Seed Phrase'}
                        </Button>
                      </div>
                      <Button
                        variant="outline"
                        type="button"
                        onClick={handleNextStep}
                        className="w-full text-white border-gray-300 hover:bg-white/20 disabled:text-slate-500 disabled:border-slate-600 disabled:hover:bg-transparent"
                        disabled={!isSeedPhraseGenerated} >
                        Tôi Đã Lưu, Tiếp Tục <ArrowRight className="ml-2 h-5 w-5" />
                      </Button>
                    </div>
                  )}
                </div>
              )}

              {/* Step 3: Verify Seed Phrase */}
              {currentStep === 3 && (
                <div className="space-y-4">
                  <h3 className="text-lg font-semibold text-slate-100">3. Xác Nhận Seed Phrase</h3>
                  <div>
                    <Label htmlFor="verifySeedPhrase" className="block text-sm font-medium text-slate-300 mb-1">
                      Nhập lại chính xác Seed Phrase bạn vừa lưu <span className="text-red-500">*</span>
                    </Label>
                    <Textarea
                      id="verifySeedPhrase"
                      value={seedPhraseVerificationInput}
                      onChange={(e) => setSeedPhraseVerificationInput(e.target.value)}
                      placeholder="Nhập lại chính xác Seed Phrase bạn vừa lưu..."
                      required
                      className="bg-slate-700 border-slate-600 text-slate-100 placeholder-slate-400 focus:ring-primary focus:border-primary h-24 resize-none"
                      rows={3}
                    />
                    {!isSeedPhraseVerified && (
                      <Button
                        type="button"
                        onClick={handleVerifySeedPhrase}
                        className="mt-2 bg-orange-500 hover:bg-orange-600 text-white disabled:bg-orange-500/70 disabled:hover:bg-orange-500/70"
                        disabled={seedPhraseVerificationInput.trim() === ''}>
                        <ShieldCheck className="mr-2 h-4 w-4" />
                        Xác Nhận
                      </Button>
                    )}
                  </div>

                  {isSeedPhraseVerified && (
                    <div className="space-y-4">
                      <p className="text-green-400 text-sm flex items-center">
                        <ShieldCheck className="mr-2 h-5 w-5" /> Seed Phrase đã được xác nhận thành công!
                      </p>
                      <Button
                        variant="outline"
                        type="button"
                        onClick={handleNextStep}
                        className="w-full text-white border-gray-300 hover:bg-white/20">
                        Tiếp Tục <ArrowRight className="ml-2 h-5 w-5" />
                      </Button>
                      <div className="mt-4 p-3 bg-red-900/50 border border-red-700 rounded-md">
                        <Label className="block text-sm font-bold text-red-400 mb-1">
                          KHÓA RIÊNG TƯ (Private Key) - CHỈ DÀNH CHO DEMO
                        </Label>
                        <Input
                          readOnly
                          value={generatedPrivateKey}
                          className="bg-slate-900 border-slate-700 text-red-300 font-mono text-xs select-all"
                        />
                        <p className="mt-1 text-xs text-red-300">
                          <strong>CẢNH BÁO TUYỆT ĐỐI:</strong> Trong phiên bản thử nghiệm này, bạn cần Khóa Riêng Tư để đăng nhập. Hãy sao chép và lưu trữ nó cực kỳ cẩn thận.
                          <strong>Trong ứng dụng thực tế, Khóa Riêng Tư KHÔNG BAO GIỜ được hiển thị hoặc xử lý theo cách này.</strong> Việc mất Khóa Riêng Tư đồng nghĩa với mất toàn bộ tài sản.
                        </p>
                      </div>
                    </div>
                  )}
                </div>
              )}

              {/* Step 4: Enter Email */}
              {currentStep === 4 && (
                <div className="space-y-4">
                  <h3 className="text-lg font-semibold text-slate-100">4. Nhập Địa Chỉ Email</h3>
                  <div>
                    <Label htmlFor="email" className="block text-sm font-medium text-slate-300 mb-1">
                      Email để nhận trạng thái duyệt EKYC <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id="email"
                      type="email"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      placeholder="Ví dụ: emailcuaban@example.com"
                      required
                      className="bg-slate-700 border-slate-600 text-slate-100 placeholder-slate-400 focus:ring-primary focus:border-primary"
                    />
                    <p className="mt-1 text-xs text-slate-400">
                      Chúng tôi sẽ gửi thông báo về trạng thái xác minh tài khoản của bạn đến địa chỉ email này.
                    </p>
                  </div>
                  <Button
                    variant="outline"
                    type="button"
                    onClick={handleNextStep}
                    className="w-full text-white border-gray-300 hover:bg-white/20 disabled:text-slate-500 disabled:border-slate-600 disabled:hover:bg-transparent"
                    disabled={!email || !email.includes('@')}>
                    Hoàn Thành <ArrowRight className="ml-2 h-5 w-5" />
                  </Button>
                </div>
              )}

              {/* Step 5: Completion & Submit */}
              {currentStep === 5 && (
                <div className="space-y-6 text-center">
                  <h3 className="text-lg font-semibold text-slate-100">5. Hoàn Tất Đăng Ký</h3>
                  <p className="text-slate-300">
                    Xin chúc mừng! Bạn đã hoàn thành các bước cần thiết.
                    Nhấn nút "Đăng Ký Tài Khoản" dưới đây để gửi thông tin và hoàn tất quá trình.
                    Trạng thái duyệt EKYC sẽ được thông báo qua email bạn đã cung cấp.
                  </p>
                  <div className="p-3 bg-red-900/50 border border-red-700 rounded-md text-left">
                    <Label className="block text-sm font-bold text-red-400 mb-1">
                      KHÓA RIÊNG TƯ CỦA BẠN (DÙNG ĐỂ ĐĂNG NHẬP DEMO)
                    </Label>
                    <Input
                      readOnly
                      value={generatedPrivateKey}
                      className="bg-slate-900 border-slate-700 text-red-300 font-mono text-xs select-all"
                    />
                    <p className="mt-1 text-xs text-red-300">
                      <strong>TUYỆT ĐỐI QUAN TRỌNG:</strong> Vui lòng sao chép và lưu trữ Khóa Riêng Tư này ngay bây giờ. Đây là cách duy nhất để đăng nhập vào tài khoản demo của bạn.
                      <strong>Mất Khóa Riêng Tư = Mất Quyền Truy Cập (trong demo) hoặc Mất Tài Sản (trong ứng dụng thật).</strong>
                    </p>
                  </div>

                  <Button
                    type="submit" // This is the main submit for the form
                    variant="outline"
                    className="w-full text-white border-gray-300 hover:bg-white/20 disabled:text-slate-500 disabled:border-slate-600 disabled:hover:bg-transparent"
                    disabled={isLoading}
                  >
                    {isLoading ? "Đang xử lý..." : "Đăng Ký Tài Khoản"}
                  </Button>
                </div>
              )}
            </form>
          </CardContent>
          <CardFooter className="mt-6 text-center flex flex-col items-center">
            <p className="text-sm text-slate-400">
              Đã có tài khoản?{' '}
              <Link href="/login" className="font-medium text-primary hover:text-primary/80 hover:underline">
                <LogIn className="inline-block mr-1 h-4 w-4" />
                Đăng nhập tại đây
              </Link>
            </p>
          </CardFooter>
        </Card>
      </main>
    </>
  );
}
