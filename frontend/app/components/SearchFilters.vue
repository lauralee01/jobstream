<script setup>
const props = defineProps({
  modelValue: {
    type: Object,
    required: true
  },
  categories: {
    type: Array,
    default: () => []
  },
  platforms: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'search', 'clear'])

const sortOptions = [
  { label: 'Newest first', sortBy: 'posted_at', sortOrder: 'desc' },
  { label: 'Oldest first', sortBy: 'posted_at', sortOrder: 'asc' }
]

const sortLabels = computed(() => sortOptions.map(opt => opt.label))

const currentSortLabel = computed(() => {
  const current = `${props.modelValue.sortBy}:${props.modelValue.sortOrder}`
  const match = sortOptions.find(
    (option) => `${option.sortBy}:${option.sortOrder}` === current
  )
  return match?.label || sortOptions[0].label
})

const handleSortChange = (label) => {
  const selected = sortOptions.find(opt => opt.label === label)
  if (selected) {
    emit('update:modelValue', {
      ...props.modelValue,
      sortBy: selected.sortBy,
      sortOrder: selected.sortOrder
    })
  }
}

const applyFilters = () => {
  emit('search')
}

const clearFilters = () => {
  emit('clear')
}

const updateField = (patch) => {
  emit('update:modelValue', { ...props.modelValue, ...patch })
}

const togglePlatform = (id) => {
  const platforms = [...(props.modelValue.platforms || [])]
  const index = platforms.indexOf(id)
  if (index === -1) {
    platforms.push(id)
  } else {
    platforms.splice(index, 1)
  }
  updateField({ platforms })
}
</script>

<template>
  <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-none rounded-xl p-6 shadow-sm sticky top-24 transition-colors duration-300">
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-lg font-bold text-gray-900 dark:text-white">Filters</h2>
      <button
        type="button"
        class="text-sm text-gray-900 dark:text-white hover:text-blue-800 dark:hover:text-blue-300 font-medium transition-colors"
        @click="clearFilters"
      >
        Clear all
      </button>
    </div>

    <div class="mb-6">
      <CustomSelect
        :model-value="currentSortLabel"
        label="Sort"
        :options="sortLabels"
        @update:model-value="handleSortChange"
      /> 

    </div>

    <div class="mb-6">
      <CustomSelect
        :model-value="modelValue.category"
        label="Category"
        placeholder="All Categories"
        :options="['All Categories', ...categories]"
        @update:model-value="updateField({ category: $event === 'All Categories' ? '' : $event })"
      />
    </div>

    <div class="mb-6 flex items-center">
      <input
        id="remote-only"
        :checked="modelValue.remote"
        type="checkbox"
        class="w-4 h-4 text-blue-600 bg-gray-100 dark:bg-gray-900 border-gray-300 dark:border-gray-700 rounded focus:ring-blue-500"
        @change="updateField({ remote: $event.target.checked })"
      >
      <label for="remote-only" class="ml-2 text-sm font-medium text-gray-700 dark:text-gray-300 select-none cursor-pointer">
        Remote Only
      </label>
    </div>

    <div class="mb-6">
      <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">Minimum Salary (USD)</label>
      <div class="relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
        </div>
        <input
          :value="modelValue.salaryMin"
          type="number"
          placeholder="e.g. 80000"
          class="bg-gray-50 dark:bg-gray-900 border border-gray-300 dark:border-gray-700 text-gray-900 dark:text-white text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full pl-7 p-2.5 outline-none transition-colors"
          @input="updateField({ salaryMin: $event.target.value })"
        >
      </div>
    </div>

    <div class="mb-6">
      <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">Sources</label>
      <div class="space-y-2">
        <label v-for="platform in platforms" :key="platform.id" class="flex items-center group cursor-pointer">
          <input
            type="checkbox"
            :checked="modelValue.platforms?.includes(platform.id)"
            class="w-4 h-4 text-blue-600 bg-gray-100 dark:bg-gray-900 border-gray-300 dark:border-gray-700 rounded focus:ring-blue-500"
            @change="togglePlatform(platform.id)"
          >
          <span class="ml-2 text-sm text-gray-600 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white transition-colors">
            {{ platform.name }}
          </span>
        </label>
      </div>
    </div>

    <button
      type="button"
      class="w-full bg-gray-900 dark:bg-gray-700 text-white font-bold py-3 px-4 rounded-lg hover:bg-blue-700 transition-all active:scale-[0.98]"
      @click="applyFilters"
    >
      Apply Filters
    </button>
  </div>
</template>
