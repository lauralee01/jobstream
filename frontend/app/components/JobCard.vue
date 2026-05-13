<script setup>
defineProps({
  job: {
    type: Object,
    required: true
  }
})

// Format date helper
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now - date)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return 'Yesterday'
  return `${diffDays} days ago`
}

// Map platform to colors/icons (example)
const getPlatformStyle = (platform) => {
  const p = platform?.toLowerCase()
  if (p === 'linkedin') return 'bg-blue-100 text-blue-700'
  if (p === 'indeed') return 'bg-blue-50 text-blue-800'
  if (p === 'glassdoor') return 'bg-green-100 text-green-700'
  return 'bg-gray-100 text-gray-700'
}
</script>

<template>
  <div class="bg-white border border-gray-200 rounded-xl p-6 hover:shadow-md transition-shadow">
    <div class="flex justify-between items-start mb-4">
      <div class="flex items-center gap-2">
        <span 
          class="px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase tracking-wider"
          :class="getPlatformStyle(job.platform)"
        >
          {{ job.platform }}
        </span>
        <span class="text-xs text-gray-500">{{ formatDate(job.posted_at) }}</span>
      </div>
      <button class="text-gray-400 hover:text-blue-600">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
        </svg>
      </button>
    </div>

    <h3 class="text-xl font-bold text-gray-900 mb-1">
      {{ job.title }}
    </h3>
    <p class="text-lg text-blue-600 font-medium mb-3">
      {{ job.company }}
    </p>

    <div class="flex flex-wrap gap-4 text-sm text-gray-600 mb-6">
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
        class="bg-blue-600 text-white px-5 py-2.5 rounded-lg font-semibold hover:bg-blue-700 transition-colors inline-flex items-center gap-2"
      >
        View Details
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
        </svg>
      </a>
    </div>
  </div>
</template>
