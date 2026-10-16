import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from '@/context/AuthContext';
import { ProtectedRoute } from '@/components/layout/ProtectedRoute';
import { AppLayout } from '@/components/layout/AppLayout';

// Pages
import { LoginPage } from '@/pages/LoginPage';
import { DashboardPage } from '@/pages/DashboardPage';
import { POSPage } from '@/pages/orders/POSPage';
import { OrderListPage } from '@/pages/orders/OrderListPage';
import { OrderDetailPage } from '@/pages/orders/OrderDetailPage';
import { CustomerListPage } from '@/pages/customers/CustomerListPage';
import { ServiceListPage } from '@/pages/services/ServiceListPage';
import { ExpenseListPage } from '@/pages/expenses/ExpenseListPage';
import { ReportPage } from '@/pages/reports/ReportPage';
import { AIAssistantPage } from '@/pages/ai/AIAssistantPage';
import { SettingsPage } from '@/pages/settings/SettingsPage';
import { UsersPage } from '@/pages/settings/UsersPage';
import { AuditLogsPage } from '@/pages/settings/AuditLogsPage';

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          {/* Public Login Route */}
          <Route path="/login" element={<LoginPage />} />

          {/* Protected Routes inside AppLayout */}
          <Route element={<ProtectedRoute />}>
            <Route element={<AppLayout />}>
              {/* Default redirect to Dashboard */}
              <Route path="/" element={<Navigate to="/dashboard" replace />} />

              {/* Dashboard (Owner & Kasir) */}
              <Route
                path="/dashboard"
                element={<DashboardPage />}
              />

              {/* POS / Buat Pesanan (Owner & Kasir) */}
              <Route
                element={<ProtectedRoute allowedRoles={['owner', 'kasir']} />}
              >
                <Route path="/orders/new" element={<POSPage />} />
              </Route>

              {/* Orders Management (Owner, Kasir, Operator) */}
              <Route path="/orders" element={<OrderListPage />} />
              <Route path="/orders/:id" element={<OrderDetailPage />} />

              {/* Customer Directory (Owner & Kasir) */}
              <Route
                element={<ProtectedRoute allowedRoles={['owner', 'kasir']} />}
              >
                <Route path="/customers" element={<CustomerListPage />} />
              </Route>

              {/* Owner-Only Routes */}
              <Route element={<ProtectedRoute allowedRoles={['owner']} />}>
                <Route path="/services" element={<ServiceListPage />} />
                <Route path="/expenses" element={<ExpenseListPage />} />
                <Route path="/reports" element={<ReportPage />} />
                <Route path="/ai" element={<AIAssistantPage />} />
                <Route path="/settings" element={<SettingsPage />} />
                <Route path="/settings/users" element={<UsersPage />} />
                <Route path="/settings/audit-logs" element={<AuditLogsPage />} />
              </Route>
            </Route>
          </Route>

          {/* Catch-all */}
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;
