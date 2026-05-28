<template>
  <NuxtLayout>
    <div class="min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 transition-colors duration-300">
    <NuxtRouteAnnouncer />
    <NuxtPage />
  </div>
  </NuxtLayout>
</template>

<style>
/* Global styles */
body {
  @apply font-sans antialiased;
}
</style>
