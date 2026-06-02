<script setup>
defineProps({
  isSyncing: {
    type: Boolean,
    default: false
  }
})

defineEmits(['sync'])
</script>

<template>
  <header class="bg-white dark:bg-gray-900 border-b border-gray-200 dark:border-gray-800 sticky top-0 z-50 transition-colors duration-300">
    <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
      <div class="flex items-center gap-8">
        <NuxtLink to="/" class="text-2xl font-black text-gray-900 dark:text-white tracking-tight">
          JobStream
        </NuxtLink>
        <nav class="hidden md:flex items-center gap-6">
          <NuxtLink
            to="/jobs"
            class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
            active-class="!text-gray-900 dark:!text-white font-semibold border-b-2 border-gray-900 dark:border-gray-700 pb-1"
          >
            Browse Jobs
          </NuxtLink>
          <NuxtLink to="#" class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors">
            Companies
          </NuxtLink>
          <!-- <NuxtLink to="#" class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors">
            Salaries
          </NuxtLink> -->
        </nav>
      </div>
      <div class="flex items-center gap-4">
        <ClientOnly>
          <ColorModeToggle />
        </ClientOnly>
        <button
          type="button"
          :disabled="isSyncing"
          class="inline-flex items-center gap-2 whitespace-nowrap text-sm font-medium text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 px-4 py-2 rounded-lg transition-colors disabled:opacity-60"
          @click="$emit('sync')"
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
          <span>{{ isSyncing ? 'Syncing...' : 'Sync Jobs' }}</span>
        </button>
      </div>
    </div>
  </header>
</template>
