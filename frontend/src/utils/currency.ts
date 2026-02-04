const formatterCache = new Map<string, Intl.NumberFormat>();

function getFormatter(mnemonic: string): Intl.NumberFormat {
  let formatter = formatterCache.get(mnemonic);
  if (!formatter) {
    try {
      formatter = new Intl.NumberFormat(undefined, {
        style: 'currency',
        currency: mnemonic,
      });
    } catch {
      // Fallback for unrecognized currency codes
      formatter = new Intl.NumberFormat(undefined, {
        style: 'decimal',
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    }
    formatterCache.set(mnemonic, formatter);
  }
  return formatter;
}

export function formatCurrency(amount: string | number, mnemonic?: string): string {
  const num = typeof amount === 'string' ? parseFloat(amount) : amount;
  if (!mnemonic) {
    return num.toFixed(2);
  }
  return getFormatter(mnemonic).format(num);
}
