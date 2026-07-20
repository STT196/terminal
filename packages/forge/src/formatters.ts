/**
 * Format currency from cents to dollars
 */
export function formatCurrency(cents: number | string | null | undefined): string {
  const amount = typeof cents === "string" ? parseInt(cents) : cents || 0;
  return `$${(amount / 100).toFixed(2)}`;
}
  if (!address) return "N/A";
  const parts = [
    address.name,
    address.street1,
    address.street2,
    `${address.city}, ${address.province} ${address.zip}`,
    address.country,
  ].filter(Boolean);
  return parts.join("\n");
}

/**
 * Format shipping address inline (name, city, province)
 */
export function formatShippingAddressInline(address: any): string {
  if (!address) return "N/A";
  const parts = [
    address.name,
    address.city,
    address.province,
  ].filter(Boolean);
  return parts.length > 0 ? parts.join(", ") : "N/A";
}

/**
 * Format card last4 with masking
 */
export function formatCardLast4(last4: string | null | undefined): string {
  return last4 ? `****${last4}` : "N/A";
}

/**
 * Format card expiration date
 */
export function formatCardExpiration(
  month: number | null | undefined,
  year: number | null | undefined,
): string {
  if (!month || !year) return "N/A";
  return `${month}/${year}`;
}

/**
 * Format boolean or truthy value as Yes/No
 */
export function formatBoolean(value: boolean | null | undefined): string {
  return value ? "Yes" : "No";
}

