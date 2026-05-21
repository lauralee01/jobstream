<script setup>
const { fetchPlatforms } = useJobs()
const { data: fetchedPlatforms } = await fetchPlatforms()
const platforms = computed(() => {
  if(!fetchedPlatforms.value) {
    return []
  }
  return fetchedPlatforms.value.map(p => ({
    name: p
  }))
})

</script>

<template>
  <div class="min-h-screen bg-white dark:bg-gray-900 transition-colors duration-300">
    <!-- Simple Hero Section -->
    <header class="bg-white dark:bg-gray-900 border-b border-gray-100 dark:border-gray-800 transition-colors duration-300">
      <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
        <NuxtLink to="/" class="text-2xl font-black text-blue-600 tracking-tight">JobStream</NuxtLink>
        <div class="flex items-center gap-4">
          <ColorModeToggle />
          <NuxtLink 
            to="/jobs" 
            class="bg-blue-600 text-white px-6 py-2 rounded-lg font-bold hover:bg-blue-700 transition-all shadow-lg shadow-blue-200 dark:shadow-none"
          >
            Find Jobs
          </NuxtLink>
        </div>
      </div>
    </header>

    <main class="max-w-4xl mx-auto px-4 py-24 text-center">
      <h1 class="text-5xl md:text-7xl font-black text-gray-900 dark:text-white mb-6 tracking-tight leading-tight">
        Your Next Career Move, <br/> 
        <span class="text-blue-600">All in One Stream.</span>
      </h1>
      <p class="text-xl text-gray-600 dark:text-gray-400 mb-10 max-w-2xl mx-auto">
        Aggregate job listings from {{ platforms.map(p => p.name).join(', ') }}. 
        Filter, search, and apply from a single, fast interface.
      </p>
      
      <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
        <NuxtLink 
          to="/jobs" 
          class="w-full sm:w-auto bg-gray-900 dark:bg-blue-600 text-white text-lg px-10 py-4 rounded-xl font-bold hover:bg-gray-800 dark:hover:bg-blue-700 transition-all shadow-xl shadow-gray-200 dark:shadow-none"
        >
          Browse All Jobs
        </NuxtLink>
        <NuxtLink 
          to="/jobs?work_model=remote" 
          class="w-full sm:w-auto bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-lg px-10 py-4 rounded-xl font-bold border-2 border-gray-100 dark:border-gray-700 hover:border-gray-200 dark:hover:border-gray-600 transition-all"
        >
          Remote Only
        </NuxtLink>
      </div>
    </main>
  </div>
</template>
