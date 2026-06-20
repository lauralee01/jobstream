<script setup>
definePageMeta({ layout: 'default' })

const { platforms } = useJobMetadata()

const platformNames = computed(() =>
  (platforms.value || []).join(', ') || 'top job boards'
)
</script>

<template>
  <header class="bg-white dark:bg-gray-900 border-b border-gray-100 dark:border-gray-800">
    <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
      <NuxtLink to="/" class="text-2xl font-black text-gray-900 dark:text-white tracking-tight">
        JobStream
      </NuxtLink>
      <div class="flex items-center gap-4">
        <ClientOnly>
          <ColorModeToggle />
        </ClientOnly>
        <NuxtLink
          to="/jobs"
          class="bg-gray-900 dark:bg-gray-700 text-white px-6 py-2 rounded-lg font-bold hover:bg-blue-700 transition-all"
        >
          Find Jobs
        </NuxtLink>
      </div>
    </div>
  </header>

  <main class="max-w-4xl mx-auto px-4 py-24 text-center overflow-hidden">
    <h1 class="text-5xl md:text-7xl font-black text-gray-900 dark:text-white mb-6 tracking-tight leading-tight">
      <span class="animate-fade-in-up delay-100 inline-block">Your Next Career Move,</span> <br>
      <span class="animate-fade-in-up delay-300 inline-block">
        <span class="gradient-text">All in One Stream.</span>
      </span>
    </h1>
    <p class="animate-fade-in-up delay-500 text-xl text-gray-600 dark:text-gray-400 mb-10 max-w-2xl mx-auto">
      Aggregate job listings from {{ platformNames }}.
      Filter, search, and apply from a single, fast interface.
    </p>

    <div class="animate-fade-in-up delay-700 flex flex-col sm:flex-row items-center justify-center gap-4">
      <NuxtLink
        to="/jobs"
        class="w-full sm:w-auto bg-gray-900 dark:bg-blue-600 text-white text-lg px-10 py-4 rounded-xl font-bold hover:bg-gray-800 dark:hover:bg-blue-700 transition-all shadow-xl"
      >
        Browse All Jobs
      </NuxtLink>
      <NuxtLink
        to="/jobs?remote=true"
        class="w-full sm:w-auto bg-white dark:bg-gray-800 text-gray-900 dark:text-white text-lg px-10 py-4 rounded-xl font-bold border-2 border-gray-100 dark:border-gray-700 hover:border-gray-200 transition-all"
      >
        Remote Only
      </NuxtLink>
    </div>
  </main>
</template>

<style scoped>
.animate-fade-in-up {
  opacity: 0;
  transform: translateY(30px);
  animation: fadeInUp 0.9s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.delay-100 {
  animation-delay: 100ms;
}

.delay-300 {
  animation-delay: 300ms;
}

.delay-500 {
  animation-delay: 500ms;
}

.delay-700 {
  animation-delay: 700ms;
}

.gradient-text {
  --gradient-from: #2563eb;
  --gradient-via: #7c3aed;
  --gradient-to: #06b6d4;

  background: linear-gradient(
    to right,
    var(--gradient-from),
    var(--gradient-via),
    var(--gradient-to),
    var(--gradient-from)
  );
  background-size: 200% auto;
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  -webkit-text-fill-color: transparent;
  animation: gradientFlow 6s linear infinite;
  display: inline-block;
}

.dark .gradient-text {
  --gradient-from: #60a5fa;
  --gradient-via: #c084fc;
  --gradient-to: #22d3ee;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes gradientFlow {
  0% {
    background-position: 0% center;
  }
  100% {
    background-position: 200% center;
  }
}
</style>
