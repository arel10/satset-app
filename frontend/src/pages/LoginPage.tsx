import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { AlertCircle, X, Shield, Phone, Sparkles, HelpCircle } from 'lucide-react';

export const LoginPage: React.FC = () => {
  const { login, isAuthenticated, user } = useAuth();
  const navigate = useNavigate();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // Modals / Drawers state
  const [showForgotModal, setShowForgotModal] = useState(false);
  const [showMenuDrawer, setShowMenuDrawer] = useState(false);

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
      setError('Silakan masukkan email / username dan password.');
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
    <div className="min-h-screen min-h-[100dvh] w-full flex flex-col lg:flex-row bg-[#F9CB43] selection:bg-black selection:text-white font-sans">
      
      {/* ==================================================================== */}
      {/* LEFT COLUMN: Warm Yellow Panel with Zen Meditating Woman & Doodles   */}
      {/* ==================================================================== */}
      <div className="w-full lg:w-1/2 bg-[#F9CB43] flex items-center justify-center p-6 sm:p-10 lg:p-14 relative overflow-hidden select-none min-h-[360px] sm:min-h-[460px] lg:min-h-screen">
        
        {/* Top brand tag for mobile view */}
        <div className="absolute top-4 left-5 lg:hidden flex items-center gap-1.5 z-10">
          <span className="font-extrabold text-xs tracking-wider text-black">SATSET</span>
          <span className="text-[9px] uppercase font-bold tracking-widest px-1.5 py-0.5 rounded-full bg-black text-[#F9CB43]">
            Laundry
          </span>
        </div>

        {/* Vector Illustration matching reference design */}
        <div className="w-full max-w-[440px] sm:max-w-[480px] lg:max-w-[500px] flex items-center justify-center">
          <svg
            viewBox="0 0 540 640"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className="w-full h-auto drop-shadow-sm pointer-events-none"
          >
            {/* ----------------- Surrounding Geometric Doodles ----------------- */}
            
            {/* Top-left diagonal slash */}
            <line x1="85" y1="75" x2="98" y2="98" stroke="#000000" strokeWidth="2.5" strokeLinecap="round" />

            {/* Top-center hollow circle */}
            <circle cx="180" cy="55" r="7.5" stroke="#000000" strokeWidth="2" fill="none" />

            {/* Top-right vertical dash */}
            <line x1="315" y1="85" x2="315" y2="108" stroke="#000000" strokeWidth="2.5" strokeLinecap="round" />

            {/* Top-left tilted outline diamond */}
            <rect x="155" y="135" width="16" height="16" rx="2" transform="rotate(45 163 143)" stroke="#000000" strokeWidth="2" fill="none" />

            {/* Top-right small solid square */}
            <rect x="250" y="185" width="12" height="12" fill="#000000" />

            {/* Middle-right horizontal dash */}
            <line x1="365" y1="198" x2="390" y2="198" stroke="#000000" strokeWidth="2.5" strokeLinecap="round" />

            {/* Left-center solid black dot */}
            <circle cx="110" cy="220" r="10.5" fill="#000000" />

            {/* Middle-left hollow circle */}
            <circle cx="48" cy="275" r="8.5" stroke="#000000" strokeWidth="2" fill="none" />

            {/* Middle-left slash */}
            <line x1="105" y1="380" x2="116" y2="396" stroke="#000000" strokeWidth="2.5" strokeLinecap="round" />

            {/* Solid dot near ground */}
            <circle cx="28" cy="445" r="7" fill="#000000" />

            {/* Bottom-left hollow circle */}
            <circle cx="72" cy="548" r="9" stroke="#000000" strokeWidth="2" fill="none" />

            {/* Bottom-left horizontal dash */}
            <line x1="58" y1="635" x2="82" y2="635" stroke="#000000" strokeWidth="2.5" strokeLinecap="round" />

            {/* Bottom-left slash */}
            <line x1="152" y1="585" x2="164" y2="600" stroke="#000000" strokeWidth="2.5" strokeLinecap="round" />

            {/* Bottom-center hollow circle */}
            <circle cx="225" cy="610" r="7" stroke="#000000" strokeWidth="2" fill="none" />

            {/* Bottom-right small solid square */}
            <rect x="330" y="530" width="9" height="9" fill="#000000" />

            {/* Bottom-right tilted outline diamond */}
            <rect x="305" y="605" width="16" height="16" rx="2" transform="rotate(45 313 613)" stroke="#000000" strokeWidth="2" fill="none" />

            {/* Upper-right small solid square */}
            <rect x="300" y="258" width="7" height="7" fill="#000000" />

            {/* ----------------- Ground Line ----------------- */}
            <line x1="45" y1="486" x2="495" y2="486" stroke="#000000" strokeWidth="2.5" strokeLinecap="round" />

            {/* ----------------- Meditating Woman Figure ----------------- */}
            
            {/* Crossed Legs (Lotus Pose) Outline & Fill */}
            <g>
              {/* Left thigh & knee fold */}
              <path
                d="M 125 440 C 95 450, 75 465, 75 478 C 75 488, 105 498, 185 498 C 240 498, 290 482, 345 455"
                fill="#F9CB43"
                stroke="#000000"
                strokeWidth="2.5"
                strokeLinecap="round"
              />

              {/* Right thigh & knee fold */}
              <path
                d="M 345 455 C 385 462, 400 472, 400 480 C 400 490, 360 500, 275 500 C 195 500, 130 485, 95 465"
                fill="#F9CB43"
                stroke="#000000"
                strokeWidth="2.5"
                strokeLinecap="round"
              />

              {/* Inner lap crease folds */}
              <path
                d="M 170 445 C 205 465, 260 465, 305 450"
                stroke="#000000"
                strokeWidth="2"
                strokeLinecap="round"
                fill="none"
              />
              <path
                d="M 195 465 C 220 482, 265 478, 295 465"
                stroke="#000000"
                strokeWidth="2"
                strokeLinecap="round"
                fill="none"
              />

              {/* Bare feet contours */}
              {/* Right foot resting over left shin */}
              <path
                d="M 130 468 C 115 476, 95 485, 90 495 C 85 506, 96 507, 108 503 C 122 497, 132 485, 140 472"
                fill="#F9CB43"
                stroke="#000000"
                strokeWidth="2"
                strokeLinecap="round"
              />
              {/* Left foot resting over right thigh */}
              <path
                d="M 295 460 C 315 460, 335 466, 340 478 C 338 488, 322 488, 302 482"
                fill="#F9CB43"
                stroke="#000000"
                strokeWidth="2"
                strokeLinecap="round"
              />
            </g>

            {/* Black Long-Sleeved Sweater / Top */}
            <path
              d="M 205 320 C 190 322, 165 330, 155 342 C 146 355, 168 405, 185 448 L 245 448 C 262 405, 284 355, 275 342 C 265 330, 240 322, 225 320 C 218 335, 212 335, 205 320 Z"
              fill="#000000"
            />

            {/* Left Arm and Chin Mudra Hand */}
            <g>
              {/* Left arm sleeve */}
              <path
                d="M 155 342 C 138 360, 115 395, 108 430 C 105 445, 116 448, 126 440 C 134 430, 148 395, 156 365 Z"
                fill="#000000"
              />
              {/* Left forearm bending upwards */}
              <path
                d="M 108 430 C 112 405, 120 380, 128 360 L 138 364 C 132 384, 124 410, 124 438 Z"
                fill="#000000"
              />
              {/* Chin mudra hand (circle thumb & index, 3 fingers up) */}
              <path
                d="M 126 358 C 122 352, 112 350, 110 356 C 108 362, 116 366, 122 363"
                stroke="#000000"
                strokeWidth="2"
                fill="#F9CB43"
              />
              <path
                d="M 112 352 C 110 345, 108 338, 108 332 C 108 330, 112 330, 113 334 L 117 350"
                stroke="#000000"
                strokeWidth="1.8"
                strokeLinecap="round"
                fill="none"
              />
              <path
                d="M 117 350 C 116 342, 116 335, 118 330 C 119 329, 122 330, 122 334 L 122 352"
                stroke="#000000"
                strokeWidth="1.8"
                strokeLinecap="round"
                fill="none"
              />
              <path
                d="M 122 352 C 123 345, 125 339, 128 335 C 129 334, 131 336, 130 340 L 127 355"
                stroke="#000000"
                strokeWidth="1.8"
                strokeLinecap="round"
                fill="none"
              />
            </g>

            {/* Right Arm and Chin Mudra Hand */}
            <g>
              {/* Right arm sleeve */}
              <path
                d="M 275 342 C 292 360, 315 395, 322 430 C 325 445, 314 448, 304 440 C 296 430, 282 395, 274 365 Z"
                fill="#000000"
              />
              {/* Right forearm bending upwards */}
              <path
                d="M 322 430 C 318 405, 310 380, 302 360 L 292 364 C 298 384, 306 410, 306 438 Z"
                fill="#000000"
              />
              {/* Chin mudra hand */}
              <path
                d="M 304 358 C 308 352, 318 350, 320 356 C 322 362, 314 366, 308 363"
                stroke="#000000"
                strokeWidth="2"
                fill="#F9CB43"
              />
              <path
                d="M 318 352 C 320 345, 322 338, 322 332 C 322 330, 318 330, 317 334 L 313 350"
                stroke="#000000"
                strokeWidth="1.8"
                strokeLinecap="round"
                fill="none"
              />
              <path
                d="M 313 350 C 314 342, 314 335, 312 330 C 311 329, 308 330, 308 334 L 308 352"
                stroke="#000000"
                strokeWidth="1.8"
                strokeLinecap="round"
                fill="none"
              />
              <path
                d="M 308 352 C 307 345, 305 339, 302 335 C 301 334, 299 336, 300 340 L 303 355"
                stroke="#000000"
                strokeWidth="1.8"
                strokeLinecap="round"
                fill="none"
              />
            </g>

            {/* Neck */}
            <path
              d="M 207 325 L 207 290 C 211 288, 219 288, 223 290 L 223 325 Z"
              fill="#F9CB43"
            />

            {/* Serene Face & Jawline */}
            <path
              d="M 201 250 C 201 285, 210 298, 215 298 C 220 298, 229 285, 229 250 Z"
              fill="#F9CB43"
              stroke="#000000"
              strokeWidth="1.5"
            />

            {/* Closed Curved Peaceful Eyes (Zen expression) */}
            <path
              d="M 206 256 C 208 261, 212 261, 214 256"
              stroke="#000000"
              strokeWidth="2.2"
              strokeLinecap="round"
              fill="none"
            />
            <path
              d="M 218 256 C 220 261, 224 261, 226 256"
              stroke="#000000"
              strokeWidth="2.2"
              strokeLinecap="round"
              fill="none"
            />

            {/* Gentle Smile */}
            <path
              d="M 212 274 Q 215 279 218 274"
              stroke="#000000"
              strokeWidth="1.8"
              strokeLinecap="round"
              fill="none"
            />

            {/* Voluminous Blowing Black Hair (Flowing to the right) */}
            <path
              d="M 202 250 C 192 225, 204 195, 226 195 C 246 195, 258 212, 264 235 C 280 242, 310 246, 355 258 C 395 270, 408 295, 396 322 C 382 345, 345 352, 298 340 C 275 334, 256 320, 246 300 C 238 298, 230 304, 226 312 C 224 280, 221 245, 221 240 C 215 240, 208 244, 202 250 Z"
              fill="#000000"
            />
            
            {/* Left side hair strand framing face */}
            <path
              d="M 203 230 C 190 242, 182 265, 185 295 C 188 318, 198 335, 206 345 C 201 328, 200 302, 202 260 Z"
              fill="#000000"
            />
          </svg>
        </div>
      </div>

      {/* ==================================================================== */}
      {/* RIGHT COLUMN: Pure White Editorial Login Form                        */}
      {/* ==================================================================== */}
      <div className="w-full lg:w-1/2 bg-white flex flex-col justify-between p-6 sm:p-10 lg:p-16 relative min-h-[520px] lg:min-h-screen">
        
        {/* Top Header Row with Hamburger Menu */}
        <div className="w-full flex items-center justify-between">
          {/* Subtle Desktop SatSet Branding */}
          <div className="hidden lg:flex items-center gap-2">
            <span className="font-extrabold text-sm tracking-wider text-black">SATSET</span>
            <span className="text-[10px] uppercase font-bold tracking-widest px-2 py-0.5 rounded-full bg-neutral-100 text-neutral-800 border border-neutral-200">
              Laundry POS
            </span>
          </div>

          <div className="lg:hidden" />

          {/* Three-Line Hamburger Menu Button (Matches mockup exactly) */}
          <button
            type="button"
            onClick={() => setShowMenuDrawer(true)}
            className="p-2 -mr-2 text-black hover:opacity-70 transition-opacity cursor-pointer rounded-lg focus:outline-none"
            aria-label="Buka menu aplikasi"
          >
            <div className="w-6 flex flex-col items-end gap-1.5">
              <span className="w-6 h-[2.5px] bg-black rounded-full" />
              <span className="w-6 h-[2.5px] bg-black rounded-full" />
              <span className="w-6 h-[2.5px] bg-black rounded-full" />
            </div>
          </button>
        </div>

        {/* Center Container: Main Form with Editorial Typography */}
        <div className="w-full max-w-[360px] sm:max-w-[400px] mx-auto my-auto py-6 sm:py-8">
          
          {/* Editorial Display Heading */}
          <div className="text-center">
            <h1 className="font-serif text-3xl sm:text-4xl lg:text-[44px] font-bold text-black tracking-tight leading-tight">
              Hey Beautiful!
            </h1>
            <p className="font-serif italic text-xl sm:text-2xl text-neutral-900 mt-1 sm:mt-2 mb-7 sm:mb-9 font-normal">
              Login please
            </p>
          </div>

          {/* Error Alert Display */}
          {error && (
            <div className="mb-5 rounded-2xl bg-rose-50 border border-rose-200 p-3.5 flex items-start gap-2.5 text-xs text-rose-800 animate-fadeIn">
              <AlertCircle className="w-4 h-4 text-rose-600 shrink-0 mt-0.5" />
              <div className="flex-1">
                <span>{error}</span>
              </div>
              <button
                type="button"
                onClick={() => setError(null)}
                className="text-rose-400 hover:text-rose-600 text-xs"
              >
                ✕
              </button>
            </div>
          )}

          {/* Login Form */}
          <form onSubmit={handleSubmit} className="space-y-4 sm:space-y-5">
            {/* Email / Username Field */}
            <div>
              <label
                htmlFor="identifier"
                className="block text-xs sm:text-sm font-normal text-neutral-800 mb-1.5 ml-1 font-sans"
              >
                Email
              </label>
              <input
                id="identifier"
                type="text"
                required
                autoFocus
                value={username}
                onChange={e => setUsername(e.target.value)}
                placeholder="Email atau username"
                className="w-full rounded-full border border-black/80 sm:border-[1.8px] sm:border-black bg-white px-5 sm:px-6 py-2.5 sm:py-3 text-sm text-neutral-900 placeholder:text-neutral-400 focus:outline-none focus:ring-2 focus:ring-black/10 focus:border-black transition-all"
              />
            </div>

            {/* Password Field */}
            <div>
              <label
                htmlFor="password"
                className="block text-xs sm:text-sm font-normal text-neutral-800 mb-1.5 ml-1 font-sans"
              >
                Password
              </label>
              <div className="relative flex items-center">
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  required
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  placeholder="Kata sandi akun"
                  className="w-full rounded-full border border-black/80 sm:border-[1.8px] sm:border-black bg-white pl-5 sm:pl-6 pr-12 py-2.5 sm:py-3 text-sm text-neutral-900 placeholder:text-neutral-400 focus:outline-none focus:ring-2 focus:ring-black/10 focus:border-black transition-all"
                />
                
                {/* Custom Eye Toggle matching the reference icon */}
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-4 p-1.5 text-black hover:opacity-70 transition-opacity focus:outline-none cursor-pointer"
                  aria-label={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
                >
                  <svg
                    className="w-5 h-5 text-black"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="1.8"
                  >
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                    {showPassword ? (
                      <line x1="2" y1="2" x2="22" y2="22" stroke="currentColor" strokeWidth="2" />
                    ) : (
                      <circle cx="12" cy="12" r="3.2" fill="currentColor" stroke="none" />
                    )}
                  </svg>
                </button>
              </div>
            </div>

            {/* Login Button (Black Pill matching mockup) */}
            <div className="pt-3 text-center">
              <button
                type="submit"
                disabled={isLoading}
                className="inline-flex items-center justify-center rounded-full bg-black hover:bg-neutral-800 active:scale-95 text-white text-sm font-medium px-10 sm:px-12 py-2.5 sm:py-3 min-w-[130px] transition-all cursor-pointer disabled:opacity-60 disabled:cursor-not-allowed shadow-sm"
              >
                {isLoading ? (
                  <span className="flex items-center gap-2">
                    <span className="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    <span>Masuk...</span>
                  </span>
                ) : (
                  <span>Login</span>
                )}
              </button>
            </div>

            {/* Forget Password Link */}
            <div className="text-center pt-1">
              <button
                type="button"
                onClick={() => setShowForgotModal(true)}
                className="text-xs sm:text-sm text-neutral-900 hover:opacity-75 underline underline-offset-4 cursor-pointer font-sans transition-opacity"
              >
                Forget Password?
              </button>
            </div>
          </form>

          {/* Quick Demo Login Chips (Built-in SatSet Testing Support) */}
          <div className="mt-8 pt-5 border-t border-neutral-100 text-center">
            <p className="text-[10px] sm:text-[11px] text-neutral-400 mb-2.5 font-medium tracking-wider uppercase">
              Pilih Akun Demo (SatSet POS)
            </p>
            <div className="flex flex-wrap justify-center gap-2">
              <button
                type="button"
                onClick={() => handleQuickFill('owner', 'owner123')}
                className="text-xs px-3.5 py-1 rounded-full border border-black/25 hover:border-black hover:bg-neutral-50 active:scale-95 text-neutral-800 transition-all cursor-pointer font-medium"
              >
                👑 Owner
              </button>
              <button
                type="button"
                onClick={() => handleQuickFill('kasir', 'kasir123')}
                className="text-xs px-3.5 py-1 rounded-full border border-black/25 hover:border-black hover:bg-neutral-50 active:scale-95 text-neutral-800 transition-all cursor-pointer font-medium"
              >
                💼 Kasir
              </button>
              <button
                type="button"
                onClick={() => handleQuickFill('operator', 'operator123')}
                className="text-xs px-3.5 py-1 rounded-full border border-black/25 hover:border-black hover:bg-neutral-50 active:scale-95 text-neutral-800 transition-all cursor-pointer font-medium"
              >
                ⚙️ Operator
              </button>
            </div>
          </div>

        </div>

        {/* Footer info */}
        <div className="w-full text-center py-2 text-[11px] text-neutral-400">
          <span>&copy; {new Date().getFullYear()} SatSet Laundry — Sistem Manajemen Terpadu</span>
        </div>

      </div>

      {/* ==================================================================== */}
      {/* MODAL: Lupa Password                                                 */}
      {/* ==================================================================== */}
      {showForgotModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-xs animate-fadeIn">
          <div className="w-full max-w-sm bg-white rounded-3xl p-6 shadow-2xl border border-neutral-100 text-center relative animate-scaleUp">
            <button
              type="button"
              onClick={() => setShowForgotModal(false)}
              className="absolute top-4 right-4 p-1.5 text-neutral-400 hover:text-black rounded-full transition-colors cursor-pointer"
            >
              <X className="w-5 h-5" />
            </button>

            <div className="w-12 h-12 rounded-full bg-[#F9CB43] text-black mx-auto flex items-center justify-center mb-4">
              <Shield className="w-6 h-6" />
            </div>

            <h3 className="font-serif text-xl font-bold text-neutral-900 mb-2">
              Lupa Kata Sandi?
            </h3>
            
            <p className="text-xs sm:text-sm text-neutral-600 leading-relaxed mb-6">
              Untuk menjaga integritas operasional dan keamanan data transaksi laundry, pembaruan password dilakukan langsung oleh <strong>Owner</strong> atau Administrator SatSet melalui menu Manajemen Pengguna.
            </p>

            <button
              type="button"
              onClick={() => setShowForgotModal(false)}
              className="w-full rounded-full bg-black hover:bg-neutral-800 text-white text-sm font-medium py-2.5 transition-all cursor-pointer"
            >
              Mengerti &amp; Tutup
            </button>
          </div>
        </div>
      )}

      {/* ==================================================================== */}
      {/* DRAWER: Top-Right Hamburger Menu                                     */}
      {/* ==================================================================== */}
      {showMenuDrawer && (
        <div className="fixed inset-0 z-50 flex justify-end bg-black/40 backdrop-blur-xs animate-fadeIn">
          <div className="w-full max-w-xs sm:max-w-sm bg-white h-full shadow-2xl p-6 sm:p-8 flex flex-col justify-between animate-slideLeft">
            
            <div>
              <div className="flex items-center justify-between pb-6 border-b border-neutral-100">
                <div className="flex items-center gap-2">
                  <span className="font-extrabold text-base tracking-wider text-black">SATSET</span>
                  <span className="text-[10px] uppercase font-bold tracking-widest px-2 py-0.5 rounded-full bg-[#F9CB43] text-black">
                    v1.0 Live
                  </span>
                </div>
                <button
                  type="button"
                  onClick={() => setShowMenuDrawer(false)}
                  className="p-1.5 text-neutral-400 hover:text-black rounded-full transition-colors cursor-pointer"
                >
                  <X className="w-5 h-5" />
                </button>
              </div>

              <div className="mt-6 space-y-4 text-left">
                <div className="p-4 rounded-2xl bg-[#F9CB43]/15 border border-[#F9CB43]/30">
                  <div className="flex items-center gap-2 text-xs font-semibold text-neutral-900 mb-1">
                    <Sparkles className="w-4 h-4 text-neutral-900" />
                    <span>Sistem Laundry Mandiri</span>
                  </div>
                  <p className="text-xs text-neutral-700 leading-relaxed">
                    Aplikasi kasir dan manajemen operasional laundry terintegrasi dengan pemantauan pengerjaan dan pelaporan keuangan real-time.
                  </p>
                </div>

                <div className="pt-2 space-y-3">
                  <div className="flex items-start gap-3 p-2.5 rounded-xl hover:bg-neutral-50 transition-colors">
                    <Shield className="w-4 h-4 text-neutral-700 mt-0.5 shrink-0" />
                    <div>
                      <p className="text-xs font-medium text-neutral-900">3 Hak Akses Peran</p>
                      <p className="text-[11px] text-neutral-500">Owner, Kasir POS, dan Operator Cucian.</p>
                    </div>
                  </div>

                  <div className="flex items-start gap-3 p-2.5 rounded-xl hover:bg-neutral-50 transition-colors">
                    <HelpCircle className="w-4 h-4 text-neutral-700 mt-0.5 shrink-0" />
                    <div>
                      <p className="text-xs font-medium text-neutral-900">Bantuan &amp; Panduan</p>
                      <p className="text-[11px] text-neutral-500">Panduan alur pesanan lunas dan cetak nota thermal.</p>
                    </div>
                  </div>

                  <div className="flex items-start gap-3 p-2.5 rounded-xl hover:bg-neutral-50 transition-colors">
                    <Phone className="w-4 h-4 text-neutral-700 mt-0.5 shrink-0" />
                    <div>
                      <p className="text-xs font-medium text-neutral-900">Dukungan Teknis</p>
                      <p className="text-[11px] text-neutral-500">Hubungi teknisi internal untuk kendala koneksi printer.</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div className="pt-6 border-t border-neutral-100 text-center">
              <button
                type="button"
                onClick={() => setShowMenuDrawer(false)}
                className="w-full rounded-full bg-neutral-100 hover:bg-neutral-200 text-neutral-900 text-xs font-medium py-2.5 transition-colors cursor-pointer"
              >
                Tutup Menu
              </button>
            </div>

          </div>
        </div>
      )}

    </div>
  );
};
