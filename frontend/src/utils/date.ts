/**
 * Date and time formatters for Asia/Jakarta timezone (WIB).
 */

export function formatJakartaDate(dateStr: string | null | undefined): string {
  if (!dateStr) return '-';
  try {
    const d = new Date(dateStr);
    return new Intl.DateTimeFormat('id-ID', {
      timeZone: 'Asia/Jakarta',
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    }).format(d);
  } catch {
    return dateStr;
  }
}

export function formatJakartaDateTime(dateStr: string | null | undefined): string {
  if (!dateStr) return '-';
  try {
    const d = new Date(dateStr);
    const dateFormatted = new Intl.DateTimeFormat('id-ID', {
      timeZone: 'Asia/Jakarta',
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    }).format(d);
    return `${dateFormatted} WIB`;
  } catch {
    return dateStr;
  }
}

export function formatJakartaTimeOnly(dateStr: string | null | undefined): string {
  if (!dateStr) return '-';
  try {
    const d = new Date(dateStr);
    const timeFormatted = new Intl.DateTimeFormat('id-ID', {
      timeZone: 'Asia/Jakarta',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    }).format(d);
    return `${timeFormatted} WIB`;
  } catch {
    return dateStr;
  }
}

/**
 * Returns true if date has passed current time
 */
export function isOverdue(estimatedCompletionStr: string | null | undefined): boolean {
  if (!estimatedCompletionStr) return false;
  return new Date(estimatedCompletionStr).getTime() < Date.now();
}
