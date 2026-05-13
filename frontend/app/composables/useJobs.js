export const useJobs = () => {
  const config = useRuntimeConfig()
  // Base URL for the Go backend
  // In a real app, this would be in nuxt.config.ts or .env
  const API_BASE = 'http://localhost:8080/api/v1'

  const fetchJobs = async (params = {}) => {
    // Format parameters for the backend
    const query = {
      q: params.keyword || undefined,
      location: params.location || undefined,
      category: params.category || undefined,
      platforms: params.platforms?.join(',') || undefined,
      work_model: params.remote ? 'remote' : undefined,
      salary_min: params.salaryMin || undefined,
      page: params.page || 1,
      limit: 20
    }

    // Use Nuxt's useFetch for data fetching
    // This allows for SSR support automatically
    const { data, pending, error, refresh } = await useFetch(`${API_BASE}/jobs`, {
      query,
      key: `jobs-${JSON.stringify(query)}`, // Unique key for caching based on query
      watch: false // We manually trigger refresh or navigation
    })

    return {
      jobs: computed(() => data.value?.data || []),
      metadata: computed(() => data.value?.metadata || {}),
      pending,
      error,
      refresh
    }
  }

  // Helper for syncing jobs (triggering the backend fetcher)
  const syncJobs = async () => {
    const { error } = await useFetch(`${API_BASE}/jobs/sync`, { method: 'POST' })
    return { error }
  }

  return {
    fetchJobs,
    syncJobs
  }
}
