<template>
  <div class="flex flex-col gap-4">
    <div
      v-for="i in 5"
      :key="i"
      class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl px-6 h-48 animate-pulse shadow-sm"
    >
      <div class="flex justify-between mb-4 mt-3">
        <div class="h-4 w-24 bg-gray-200 dark:bg-gray-700 rounded-full" />
        <div class="h-4 w-6 bg-gray-200 dark:bg-gray-700 rounded" />
      </div>
      <div class="h-8 w-64 bg-gray-200 dark:bg-gray-700 rounded mb-4" />
      <div class="h-6 w-32 bg-gray-200 dark:bg-gray-700 rounded mb-4" />
      <div class="h-10 w-32 bg-gray-200 dark:bg-gray-700 rounded" />
    </div>
  </div>
</template>
