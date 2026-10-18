import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Lock, User as UserIcon, AlertCircle, Eye, EyeOff, Sparkles, Shirt, Waves } from 'lucide-react';

export const LoginPage: React.FC = () => {
  const { login, isAuthenticated, user } = useAuth();
  const navigate = useNavigate();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // If already logged in, redirect based on role
  React.useEffect(() => {
    if (isAuthenticated && user) {
      if (user.role === 'operator') {
        navigate('/orders');
      } else {
        navigate('/dashboard');
      }
    }
  }, [isAuthenticated, user, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!username.trim() || !password) {
      setError('Silakan masukkan nama pengguna dan kata sandi.');
      return;
    }

    setIsLoading(true);
    try {
      await login(username.trim(), password);
    } catch (err: any) {
      setError(err.message || 'Login gagal. Periksa username dan password Anda.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleQuickFill = (u: string, p: string) => {
    setUsername(u);
    setPassword(p);
    setError(null);
  };

  return (
    <div className="min-h-screen min-h-[100dvh] bg-gradient-to-br from-[#122822] via-[#1a3e35] to-[#0d1e19] flex items-center justify-center p-3.5 sm:p-6 lg:p-8 relative overflow-x-hidden font-sans">
      {/* Ambient background bokeh orbs matching reference aesthetic */}
      <div className="absolute top-6 left-10 w-28 sm:w-44 h-28 sm:h-44 rounded-full bg-emerald-400/10 blur-2xl pointer-events-none" />
      <div className="absolute top-1/4 right-10 sm:right-20 w-48 sm:w-72 h-48 sm:h-72 rounded-full bg-teal-300/10 blur-3xl pointer-events-none" />
      <div className="absolute bottom-10 left-1/4 w-44 sm:w-64 h-44 sm:h-64 rounded-full bg-emerald-500/10 blur-3xl pointer-events-none" />
      <div className="absolute -bottom-10 -right-10 w-64 sm:w-80 h-64 sm:h-80 rounded-full bg-teal-600/15 blur-3xl pointer-events-none" />

      {/* Floating subtle bokeh dots */}
      <div className="absolute top-12 left-1/4 w-6 sm:w-8 h-6 sm:h-8 rounded-full bg-white/5 pointer-events-none" />
      <div className="absolute top-36 left-8 sm:left-12 w-10 sm:w-14 h-10 sm:h-14 rounded-full bg-white/5 pointer-events-none" />
      <div className="absolute bottom-20 right-1/4 w-8 sm:w-12 h-8 sm:h-12 rounded-full bg-white/5 pointer-events-none" />

      {/* Main Container Card: Full Responsive layout across Mobile, Tablet, and Desktop */}
      <div className="w-full max-w-5xl bg-[#1b3e34] rounded-2xl sm:rounded-3xl lg:rounded-[36px] shadow-2xl shadow-black/40 overflow-hidden flex flex-col lg:flex-row border border-emerald-700/30 relative z-10 my-auto">
        
        {/* ==================================================================== */}
        {/* LEFT COLUMN: Clean White Section with Organic Curve & Illustration   */}
        {/* ==================================================================== */}
        <div className="w-full lg:w-[50%] bg-white p-5 sm:p-8 lg:p-10 flex flex-col justify-between relative overflow-hidden">
          
          {/* Organic curved shape flowing into the right panel on desktop */}
          <div className="absolute -right-16 top-0 bottom-0 w-32 bg-white hidden lg:block rounded-r-[60px] pointer-events-none" />
          
          {/* Top Brand Logo */}
          <div className="relative z-10 flex items-center justify-between lg:justify-start gap-3">
            <div className="flex items-center gap-2.5 sm:gap-3">
              <div className="w-9 h-9 sm:w-11 sm:h-11 rounded-xl sm:rounded-2xl bg-gradient-to-tr from-teal-700 to-emerald-600 flex items-center justify-center text-white shadow-md shadow-teal-700/20 shrink-0">
                <Shirt className="w-5 h-5 sm:w-6 sm:h-6" />
              </div>
              <div>
                <div className="flex items-center gap-1.5">
                  <span className="font-extrabold text-sm sm:text-base tracking-wider text-slate-800">SATSET</span>
                  <span className="text-[9px] sm:text-[10px] uppercase font-bold tracking-widest px-1.5 py-0.5 rounded bg-teal-50 text-teal-700 border border-teal-200">
                    Laundry
                  </span>
                </div>
                <p className="text-[10px] sm:text-[11px] font-medium text-slate-400">Sistem POS &amp; Manajemen Terpadu</p>
              </div>
            </div>

            {/* Mobile badge indicator */}
            <span className="lg:hidden text-[10px] text-teal-700 bg-teal-50 border border-teal-200/60 px-2 py-0.5 rounded-full font-medium">
              v1.0 Live
            </span>
          </div>

          {/* Center Illustration: Visible on tablets and desktop, gracefully collapsed on small mobile */}
          <div className="relative z-10 my-3 sm:my-6 lg:my-8 hidden sm:flex flex-col items-center justify-center">
            {/* Ambient decorative soft shapes & bubbles */}
            <div className="relative w-full max-w-[240px] sm:max-w-[280px] lg:max-w-[320px] aspect-square flex items-center justify-center">
              <div className="absolute inset-2 sm:inset-4 rounded-[32px] sm:rounded-[40px] bg-gradient-to-br from-teal-50 via-emerald-50/70 to-slate-100/80 -rotate-3" />
              
              {/* Floating decorative bubbles */}
              <div className="absolute top-2 left-4 sm:top-4 sm:left-6 w-7 h-7 sm:w-9 sm:h-9 rounded-full bg-teal-100/60 flex items-center justify-center text-teal-500 animate-pulse">
                <Sparkles className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
              </div>
              <div className="absolute top-8 right-6 sm:top-10 sm:right-8 w-9 h-9 sm:w-11 sm:h-11 rounded-full bg-emerald-100/50 flex items-center justify-center text-emerald-600">
                <Waves className="w-4 h-4 sm:w-5 sm:h-5" />
              </div>
              <div className="absolute bottom-6 left-6 w-5 h-5 rounded-full bg-teal-200/40" />
              <div className="absolute bottom-10 right-4 w-7 h-7 rounded-full bg-emerald-200/50" />

              {/* Vector Art: Washing Machine & Fresh Linen */}
              <svg className="w-40 h-40 sm:w-48 sm:h-48 lg:w-52 lg:h-52 relative z-10 drop-shadow-md" viewBox="0 0 200 200" fill="none" xmlns="http://www.w3.org/2000/svg">
                {/* Washing Machine Body */}
                <rect x="45" y="35" width="110" height="135" rx="18" fill="#1b3e34" />
                <rect x="52" y="42" width="96" height="121" rx="14" fill="#244d41" />
                
                {/* Control Panel */}
                <rect x="56" y="48" width="88" height="24" rx="8" fill="#1b3e34" />
                <circle cx="68" cy="60" r="4" fill="#52948e" />
                <circle cx="80" cy="60" r="3" fill="#88c5be" />
                <rect x="94" y="56" width="42" height="8" rx="4" fill="#122a23" />
                <circle cx="128" cy="60" r="2" fill="#34d399" />

                {/* Machine Drum Outer Ring */}
                <circle cx="100" cy="115" r="42" fill="#1a3a31" stroke="#366b5e" strokeWidth="3" />
                <circle cx="100" cy="115" r="34" fill="url(#waterGrad)" />
                
                {/* Water Waves */}
                <path d="M72 120 C82 108, 92 128, 106 116 C118 106, 126 122, 128 118" stroke="#ffffff" strokeWidth="3" strokeLinecap="round" strokeOpacity="0.8" />
                <path d="M78 128 C88 120, 98 134, 112 124 C120 118, 124 126, 125 124" stroke="#88c5be" strokeWidth="2.5" strokeLinecap="round" strokeOpacity="0.6" />
                
                {/* Clean Sparkles */}
                <circle cx="95" cy="100" r="2" fill="#ffffff" />
                <circle cx="112" cy="108" r="1.5" fill="#ffffff" />
                <circle cx="86" cy="118" r="1.5" fill="#ffffff" />

                {/* Base Shadow */}
                <ellipse cx="100" cy="178" rx="65" ry="7" fill="#000000" fillOpacity="0.1" />

                {/* Folded Towels next to machine */}
                <rect x="135" y="132" width="34" height="10" rx="4" fill="#52948e" />
                <rect x="133" y="142" width="38" height="11" rx="4" fill="#88c5be" />
                <rect x="130" y="153" width="44" height="13" rx="5" fill="#e2e8f0" stroke="#cbd5e1" strokeWidth="1" />

                <defs>
                  <linearGradient id="waterGrad" x1="70" y1="85" x2="130" y2="145" gradientUnits="userSpaceOnUse">
                    <stop stopColor="#38bdf8" stopOpacity="0.85" />
                    <stop offset="0.5" stopColor="#0284c7" stopOpacity="0.9" />
                    <stop offset="1" stopColor="#0f766e" />
                  </linearGradient>
                </defs>
              </svg>
            </div>

            <p className="text-center text-xs font-medium text-slate-500 mt-2 max-w-xs">
              Mulai hari dengan alur operasional laundry yang cepat, rapi, dan otomatis.
            </p>
          </div>

          {/* Bottom Copyright on Desktop */}
          <div className="relative z-10 pt-2 border-t border-slate-100 hidden lg:flex items-center justify-between text-[11px] text-slate-400">
            <span>&copy; {new Date().getFullYear()} SatSet Laundry</span>
            <span className="text-[10px] text-teal-600 bg-teal-50 px-2 py-0.5 rounded-full font-medium">
              v1.0 Production
            </span>
          </div>
        </div>

        {/* ==================================================================== */}
        {/* RIGHT COLUMN: Forest Green Login Form (Responsive on all viewports)  */}
        {/* ==================================================================== */}
        <div className="w-full lg:w-[50%] bg-[#1b3e34] p-5 sm:p-8 lg:p-11 flex flex-col justify-between text-white relative">
          
          {/* Subtle background glow on right panel */}
          <div className="absolute top-0 right-0 w-36 sm:w-48 h-36 sm:h-48 bg-teal-400/10 rounded-full blur-3xl pointer-events-none" />

          {/* Form Header */}
          <div>
            <h1 className="text-2xl sm:text-3xl lg:text-4xl font-bold tracking-tight text-white">
              Login
            </h1>
            <p className="text-xs sm:text-sm text-emerald-200/70 mt-1">
              Masukkan kredensial akun Anda untuk masuk ke sistem.
            </p>

            {/* Error Message Alert */}
            {error && (
              <div className="mt-3 sm:mt-4 rounded-xl bg-rose-950/70 border border-rose-500/40 p-3 flex items-start gap-2.5 text-xs text-rose-200 animate-fadeIn">
                <AlertCircle className="w-4 h-4 text-rose-400 shrink-0 mt-0.5" />
                <span>{error}</span>
              </div>
            )}

            {/* Login Form */}
            <form onSubmit={handleSubmit} className="mt-5 sm:mt-6 space-y-3.5 sm:space-y-4">
              {/* Username Field */}
              <div>
                <label className="block text-xs sm:text-sm font-medium text-emerald-100/90 mb-1 sm:mb-1.5">
                  Username
                </label>
                <div className="relative flex items-center">
                  <div className="absolute left-3.5 sm:left-4 text-emerald-300/50 pointer-events-none">
                    <UserIcon className="w-4 h-4" />
                  </div>
                  <input
                    type="text"
                    required
                    autoFocus
                    value={username}
                    onChange={e => setUsername(e.target.value)}
                    placeholder="Enter your username"
                    className="w-full rounded-full bg-[#132d26] border border-emerald-700/50 text-white placeholder-emerald-200/30 text-base sm:text-sm pl-10 sm:pl-11 pr-4 py-2.5 sm:py-3.5 focus:outline-none focus:ring-2 focus:ring-teal-400 focus:border-transparent transition-all min-h-[44px]"
                  />
                </div>
              </div>

              {/* Password Field */}
              <div>
                <label className="block text-xs sm:text-sm font-medium text-emerald-100/90 mb-1 sm:mb-1.5">
                  Password
                </label>
                <div className="relative flex items-center">
                  <div className="absolute left-3.5 sm:left-4 text-emerald-300/50 pointer-events-none">
                    <Lock className="w-4 h-4" />
                  </div>
                  <input
                    type={showPassword ? 'text' : 'password'}
                    required
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                    placeholder="Enter your password"
                    className="w-full rounded-full bg-[#132d26] border border-emerald-700/50 text-white placeholder-emerald-200/30 text-base sm:text-sm pl-10 sm:pl-11 pr-11 py-2.5 sm:py-3.5 focus:outline-none focus:ring-2 focus:ring-teal-400 focus:border-transparent transition-all min-h-[44px]"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    tabIndex={-1}
                    className="absolute right-2 sm:right-3 p-2 text-emerald-300/50 hover:text-emerald-200 transition-colors rounded-full focus:outline-none"
                    aria-label={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
                  >
                    {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
                
                {/* Forgot Password Link matching reference */}
                <div className="mt-1.5 sm:mt-2 text-right">
                  <span className="text-[11px] text-emerald-300/60 hover:text-emerald-200 cursor-pointer transition-colors">
                    Lupa Password?
                  </span>
                </div>
              </div>

              {/* Action Button: Styled like reference "Login to Wifi" */}
              <div className="pt-1 sm:pt-2">
                <button
                  type="submit"
                  disabled={isLoading}
                  className="w-full rounded-full bg-[#52948e] hover:bg-[#5da7a0] active:scale-[0.98] text-white font-medium text-sm sm:text-base py-3 sm:py-3.5 px-6 min-h-[46px] shadow-lg shadow-black/20 transition-all flex items-center justify-center gap-2 cursor-pointer disabled:opacity-60 disabled:cursor-not-allowed"
                >
                  {isLoading ? (
                    <>
                      <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                      <span>Memverifikasi...</span>
                    </>
                  ) : (
                    <span>Masuk ke SatSet</span>
                  )}
                </button>
              </div>
            </form>

            {/* Quick Login Helper Chips (Convenience for POS/Laundry testing) */}
            <div className="mt-4 sm:mt-5 pt-3.5 sm:pt-4 border-t border-emerald-800/40">
              <p className="text-[10px] sm:text-[11px] text-emerald-200/60 mb-2 font-medium">
                Pilih akun cepat (Demo):
              </p>
              <div className="flex flex-wrap gap-1.5 sm:gap-2">
                <button
                  type="button"
                  onClick={() => handleQuickFill('owner', 'owner123')}
                  className="text-[10px] sm:text-[11px] bg-emerald-900/60 hover:bg-emerald-800/80 active:scale-95 text-emerald-200 border border-emerald-700/50 px-3 py-1.5 sm:py-1 rounded-full transition-all cursor-pointer"
                >
                  👑 Owner
                </button>
                <button
                  type="button"
                  onClick={() => handleQuickFill('kasir', 'kasir123')}
                  className="text-[10px] sm:text-[11px] bg-emerald-900/60 hover:bg-emerald-800/80 active:scale-95 text-emerald-200 border border-emerald-700/50 px-3 py-1.5 sm:py-1 rounded-full transition-all cursor-pointer"
                >
                  💼 Kasir
                </button>
                <button
                  type="button"
                  onClick={() => handleQuickFill('operator', 'operator123')}
                  className="text-[10px] sm:text-[11px] bg-emerald-900/60 hover:bg-emerald-800/80 active:scale-95 text-emerald-200 border border-emerald-700/50 px-3 py-1.5 sm:py-1 rounded-full transition-all cursor-pointer"
                >
                  ⚙️ Operator
                </button>
              </div>
            </div>
          </div>

          {/* Bottom Support / Terms Links */}
          <div className="mt-5 sm:mt-6 pt-3 flex flex-col sm:flex-row items-center justify-between text-center sm:text-left text-[11px] text-emerald-300/50 gap-2 border-t border-emerald-800/30 lg:border-t-0">
            <span className="hover:underline cursor-pointer">Syarat &amp; Ketentuan</span>
            <span className="hover:underline cursor-pointer">Pusat Bantuan SatSet</span>
          </div>

        </div>
      </div>
    </div>
  );
};
