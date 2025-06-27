import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { ParticlesWrapper } from "@/components/ParticlesWrapper";
import { Providers } from './providers'

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});


export const metadata: Metadata = {
  // metadataBase: new URL('https://yourdomain.com'), // Optional: Set a base URL if not set in root layout
  title: "Fichain - Tương Lai Của Ngân Hàng: An Toàn, Minh Bạch, Phi Tập Trung",
  description: "Quản lý Tiền Mã Hóa và Tiền Pháp Định Liền Mạch với Fichain. Trải nghiệm sức mạnh của DeFi với sự đơn giản của ngân hàng truyền thống. Kết nối ngân hàng truyền thống với các công cụ DeFi hiện đại.",
  keywords: [
    "Fichain",
    "ngân hàng",
    "blockchain",
    "DeFi",
    "tài chính phi tập trung",
    "an toàn",
    "minh bạch",
    "tiền mã hóa",
    "cryptocurrency",
    "tiền pháp định",
    "ví đa tài sản",
    "chuyển tiền toàn cầu",
    "core banking",
    "eKYC",
    "fintech",
    "web3 banking",
    "ngân hàng số"
  ],
  authors: [{ name: "Fichain Team", url: "https://yourdomain.com/team" }], // Adjust URL if you have a team page
  //manifest: "/site.webmanifest", // If you have a PWA manifest

  // Open Graph (for social sharing like Facebook, LinkedIn)
  openGraph: {
    title: "Fichain - Tương Lai Của Ngân Hàng: An Toàn, Minh Bạch, Phi Tập Trung",
    description: "Quản lý Tiền Mã Hóa và Tiền Pháp Định Liền Mạch với Fichain. Trải nghiệm sức mạnh của DeFi với sự đơn giản của ngân hàng truyền thống.",
    url: "https://yourdomain.com", // Replace with your actual page URL
    siteName: "Fichain",
    images: [
      {
        url: "/hero.png", // Path to your hero image in the public folder
        width: 1200,      // Recommended width
        height: 630,     // Recommended height
        alt: "Nền tảng Fichain công nghệ blockchain",
      },
      // You can add more images if needed
      // {
      //   url: "/og-image-alternative.png",
      //   width: 800,
      //   height: 600,
      //   alt: "Another image for Fichain",
      // },
    ],
    locale: "vi_VN", // Specify the locale
    type: "website",
  },

  // Twitter Card (for sharing on Twitter)
  twitter: {
    card: "summary_large_image", // Use "summary_large_image" if you have a prominent image
    title: "Fichain - Tương Lai Của Ngân Hàng: An Toàn, Minh Bạch, Phi Tập Trung",
    description: "Quản lý Tiền Mã Hóa và Tiền Pháp Định Liền Mạch với Fichain.",
    // site: "@fichain_handle", // Your Twitter handle (optional)
    // creator: "@creator_handle", // Twitter handle of the content creator (optional)
    images: ["/hero.png"], // Path to your image, same as og:image
  },

  // Icons
  icons: {
    icon: "/favicon.ico", // Standard favicon
  },

  // Robots (SEO directives for crawlers)
  robots: {
    index: true,  // Allow indexing of this page
    follow: true, // Allow crawlers to follow links from this page
    nocache: false, // Allow caching (can be true if content changes very frequently)
    googleBot: { // Specific directives for GoogleBot
      index: true,
      follow: true,
      noimageindex: false, // Allow Google to index images on this page
      'max-video-preview': -1,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },

};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased bg-slate-900 text-slate-200`}
      >
        <ParticlesWrapper>{children}</ParticlesWrapper>
      </body>
    </html>
  );
}
