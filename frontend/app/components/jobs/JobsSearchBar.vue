<script setup>
const props = defineProps({
  modelValue: { type: Object, required: true }
})
const draft = defineModel({ type: Object, required: true })

const emit = defineEmits(['update:modelValue', 'search'])

// Track previous values to detect when clearing happens
const prevKeyword = ref(props.modelValue.keyword)
const prevLocation = ref(props.modelValue.location)

// When keyword is cleared
watch(
  () => props.modelValue.keyword,
  (newValue) => {
    if (!newValue?.trim() && prevKeyword.value?.trim()) {
      emit('search')
    }
    prevKeyword.value = newValue
  }
)

// When location is cleared
watch(
  () => props.modelValue.location,
  (newValue) => {
    if (!newValue?.trim() && prevLocation.value?.trim()) {
      emit('search')
    }
    prevLocation.value = newValue
  }
)
</script>

<template>
  <div class="bg-white dark:bg-gray-800 p-4 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 mb-8 flex flex-col sm:flex-row gap-4 transition-colors duration-300">
    <div class="flex-grow relative">
      <span class="absolute inset-y-0 left-3 flex items-center text-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </span>
      <input
        v-model="draft.keyword"
        type="text"
        placeholder="Job title, keywords or company"
        class="w-full pl-10 pr-4 py-2.5 bg-gray-50 dark:bg-gray-900 border border-gray-300 dark:border-gray-700 text-gray-900 text-base dark:text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
        @keyup.enter="$emit('search')"
      >
    </div>
    <div class="flex-grow relative">
      <span class="absolute inset-y-0 left-3 flex items-center text-gray-400">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
      </span>
      <input
        v-model="draft.location"
        type="text"
        placeholder="City, state or remote"
        class="w-full pl-10 pr-4 py-2.5 bg-gray-50 dark:bg-gray-900 border border-gray-300 dark:border-gray-700 text-base text-gray-900 dark:text-white rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
        @keyup.enter="$emit('search')"
      >
    </div>
    <button
      type="button"
      class="bg-gray-900 dark:bg-gray-700 text-white px-8 py-2.5 rounded-lg font-bold hover:bg-blue-700 transition-all active:scale-[0.98] shadow-lg shadow-blue-200 dark:shadow-none"
      @click="$emit('search')"
    >
      Search
    </button>
  </div>
</template>
