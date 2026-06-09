<script setup>
definePageMeta({ layout: 'default' })

const { fetchJobs, syncJobs } = useJobs()
const { filters, draft, applyFilters, setPage, clearAll } = useJobFilters()
const { categories, platformOptions, refreshMetadata } = useJobMetadata()

const { jobs, metadata, pending, error, refresh } = fetchJobs(filters)

const isSyncing = ref(false)

const showSearch = ref(false)
const showFilters = ref(false)

const hasActiveSearch = computed(() => {
  return !!(draft.value.keyword?.trim() || draft.value.location?.trim())
})

const activeFiltersCount = computed(() => {
  let count = 0
  if (draft.value.category) count++
  if (draft.value.remote) count++
  if (draft.value.salaryMin) count++
  if (draft.value.platforms?.length) count += draft.value.platforms.length
  return count
})

const handleSearch = () => {
  applyFilters()
  showSearch.value = false
}

const handleApplyFilters = () => {
  applyFilters()
  showFilters.value = false
}

const handleSync = async () => {
  isSyncing.value = true
  try {
    await syncJobs()
    await refreshMetadata()
    await refreshNuxtData(['jobs-list'])
  } finally {
    isSyncing.value = false
  }
}
</script>

<template>
  <JobsHeader :is-syncing="isSyncing" @sync="handleSync" />

  <main class="flex-grow max-w-7xl mx-auto px-4 py-8 w-full">
    <div class="lg:hidden flex gap-4 mb-6">
      <button
        type="button"
        @click="showSearch = !showSearch"
        class="flex-1 inline-flex items-center justify-center gap-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 py-3 px-4 rounded-xl text-gray-700 dark:text-gray-200 font-semibold shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 transition-all focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
        :class="{ 'border-blue-500 ring-2 ring-blue-500 ring-offset-2 dark:ring-offset-gray-900': showSearch }"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" :class="showSearch ? 'text-blue-500' : 'text-gray-400 dark:text-gray-500'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <span>Search</span>
        <span v-if="hasActiveSearch" class="w-2 h-2 rounded-full bg-blue-500"></span>
      </button>

      <button
        type="button"
        @click="showFilters = !showFilters"
        class="flex-1 inline-flex items-center justify-center gap-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 py-3 px-4 rounded-xl text-gray-700 dark:text-gray-200 font-semibold shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 transition-all focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
        :class="{ 'border-blue-500 ring-2 ring-blue-500 ring-offset-2 dark:ring-offset-gray-900': showFilters }"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" :class="showFilters ? 'text-blue-500' : 'text-gray-400 dark:text-gray-500'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 8.293A1 1 0 013 7.586V4z" />
        </svg>
        <span>Filters</span>
        <span v-if="activeFiltersCount > 0" class="inline-flex items-center justify-center px-2 py-0.5 text-xs font-bold leading-none text-white bg-blue-500 rounded-full">
          {{ activeFiltersCount }}
        </span>
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
      <aside 
        class="lg:col-span-1"
        :class="{ 'hidden lg:block': !showFilters, 'block': showFilters }"
      >
        <SearchFilters
          v-model="draft"
          :categories="categories"
          :platforms="platformOptions"
          @search="handleApplyFilters"
          @clear="clearAll"
        />
      </aside>

      <div class="lg:col-span-3">
        <div :class="{ 'hidden lg:block': !showSearch, 'block': showSearch }">
          <JobsSearchBar v-model="draft" @search="handleSearch" />
        </div>

        <JobsResults
          :jobs="jobs"
          :metadata="metadata"
          :page="filters.page"
          :pending="pending"
          :error="error"
          @retry="applyFilters"
          @prev-page="setPage(filters.page - 1)"
          @next-page="setPage(filters.page + 1)"
        />
      </div>
    </div>
  </main>

  <footer class="bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800 py-12 mt-12">
    <div class="max-w-7xl mx-auto px-4 text-center">
      <p class="text-gray-500 text-sm">© 2026 JobStream Aggregator. All rights reserved.</p>
    </div>
  </footer>
</template>
