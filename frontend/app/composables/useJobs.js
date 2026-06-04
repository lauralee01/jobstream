import { buildJobsApiQuery } from '~/utils/filterParams'

export const useJobs = () => {
  const config = useRuntimeConfig()
  const API_BASE = config.public.apiBase

  const fetchJobs = (filtersRef) => {
    const query = computed(() => buildJobsApiQuery(filtersRef.value))

    const { data, pending, error, refresh } = useFetch(`${API_BASE}/jobs`, {
      query,
      watch: [query],
      default: () => ({ data: [], metadata: {} })
    })

    return {
      jobs: computed(() => data.value?.data || []),
      metadata: computed(() => data.value?.metadata || {}),
      pending,
      error,
      refresh
    }
  }

  /** Imperative POST — use $fetch, not useFetch. */
  const syncJobs = () =>
    $fetch(`${API_BASE}/jobs/sync`, { method: 'POST' })

  return {
    fetchJobs,
    syncJobs
  }
}
