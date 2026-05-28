// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss', '@nuxtjs/color-mode'],
  colorMode: {
    classSuffix: ''
  },
  css: ['~/assets/css/main.css'],

  runtimeConfig: {
    public: {
      // In dev, prefer same-origin proxy (/api/v1) to avoid CORS round-trips.
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api/v1'
    }
  },

  nitro: {
    devProxy: {
      '/api/v1': {
        target: 'http://localhost:8080/api/v1',
        changeOrigin: true
      }
    }
  },

  routeRules: {
    '/': { prerender: false },
    '/jobs': { ssr: true }
  }
})
