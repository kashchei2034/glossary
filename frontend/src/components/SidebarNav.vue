<template>
  <aside class="flex flex-col h-full bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 w-64 shrink-0 select-none">
    <!-- Header Brand -->
    <div class="p-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between">
      <router-link to="/" class="flex items-center gap-2.5 group">
        <div class="w-8 h-8 rounded-xl bg-brand-500 text-white flex items-center justify-center font-bold text-lg shadow-md group-hover:scale-105 transition-transform">
          K
        </div>
        <div>
          <h1 class="text-sm font-bold text-slate-900 dark:text-white leading-tight">Help Desk</h1>
          <p class="text-[11px] text-slate-400">Knowledge Hub</p>
        </div>
      </router-link>
    </div>

    <!-- Navigation Menu List -->
    <div class="flex-1 overflow-y-auto p-3 space-y-4">
      <!-- Main Links -->
      <div class="space-y-1">
        <router-link
          to="/"
          class="flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-semibold transition-all"
          :class="route.path === '/' && !categoryStore.selectedCategorySlug ? 'bg-brand-50 dark:bg-brand-950/50 text-brand-600 dark:text-brand-400 font-bold' : 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400'"
          @click="resetFilters"
        >
          <Icon name="book-open" class="w-4 h-4" />
          <span>All Documents</span>
        </router-link>

        <router-link
          to="/admin"
          class="flex items-center gap-2.5 px-3 py-2 rounded-xl text-xs font-semibold transition-all"
          :class="route.path === '/admin' ? 'bg-brand-50 dark:bg-brand-950/50 text-brand-600 dark:text-brand-400 font-bold' : 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400'"
        >
          <Icon name="sliders" class="w-4 h-4" />
          <span>Document Manager</span>
        </router-link>
      </div>

      <!-- Categories Section -->
      <div>
        <div class="flex items-center justify-between px-3 mb-2 text-[11px] font-semibold uppercase tracking-wider text-slate-400">
          <span>Categories</span>
          <Icon name="folder-tree" class="w-3.5 h-3.5" />
        </div>
        <div class="space-y-1">
          <template v-for="cat in categoryStore.categories" :key="cat.id">
            <!-- Parent Category item -->
            <div class="space-y-1">
              <button
                @click="selectCategory(cat.slug)"
                class="w-full flex items-center justify-between px-3 py-1.5 rounded-xl text-xs transition-all text-left group"
                :class="categoryStore.selectedCategorySlug === cat.slug ? 'bg-slate-100 dark:bg-slate-800 text-brand-600 dark:text-brand-400 font-semibold' : 'hover:bg-slate-50 dark:hover:bg-slate-800/60 text-slate-700 dark:text-slate-300'"
              >
                <div class="flex items-center gap-2 truncate">
                  <Icon name="folder" class="w-4 h-4 text-slate-400 group-hover:text-brand-500 transition-colors shrink-0" />
                  <span class="truncate">{{ cat.name }}</span>
                </div>
                <span v-if="cat.doc_count && cat.doc_count > 0" class="text-[10px] px-1.5 py-0.5 rounded-full bg-slate-200 dark:bg-slate-800 text-slate-500">
                  {{ cat.doc_count }}
                </span>
              </button>

              <!-- Child Categories -->
              <div v-if="cat.children && cat.children.length > 0" class="pl-4 space-y-1 border-l border-slate-200 dark:border-slate-800 ml-4">
                <button
                  v-for="sub in cat.children"
                  :key="sub.id"
                  @click="selectCategory(sub.slug)"
                  class="w-full flex items-center justify-between px-2.5 py-1 rounded-lg text-[11px] transition-all text-left"
                  :class="categoryStore.selectedCategorySlug === sub.slug ? 'text-brand-600 dark:text-brand-400 font-bold bg-brand-50/50 dark:bg-brand-950/30' : 'text-slate-500 hover:text-slate-900 dark:hover:text-white'"
                >
                  <span class="truncate">{{ sub.name }}</span>
                  <span v-if="sub.doc_count" class="text-[9px] opacity-70">{{ sub.doc_count }}</span>
                </button>
              </div>
            </div>
          </template>
        </div>
      </div>

      <!-- Popular Tags -->
      <div v-if="categoryStore.tags.length > 0">
        <div class="flex items-center justify-between px-3 mb-2 text-[11px] font-semibold uppercase tracking-wider text-slate-400">
          <span>Filter by Tags</span>
          <Icon name="tag" class="w-3.5 h-3.5" />
        </div>
        <div class="flex flex-wrap gap-1 px-2">
          <button
            v-for="t in categoryStore.tags"
            :key="t.id"
            @click="toggleTag(t.slug)"
            class="px-2 py-1 rounded-lg text-[11px] font-medium transition-all"
            :class="categoryStore.selectedTagSlug === t.slug ? 'bg-brand-500 text-white shadow-xs' : 'bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-600 dark:text-slate-400'"
          >
            #{{ t.name }}
          </button>
        </div>
      </div>
    </div>

    <!-- Workspace Create Button Footer -->
    <div class="p-3 border-t border-slate-200 dark:border-slate-800">
      <router-link
        to="/editor/new"
        class="w-full flex items-center justify-center gap-2 py-2.5 px-4 rounded-xl bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold shadow-md shadow-brand-500/20 transition-all hover:scale-[1.02] active:scale-[0.98]"
      >
        <Icon name="plus" class="w-4 h-4" />
        <span>New Document</span>
      </router-link>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { useCategoryStore } from '@/stores/categoryStore';
import { useDocStore } from '@/stores/docStore';
import Icon from '@/components/Icon.vue';

const route = useRoute();
const router = useRouter();
const categoryStore = useCategoryStore();
const docStore = useDocStore();

function selectCategory(slug: string) {
  categoryStore.setCategory(slug);
  docStore.fetchDocuments(slug, categoryStore.selectedTagSlug);
  if (route.path !== '/') router.push('/');
}

function toggleTag(slug: string) {
  if (categoryStore.selectedTagSlug === slug) {
    categoryStore.setTag('');
  } else {
    categoryStore.setTag(slug);
  }
  docStore.fetchDocuments(categoryStore.selectedCategorySlug, categoryStore.selectedTagSlug);
  if (route.path !== '/') router.push('/');
}

function resetFilters() {
  categoryStore.setCategory('');
  categoryStore.setTag('');
  docStore.fetchDocuments();
}
</script>
