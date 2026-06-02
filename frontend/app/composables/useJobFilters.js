import {
  emptyFilters,
  filtersFromQuery,
  normalizeFilters,
  queryFromFilters
} from '~/utils/filterParams'

/**
 * URL is the source of truth for applied filters.
 * `draft` is bound to form controls; `applyFilters` commits draft → URL.
 */
export const useJobFilters = () => {
  const route = useRoute()
  const router = useRouter()

  const filters = computed(() => filtersFromQuery(route.query))
  const draft = ref({ ...filters.value })

  watch(
    filters,
    (next) => {
      draft.value = { ...next }
    },
    { deep: true }
  )

  const applyFilters = async (partial = draft.value) => {
    const next = normalizeFilters({ ...draft.value, ...partial })
    draft.value = next
    await router.replace({ query: queryFromFilters(next) })
  }

  const setPage = async (page) => {
    const next = normalizeFilters({ ...filters.value, page }, { resetPage: false })
    draft.value = next
    await router.replace({ query: queryFromFilters(next) })
  }

  const clearAll = () => applyFilters(emptyFilters())

  return {
    filters,
    draft,
    applyFilters,
    setPage,
    clearAll
  }
}
