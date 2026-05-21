<script setup>
const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      keyword: '',
      location: '',
      platforms: [],
      remote: false,
      salaryMin: '',
      category: ''
    })
  }
})

const emit = defineEmits(['update:modelValue', 'search'])

const { fetchCategories, fetchPlatforms } = useJobs()
const { data: fetchedCategories } = await fetchCategories()
const categories = computed(() => fetchedCategories.value || [])

const { data: fetchedPlatforms } = await fetchPlatforms()
const platforms = computed(() => {
  return (fetchedPlatforms.value || []).map(p => ({
    id: p,
    name: p
  }))
})

const filters = ref({ ...props.modelValue })

// Sync local ref with parent modelValue
watch(() => props.modelValue, (newVal) => {
  if (JSON.stringify(newVal) !== JSON.stringify(filters.value)) {
    filters.value = { ...newVal }
  }
}, { deep: true })

// Emit changes
watch(filters, (newVal) => {
  if (JSON.stringify(newVal) !== JSON.stringify(props.modelValue)) {
    emit('update:modelValue', newVal)
  }
}, { deep: true })

const clearFilters = () => {
  const cleared = {
    keyword: '',
    location: '',
    platforms: [],
    remote: false,
    salaryMin: '',
    category: ''
  }
  filters.value = cleared
  emit('update:modelValue', cleared)
  emit('search')
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-6 shadow-sm sticky top-24 transition-colors duration-300">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-lg font-bold text-gray-900 dark:text-white">Filters</h2>
      <button 
        @click="clearFilters"
        class="text-sm text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 font-medium transition-colors"
      >
        Clear all
      </button>
    </div>

    <!-- Category -->
    <div class="mb-6">
      <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">Category</label>
      <select 
        v-model="filters.category"
        class="w-full bg-gray-50 dark:bg-gray-900 border border-gray-300 dark:border-gray-700 text-gray-900 dark:text-white text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 outline-none transition-colors"
      >
        <option value="">All Categories</option>
        <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
      </select>
    </div>

    <!-- Remote Checkbox -->
    <div class="mb-6 flex items-center">
      <input 
        id="remote-only" 
        v-model="filters.remote"
        type="checkbox" 
        class="w-4 h-4 text-blue-600 bg-gray-100 dark:bg-gray-900 border-gray-300 dark:border-gray-700 rounded focus:ring-blue-500 transition-colors"
      >
      <label for="remote-only" class="ml-2 text-sm font-medium text-gray-700 dark:text-gray-300 select-none cursor-pointer hover:text-gray-900 dark:hover:text-white transition-colors">
        Remote Only
      </label>
    </div>

    <!-- Salary Range -->
    <div class="mb-6">
      <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">Minimum Salary (USD)</label>
      <div class="relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
        </div>
        <input 
          v-model="filters.salaryMin"
          type="number" 
          placeholder="e.g. 80000"
          class="bg-gray-50 dark:bg-gray-900 border border-gray-300 dark:border-gray-700 text-gray-900 dark:text-white text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full pl-7 p-2.5 outline-none transition-colors"
        >
      </div>
    </div>

    <!-- Platforms -->
    <div class="mb-6">
      <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">Sources</label>
      <div class="space-y-2">
        <label v-for="platform in platforms" :key="platform.id" class="flex items-center group cursor-pointer">
          <input 
            v-model="filters.platforms"
            type="checkbox" 
            :value="platform.id"
            class="w-4 h-4 text-blue-600 bg-gray-100 dark:bg-gray-900 border-gray-300 dark:border-gray-700 rounded focus:ring-blue-500 transition-colors"
          >
          <span class="ml-2 text-sm text-gray-600 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white transition-colors">{{ platform.name }}</span>
        </label>
      </div>
    </div>

    <button 
      @click="$emit('search')"
      class="w-full bg-blue-600 text-white font-bold py-3 px-4 rounded-lg hover:bg-blue-700 transition-all shadow-lg shadow-blue-200 dark:shadow-none active:scale-[0.98]"
    >
      Apply Filters
    </button>
  </div>
</template>
