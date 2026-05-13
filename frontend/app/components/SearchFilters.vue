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

const platforms = [
  { id: 'linkedin', name: 'LinkedIn' },
  { id: 'indeed', name: 'Indeed' },
  { id: 'glassdoor', name: 'Glassdoor' },
  { id: 'handshake', name: 'Handshake' }
]

const categories = [
  'Software Engineer',
  'Designer',
  'Marketing',
  'Sales',
  'Accountant',
  'Product Manager'
]

const filters = ref({ ...props.modelValue })

// Sync local ref with parent modelValue
watch(() => props.modelValue, (newVal) => {
  filters.value = { ...newVal }
}, { deep: true })

// Emit changes
watch(filters, (newVal) => {
  emit('update:modelValue', newVal)
}, { deep: true })

const clearFilters = () => {
  filters.value = {
    keyword: '',
    location: '',
    platforms: [],
    remote: false,
    salaryMin: '',
    category: ''
  }
  emit('search')
}
</script>

<template>
  <div class="bg-white border border-gray-200 rounded-xl p-6 shadow-sm sticky top-24">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-lg font-bold text-gray-900">Filters</h2>
      <button 
        @click="clearFilters"
        class="text-sm text-blue-600 hover:text-blue-800 font-medium"
      >
        Clear all
      </button>
    </div>

    <!-- Category -->
    <div class="mb-6">
      <label class="block text-sm font-semibold text-gray-700 mb-2">Category</label>
      <select 
        v-model="filters.category"
        class="w-full bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5"
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
        class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 transition-colors"
      >
      <label for="remote-only" class="ml-2 text-sm font-medium text-gray-700 select-none cursor-pointer">
        Remote Only
      </label>
    </div>

    <!-- Salary Range -->
    <div class="mb-6">
      <label class="block text-sm font-semibold text-gray-700 mb-2">Minimum Salary (USD)</label>
      <div class="relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <span class="text-gray-500 sm:text-sm">$</span>
        </div>
        <input 
          v-model="filters.salaryMin"
          type="number" 
          placeholder="e.g. 80000"
          class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full pl-7 p-2.5"
        >
      </div>
    </div>

    <!-- Platforms -->
    <div class="mb-6">
      <label class="block text-sm font-semibold text-gray-700 mb-2">Sources</label>
      <div class="space-y-2">
        <label v-for="platform in platforms" :key="platform.id" class="flex items-center group cursor-pointer">
          <input 
            v-model="filters.platforms"
            type="checkbox" 
            :value="platform.id"
            class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 transition-colors"
          >
          <span class="ml-2 text-sm text-gray-600 group-hover:text-gray-900 transition-colors">{{ platform.name }}</span>
        </label>
      </div>
    </div>

    <button 
      @click="$emit('search')"
      class="w-full bg-blue-600 text-white font-bold py-3 px-4 rounded-lg hover:bg-blue-700 transition-all shadow-lg shadow-blue-200 active:scale-[0.98]"
    >
      Apply Filters
    </button>
  </div>
</template>
