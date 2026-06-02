<script setup>
import { ref, computed } from 'vue'
import { onClickOutside } from '@vueuse/core'

const props = defineProps({
  modelValue: {
    type: String,
    required: true
  },
  options: {
    type: Array,
    required: true
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Select an option'
  }
})

const emit = defineEmits(['update:modelValue'])

const dropdownRef = ref(null)

const isOpen = ref(false)
const searchQuery = ref('')

onClickOutside(dropdownRef, () => {
  isOpen.value = false
  searchQuery.value = ''
})

const filteredOptions = computed(() => {
  if (!searchQuery.value) return props.options

  return props.options.filter(opt =>
    opt.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const selectedLabel = computed(() => {
  return props.options.find(opt => opt === props.modelValue)
    || props.placeholder
})

const selectOption = (option) => {
  emit('update:modelValue', option)
  isOpen.value = false
  searchQuery.value = ''
}

const toggleDropdown = () => {
  isOpen.value = !isOpen.value

  if (isOpen.value) {
    searchQuery.value = ''
  }
}
</script>

<template>
  <div class="relative w-full" ref="dropdownRef">
    <!-- Label -->
    <label v-if="label" class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
      {{ label }}
    </label>

    <!-- Main Button -->
    <button
      @click="toggleDropdown"
      class="w-full bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 text-gray-900 dark:text-white text-sm rounded-lg p-3 flex items-center justify-between transition-all duration-200 hover:border-gray-300 dark:hover:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
    >
      <span class="font-medium">{{ selectedLabel }}</span>
      <svg
        class="w-3 h-3 transition-transform duration-300"
        :class="{ 'rotate-180': isOpen }"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
      </svg>
    </button>

    <Transition
      enter-active-class="transition ease-out duration-150"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in duration-100"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div
        v-if="isOpen"
        class="absolute top-full left-0 right-0 mt-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg z-[60]"
      >
        <div class="p-3 border-b border-gray-200 dark:border-gray-700 sticky top-0 bg-white dark:bg-gray-800">
          <input
            v-model="searchQuery"
            type="text"
            @click.stop
            placeholder="Search..."
            class="w-full bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 text-gray-900 dark:text-white text-sm rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        <div class="max-h-72 overflow-y-auto">
          <button
            type="button"
            v-for="(option, index) in filteredOptions"
            :key="index"
            @click="selectOption(option)"
            :class="[
              'w-full text-left px-4 py-3 transition-colors duration-150 flex items-center',
              modelValue === option
                ? 'bg-blue-50 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 border-l-4 border-blue-500'
                : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700/50 border-l-4 border-transparent'
            ]"
          >
            <span class="font-medium">{{ option }}</span>
          </button>

          <!-- Empty State -->
          <div v-if="filteredOptions.length === 0" class="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
            No options found
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
/* Smooth scrollbar styling */
div::-webkit-scrollbar {
  width: 6px;
}

div::-webkit-scrollbar-track {
  background: transparent;
}

div::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

div::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

:root.dark div::-webkit-scrollbar-thumb {
  background: #475569;
}

:root.dark div::-webkit-scrollbar-thumb:hover {
  background: #64748b;
}
</style>