<script setup>
definePageMeta({ layout: 'default' })

const { platforms } = useJobMetadata()

const platformNames = computed(() =>
  (platforms.value || []).join(', ') || 'top job boards'
)
</script>

<template>
  <header class="bg-white dark:bg-gray-900 border-b border-gray-100 dark:border-gray-800">
    <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
      <NuxtLink to="/" class="text-2xl font-black text-gray-900 dark:text-white tracking-tight">
        JobStream
      </NuxtLink>
      <div class="flex items-center gap-4">
        <ClientOnly>
          <ColorModeToggle />
        </ClientOnly>
        <NuxtLink
          to="/jobs"
          class="bg-gray-900 dark:bg-gray-700 text-white px-6 py-2 rounded-lg font-bold hover:bg-blue-700 transition-all"
        >
          Find Jobs
        </NuxtLink>
      </div>
    </div>
  </header>

  <main class="max-w-4xl mx-auto px-4 py-24 text-center">
    <h1 class="text-5xl md:text-7xl font-black text-gray-900 dark:text-white mb-6 tracking-tight leading-tight">
      Your Next Career Move, <br>
      <span class="text-blue-600">All in One Stream.</span>
    </h1>
    <p class="text-xl text-gray-600 dark:text-gray-400 mb-10 max-w-2xl mx-auto">
      Aggregate job listings from {{ platformNames }}.
      Filter, search, and apply from a single, fast interface.
    </p>

    <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
      <NuxtLink
        to="/jobs"
        class="w-full sm:w-auto bg-gray-900 dark:bg-blue-600 text-white text-lg px-10 py-4 rounded-xl font-bold hover:bg-gray-800 dark:hover:bg-blue-700 transition-all shadow-xl"
      >
        Browse All Jobs
      </NuxtLink>
      <NuxtLink
        to="/jobs?work_model=remote"
        class="w-full sm:w-auto bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-lg px-10 py-4 rounded-xl font-bold border-2 border-gray-100 dark:border-gray-700 hover:border-gray-200 transition-all"
      >
        Remote Only
      </NuxtLink>
    </div>
  </main>
</template>
