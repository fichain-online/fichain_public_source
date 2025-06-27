'use client'

import React from 'react';

// --- UI and Icon Imports ---
import {
  ArrowRightLeft, Search, PiggyBank, Coins, ReceiptText,
  FilePlus2, ClipboardCheck, Gem, GitCompareArrows, CreditCard
} from 'lucide-react';

// (Dapp interface and dappsData array remain exactly the same)
interface Dapp {
  id: number;
  title: string;
  description: string;
  icon: React.ElementType;
  href: string;
}

const dappsData: Dapp[] = [
    // ... your dapps data is unchanged
    { id: 1, title: 'Chuyển Khoản', description: 'Thực hiện chuyển và nhận VNĐ một cách an toàn và nhanh chóng trên blockchain.', icon: ArrowRightLeft, href: '/dapps/transfer', },
    { id: 2, title: 'Trình Khám Phá (Explorer)', description: 'Theo dõi các giao dịch, địa chỉ ví và thông tin chi tiết của các khối.', icon: Search, href: '/dapps/explorer', },
    { id: 3, title: 'Gửi Tiết Kiệm', description: 'Gửi tài sản của bạn để nhận lãi suất hấp dẫn và an toàn với smart contract.', icon: PiggyBank, href: '/dapps/saving', },
    { id: 4, title: 'Đầu Tư Vàng', description: 'Mua, bán và lưu trữ vàng được mã hóa (tokenized gold) trên blockchain.', icon: Gem, href: '/dapps/gold-invest', },
    { id: 5, title: 'Đầu Tư Token', description: 'Giao dịch và đầu tư vào các loại token tiềm năng trong hệ sinh thái.', icon: Coins, href: '#', },
    { id: 6, title: 'Thanh Toán Tiện Ích', description: 'Chi trả hóa đơn điện, nước, internet, v.v. một cách nhanh chóng và minh bạch.', icon: ReceiptText, href: '/dapps/service-bills', },
    { id: 7, title: 'Tạo Hóa Đơn', description: 'Tạo và gửi các hóa đơn thanh toán điện tử cho đối tác hoặc khách hàng.', icon: FilePlus2, href: '/dapps/create-invoice', },
    { id: 8, title: 'Thanh Toán Hóa Đơn', description: 'Thanh toán các hóa đơn bạn đã nhận một cách dễ dàng và an toàn.', icon: ClipboardCheck, href: '/dapps/pay-invoice', },
    { id: 9, title: 'Cầu Nối (Bridge)', description: 'Di chuyển tài sản của bạn giữa blockchain này và các blockchain khác.', icon: GitCompareArrows, href: '/dapps/bridge', },
    { id: 10, title: 'Dịch vụ thẻ', description: 'Thẻ kỹ thuật số dành cho các giao dịch blockchain an toàn và tức thì.', icon: CreditCard, href: '/dapps/card', }
];

// (DappCard component remains exactly the same)
const DappCard = ({ title, description, icon: Icon, href }: Omit<Dapp, 'id'>) => (
    <a href={href} className="group block bg-gray-800 p-6 rounded-xl shadow-lg hover:bg-gray-700 transition-all duration-300 transform hover:-translate-y-1">
        <div className="flex items-start">
            <div className="flex-shrink-0">
                <Icon className="h-10 w-10 text-cyan-400 group-hover:text-cyan-300 transition-colors" />
            </div>
            <div className="ml-4">
                <h3 className="text-xl font-bold text-white">{title}</h3>
                <p className="mt-2 text-gray-400">{description}</p>
            </div>
        </div>
    </a>
);



// --- HomePage Component is now much cleaner ---
export default function DappsHomePage() {
  // No more useEffect for auth or connection! The layout handles it all.
  
  return (
    <div className="container mx-auto px-4 py-16 sm:py-24">
      <div className="text-center mb-16">
        <h1 className="text-4xl sm:text-5xl font-extrabold tracking-tight text-white">
          Hệ sinh thái dApps
        </h1>
        <p className="mt-4 max-w-2xl mx-auto text-lg text-gray-400">
          Khám phá các ứng dụng phi tập trung được xây dựng trên Blockchain của chúng tôi.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {dappsData.map((dapp) => (
          <DappCard
            key={dapp.id}
            title={dapp.title}
            description={dapp.description}
            icon={dapp.icon}
            href={dapp.href}
          />
        ))}
      </div>
    </div>
  );
}
