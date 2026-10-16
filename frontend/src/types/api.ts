// API Response Wrappers
export interface ApiResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
  error?: string;
  meta?: PaginationMeta;
}

export interface PaginationMeta {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

// User & Auth
export type UserRole = 'owner' | 'kasir' | 'operator';

export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  role: UserRole;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

// Customer
export interface Customer {
  id: number;
  name: string;
  phone: string;
  address?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

// Laundry Service
export interface LaundryService {
  id: number;
  name: string;
  category: string;
  unit: string;
  price_per_unit: number;
  estimated_duration: number; // in hours
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

// Order Status Enums
export type WorkStatus =
  | 'DITERIMA'
  | 'DICUCI'
  | 'DIKERINGKAN'
  | 'DISETRIKA'
  | 'DIKEMAS'
  | 'SIAP_DIAMBIL'
  | 'SELESAI'
  | 'DIBATALKAN';

export type PaymentStatus = 'BELUM_LUNAS' | 'LUNAS';

export type PaymentMethod = 'TUNAI' | 'TRANSFER' | 'QRIS';

export interface OrderItem {
  id: number;
  order_id: number;
  service_id: number;
  service_name: string;
  service_unit: string;
  price_per_unit: number;
  quantity: number;
  subtotal: number;
}

export interface OrderStatusHistory {
  id: number;
  order_id: number;
  from_status: WorkStatus;
  to_status: WorkStatus;
  changed_by: number;
  notes?: string;
  created_at: string;
  changer?: {
    id: number;
    full_name: string;
    username: string;
  };
}

export interface Payment {
  id: number;
  order_id: number;
  amount: number;
  payment_method: PaymentMethod;
  paid_at: string;
  received_by: number;
  notes?: string;
  is_cancelled: boolean;
  cancelled_by?: number;
  cancelled_at?: string;
  cancel_reason?: string;
  receiver?: {
    id: number;
    full_name: string;
    username: string;
  };
}

export interface Order {
  id: number;
  order_number: string;
  customer_id: number;
  created_by: number;
  order_date: string;
  estimated_completion: string;
  work_status: WorkStatus;
  payment_status: PaymentStatus;
  subtotal: number;
  discount: number;
  total: number;
  notes?: string;
  cancelled_by?: number;
  cancelled_at?: string;
  cancel_reason?: string;
  customer?: Customer;
  creator?: {
    id: number;
    full_name: string;
    username: string;
  };
  items?: OrderItem[];
  payments?: Payment[];
  status_history?: OrderStatusHistory[];
  created_at: string;
  updated_at: string;
}

// POS Order Creation Request
export interface CreateOrderItemInput {
  service_id: number;
  quantity: number;
}

export interface CreateOrderInput {
  customer_id: number;
  items: CreateOrderItemInput[];
  discount?: number;
  notes?: string;
  estimated_completion?: string;
}

// Expenses
export interface ExpenseCategory {
  id: number;
  name: string;
  description?: string;
  is_active: boolean;
  created_at?: string;
}

export interface Expense {
  id: number;
  category_id: number;
  amount: number;
  expense_date: string;
  description: string;
  vendor?: string;
  receipt_url?: string;
  recorded_by: number;
  is_cancelled: boolean;
  cancel_reason?: string;
  cancelled_by?: number;
  cancelled_at?: string;
  category?: ExpenseCategory;
  recorder?: {
    id: number;
    full_name: string;
    username: string;
  };
  created_at: string;
}

// Dashboards
export interface OwnerDashboardData {
  period: string;
  revenue?: {
    total_orders: number;
    total_revenue: number;
    total_discount: number;
  };
  cash_received: number;
  total_expenses: number;
  estimated_profit: number;
  active_orders: number;
  overdue_orders: number;
  ready_for_pickup: number;
  unpaid_orders: number;
  top_services?: Array<{
    service_id: number;
    service_name: string;
    service_unit: string;
    order_count: number;
    total_quantity: number;
    total_revenue: number;
  }>;
  daily_revenue?: Array<{
    date: string;
    order_count: number;
    revenue: number;
  }>;
  orders_by_status?: Array<{
    work_status?: string;
    status?: string;
    count: number;
  }>;
}

export interface KasirDashboardData {
  today_orders: number;
  ready_for_pickup: number;
  unpaid_orders: number;
  today_cash_received: number;
}

// Reports
export interface BusinessReportData {
  period_start: string;
  period_end: string;
  total_orders: number;
  gross_revenue: number;
  total_discount: number;
  net_revenue: number;
  total_received: number;
  unpaid_amount: number;
  total_expenses: number;
  net_profit: number;
  revenue_by_service: Array<{
    service_id: number;
    service_name: string;
    quantity_total: number;
    revenue_total: number;
  }>;
  expenses_by_category: Array<{
    category_id: number;
    category_name: string;
    amount_total: number;
  }>;
  daily_trends: Array<{
    date: string;
    order_count: number;
    revenue: number;
    cash_received: number;
    expenses: number;
  }>;
}

// AI Assistant
export interface AIMessage {
  id: number;
  role: 'user' | 'model';
  content: string;
  tool_calls_json?: string;
  created_at: string;
}

export interface AIConversation {
  id: number;
  title: string;
  messages?: AIMessage[];
  created_at: string;
  updated_at: string;
}

export interface AIChatResponse {
  conversation_id: number;
  message?: string;
  answer?: string;
  tool_calls_made?: number;
  tokens_used?: number;
  tool_calls?: Array<{
    name: string;
    arguments: Record<string, any>;
    result?: any;
  }>;
}

// Audit Log
export interface AuditLog {
  id: number;
  user_id: number;
  action: string;
  entity_type: string;
  entity_id: number;
  old_values?: Record<string, any>;
  new_values?: Record<string, any>;
  ip_address: string;
  created_at: string;
  user?: {
    id: number;
    full_name: string;
    username: string;
    role: string;
  };
}
