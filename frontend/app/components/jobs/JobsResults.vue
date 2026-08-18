<script setup>

defineProps({
  jobs: {
    type: Array,
    default: () => []
  },
  metadata: {
    type: Object,
    default: () => ({})
  },
  page: {
    type: Number,
    default: 1
  },
  pending: Boolean,
  error: [Object, Error, null]
})

const emit = defineEmits(['retry', 'prev-page', 'next-page'])

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const handlePrevPage = () => {
  emit('prev-page')
  scrollToTop()
}

const handleNextPage = () => {
  emit('next-page')
  scrollToTop()
}
</script>

<template>
  <!-- Full skeleton loader for initial load when no data exists yet -->
  <JobsJobListSkeleton v-if="pending && (!jobs || jobs.length === 0)" />

  <div
    v-else-if="error"
    class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-900 text-red-700 dark:text-red-400 p-8 rounded-xl text-center shadow-sm"
  >
    <h3 class="text-xl font-bold mb-2">Something went wrong</h3>
    <p class="mb-4">We couldn't load the jobs right now. Please try again later.</p>
    <button
      type="button"
      class="bg-red-600 text-white px-6 py-2 rounded-lg font-bold hover:bg-red-700 transition-colors"
      @click="$emit('retry')"
    >
      Retry
    </button>
  </div>

  <div v-else class="space-y-6 relative">
    <!-- Subtle top loading indicator bar during background refetches -->
    <div
      v-if="pending && jobs.length > 0"
      class="absolute -top-3 left-0 right-0 h-1 bg-gradient-to-r from-blue-500 via-purple-500 to-cyan-500 rounded-full animate-pulse z-10"
    />

    <div
      v-if="jobs.length > 0"
      class="flex flex-col gap-6 transition-opacity duration-200"
      :class="{ 'opacity-60 pointer-events-none': pending }"
    >
      <JobCard v-for="job in jobs" :key="job.id" :job="job" />

      <div
        v-if="metadata?.total_pages > 1"
        class="flex justify-center items-center gap-4 mt-8 py-4"
      >
        <button
          type="button"
          :disabled="page <= 1"
          class="p-2 rounded-lg border border-gray-300 dark:border-gray-700 disabled:opacity-30 disabled:cursor-not-allowed bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-all"
          @click="handlePrevPage"
        >
          <span class="sr-only">Previous page</span>
          ‹
        </button>
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
          Page {{ page }} of {{ metadata.total_pages }}
        </span>
        <button
          type="button"
          :disabled="page >= metadata.total_pages"
          class="p-2 rounded-lg border border-gray-300 dark:border-gray-700 disabled:opacity-30 disabled:cursor-not-allowed bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-all"
          @click="handleNextPage"
        >
          <span class="sr-only">Next page</span>
          ›
        </button>
      </div>
    </div>

    <div
      v-else
      class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-12 text-center shadow-sm"
    >
      <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">No jobs found</h3>
      <p class="text-gray-600 dark:text-gray-400 mb-6">
        Try adjusting your filters or search terms to find what you're looking for.
      </p>
      <button
        type="button"
        class="inline-flex items-center gap-2 bg-gray-900 dark:bg-gray-700 text-white font-semibold px-5 py-2.5 rounded-lg hover:bg-blue-700 transition-colors"
        @click="$emit('retry')"
      >
        Reset Search / Retry
      </button>
    </div>
  </div>
</template>
