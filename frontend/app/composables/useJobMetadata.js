/**
 * Fetch categories and platforms once per session (shared cache keys).
 * Call from page setup so sidebar does not block on its own fetches.
 */
export const useJobMetadata = () => {
  const config = useRuntimeConfig()
  const API_BASE = config.public.apiBase

  const { data: categories, pending: categoriesPending } = useFetch(
    `${API_BASE}/jobs/categories`,
    { key: 'categories', default: () => [] }
  )

  const { data: platforms, pending: platformsPending } = useFetch(
    `${API_BASE}/jobs/platforms`,
    { key: 'platforms', default: () => [] }
  )

  const platformOptions = computed(() =>
    (platforms.value || []).map((name) => ({ id: name, name }))
  )

  const metadataPending = computed(
    () => categoriesPending.value || platformsPending.value
  )

  const refreshMetadata = () => refreshNuxtData(['categories', 'platforms'])

  return {
    categories,
    platforms,
    platformOptions,
    metadataPending,
    refreshMetadata
  }
}
