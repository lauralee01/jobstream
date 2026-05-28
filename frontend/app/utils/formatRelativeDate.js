/** Human-readable relative time (e.g. "2 days ago"). */
export const formatRelativeDate = (dateString) => {
  if (!dateString) return ''

  const date = new Date(dateString)
  if (Number.isNaN(date.getTime())) return ''

  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays < 0) return 'Today'
  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return 'Yesterday'
  return `${diffDays} days ago`
}
