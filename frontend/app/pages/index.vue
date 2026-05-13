<script setup>
// Redirect home to /jobs automatically for now, 
// or show a simple hero section.
// navigateTo('/jobs')
</script>

<template>
  <div class="min-h-screen bg-white">
    <!-- Simple Hero Section -->
    <header class="bg-white border-b border-gray-100">
      <div class="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
        <NuxtLink to="/" class="text-2xl font-black text-blue-600 tracking-tight">JobStream</NuxtLink>
        <NuxtLink 
          to="/jobs" 
          class="bg-blue-600 text-white px-6 py-2 rounded-lg font-bold hover:bg-blue-700 transition-all"
        >
          Find Jobs
        </NuxtLink>
      </div>
    </header>

    <main class="max-w-4xl mx-auto px-4 py-24 text-center">
      <h1 class="text-6xl font-black text-gray-900 mb-6 tracking-tight">
        Your Next Career Move, <br/> 
        <span class="text-blue-600">All in One Stream.</span>
      </h1>
      <p class="text-xl text-gray-600 mb-10 max-w-2xl mx-auto">
        Aggregate job listings from LinkedIn, Indeed, Glassdoor and more. 
        Filter, search, and apply from a single, fast interface.
      </p>
      
      <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
        <NuxtLink 
          to="/jobs" 
          class="w-full sm:w-auto bg-gray-900 text-white text-lg px-10 py-4 rounded-xl font-bold hover:bg-gray-800 transition-all shadow-xl shadow-gray-200"
        >
          Browse All Jobs
        </NuxtLink>
        <NuxtLink 
          to="/jobs?work_model=remote" 
          class="w-full sm:w-auto bg-white text-gray-900 text-lg px-10 py-4 rounded-xl font-bold border-2 border-gray-100 hover:border-gray-200 transition-all"
        >
          Remote Only
        </NuxtLink>
      </div>

      <div class="mt-24 grid grid-cols-2 md:grid-cols-4 gap-8 opacity-40 grayscale">
        <div class="font-black text-2xl italic">LinkedIn</div>
        <div class="font-black text-2xl italic">Indeed</div>
        <div class="font-black text-2xl italic">Glassdoor</div>
        <div class="font-black text-2xl italic">Handshake</div>
      </div>
    </main>
  </div>
</template>
