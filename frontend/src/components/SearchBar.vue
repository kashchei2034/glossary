<template>
  <div>
    <!-- Command Palette Trigger Button -->
    <button
      @click="openModal"
      class="flex items-center justify-between gap-2 px-2.5 sm:px-3 py-1.5 rounded-xl bg-slate-100 dark:bg-slate-800/80 hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-500 dark:text-slate-400 text-xs font-medium border border-slate-200 dark:border-slate-700/60 transition-all w-full max-w-xs shadow-xs"
    >
      <div class="flex items-center gap-2 min-w-0">
        <Icon name="search" class="w-4 h-4 text-slate-400 shrink-0" />
        <span class="truncate text-left hidden xs:inline">Search docs...</span>
        <span class="truncate text-left xs:hidden">Search...</span>
      </div>
      <kbd class="hidden sm:inline-block px-1.5 py-0.5 text-[10px] font-mono bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded shadow-xs text-slate-400 shrink-0">⌘K</kbd>
    </button>

    <!-- Modal Dialog -->
    <Teleport to="body">
      <Transition name="fade">
        <div
          v-if="docStore.isSearchOpen"
          class="fixed inset-0 z-50 flex items-start justify-center pt-12 sm:pt-24 px-3 sm:px-4 bg-slate-950/60 backdrop-blur-sm"
          @click.self="closeModal"
        >
          <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-2xl w-full max-w-2xl overflow-hidden flex flex-col max-h-[85vh]">
            <!-- Search Bar Input Header -->
            <div class="p-3.5 sm:p-4 border-b border-slate-200 dark:border-slate-800 flex items-center gap-2.5 sm:gap-3">
              <Icon name="search" class="w-5 h-5 text-brand-500 shrink-0" />
              <input
                ref="searchInputRef"
                v-model="query"
                type="text"
                placeholder="Search documents by title, tags, or content..."
                class="flex-1 bg-transparent border-0 text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-0 text-sm sm:text-base"
                @input="handleInput"
                @keydown.down.prevent="navigateResults(1)"
                @keydown.up.prevent="navigateResults(-1)"
                @keydown.enter.prevent="selectActive"
                @keydown.esc="closeModal"
              />
              <Icon v-if="docStore.searchLoading" name="loader" class="w-4 h-4 text-slate-400 animate-spin" />
              <button @click="closeModal" class="p-1 rounded-lg text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800">
                <Icon name="x" class="w-5 h-5" />
              </button>
            </div>

            <!-- Results List -->
            <div class="flex-1 overflow-y-auto p-2">
              <div v-if="docStore.searchResults.length > 0" class="space-y-1">
                <div
                  v-for="(res, idx) in docStore.searchResults"
                  :key="res.id"
                  @click="navigateTo(res.slug)"
                  @mouseenter="selectedIndex = idx"
                  class="p-3 rounded-xl cursor-pointer transition-all flex flex-col gap-1 border"
                  :class="
                    selectedIndex === idx
                      ? 'bg-brand-50/70 dark:bg-brand-950/30 border-brand-200 dark:border-brand-800/60'
                      : 'border-transparent hover:bg-slate-50 dark:hover:bg-slate-800/50'
                  "
                >
                  <div class="flex items-center justify-between gap-2">
                    <h4 class="text-sm font-semibold text-slate-900 dark:text-slate-100 flex items-center gap-2">
                      <Icon name="file" class="w-4 h-4 text-brand-500 shrink-0" />
                      <span>{{ res.title }}</span>
                    </h4>
                    <span v-if="res.category" class="px-2 py-0.5 text-[10px] font-medium rounded-full bg-slate-100 dark:bg-slate-800 text-slate-500 shrink-0">
                      {{ res.category }}
                    </span>
                  </div>
                  <p class="text-xs text-slate-500 dark:text-slate-400 line-clamp-2" v-html="res.snippet"></p>
                </div>
              </div>

              <!-- Empty Results -->
              <div v-else-if="query.trim() && !docStore.searchLoading" class="py-12 text-center text-slate-400 text-sm">
                No matching documents found for "<span class="font-medium text-slate-600 dark:text-slate-300">{{ query }}</span>"
              </div>

              <!-- Search Hints -->
              <div v-else-if="!query.trim()" class="py-8 px-4 text-center text-xs text-slate-400">
                Type keywords, topics, or code snippets to search instantly using PostgreSQL full-text index.
              </div>
            </div>

            <!-- Footer Keyboard Nav Hints -->
            <div class="px-4 py-2 bg-slate-50 dark:bg-slate-950 border-t border-slate-200 dark:border-slate-800 text-[11px] text-slate-400 flex items-center justify-between">
              <div class="flex items-center gap-3">
                <span><kbd class="font-mono bg-slate-200 dark:bg-slate-800 px-1 rounded">↑↓</kbd> Navigate</span>
                <span><kbd class="font-mono bg-slate-200 dark:bg-slate-800 px-1 rounded">↵</kbd> Select</span>
                <span><kbd class="font-mono bg-slate-200 dark:bg-slate-800 px-1 rounded">ESC</kbd> Close</span>
              </div>
              <span>Instant Search</span>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { useDocStore } from '@/stores/docStore';
import Icon from '@/components/Icon.vue';

const docStore = useDocStore();
const router = useRouter();
const query = ref('');
const selectedIndex = ref(0);
const searchInputRef = ref<HTMLInputElement | null>(null);

let debounceTimer: any = null;

function handleInput() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(() => {
    docStore.searchDocs(query.value);
    selectedIndex.value = 0;
  }, 200);
}

function openModal() {
  docStore.isSearchOpen = true;
  query.value = '';
  docStore.searchResults = [];
  nextTick(() => {
    searchInputRef.value?.focus();
  });
}

function closeModal() {
  docStore.isSearchOpen = false;
}

function navigateResults(delta: number) {
  if (docStore.searchResults.length === 0) return;
  const max = docStore.searchResults.length - 1;
  selectedIndex.value = Math.max(0, Math.min(max, selectedIndex.value + delta));
}

function selectActive() {
  if (docStore.searchResults[selectedIndex.value]) {
    navigateTo(docStore.searchResults[selectedIndex.value].slug);
  }
}

function navigateTo(slug: string) {
  closeModal();
  router.push(`/doc/${slug}`);
}

function handleKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
    e.preventDefault();
    if (docStore.isSearchOpen) {
      closeModal();
    } else {
      openModal();
    }
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown);
});
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
