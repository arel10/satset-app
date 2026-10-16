/**
 * Format integer amount to Indonesian Rupiah.
 * Example: 15000 -> "Rp 15.000"
 */
export function formatRupiah(amount: number | null | undefined): string {
  if (amount == null || isNaN(amount)) {
    return 'Rp 0';
  }
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
}

/**
 * Format decimal quantity (e.g. 3.5 kg or 5 pcs)
 */
export function formatQuantity(qty: number, unit: string = ''): string {
  const formatted = new Intl.NumberFormat('id-ID', {
    maximumFractionDigits: 2,
  }).format(qty);
  return unit ? `${formatted} ${unit}` : formatted;
}
