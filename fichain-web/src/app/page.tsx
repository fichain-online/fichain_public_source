'use client'
import Image from "next/image";
import Head from "next/head";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import {
  ShieldCheck, Users, DollarSign, TrendingUp, CheckCircle, HelpCircle, UsersRound,
  Building, HeartHandshake, Briefcase, Linkedin, Twitter, MessageCircle, ArrowRight,
  LockKeyhole, FileText, Info, BarChart3, WalletCards, Send, Sparkles, Menu, X, LayoutDashboard,
  // NEW ICONS for B2B & Tech focus
  Landmark, BriefcaseBusiness, Replace, CodeXml, BookOpen, Layers, Zap, Cpu
} from "lucide-react";

import Link from 'next/link';
// Import for Framer Motion
import { motion } from "framer-motion";

// Import Header
import Header from "@/components/header"; // Adjust path if needed

// Framer Motion variants for section animations
const sectionVariants = {
  hidden: { opacity: 0, y: 50 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.6, ease: "easeOut" }
  },
};

export default function Home() {
  // NOTE: Smooth scrolling is handled natively by CSS `scroll-behavior: smooth`
  // which should be set in your globals.css on the `html` element for best practice.

  return (
    <>
      <Head>
        <title>Fichain: Blockchain Layer 1 Chuyên Dụng cho Ngân Hàng | Công Nghệ Lõi Việt Nam</title>
        <meta name="description" content="Fichain - Nền tảng blockchain Layer 1 hiệu suất cao do Việt Nam phát triển và làm chủ, cung cấp hạ tầng an toàn và tuân thủ cho các ứng dụng tài chính - ngân hàng thế hệ mới." />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <Header />

      <div className="space-y-12 md:space-y-16 bg-transparent text-slate-200 relative z-0">

        {/* UPDATED: Hero Section with new messaging */}
        <motion.section
          id="hero"
          className="relative text-center isolate overflow-hidden pt-20"
          initial="hidden"
          animate="visible"
          variants={sectionVariants}
        >
          <div className="absolute inset-0 z-[-1]">
            <Image
              src="/hero.png" // Consider a more abstract/tech background image
              alt="Nền tảng Fichain công nghệ blockchain"
              layout="fill"
              objectFit="cover"
              priority
            />
            <div className="absolute inset-0 bg-black/75"></div>
          </div>
          <div className="container mx-auto px-4 py-24 md:py-32 lg:py-40 xl:py-48 relative z-10 space-y-8">
            <h1 className="text-4xl sm:text-5xl lg:text-6xl xl:text-7xl font-bold tracking-tight text-white">
              Fichain: Nền Tảng Blockchain Layer 1
              <br className="hidden sm:block" />
              <span className="text-cyan-400 hover:text-cyan-300 transition-colors duration-300">
                Kiến Tạo Tương Lai Ngân Hàng Số
              </span>
            </h1>
            <p className="text-lg sm:text-xl text-gray-200 max-w-4xl mx-auto">
              Nền tảng blockchain hiệu suất cao, an toàn và tuân thủ, do người Việt Nam phát triển và làm chủ hoàn toàn công nghệ, được thiết kế chuyên biệt cho các nghiệp vụ lõi của ngành tài chính - ngân hàng.
            </p>
            <div className="pt-4 md:pt-2 space-y-4 sm:space-y-0 sm:flex sm:justify-center sm:space-x-4">
              <Button
                size="lg"
                className="w-full sm:w-auto text-lg px-8 py-3 bg-primary hover:bg-primary/90 text-primary-foreground"
              >
                Liên Hệ Hợp Tác
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
              <Button
                size="lg"
                variant="outline"
                className="w-full sm:w-auto text-lg px-8 py-3 border-gray-400 text-gray-100 hover:bg-white/10 hover:text-white hover:border-white transition-colors duration-300"
              >
                Xem Tài Liệu Kỹ Thuật
                <BookOpen className="ml-2 h-5 w-5" />
              </Button>
            </div>
            <p className="text-sm text-gray-300 pt-6 md:pt-4 flex items-center justify-center">
              <span className="text-3xl mr-2">🇻🇳</span>
              Công nghệ lõi được phát triển và làm chủ bởi đội ngũ kỹ sư Việt Nam.
            </p>
          </div>
        </motion.section>

        {/* NEW: Platform Core Section */}
        <motion.section
          id="platform-core"
          className="container mx-auto px-4"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.2 }}
          variants={sectionVariants}
        >
          <div className="bg-slate-800/80 backdrop-blur-md rounded-lg shadow-xl p-8 md:p-10">
            <div className="text-center space-y-3 mb-10">
              <h2 className="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-100">Nền Tảng Blockchain Layer 1 Chuyên Dụng</h2>
              <p className="text-lg text-slate-300 max-w-3xl mx-auto">
                Fichain được xây dựng từ gốc để đáp ứng các yêu cầu khắt khe nhất của ngành tài chính.
              </p>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8 text-center">
              {[
                { icon: Zap, title: "Hiệu Suất Vượt Trội", description: "Xử lý hàng nghìn giao dịch mỗi giây (TPS) với độ trễ thấp, đảm bảo các hoạt động thanh toán và giao dịch diễn ra tức thì." },
                { icon: CodeXml, title: "Tương Thích EVM", description: "Hỗ trợ đầy đủ Ethereum Virtual Machine, cho phép triển khai nhanh chóng các hợp đồng thông minh và DApps từ hệ sinh thái Ethereum." },
                { icon: ShieldCheck, title: "Bảo Mật & Tuân Thủ", description: "Thiết kế với cơ chế đồng thuận PoA/PoS và khả năng tạo mạng riêng (private network) để đáp ứng các quy định của ngân hàng." },
              ].map((item, index) => (
                <motion.div
                  key={item.title}
                  className="p-6"
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ duration: 0.5, delay: index * 0.15 }}
                >
                  <div className="bg-primary/20 text-primary rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
                    <item.icon className="h-8 w-8" />
                  </div>
                  <h3 className="text-xl font-semibold mb-2 text-slate-100">{item.title}</h3>
                  <p className="text-slate-300">{item.description}</p>
                </motion.div>
              ))}
            </div>
          </div>
        </motion.section>

        {/* UPDATED: Problem/Solution Section reframed for B2B */}
        <motion.section
          id="problem-solution"
          className="py-12 md:py-16 space-y-8 container mx-auto px-4 bg-slate-800/80 backdrop-blur-md rounded-lg shadow-xl"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.2 }}
          variants={sectionVariants}
        >
          <div className="text-center space-y-3 mb-10">
            <h2 className="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-100">Hạ Tầng Cho Kỷ Nguyên Số Của Ngân Hàng</h2>
            <p className="text-lg text-slate-300 max-w-3xl mx-auto">
              Vượt qua rào cản của hệ thống kế thừa (legacy systems) và nắm bắt cơ hội từ công nghệ blockchain một cách an toàn và hiệu quả.
            </p>
          </div>
          <div className="grid md:grid-cols-2 gap-8 items-center">
            <div className="space-y-4">
              <h3 className="text-2xl font-semibold text-primary">Thách Thức Của Ngành Ngân Hàng</h3>
              <p className="text-slate-300">
                Các ngân hàng đối mặt với áp lực hiện đại hóa nhưng bị cản trở bởi hạ tầng công nghệ cũ, chi phí tích hợp cao và những lo ngại về bảo mật, tuân thủ khi áp dụng blockchain. Việc xây dựng một blockchain riêng đòi hỏi nguồn lực khổng lồ và chuyên môn sâu.
              </p>
              <h3 className="text-2xl font-semibold text-primary mt-6">Giải Pháp Hạ Tầng Từ Fichain</h3>
              <p className="text-slate-300">
                Fichain cung cấp một nền tảng Layer 1 sẵn sàng sử dụng (turnkey), giúp ngân hàng và tổ chức tài chính nhanh chóng phát triển và triển khai các sản phẩm mới. Chúng tôi loại bỏ gánh nặng xây dựng hạ tầng, cho phép đối tác tập trung vào việc tạo ra giá trị kinh doanh.
              </p>
            </div>
            <div className="relative w-full h-64 sm:h-80 rounded-lg shadow-lg overflow-hidden">
              <Image src="/Financial-services-2.webp" alt="Hình ảnh Kết nối Tài chính" layout="fill" objectFit="cover" />
            </div>
          </div>
        </motion.section>

        {/* UPDATED: Features Section focusing on platform capabilities */}
        <motion.section
          id="features"
          className="py-12 md:py-16 bg-slate-800/60 backdrop-blur-md rounded-lg shadow-xl"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.1 }}
          variants={sectionVariants}
        >
          <div className="container mx-auto px-4">
            <div className="text-center space-y-3 mb-10">
              <h2 className="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-100">Năng Lực Vượt Trội Của Nền Tảng</h2>
              <p className="text-lg text-slate-300 max-w-2xl mx-auto">
                Các tính năng được thiết kế để xây dựng những ứng dụng tài chính mạnh mẽ nhất.
              </p>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {[
                {
                  icon: Layers,
                  title: "PoSA: Cơ Chế Đồng Thuận cho Liên Minh Ngân Hàng",
                  description: "Cơ chế PoSA yêu cầu các validator phải là những ngân hàng có định danh và uy tín (Authority), đồng thời phải ký quỹ một lượng tài sản lớn để tham gia xác thực (Stake). Sự kết hợp này tạo ra một mạng lưới vừa hiệu suất cao, vừa có độ tin cậy và bảo mật cấp độ doanh nghiệp, lý tưởng cho các hoạt động liên ngân hàng."
                },
                { icon: Cpu, title: "Chi Phí Giao Dịch Thấp & Ổn Định", description: "Cấu trúc phí được tối ưu hóa để đảm bảo các giao dịch vi mô (micro-transactions) và hoạt động nghiệp vụ diễn ra hiệu quả." },
                { icon: LockKeyhole, title: "Bảo Mật Cấp Tổ Chức", description: "Nền tảng được xây dựng với các tiêu chuẩn bảo mật cao nhất, hỗ trợ HSM và các giải pháp mã hóa đầu cuối." },
                { icon: LayoutDashboard, title: "Bộ Công Cụ Phát Triển (SDKs)", description: "Cung cấp đầy đủ API/SDK giúp các nhà phát triển dễ dàng tích hợp và xây dựng ứng dụng trên nền tảng Fichain." },
                { icon: BarChart3, title: "Hệ Thống Giám Sát Toàn Diện", description: "Cung cấp các công cụ giám sát (monitoring) và phân tích on-chain theo thời gian thực, đảm bảo tính minh bạch và ổn định." },
                { icon: UsersRound, title: "Hỗ Trợ Kỹ Thuật Chuyên Sâu", description: "Đội ngũ kỹ sư lõi của chúng tôi sẵn sàng hỗ trợ đối tác trong quá trình tích hợp và vận hành." },
              ].map((feature, index) => (
                <motion.div
                  key={feature.title}
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ duration: 0.5, delay: index * 0.1 }}
                >
                  <Card className="bg-slate-700/70 backdrop-blur-sm hover:bg-slate-600/80 transition-colors duration-300 shadow-lg h-full border-slate-600">
                    <CardHeader>
                      <div className="bg-primary/20 text-primary p-3 rounded-lg w-fit mb-4">
                        <feature.icon className="h-8 w-8" />
                      </div>
                      <CardTitle className="text-slate-100">{feature.title}</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <p className="text-slate-300">{feature.description}</p>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </div>
          </div>
        </motion.section>

        {/* UPDATED: How It Works Section for B2B clients */}
        <motion.section
          id="how-it-works"
          className="py-12 md:py-16 container mx-auto px-4 bg-slate-800/80 backdrop-blur-md rounded-lg shadow-xl"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.2 }}
          variants={sectionVariants}
        >
          <div className="text-center space-y-3 mb-10">
            <h2 className="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-100">Triển Khai Giải Pháp Của Bạn Trên Fichain</h2>
            <p className="text-lg text-slate-300 max-w-xl mx-auto">
              Quy trình hợp tác đơn giản và hiệu quả để đưa ý tưởng của bạn vào thực tế.
            </p>
          </div>
          <div className="grid md:grid-cols-3 gap-8 text-center">
            {[
              { step: 1, title: "Liên Hệ & Tư Vấn", description: "Trao đổi với đội ngũ của chúng tôi để phân tích nhu cầu và thiết kế giải pháp kiến trúc phù hợp nhất cho bài toán của bạn.", icon: MessageCircle },
              { step: 2, title: "Tích Hợp & Phát Triển", description: "Sử dụng bộ SDK và tài liệu kỹ thuật đầy đủ của chúng tôi để phát triển và thử nghiệm các hợp đồng thông minh, ứng dụng.", icon: CodeXml },
              { step: 3, title: "Triển Khai & Vận Hành", description: "Triển khai giải pháp của bạn trên mạng chính (mainnet) của Fichain với sự hỗ trợ vận hành và giám sát liên tục từ chúng tôi.", icon: Send },
            ].map((item, index) => (
              <motion.div
                key={item.step}
                className="p-6 border border-slate-700 rounded-lg shadow-md hover:shadow-lg transition-shadow bg-slate-700/70 backdrop-blur-sm hover:bg-slate-600/80"
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ duration: 0.5, delay: index * 0.15 }}
              >
                <div className="bg-primary text-primary-foreground rounded-full w-12 h-12 flex items-center justify-center mx-auto mb-4">
                  <item.icon className="h-6 w-6" />
                </div>
                <h3 className="text-xl font-semibold mb-2 text-slate-100">Bước {item.step}: {item.title}</h3>
                <p className="text-slate-300">{item.description}</p>
              </motion.div>
            ))}
          </div>
        </motion.section>

        {/* UPDATED: Use Cases Section for Banking */}
        <motion.section
          id="use-cases"
          className="py-12 md:py-16 bg-slate-800/60 backdrop-blur-md rounded-lg shadow-xl"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.1 }}
          variants={sectionVariants}
        >
          <div className="container mx-auto px-4">
            <div className="text-center space-y-3 mb-10">
              <h2 className="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-100">Ứng Dụng Thực Tiễn Cho Ngành Ngân Hàng</h2>
              <p className="text-lg text-slate-300 max-w-xl mx-auto">
                Mở khóa các mô hình kinh doanh mới và tối ưu hóa hoạt động hiện tại.
              </p>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
              {[
                { icon: Replace, title: "Số Hóa Tài Sản & Chứng Khoán", description: "Phát hành và giao dịch các loại tài sản được mã hóa (tokenized assets) như trái phiếu, cổ phần, bất động sản một cách minh bạch và hiệu quả." },
                { icon: Landmark, title: "Thanh Toán Xuyên Biên Giới", description: "Xây dựng hệ thống chuyển tiền và thanh toán quốc tế tức thì, chi phí thấp, giảm sự phụ thuộc vào các mạng lưới trung gian." },
                { icon: BriefcaseBusiness, title: "Tài Chính Doanh Nghiệp & Tín Dụng", description: "Tạo lập các nền tảng cho vay ngang hàng (P2P Lending), tài trợ thương mại, và quản lý chuỗi cung ứng trên blockchain." },
              ].map((useCase, index) => (
                <motion.div
                  key={useCase.title}
                  initial={{ opacity: 0, y: 20 }}
                  whileInView={{ opacity: 1, y: 0 }}
                  viewport={{ once: true }}
                  transition={{ duration: 0.5, delay: index * 0.1 }}
                >
                  <Card className="text-center bg-slate-700/70 backdrop-blur-sm hover:bg-slate-600/80 transition-colors duration-300 shadow-lg h-full border-slate-600">
                    <CardHeader>
                      <div className="mx-auto bg-primary/20 text-primary p-3 rounded-lg w-fit mb-3">
                        <useCase.icon className="h-7 w-7" />
                      </div>
                      <CardTitle className="text-slate-100">{useCase.title}</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <p className="text-slate-300">{useCase.description}</p>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </div>
          </div>
        </motion.section>
        
        {/* Trust & Security Section (Largely remains the same, as it's still relevant and strong) */}
        <motion.section
          id="trust-security"
          className="py-12 md:py-16 container mx-auto px-4 bg-slate-800/80 backdrop-blur-md rounded-lg shadow-xl"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.2 }}
          variants={sectionVariants}
        >
          <div className="text-center space-y-3 mb-10">
            <h2 className="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-100">Sự Tin Tưởng Của Bạn, Ưu Tiên Của Chúng Tôi</h2>
            <p className="text-lg text-slate-300 max-w-xl mx-auto">
              Được xây dựng với bảo mật cấp tổ chức và tuân thủ quy định là cốt lõi.
            </p>
          </div>
          <div className="grid md:grid-cols-2 gap-8 items-start">
            <div className="space-y-6">
              <ul className="space-y-3">
                 {[
                  { icon: ShieldCheck, text: "Công nghệ mã hóa 256-bit tiên tiến và hỗ trợ ví lạnh, HSM cho các giải pháp lưu ký an toàn." },
                  { icon: Building, text: "Kiến trúc cho phép tuân thủ đầy đủ các quy định của Ngân hàng Nhà nước, KYC/AML." },
                  { icon: CheckCircle, text: "Được kiểm toán độc lập bởi các công ty an ninh mạng hàng đầu và cộng đồng." },
                  { icon: Info, text: "Minh bạch có kiểm soát: dữ liệu giao dịch có thể được phân quyền truy cập theo yêu cầu của cơ quan quản lý." },
                ].map((item, index) => (
                  <motion.li
                    key={item.text}
                    className="flex items-start"
                    initial={{ opacity: 0, x: -20 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    viewport={{ once: true }}
                    transition={{ duration: 0.4, delay: index * 0.1 }}
                  >
                    <item.icon className="h-6 w-6 text-green-400 mr-3 mt-1 flex-shrink-0" />
                    <span className="text-slate-300">{item.text}</span>
                  </motion.li>
                ))}
              </ul>
            </div>
            <Card className="bg-slate-700/50 backdrop-blur-sm p-6 border-slate-600 shadow-lg">
              <CardContent className="text-center space-y-4">
                 <Landmark className="h-16 w-16 mx-auto text-primary" />
                <blockquote className="text-lg italic text-slate-200">
                  "Fichain cung cấp một hạ tầng blockchain mạnh mẽ và linh hoạt, cho phép chúng tôi nhanh chóng thử nghiệm và triển khai các dịch vụ tài chính số mới mà không cần đầu tư lớn vào việc xây dựng từ đầu."
                </blockquote>
                <p className="font-semibold text-slate-400">– Giám Đốc Khối Chuyển Đổi Số, Một Ngân Hàng Đối Tác</p>
              </CardContent>
            </Card>
          </div>
        </motion.section>


        {/* FAQ Section (Kept for general info) */}
        {/* You can update these questions to be more B2B/tech-focused if needed */}
        <motion.section
          id="faq"
          className="py-12 md:py-16 container mx-auto px-4 bg-slate-800/80 backdrop-blur-md rounded-lg shadow-xl"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.1 }}
          variants={sectionVariants}
        >
          <div className="text-center space-y-3 mb-10">
            <HelpCircle className="h-12 w-12 text-primary mx-auto" />
            <h2 className="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-100">Câu Hỏi Thường Gặp</h2>
          </div>
          <Accordion type="single" collapsible className="w-full max-w-3xl mx-auto">
             {[
              { q: "Tại sao lại chọn xây dựng một Layer 1 riêng thay vì dùng các Layer 1 có sẵn?", a: "Các Layer 1 phổ biến thường là mạng không cần cấp phép (permissionless) và có chi phí giao dịch biến động, không phù hợp cho các nghiệp vụ ngân hàng đòi hỏi tính riêng tư, ổn định và tuân thủ. Fichain được thiết kế để giải quyết những vấn đề này, cung cấp một môi trường được tối ưu hóa cho ngành tài chính." },
              { q: "Fichain có phải là một private blockchain không?", a: "Fichain là một nền tảng linh hoạt. Nó có thể được triển khai dưới dạng một mạng công cộng (public), một mạng riêng tư (private) cho một ngân hàng hoặc một liên minh (consortium) cho nhiều tổ chức, tùy thuộc vào yêu cầu của bài toán cụ thể." },
              { q: "Làm thế nào để các nhà phát triển có thể bắt đầu xây dựng trên Fichain?", a: "Chúng tôi cung cấp tài liệu kỹ thuật chi tiết, bộ công cụ phát triển (SDK), và một mạng thử nghiệm (testnet) công khai. Do Fichain tương thích EVM, các nhà phát triển có kinh nghiệm với Solidity và Ethereum có thể bắt đầu một cách nhanh chóng. Hãy truy cập mục 'Tài Liệu' của chúng tôi để biết thêm chi tiết." },
              { q: "Mô hình kinh doanh và chi phí khi sử dụng nền tảng Fichain là gì?", a: "Chúng tôi có các mô hình hợp tác linh hoạt, bao gồm phí giấy phép, phí giao dịch hoặc chia sẻ doanh thu tùy thuộc vào quy mô và loại hình dự án. Vui lòng liên hệ với đội ngũ của chúng tôi để được tư vấn cụ thể." },
            ].map((faq, i) => (
              <AccordionItem value={`item-${i + 1}`} key={i} className="border-b-slate-700">
                <AccordionTrigger className="text-lg text-left hover:no-underline text-slate-200 hover:text-slate-100">{faq.q}</AccordionTrigger>
                <AccordionContent className="text-base text-slate-300">{faq.a}</AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        </motion.section>

        {/* UPDATED: Team Section with a stronger headline */}
        <motion.section
          id="team"
          className="py-12 md:py-16 container mx-auto px-4 bg-slate-800/80 backdrop-blur-md rounded-lg shadow-xl"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.05 }}
          variants={sectionVariants}
        >
          <div className="text-center space-y-3 mb-10">
            <UsersRound className="h-12 w-12 text-primary mx-auto" />
            <h2 className="text-3xl sm:text-4xl font-semibold tracking-tight text-slate-100">Đội Ngũ Chuyên Gia Việt Nam Làm Chủ Công Nghệ Lõi</h2>
            <p className="text-lg text-slate-300 max-w-2xl mx-auto">
              Sự kết hợp giữa kinh nghiệm chuyên sâu về blockchain Layer 1, Core Banking và Khoa học dữ liệu tạo nên sức mạnh của Fichain.
            </p>
          </div>
          {/* The team member grid is already very strong and detailed, no changes needed here */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 gap-8 items-stretch">
            {[
              { name: "Hiếu Phan", role: "Trưởng nhóm", field: "Phát triển Layer 1", experience: "Chuyên sâu về thiết kế, xây dựng và tối ưu core blockchain engine (transaction pool, consensus, P2P, state machine), EVM. Có kinh nghiệm triển khai nhiều cơ chế đồng thuận (PoW, PoS, DPoS, PoA), tối ưu hiệu suất cao (hàng nghìn TPS), vận hành mainnet Layer 1 (>2 năm), bảo mật (slashing, BLS, HSM), giám sát (Prometheus, Grafana), CI/CD, xử lý sự cố và cố vấn kỹ thuật.", img: "/team/hieu.png", linkedin: "#" },
              { name: "Hân Phạm", role: "Lập trình viên Layer 1 & DevOps", field: "Phát triển Layer 1, DevOps", experience: "Phát triển và bảo trì các thành phần cốt lõi blockchain (Golang), triển khai/kiểm thử node, tối ưu hiệu suất. Thiết lập CI/CD (AWS/GCP). Phát triển smart contract (Solidity, ERC-20/721/1155), DApp frontend (React/Next.js, Web3.js/Ethers.js). Nghiên cứu Layer 2, cross-chain, oracle.", img: "/team/han.jpg", linkedin: "#" },
              { name: "Nguy Nguyễn", role: "Lập trình viên DApps", field: "Phát triển DApp", experience: "Chuyên sâu thiết kế, phát triển smart contract (Solidity, Hardhat, tối ưu gas, bảo mật). Phát triển frontend DApp (React, Web3.js/Ethers.js) và backend (Golang). Kinh nghiệm với GameFi (P2E, NFT Staking, DAO, Marketplace), token standards (ERC-20/721/1155), DEX, Cross-chain Bridge, và xử lý dữ liệu on-chain/off-chain (The Graph).", img: "/team/nguy.jpg", linkedin: "#" },
              { name: "Hải Trần", role: "Lập trình viên Core Banking", field: "Phát triển Core Banking", experience: "Tích hợp hệ thống core banking quy mô lớn với các hệ thống vệ tinh (eKYC, CRM, CIC, ví điện tử). Thành thạo các giao thức tích hợp (SOAP, REST, MQ, SFTP, ISO 8583). Kinh nghiệm làm việc với đối tác, vận hành hệ thống thanh toán 24/7, bảo mật dữ liệu (mã hóa, 2FA, SSL/TLS), API Gateway, và tuân thủ quy định NHNN.", img: "/team/hai.png", linkedin: "#" },
              { name: "Hùng Hà", role: "Chuyên viên Phân tích Dữ liệu Ngân hàng", field: "Phân tích Dữ liệu Ngân hàng", experience: "Tích hợp hệ thống core banking. Thiết kế, tối ưu và vận hành CSDL lớn (PostgreSQL, MySQL, MongoDB, Cassandra) cho ứng dụng ngân hàng. Xây dựng quy trình ETL/ELT. Phát triển và triển khai pipeline AI/ML (Python, scikit-learn, XGBoost, TensorFlow/PyTorch, Airflow, MLflow) cho các bài toán tài chính (phân loại rủi ro, phát hiện gian lận, scoring tín dụng).", img: "/team/hung.png", linkedin: "#" },
            ].map((member, index) => (
              <motion.div
                key={member.name}
                className="flex flex-col text-center p-4 border border-slate-700 rounded-lg hover:shadow-xl transition-shadow bg-slate-700/70 backdrop-blur-sm hover:bg-slate-600/80 h-full shadow-lg"
                initial={{ opacity: 0, scale: 0.9 }}
                whileInView={{ opacity: 1, scale: 1 }}
                viewport={{ once: true }}
                transition={{ duration: 0.4, delay: index * 0.1 }}
              >
                <div className="relative mx-auto w-28 h-28 sm:w-32 sm:h-32 rounded-full overflow-hidden mb-4 shadow-md flex-shrink-0">
                  <Image src={member.img} alt={member.name} layout="fill" objectFit="cover" />
                </div>
                <p className="font-semibold text-lg text-slate-100">{member.name}</p>
                <p className="text-sm text-primary font-medium my-1">{member.role}</p>
                <p className="text-xs text-cyan-400 uppercase tracking-wider mb-2 font-semibold">{member.field}</p>
                <div className="text-left mt-2 flex-grow">
                  <p className="text-xs text-slate-400 leading-relaxed">
                    <span className="font-semibold text-slate-300">Kinh nghiệm nổi bật:</span> {member.experience}
                  </p>
                </div>
                <a
                  href={member.linkedin || '#'}
                  target="_blank"
                  rel="noopener noreferrer"
                  aria-label={`Hồ sơ LinkedIn của ${member.name}`}
                  className="text-slate-400 hover:text-primary mt-3 inline-block self-center pt-2 flex-shrink-0"
                >
                  <Linkedin className="h-5 w-5" />
                </a>
              </motion.div>
            ))}
          </div>
        </motion.section>

        {/* UPDATED: Final Call to Action for B2B */}
        <motion.section
          id="cta"
          className="py-16 md:py-24 text-center bg-gradient-to-r from-primary to-primary/80 text-primary-foreground rounded-lg shadow-xl"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, amount: 0.3 }}
          variants={sectionVariants}
        >
          <div className="container mx-auto px-4 space-y-6">
            <h2 className="text-3xl sm:text-4xl font-bold">Sẵn Sàng Xây Dựng Tương Lai Tài Chính Cùng Chúng Tôi?</h2>
            <p className="text-lg sm:text-xl max-w-2xl mx-auto text-primary-foreground/90">
              Hãy trở thành đối tác của Fichain và cùng nhau kiến tạo những giải pháp tài chính - ngân hàng đột phá trên nền tảng blockchain Layer 1 của người Việt.
            </p>
            <div className="space-y-4 sm:space-y-0 sm:flex sm:justify-center sm:space-x-4">
              <Button
                variant="outline"
                className="w-full sm:w-auto text-lg px-8 py-3 text-white border-gray-300 hover:bg-white/20"
              >
                Liên Hệ Hợp Tác
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
              <Button
                variant="outline"
                className="w-full sm:w-auto text-lg px-8 py-3 text-white border-gray-300 hover:bg-white/20"
              >
                Tham Gia Cộng Đồng Dev
                <MessageCircle className="ml-2 h-5 w-5" />
              </Button>
            </div>
          </div>
        </motion.section>

        {/* Footer (No major changes needed) */}
        <footer className="py-12 border-t border-slate-700 bg-slate-900">
            <div className="container mx-auto px-4 text-center sm:text-left">
            <div className="grid sm:grid-cols-2 md:grid-cols-4 gap-8 mb-8">
              <div>
                <h3 className="text-xl font-bold mb-2 text-primary">Fichain</h3>
                <p className="text-sm text-slate-400">Nền tảng Blockchain Layer 1 cho ngành Tài chính - Ngân hàng.</p>
              </div>
              <div>
                <h4 className="font-semibold mb-3 text-slate-200">Sản phẩm</h4>
                <ul className="space-y-2 text-sm">
                  <li><a href="#platform-core" className="text-slate-300 hover:text-primary">Nền Tảng</a></li>
                  <li><a href="#features" className="text-slate-300 hover:text-primary">Tính Năng</a></li>
                  <li><a href="#use-cases" className="text-slate-300 hover:text-primary">Ứng Dụng</a></li>
                  <li><a href="#faq" className="text-slate-300 hover:text-primary">FAQ</a></li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold mb-3 text-slate-200">Dành cho Nhà phát triển</h4>
                <ul className="space-y-2 text-sm">
                  <li><Link href="/docs" className="text-slate-300 hover:text-primary">Tài Liệu Kỹ Thuật</Link></li>
                  <li><Link href="/github" className="text-slate-300 hover:text-primary">Github</Link></li>
                  <li><Link href="/sdk" className="text-slate-300 hover:text-primary">SDKs & APIs</Link></li>
                  <li><Link href="/bug-bounty" className="text-slate-300 hover:text-primary">Bug Bounty</Link></li>
                </ul>
              </div>
              <div>
                <h4 className="font-semibold mb-3 text-slate-200">Kết Nối</h4>
                <div className="flex justify-center sm:justify-start space-x-4 mb-3">
                  <a href="https://twitter.com/fichain" target="_blank" rel="noopener noreferrer" className="text-slate-400 hover:text-primary"><Twitter className="h-6 w-6" /></a>
                  <a href="https://discord.gg/fichain" target="_blank" rel="noopener noreferrer" className="text-slate-400 hover:text-primary"><MessageCircle className="h-6 w-6" /></a>
                  <a href="https://linkedin.com/company/fichain" target="_blank" rel="noopener noreferrer" className="text-slate-400 hover:text-primary"><Linkedin className="h-6 w-6" /></a>
                </div>
                <p className="text-sm text-slate-400">Hợp tác: <a href="mailto:partner@fichain.online" className="hover:text-primary text-slate-300">partner@fichain.online</a></p>
              </div>
            </div>
            <div className="text-sm text-slate-400 pt-8 border-t border-slate-700 text-center">
              © {new Date().getFullYear()} Fichain Inc. Bảo lưu mọi quyền.
              <p className="mt-1">Hồ Chí Minh, Việt Nam</p>
            </div>
          </div>
        </footer>
      </div>
    </>
  );
}
