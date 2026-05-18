export const useJobs = () => {
  const config = useRuntimeConfig()

  const API_BASE =
    config.public.apiBase || 'http://localhost:8080/api/v1'

  // =========================================
  // Fetch Jobs
  // =========================================

  const fetchJobs = (params) => {

    const query = computed(() => ({
      keyword: params.value.keyword || undefined,
      location: params.value.location || undefined,
      category: params.value.category || undefined,
      platforms:
        params.value.platforms?.length
          ? params.value.platforms.join(',')
          : undefined,
      page: params.value.page || 1,
      limit: 20
    }))

    const {
      data,
      pending,
      error,
      refresh
    } = useFetch(
      `${API_BASE}/jobs`,
      {
        query,
        key: () => `jobs-${JSON.stringify(query.value)}`
      }
    )

    return {
      jobs: computed(() => data.value?.data || []),

      metadata: computed(
        () => data.value?.metadata || {}
      ),

      pending,

      error,

      refresh
    }
  }

  // =========================================
  // Sync Jobs
  // =========================================

  const syncJobs = async () => {

    const { data, error } = await useFetch(
      `${API_BASE}/jobs/sync`,
      {
        method: 'POST'
      }
    )

    return {
      data,
      error
    }
  }

  return {
    fetchJobs,
    syncJobs
  }
}