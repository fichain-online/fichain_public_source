import Link from 'next/link';
import Image from 'next/image';
import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Menu, X, LogIn, LogOut, LayoutDashboard } from 'lucide-react';
import { useRouter } from 'next/navigation';

// --- Zustand Store Imports ---
import { useAuthStore } from '@/stores/authStore';

const navLinks = [
  { href: '/#hero', label: 'Trang Chủ' },
  { href: '/#features', label: 'Tính Năng' },
  { href: '/#how-it-works', label: 'Cách Hoạt Động' },
  { href: '/#faq', label: 'FAQ' },
  { href: '/#team', label: 'Đội Ngũ' },
];

export default function Header() {
  const [isOpen, setIsOpen] = useState(false);
  const [isScrolled, setIsScrolled] = useState(false);
  const router = useRouter();

  const { isAuthenticated, logout } = useAuthStore();

  const handleAuth = () => {
    if (isAuthenticated) {
      logout();
      // Optional: redirect to home after logout
      router.push('/');
    } else {
      router.push('/login');
    }
  };


  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 50);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  // Close menu on route change
  useEffect(() => {
    setIsOpen(false);
  }, [router]);


  const toggleMenu = () => setIsOpen(!isOpen);
  const closeMenu = () => setIsOpen(false);

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ease-in-out 
                 ${isScrolled || isOpen ? 'bg-slate-900/95 shadow-lg backdrop-blur-md' : 'bg-transparent'}`}
    >
      <div className="container mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-20">
          {/* Logo */}
          <Link href="/" className="flex-shrink-0 flex items-center" onClick={closeMenu}>
            <Image src="/logo.png" alt="Fichain Logo" width={80} height={80} />
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex space-x-1 lg:space-x-3">
            {navLinks.map((link) => (
              <a
                key={link.label}
                href={link.href}
                onClick={closeMenu}
                className={`px-3 py-2 rounded-md text-sm font-medium transition-colors
                           ${isScrolled || isOpen ? 'text-gray-300 hover:text-white hover:bg-white/10' : 'text-gray-200 hover:text-white hover:bg-black/20'}`}
              >
                {link.label}
              </a>
            ))}
            <Link
              href="/dapps"
              className={`px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center
                         ${isScrolled || isOpen ? 'text-gray-300 hover:text-white hover:bg-white/10' : 'text-gray-200 hover:text-white hover:bg-black/20'}`}
              onClick={closeMenu}
            >
              <LayoutDashboard className="mr-1 h-4 w-4" /> DApps
            </Link>
          </nav>

          {/* Auth Button - Desktop */}
          <div className="hidden md:block">
            <Button
              onClick={handleAuth}
              variant={isScrolled || isOpen ? "outline" : "outline"}
              className={isScrolled || isOpen ? "text-white border-cyan-400 hover:bg-cyan-400 hover:text-slate-900" : "text-white border-gray-300 hover:bg-white/20"}
            >
              {isAuthenticated ? <LogOut className="mr-2 h-4 w-4" /> : <LogIn className="mr-2 h-4 w-4" />}
              {isAuthenticated ? 'Đăng Xuất' : 'Đăng Nhập'}
            </Button>
          </div>

          {/* ✅ ADDED: Mobile Menu Button (Hamburger) */}
          <div className="md:hidden flex items-center">
            <button
              onClick={toggleMenu}
              className="inline-flex items-center justify-center p-2 rounded-md text-gray-300 hover:text-white hover:bg-white/10 focus:outline-none"
              aria-controls="mobile-menu"
              aria-expanded={isOpen}
            >
              <span className="sr-only">Open main menu</span>
              {isOpen ? (
                <X className="block h-6 w-6" aria-hidden="true" />
              ) : (
                <Menu className="block h-6 w-6" aria-hidden="true" />
              )}
            </button>
          </div>

          {/* ❌ REMOVED: Redundant and misplaced mobile auth button was here */}
        </div>
      </div>

      {/* ✅ CORRECTED: Mobile Menu Panel */}
      {/* This now sits outside the main flex container, allowing it to display correctly below the header bar */}
      {isOpen && (
        <div className="md:hidden" id="mobile-menu">
          <nav className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
            {navLinks.map((link) => (
              <a
                key={link.label}
                href={link.href}
                onClick={closeMenu}
                className="block px-3 py-2 rounded-md text-base font-medium text-gray-300 hover:text-white hover:bg-white/10"
              >
                {link.label}
              </a>
            ))}
             <Link
              href="/dapps"
              className="block px-3 py-2 rounded-md text-base font-medium text-gray-300 hover:text-white hover:bg-white/10 flex items-center"
              onClick={closeMenu}
            >
              <LayoutDashboard className="mr-2 h-5 w-5" /> DApps
            </Link>
          </nav>
          <div className="px-4 pt-4 pb-3 border-t border-slate-700">
            <Button
              onClick={() => { handleAuth(); closeMenu(); }}
              variant="outline"
              className="w-full text-white border-cyan-400 hover:bg-cyan-400 hover:text-slate-900"
            >
              {isAuthenticated ? <LogOut className="mr-2 h-4 w-4" /> : <LogIn className="mr-2 h-4 w-4" />}
              {isAuthenticated ? 'Đăng Xuất' : 'Đăng Nhập'}
            </Button>
          </div>
        </div>
      )}
    </header>
  );
}
