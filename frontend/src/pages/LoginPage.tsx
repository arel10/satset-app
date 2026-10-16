import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { Lock, User as UserIcon, AlertCircle, Sparkles } from 'lucide-react';

export const LoginPage: React.FC = () => {
  const { login, isAuthenticated, user } = useAuth();
  const navigate = useNavigate();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  // If already logged in, redirect
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
      // AuthContext will update user state and trigger useEffect redirect
    } catch (err: any) {
      setError(err.message || 'Login gagal. Periksa username dan password Anda.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-teal-50/30 to-slate-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
      {/* Decorative ambient background blur */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-96 h-96 bg-teal-400/10 rounded-full blur-3xl pointer-events-none" />

      <div className="sm:mx-auto sm:w-full sm:max-w-md text-center relative z-10">
        <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[11px] font-semibold bg-teal-100/80 text-teal-800 border border-teal-200 mb-3 shadow-xs">
          <Sparkles className="w-3.5 h-3.5 text-teal-600" />
          <span>Cloud Edition • Live VPS</span>
        </div>

        <div className="w-14 h-14 rounded-2xl bg-gradient-to-tr from-teal-700 via-teal-600 to-emerald-500 text-white flex items-center justify-center font-bold text-2xl mx-auto shadow-lg shadow-teal-700/25 ring-4 ring-white">
          S
        </div>
        <h2 className="mt-4 text-2xl font-bold tracking-tight text-slate-900">
          Masuk ke SatSet
        </h2>
        <p className="mt-1 text-xs text-slate-500">
          Sistem Operasional Laundry &amp; POS Terpadu
        </p>
      </div>

      <div className="mt-7 sm:mx-auto sm:w-full sm:max-w-md px-4 relative z-10">
        <div className="bg-white/95 backdrop-blur-sm py-8 px-6 shadow-xl shadow-slate-200/60 border border-slate-200/90 rounded-2xl sm:px-10">
          <form className="space-y-4" onSubmit={handleSubmit}>
            {error && (
              <div className="rounded-lg bg-rose-50 border border-rose-200 p-3.5 flex items-start gap-2.5 text-xs text-rose-800">
                <AlertCircle className="w-4 h-4 text-rose-600 shrink-0 mt-0.5" />
                <span>{error}</span>
              </div>
            )}

            <Input
              label="Nama Pengguna (Username)"
              type="text"
              required
              autoFocus
              value={username}
              onChange={e => setUsername(e.target.value)}
              placeholder="Contoh: owner atau kasir"
              leftIcon={<UserIcon className="w-4 h-4" />}
            />

            <Input
              label="Kata Sandi"
              type="password"
              required
              value={password}
              onChange={e => setPassword(e.target.value)}
              placeholder="••••••••"
              leftIcon={<Lock className="w-4 h-4" />}
            />

            <div className="pt-2">
              <Button
                type="submit"
                variant="primary"
                size="md"
                className="w-full"
                isLoading={isLoading}
              >
                Masuk
              </Button>
            </div>
          </form>

          <div className="mt-6 border-t border-slate-100 pt-4 text-center">
            <p className="text-[11px] text-slate-400">
              SatSet Laundry Management System • Akun default: owner / owner123
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};
