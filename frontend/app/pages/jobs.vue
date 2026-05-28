<script setup>
const { fetchJobs, syncJobs } = useJobs()
const route = useRoute()
const router = useRouter()

const filtersFromQuery = (query) => ({
  keyword: query.q || '',
  location: query.location || '',
  category: query.category || '',
  platforms: query.platforms ? query.platforms.split(',') : [],
  remote: query.work_model === 'remote',
  salaryMin: query.salary_min || query.min_salary || '',
  page: parseInt(query.page) || 1
})

const queryFromFilters = (filters) => {
  const query = {}

  const keyword = filters.keyword?.trim()
  const location = filters.location?.trim()

  if (keyword) query.q = keyword
  if (location) query.location = location
  if (filters.category) query.category = filters.category
  if (filters.platforms?.length) query.platforms = filters.platforms.join(',')
  if (filters.remote) query.work_model = 'remote'
  if (filters.salaryMin) query.min_salary = filters.salaryMin
  if (filters.page > 1) query.page = String(filters.page)

  return query
}

// Filters state initialized from URL query
const searchParams = ref(filtersFromQuery(route.query))
const appliedParams = ref({ ...searchParams.value })

watch(
  () => route.query,
  (query) => {
    const nextParams = filtersFromQuery(query)
    searchParams.value = { ...nextParams }
    appliedParams.value = { ...nextParams }
  },
  { immediate: true }
)

// Computed jobs and metadata from the composable
const { jobs, metadata, pending, error, refresh } = fetchJobs(appliedParams)

// Handle search trigger
const handleSearch = (filters = searchParams.value) => {
  const nextParams = {
    keyword: filters.keyword?.trim() || '',
    location: filters.location?.trim() || '',
    category: filters.category || '',
    platforms: filters.platforms || [],
    remote: filters.remote || false,
    salaryMin: filters.salaryMin || '',
    page: 1
  }

  searchParams.value = nextParams
  appliedParams.value = nextParams

  updateUrl(nextParams)
}

// Handle pagination
const changePage = async (newPage) => {
  const nextParams = {
    ...appliedParams.value,
    page: newPage
  }

  searchParams.value = { ...nextParams }
  appliedParams.value = { ...nextParams }

  await updateUrl(nextParams)

  await refresh()
}

// Update URL to make the applied search state shareable
const updateUrl = async (filters) => {
  await router.replace({ query: queryFromFilters(filters) })
}

// When top search fields are cleared, re-apply filters so URL + API query reset.
watch(
  () => [searchParams.value.keyword, searchParams.value.location],
  ([keyword, location]) => {
    const keywordEmpty = !keyword?.trim()
    const locationEmpty = !location?.trim()
    const hadActiveSearch =
      appliedParams.value.keyword?.trim() || appliedParams.value.location?.trim()

    if (keywordEmpty && locationEmpty && hadActiveSearch) {
      handleSearch()
    }
  }
)

// Sync jobs helper
const isSyncing = ref(false)
const handleSync = async () => {
  isSyncing.value = true
  await syncJobs()
  isSyncing.value = false
  // Sync can introduce new platforms/categories; refresh cached fetches.
  await refreshNuxtData(['platforms', 'categories'])
  // Refresh the list after sync
  refresh()
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col transition-colors duration-300">
    <!-- Header -->
    <header class="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 sticky top-0 z-50 transition-colors duration-300">
      <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
        <div class="flex items-center gap-8">
          <NuxtLink to="/" class="text-2xl font-black text-gray-900 dark:text-white tracking-tight">JobStream</NuxtLink>
          <nav class="hidden md:flex items-center gap-6">
            <NuxtLink to="/jobs" class="text-gray-900 dark:text-white font-semibold border-b-2 border-gray-900 dark:border-gray-700 pb-1">Browse Jobs</NuxtLink>
            <NuxtLink to="#" class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors">Companies</NuxtLink>
            <NuxtLink to="#" class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors">Salaries</NuxtLink>
          </nav>
        </div>
        <div class="flex items-center gap-4">
          <ColorModeToggle />
          <button 
            @click="handleSync" 
            :disabled="isSyncing"
            class="inline-flex items-center gap-2 whitespace-nowrap text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 px-4 py-2 rounded-lg transition-colors"
          >
            <svg 
              xmlns="http://www.w3.org/2000/svg"
              class="h-4 w-4"
              :class="{ 'animate-spin': isSyncing }"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              />
            </svg>

            <span>
              {{ isSyncing ? 'Syncing...' : 'Sync Jobs' }}
            </span>
          </button>
          <!-- <button class="bg-gray-900 dark:bg-blue-600 text-white px-5 py-2 rounded-lg font-semibold hover:bg-gray-800 dark:hover:bg-blue-700 transition-colors">
            Sign In
          </button> -->
        </div>
      </div>
    </header>

    <main class="flex-grow max-w-7xl mx-auto px-4 py-8 w-full">
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Sidebar Filters -->
        <aside class="lg:col-span-1">
          <SearchFilters 
            v-model="searchParams"
            @search="handleSearch"
          />
        </aside>

        <!-- Main Content -->
        <div class="lg:col-span-3">
          <!-- Search Bar Top -->
          <div class="bg-white dark:bg-gray-800 p-4 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 mb-8 flex flex-col sm:flex-row gap-4 transition-colors duration-300">
            <div class="flex-grow relative">
              <span class="absolute inset-y-0 left-3 flex items-center text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </span>
              <input 
                v-model="searchParams.keyword"
                type="text" 
                placeholder="Job title, keywords, or company" 
                class="w-full pl-10 pr-4 py-2.5 bg-gray-50 dark:bg-gray-900 border border-gray-300 dark:border-gray-700 text-gray-900 dark:text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
              />
            </div>
            <div class="flex-grow relative">
              <span class="absolute inset-y-0 left-3 flex items-center text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </span>
              <input 
                v-model="searchParams.location"
                type="text" 
                placeholder="City, state, or remote" 
                class="w-full pl-10 pr-4 py-2.5 bg-gray-50 dark:bg-gray-900 border border-gray-300 dark:border-gray-700 text-gray-900 dark:text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
              />
            </div>
            <button 
              @click="handleSearch()"
              class="bg-gray-900 dark:bg-gray-700 text-white px-8 py-2.5 rounded-lg font-bold hover:bg-blue-700 transition-all active:scale-[0.98] shadow-lg shadow-blue-200 dark:shadow-none"
            >
              Search
            </button>
          </div>

          <!-- Loading/Error States -->
          <div v-if="pending" class="flex flex-col gap-4">
            <div v-for="i in 5" :key="i" class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-6 h-48 animate-pulse shadow-sm">
              <div class="flex justify-between mb-4">
                <div class="h-6 w-24 bg-gray-200 dark:bg-gray-700 rounded-full"></div>
                <div class="h-6 w-6 bg-gray-200 dark:bg-gray-700 rounded"></div>
              </div>
              <div class="h-8 w-64 bg-gray-200 dark:bg-gray-700 rounded mb-4"></div>
              <div class="h-6 w-32 bg-gray-200 dark:bg-gray-700 rounded mb-8"></div>
              <div class="h-10 w-32 bg-gray-200 dark:bg-gray-700 rounded"></div>
            </div>
          </div>

          <div v-else-if="error" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900 text-red-700 dark:text-red-400 p-8 rounded-xl text-center shadow-sm transition-colors duration-300">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto mb-4 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <h3 class="text-xl font-bold mb-2">Something went wrong</h3>
            <p class="mb-4">We couldn't load the jobs right now. Please try again later.</p>
            <button @click="handleSearch" class="bg-red-600 text-white px-6 py-2 rounded-lg font-bold hover:bg-red-700 transition-colors">
              Retry
            </button>
          </div>

          <!-- Results List -->
          <div v-else class="space-y-6">
            <div v-if="jobs?.length > 0" class="flex flex-col gap-6">
              <JobCard 
                v-for="job in jobs" 
                :key="job.id" 
                :job="job" 
              />

              <!-- Pagination -->
              <div v-if="metadata && metadata.total_pages > 1" class="flex justify-center items-center gap-4 mt-8 py-4">
                <button 
                  @click="changePage(searchParams.page - 1)"
                  :disabled="searchParams.page <= 1"
                  class="p-2 rounded-lg border border-gray-300 dark:border-gray-700 disabled:opacity-30 disabled:cursor-not-allowed bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:shadow-sm transition-all"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                  </svg>
                </button>
                
                <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                  Page {{ searchParams.page }} of {{ metadata.total_pages }}
                </span>

                <button 
                  @click="changePage(searchParams.page + 1)"
                  :disabled="searchParams.page >= metadata.total_pages"
                  class="p-2 rounded-lg border border-gray-300 dark:border-gray-700 disabled:opacity-30 disabled:cursor-not-allowed bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 hover:shadow-sm transition-all"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </button>
              </div>
            </div>

            <div v-else class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-12 text-center shadow-sm transition-colors duration-300">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto mb-4 text-gray-300 dark:text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">No jobs found</h3>
              <p class="text-gray-600 dark:text-gray-400">Try adjusting your filters or search terms to find what you're looking for.</p>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800 py-12 mt-12 transition-colors duration-300">
      <div class="max-w-7xl mx-auto px-4 text-center">
        <p class="text-gray-500 dark:text-gray-500 text-sm">© 2026 JobStream Aggregator. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.router-link-active {
  @apply text-gray-900 dark:text-white border-gray-900 dark:border-white;
}
</style>
