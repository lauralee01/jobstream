<script setup>
import { formatRelativeDate } from '~/utils/formatRelativeDate'

defineProps({
  job: {
    type: Object,
    required: true
  }
})

const platformStyles = {
  remotive: 'bg-purple-100 text-purple-800 dark:bg-purple-900/40 dark:text-purple-200',
  adzuna: 'bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-200',
  greenhouse: 'bg-green-100 text-green-800 dark:bg-green-900/40 dark:text-green-200',
  lever: 'bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-200',
  weworkremotely: 'bg-teal-100 text-teal-800 dark:bg-teal-900/40 dark:text-teal-200'
}

const getPlatformStyle = (platform) => {
  const key = platform?.toLowerCase().replace(/\s+/g, '')
  return platformStyles[key] || 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-200'
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-6 hover:shadow-md transition-all duration-300">
    <div class="flex justify-between items-start mb-4">
      <div class="flex items-center gap-2">
        <span 
          class="px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wider"
          :class="getPlatformStyle(job.platform)"
        >
          {{ job.platform }}
        </span>
        <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatRelativeDate(job.posted_at) }}</span>
      </div>
    </div>
 
    <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-1 transition-colors">
      {{ job.title }}
    </h3>
    <p class="text-lg text-gray-600 dark:text-gray-300 font-semibold mb-3">
      {{ job.company }}
    </p>
 
    <div class="flex flex-wrap gap-4 text-sm text-gray-600 dark:text-gray-400 mb-6 transition-colors">
      <div class="flex items-center gap-1">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        {{ job.location }}
      </div>
      <div v-if="job.salary" class="flex items-center gap-1">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {{ job.salary }}
      </div>
    </div>
 
    <div class="flex justify-between items-center mt-auto">
      <a 
        :href="job.url" 
        target="_blank"
        class="bg-gray-900 dark:bg-gray-700 text-white px-5 py-2.5 rounded-lg font-semibold hover:bg-gray-800 dark:hover:bg-gray-600 transition-all inline-flex items-center gap-2 shadow-lg shadow-gray-200 dark:shadow-none"
      >
        View Details
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
        </svg>
      </a>
    </div>
  </div>
</template>
