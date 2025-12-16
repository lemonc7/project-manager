export function toRFC3339(date?: string) {
  if (!date) return undefined
  return new Date(date).toISOString()
}
