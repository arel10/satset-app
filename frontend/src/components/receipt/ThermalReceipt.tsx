import React from 'react';
import { Order } from '@/types/api';
import { formatRupiah, formatQuantity } from '@/utils/currency';
import { formatJakartaDate, formatJakartaDateTime } from '@/utils/date';

export interface ThermalReceiptProps {
  order: Order;
  businessSettings?: {
    business_name?: string;
    business_address?: string;
    business_phone?: string;
    receipt_footer?: string;
    business_logo_url?: string;
  };
}

export const ThermalReceipt: React.FC<ThermalReceiptProps> = ({ order, businessSettings }) => {
  const isPaid = order.payment_status === 'LUNAS';

  return (
    <div id="thermal-receipt-area" className="p-4 bg-white text-black font-mono text-xs max-w-[320px] mx-auto border border-dashed border-slate-300">
      {/* Business Header */}
      <div className="text-center pb-2 border-b border-black">
        {businessSettings?.business_logo_url && (
          <img
            src={businessSettings.business_logo_url}
            alt="Logo"
            className="w-12 h-12 object-contain mx-auto mb-1"
          />
        )}
        <h2 className="font-bold text-sm tracking-wider uppercase">
          {businessSettings?.business_name || 'SATSET LAUNDRY'}
        </h2>
        {businessSettings?.business_address && (
          <p className="text-[10px] mt-0.5">{businessSettings.business_address}</p>
        )}
        {businessSettings?.business_phone && (
          <p className="text-[10px]">Telp/WA: {businessSettings.business_phone}</p>
        )}
      </div>

      {/* Order & Customer Details */}
      <div className="py-2 border-b border-black text-[11px] space-y-1">
        <div className="flex justify-between">
          <span>No. Nota:</span>
          <span className="font-bold">{order.order_number}</span>
        </div>
        <div className="flex justify-between">
          <span>Masuk:</span>
          <span>{formatJakartaDateTime(order.order_date || order.created_at)}</span>
        </div>
        <div className="flex justify-between">
          <span>Estimasi:</span>
          <span>{formatJakartaDateTime(order.estimated_completion)}</span>
        </div>
        <div className="flex justify-between pt-1 border-t border-dashed border-slate-400">
          <span>Pelanggan:</span>
          <span className="font-semibold">{order.customer?.name || 'Umum'}</span>
        </div>
        {order.customer?.phone && (
          <div className="flex justify-between text-[10px]">
            <span>Telp:</span>
            <span>{order.customer.phone}</span>
          </div>
        )}
      </div>

      {/* Items List */}
      <div className="py-2 border-b border-black">
        <div className="text-[11px] font-bold pb-1 flex justify-between border-b border-dashed border-slate-400">
          <span>Layanan</span>
          <span>Total</span>
        </div>
        <div className="pt-1 space-y-1.5 text-[11px]">
          {order.items?.map((item, idx) => (
            <div key={idx} className="flex flex-col">
              <span className="font-semibold">{item.service_name}</span>
              <div className="flex justify-between text-[10px] text-slate-600">
                <span>
                  {formatQuantity(item.quantity, item.service_unit)} x {formatRupiah(item.price_per_unit)}
                </span>
                <span className="font-medium text-black">{formatRupiah(item.subtotal)}</span>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Totals */}
      <div className="py-2 border-b border-black text-[11px] space-y-1">
        <div className="flex justify-between">
          <span>Subtotal:</span>
          <span>{formatRupiah(order.subtotal)}</span>
        </div>
        {order.discount > 0 && (
          <div className="flex justify-between text-rose-700">
            <span>Diskon:</span>
            <span>-{formatRupiah(order.discount)}</span>
          </div>
        )}
        <div className="flex justify-between font-bold text-xs pt-1 border-t border-dashed border-slate-400">
          <span>TOTAL:</span>
          <span>{formatRupiah(order.total)}</span>
        </div>
      </div>

      {/* Payment Status Stamp */}
      <div className="py-2 text-center border-b border-black">
        <span
          className={`inline-block px-3 py-1 text-xs font-bold border-2 rounded ${
            isPaid ? 'border-black text-black' : 'border-black text-black'
          }`}
        >
          {isPaid ? '*** LUNAS ***' : '*** BELUM LUNAS ***'}
        </span>
        {isPaid && order.payments && order.payments.length > 0 && (
          <p className="text-[10px] mt-1 text-slate-700">
            Metode: {order.payments[0].payment_method} | Kasir: {order.creator?.full_name || 'Staff'}
          </p>
        )}
      </div>

      {/* Footer Notes */}
      <div className="pt-2 text-center text-[9px] text-slate-600 space-y-0.5">
        <p>{businessSettings?.receipt_footer || 'Terima kasih atas kepercayaan Anda.'}</p>
        <p>Cucian tidak diambil &gt; 30 hari di luar tanggung jawab toko.</p>
      </div>
    </div>
  );
};
