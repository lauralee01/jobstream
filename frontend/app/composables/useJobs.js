export const useJobs = () => {
  const config = useRuntimeConfig()

  const API_BASE = config.public.apiBase

  // =========================================
  // Fetch Jobs
  // =========================================

  const fetchJobs = (params) => {

    const query = computed(() => ({
      keyword: params.value.keyword || undefined,
      location: params.value.location || undefined,
      category: params.value.category || undefined,
      // Primary key used by the current filters UI is `salaryMin`.
      // Keep `minSalary` as fallback for backwards compatibility.
      min_salary: params.value.salaryMin || params.value.minSalary || undefined,
      platforms:
        params.value.platforms?.length
          ? params.value.platforms.join(',')
          : undefined,
      remote: params.value.remote ? 'true' : undefined,
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

  // =========================================
  // Fetch Categories
  // =========================================

  const fetchCategories = () => {
    return useFetch(`${API_BASE}/jobs/categories`, {
      key: 'categories'
    })
  }

  // =========================================
  // Fetch Platforms
  // =========================================

  const fetchPlatforms = () => {
    return useFetch(`${API_BASE}/jobs/platforms`, {
      key: 'platforms'
    })
  }

  return {
    fetchJobs,
    syncJobs,
    fetchCategories,
    fetchPlatforms
  }
}