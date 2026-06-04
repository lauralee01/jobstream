<script setup>
definePageMeta({ layout: 'default' })

const { fetchJobs, syncJobs } = useJobs()
const { filters, draft, applyFilters, setPage, clearAll } = useJobFilters()
const { categories, platformOptions, refreshMetadata } = useJobMetadata()

const { jobs, metadata, pending, error, refresh } = fetchJobs(filters)

const isSyncing = ref(false)

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
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
      <aside class="lg:col-span-1">
        <SearchFilters
          v-model="draft"
          :categories="categories"
          :platforms="platformOptions"
          @search="applyFilters"
          @clear="clearAll"
        />
      </aside>

      <div class="lg:col-span-3">
        <JobsSearchBar v-model="draft" @search="applyFilters" />

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
