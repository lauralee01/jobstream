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
  <JobsJobListSkeleton v-if="pending" />

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

  <div v-else class="space-y-6">
    <div v-if="jobs.length > 0" class="flex flex-col gap-6">
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
      <p class="text-gray-600 dark:text-gray-400">
        Try adjusting your filters or search terms to find what you're looking for.
      </p>
    </div>
  </div>
</template>
