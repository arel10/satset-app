import { Button } from "@/components/ui/button";
import { Sparkles, Activity, Layers, ShieldCheck, CheckCircle2 } from "lucide-react";

export default function Home() {
  return (
    <main className="min-h-screen bg-slate-50 text-slate-900 flex flex-col">
      {/* Header */}
      <header className="border-b bg-white/80 backdrop-blur sticky top-0 z-50">
        <div className="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="h-9 w-9 rounded-lg bg-blue-600 flex items-center justify-center text-white font-bold text-lg shadow-sm">
              S
            </div>
            <span className="font-bold text-xl tracking-tight text-slate-900">
              SatSet<span className="text-blue-600">.</span>
            </span>
            <span className="ml-2 px-2 py-0.5 text-xs font-medium rounded-full bg-blue-50 text-blue-700 border border-blue-200">
              v1.0 Scaffold
            </span>
          </div>

          <nav className="flex items-center gap-4">
            <a
              href="/api/healthz"
              target="_blank"
              rel="noreferrer"
              className="text-xs text-slate-600 hover:text-slate-900 flex items-center gap-1.5 bg-slate-100 hover:bg-slate-200 px-3 py-1.5 rounded-md transition-colors"
            >
              <Activity className="size-3.5 text-emerald-600" />
              API Healthz
            </a>
            <Button size="sm" variant="default" className="bg-blue-600 hover:bg-blue-700 text-white">
              Buka Dashboard
            </Button>
          </nav>
        </div>
      </header>

      {/* Hero Section */}
      <section className="flex-1 max-w-5xl mx-auto px-4 py-16 sm:py-24 flex flex-col items-center text-center justify-center">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-50 text-blue-700 border border-blue-200 text-xs font-medium mb-6">
          <Sparkles className="size-3.5" />
          Platform Manajemen Laundry & AI Business Assistant
        </div>

        <h1 className="text-4xl sm:text-6xl font-extrabold tracking-tight text-slate-950 max-w-3xl leading-tight">
          Operasional Laundry Cepat, Rapi, & Berbasis Data
        </h1>

        <p className="mt-6 text-lg sm:text-xl text-slate-600 max-w-2xl leading-relaxed">
          Fondasi monorepo SatSet telah aktif. Siap mendukung POS kasir terintegrasi, pelacakan proses produksi, dan asisten bisnis pintar.
        </p>

        {/* Feature Cards Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 mt-14 w-full text-left">
          <div className="p-6 bg-white rounded-xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
            <div className="h-10 w-10 rounded-lg bg-blue-100 text-blue-600 flex items-center justify-center mb-4">
              <Layers className="size-5" />
            </div>
            <h3 className="font-semibold text-slate-900 text-base">Modular Monolith</h3>
            <p className="text-sm text-slate-600 mt-2">
              Backend Go + Gin dengan pemisahan handler, service, dan repository yang ketat serta aman multi-tenant.
            </p>
          </div>

          <div className="p-6 bg-white rounded-xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
            <div className="h-10 w-10 rounded-lg bg-emerald-100 text-emerald-600 flex items-center justify-center mb-4">
              <ShieldCheck className="size-5" />
            </div>
            <h3 className="font-semibold text-slate-900 text-base">Integritas Data</h3>
            <p className="text-sm text-slate-600 mt-2">
              MySQL 8 dengan transaksi database atomik, BIGINT rupiah, UTC timezone, dan log audit otomatis.
            </p>
          </div>

          <div className="p-6 bg-white rounded-xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow">
            <div className="h-10 w-10 rounded-lg bg-purple-100 text-purple-600 flex items-center justify-center mb-4">
              <Sparkles className="size-5" />
            </div>
            <h3 className="font-semibold text-slate-900 text-base">AI Assistant Ready</h3>
            <p className="text-sm text-slate-600 mt-2">
              Arsitektur adapter Gemini API yang read-only, terlindungi di sisi backend, dan decoupled dari frontend.
            </p>
          </div>
        </div>

        {/* Status Check List */}
        <div className="mt-12 w-full max-w-xl p-5 bg-white border border-slate-200 rounded-lg text-left shadow-sm">
          <h4 className="text-xs font-semibold uppercase tracking-wider text-slate-500 mb-3">
            Pemeriksaan Lingkungan Awal (Sprint 0)
          </h4>
          <ul className="space-y-2 text-sm">
            <li className="flex items-center gap-2 text-slate-700">
              <CheckCircle2 className="size-4 text-emerald-600" />
              <span>Next.js 15 (App Router) + Tailwind CSS + shadcn/ui</span>
            </li>
            <li className="flex items-center gap-2 text-slate-700">
              <CheckCircle2 className="size-4 text-emerald-600" />
              <span>Go API Engine + Gin + Zerolog logger</span>
            </li>
            <li className="flex items-center gap-2 text-slate-700">
              <CheckCircle2 className="size-4 text-emerald-600" />
              <span>Health check endpoint (/healthz &amp; /readyz)</span>
            </li>
            <li className="flex items-center gap-2 text-slate-700">
              <CheckCircle2 className="size-4 text-emerald-600" />
              <span>Multi-stage Dockerfile &amp; docker-compose stack</span>
            </li>
          </ul>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t bg-white py-6 text-center text-xs text-slate-500">
        SatSet Laundry Management System &copy; 2026. All rights reserved.
      </footer>
    </main>
  );
}
