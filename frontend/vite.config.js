import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  base: '/',
  build: {
    outDir: 'dist',
  },
  server: {
    proxy: {
      '/prod': {
        target: 'https://sdbx.it',
        changeOrigin: true,
      },
      '/dev': {
        target: 'https://sdbx.it',
        changeOrigin: true,
      },
    },
  },
});
