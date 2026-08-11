<template>
  <div class="flex flex-col h-screen w-screen overflow-hidden bg-slate-50 dark:bg-slate-900 text-slate-900 dark:text-slate-100">
    <!-- Top Global Header -->
    <header class="h-14 border-b border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md px-2.5 sm:px-4 flex items-center justify-between z-30 shrink-0 gap-2">
      <!-- Left Brand Section -->
      <div class="flex items-center gap-1.5 sm:gap-3 shrink-0">
        <!-- Mobile Drawer Toggle -->
        <button
          @click="docStore.isSidebarOpen = !docStore.isSidebarOpen"
          class="md:hidden p-1.5 sm:p-2 rounded-xl text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
          aria-label="Toggle Sidebar"
        >
          <Icon name="menu" class="w-5 h-5" />
        </button>

        <!-- Brand Icon / Title for small screens -->
        <router-link to="/" class="flex items-center gap-1.5 sm:gap-2 md:hidden">
          <div class="w-7 h-7 rounded-lg bg-brand-500 text-white font-bold text-sm flex items-center justify-center shadow-xs">K</div>
          <span class="font-bold text-xs sm:text-sm hidden xs:inline-block">Help Desk</span>
        </router-link>
      </div>

      <!-- Center Search Bar (Fluid flex-1) -->
      <div class="flex-1 max-w-md mx-1 sm:mx-4 min-w-0">
        <SearchBar />
      </div>

      <!-- Right Action Controls -->
      <div class="flex items-center gap-1 sm:gap-2 shrink-0">
        <!-- Dark/Light Theme Switcher -->
        <button
          @click="themeStore.toggleTheme"
          class="p-1.5 sm:p-2 rounded-xl text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors shrink-0"
          title="Toggle Light/Dark Theme"
        >
          <Icon v-if="themeStore.isDark" name="sun" class="w-5 h-5 text-amber-400" />
          <Icon v-else name="moon" class="w-5 h-5 text-slate-600" />
        </button>

        <!-- Desktop New Doc Button -->
        <router-link
          to="/editor/new"
          class="hidden sm:flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-semibold shadow-sm transition-all"
        >
          <Icon name="plus" class="w-4 h-4" />
          <span>New Doc</span>
        </router-link>

        <!-- Mobile New Doc Icon Button -->
        <router-link
          to="/editor/new"
          class="sm:hidden flex items-center justify-center p-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white shadow-sm transition-all shrink-0"
          title="Create New Document"
        >
          <Icon name="plus" class="w-4 h-4" />
        </router-link>
      </div>
    </header>

    <!-- Main Viewport Container (Flex 1, overflow-hidden) -->
    <div class="flex flex-1 min-h-0 min-w-0 relative overflow-hidden">
      <!-- Desktop Sidebar -->
      <div class="hidden md:block h-full shrink-0">
        <SidebarNav />
      </div>

      <!-- Mobile Off-Canvas Drawer -->
      <Teleport to="body">
        <Transition name="drawer">
          <div
            v-if="docStore.isSidebarOpen"
            class="fixed inset-0 z-50 flex md:hidden bg-slate-950/60 backdrop-blur-xs"
            @click.self="docStore.isSidebarOpen = false"
          >
            <div class="w-64 h-full shadow-2xl">
              <SidebarNav @click="docStore.isSidebarOpen = false" />
            </div>
          </div>
        </Transition>
      </Teleport>

      <!-- Router View Content Body (strict overflow handling) -->
      <main class="flex-1 flex flex-col min-h-0 min-w-0 overflow-y-auto">
        <slot />
      </main>
    </div>

    <!-- Global Toast Alert Container -->
    <Toast />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import SidebarNav from '@/components/SidebarNav.vue';
import SearchBar from '@/components/SearchBar.vue';
import Toast from '@/components/Toast.vue';
import Icon from '@/components/Icon.vue';
import { useThemeStore } from '@/stores/themeStore';
import { useCategoryStore } from '@/stores/categoryStore';
import { useDocStore } from '@/stores/docStore';

const themeStore = useThemeStore();
const categoryStore = useCategoryStore();
const docStore = useDocStore();

onMounted(() => {
  categoryStore.fetchCategories();
  categoryStore.fetchTags();
});
</script>

<style scoped>
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
  transform: translateX(-100%);
}
</style>
