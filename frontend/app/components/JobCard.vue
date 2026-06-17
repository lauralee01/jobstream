<script setup>
import { computed } from 'vue'
import { formatRelativeDate } from '~/utils/formatRelativeDate'

const props = defineProps({
  job: {
    type: Object,
    required: true
  }
})

// Base styles for each color theme
const themeColors = {
  purple: {
    badge: 'bg-purple-100 text-purple-800 dark:bg-purple-900/40 dark:text-purple-200',
    card: 'hover:border-purple-300 dark:hover:border-purple-800/80 hover:shadow-purple-500/5 dark:hover:shadow-purple-500/10',
    title: 'group-hover:text-purple-600 dark:group-hover:text-purple-400',
    button: 'group-hover:bg-purple-600 dark:group-hover:bg-purple-500 hover:!bg-purple-700 dark:hover:!bg-purple-600 group-hover:shadow-purple-100 dark:group-hover:shadow-none',
    topBar: 'bg-purple-500 dark:bg-purple-400'
  },
  amber: {
    badge: 'bg-amber-100 text-amber-800 dark:bg-amber-900/40 dark:text-amber-200',
    card: 'hover:border-amber-300 dark:hover:border-amber-800/80 hover:shadow-amber-500/5 dark:hover:shadow-amber-500/10',
    title: 'group-hover:text-amber-600 dark:group-hover:text-amber-400',
    button: 'group-hover:bg-amber-600 dark:group-hover:bg-amber-500 hover:!bg-amber-700 dark:hover:!bg-amber-600 group-hover:shadow-amber-100 dark:group-hover:shadow-none',
    topBar: 'bg-amber-500 dark:bg-amber-400'
  },
  green: {
    badge: 'bg-green-100 text-green-800 dark:bg-green-900/40 dark:text-green-200',
    card: 'hover:border-green-300 dark:hover:border-green-800/80 hover:shadow-green-500/5 dark:hover:shadow-green-500/10',
    title: 'group-hover:text-green-600 dark:group-hover:text-green-400',
    button: 'group-hover:bg-green-600 dark:group-hover:bg-green-500 hover:!bg-green-700 dark:hover:!bg-green-600 group-hover:shadow-green-100 dark:group-hover:shadow-none',
    topBar: 'bg-green-500 dark:bg-green-400'
  },
  blue: {
    badge: 'bg-blue-100 text-blue-800 dark:bg-blue-900/40 dark:text-blue-200',
    card: 'hover:border-blue-300 dark:hover:border-blue-800/80 hover:shadow-blue-500/5 dark:hover:shadow-blue-500/10',
    title: 'group-hover:text-blue-600 dark:group-hover:text-blue-400',
    button: 'group-hover:bg-blue-600 dark:group-hover:bg-blue-500 hover:!bg-blue-700 dark:hover:!bg-blue-600 group-hover:shadow-blue-100 dark:group-hover:shadow-none',
    topBar: 'bg-blue-500 dark:bg-blue-400'
  },
  teal: {
    badge: 'bg-teal-100 text-teal-800 dark:bg-teal-900/40 dark:text-teal-200',
    card: 'hover:border-teal-300 dark:hover:border-teal-800/80 hover:shadow-teal-500/5 dark:hover:shadow-teal-500/10',
    title: 'group-hover:text-teal-600 dark:group-hover:text-teal-400',
    button: 'group-hover:bg-teal-600 dark:group-hover:bg-teal-500 hover:!bg-teal-700 dark:hover:!bg-teal-600 group-hover:shadow-teal-100 dark:group-hover:shadow-none',
    topBar: 'bg-teal-500 dark:bg-teal-400'
  },
  pink: {
    badge: 'bg-pink-100 text-pink-800 dark:bg-pink-900/40 dark:text-pink-200',
    card: 'hover:border-pink-300 dark:hover:border-pink-800/80 hover:shadow-pink-500/5 dark:hover:shadow-pink-500/10',
    title: 'group-hover:text-pink-600 dark:group-hover:text-pink-400',
    button: 'group-hover:bg-pink-600 dark:group-hover:bg-pink-500 hover:!bg-pink-700 dark:hover:!bg-pink-600 group-hover:shadow-pink-100 dark:group-hover:shadow-none',
    topBar: 'bg-pink-500 dark:bg-pink-400'
  },
  yellow: {
    badge: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/40 dark:text-yellow-200',
    card: 'hover:border-yellow-300 dark:hover:border-yellow-800/80 hover:shadow-yellow-500/5 dark:hover:shadow-yellow-500/10',
    title: 'group-hover:text-yellow-600 dark:group-hover:text-yellow-400',
    button: 'group-hover:bg-yellow-600 dark:group-hover:bg-yellow-500 hover:!bg-yellow-700 dark:hover:!bg-yellow-600 group-hover:shadow-yellow-100 dark:group-hover:shadow-none',
    topBar: 'bg-yellow-500 dark:bg-yellow-400'
  }
}

// Explicit mapping for known platforms
const platformColorMap = {
  remotive: 'purple',
  ashby: 'yellow',
  adzuna: 'amber',
  greenhouse: 'green',
  lever: 'blue',
  weworkremotely: 'teal',
  workable: 'pink'
}

const colorKeys = Object.keys(themeColors)

// Dynamically compute the theme to use based on the platform name
const theme = computed(() => {
  const platformName = props.job?.platform || ''
  const key = platformName.toLowerCase().replace(/\s+/g, '')
  
  // Use explicit mapping if it exists
  if (platformColorMap[key]) {
    return themeColors[platformColorMap[key]]
  }
  
  // For any new or unknown platforms, assign a consistent color by hashing the name
  if (key) {
    let hash = 0
    for (let i = 0; i < key.length; i++) {
      hash = key.charCodeAt(i) + ((hash << 5) - hash)
    }
    const index = Math.abs(hash) % colorKeys.length
    return themeColors[colorKeys[index]]
  }
  
  // Fallback
  return themeColors.blue
})
</script>

<template>
  <div 
    class="group bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700/60 rounded-xl p-6 transition-all duration-300 hover:-translate-y-1.5 active:scale-[0.985] active:duration-100 relative overflow-hidden"
    :class="theme.card"
  >
    <!-- Subtle top bar indicator matching the platform theme -->
    <div 
      class="absolute top-0 left-0 w-full h-[4px] opacity-60 group-hover:opacity-100 transition-all duration-300"
      :class="theme.topBar"
    ></div>

    <div class="flex justify-between items-start mb-4">
      <div class="flex items-center gap-2">
        <span 
          class="px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wider transition-all duration-300"
          :class="theme.badge"
        >
          {{ job.platform }}
        </span>
        <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatRelativeDate(job.posted_at) }}</span>
      </div>
    </div>
 
    <h3 
      class="text-xl font-bold text-gray-900 dark:text-white mb-1 transition-colors duration-300"
      :class="theme.title"
    >
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
        class="bg-gray-900 dark:bg-gray-700 text-white px-5 py-2.5 rounded-lg font-semibold transition-all duration-300 inline-flex items-center gap-2 shadow-lg shadow-gray-200 dark:shadow-none"
        :class="theme.button"
      >
        View Details
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transform group-hover:translate-x-1 transition-transform duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
        </svg>
      </a>
    </div>
  </div>
</template>
