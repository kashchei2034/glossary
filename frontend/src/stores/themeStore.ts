import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useThemeStore = defineStore('themeStore', () => {
  const isDark = ref<boolean>(
    localStorage.getItem('kb-theme') === 'dark' ||
      (!('kb-theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)
  );

  function applyTheme() {
    if (isDark.value) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('kb-theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('kb-theme', 'light');
    }
  }

  function toggleTheme() {
    isDark.value = !isDark.value;
    applyTheme();
  }

  // Initialize
  applyTheme();

  return { isDark, toggleTheme, applyTheme };
});
