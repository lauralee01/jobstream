import {
  emptyFilters,
  filtersFromQuery,
  normalizeFilters,
  queryFromFilters
} from '~/utils/filterParams'

/**
 * URL is the source of truth for applied filters.
 * `draft` is bound to form controls; `applyFilters` commits draft → URL.
 * * DEBOUNCED: Keyword and location searches are debounced to prevent excessive API calls.
 */
export const useJobFilters = () => {
  const route = useRoute()
  const router = useRouter()

  const filters = computed(() => filtersFromQuery(route.query))
  const draft = ref({ ...filters.value })

  // Debounce timer for search inputs
  let debounceTimer = null
  const DEBOUNCE_DELAY = 500 // Wait 500ms after user stops typing

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


  // DEBOUNCED: Used for keyword and location inputs
  // Prevents sending an API request on every keystroke
  const applyFiltersDebounced = async (partial = draft.value) => {
    // Clear any existing timer
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }

    // Set new timer - only apply filters after user stops typing
    debounceTimer = setTimeout(async () => {
      await applyFilters(partial)
    }, DEBOUNCE_DELAY)
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
    applyFiltersDebounced,
    setPage,
    clearAll
  }
}
