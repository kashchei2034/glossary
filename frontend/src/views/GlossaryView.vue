<template>
  <AppLayout>
    <div class="flex-1 flex flex-col min-h-0 min-w-0">
      <!-- Loading State -->
      <div v-if="docStore.docLoading" class="p-6 md:p-10 max-w-4xl mx-auto w-full">
        <SkeletonLoader />
      </div>

      <!-- Active Document View -->
      <div v-else-if="docStore.activeDoc" class="flex-1 grid grid-cols-1 lg:grid-cols-12 min-h-0 min-w-0">
        <!-- Main Document Body (Col 1 to 9) -->
        <div class="lg:col-span-9 p-6 md:p-10 overflow-y-auto min-w-0 flex justify-center">
          <div class="max-w-prose w-full space-y-6">
            <!-- Breadcrumbs / Category Badge -->
            <div class="flex flex-wrap items-center justify-between gap-2 border-b border-slate-200 dark:border-slate-800 pb-4">
              <div class="flex items-center gap-2 text-xs font-medium text-slate-500">
                <span v-if="docStore.activeDoc.category" class="px-2.5 py-1 rounded-lg bg-slate-100 dark:bg-slate-800 text-brand-600 dark:text-brand-400 font-semibold flex items-center gap-1.5">
                  <Icon name="folder" class="w-3.5 h-3.5" />
                  {{ docStore.activeDoc.category.name }}
                </span>
                <span>•</span>
                <span>Updated {{ formatDate(docStore.activeDoc.updated_at) }}</span>
              </div>

              <!-- Action Edit Button -->
              <router-link
                :to="`/editor/${docStore.activeDoc.id}`"
                class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-xl border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 transition-colors shadow-xs"
              >
                <Icon name="edit" class="w-3.5 h-3.5" />
                <span>Edit Article</span>
              </router-link>
            </div>

            <!-- Title & Tags -->
            <div class="space-y-3">
              <h1 class="text-3xl sm:text-4xl font-extrabold tracking-tight text-slate-900 dark:text-white">
                {{ docStore.activeDoc.title }}
              </h1>
              <div v-if="docStore.activeDoc.tags && docStore.activeDoc.tags.length > 0" class="flex flex-wrap gap-1.5">
                <span
                  v-for="t in docStore.activeDoc.tags"
                  :key="t.id"
                  class="px-2 py-0.5 rounded-md bg-slate-100 dark:bg-slate-800 text-xs text-slate-600 dark:text-slate-400 font-medium"
                >
                  #{{ t.name }}
                </span>
              </div>
            </div>

            <!-- Rendered Markdown Content -->
            <div class="pt-4">
              <MarkdownRenderer :content="docStore.activeDoc.content_md" />
            </div>
          </div>
        </div>

        <!-- Table of Contents Sticky Sidebar (Col 10 to 12) -->
        <div class="hidden lg:block lg:col-span-3 p-6 border-l border-slate-200 dark:border-slate-800 overflow-y-auto bg-slate-50/50 dark:bg-slate-950/20">
          <div class="sticky top-6">
            <TableOfContents :content="docStore.activeDoc.content_md" />
          </div>
        </div>
      </div>

      <!-- Document List Grid View (when no single doc active or browsing category) -->
      <div v-else class="p-6 md:p-10 overflow-y-auto max-w-5xl mx-auto w-full space-y-6">
        <div class="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-4">
          <div>
            <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Knowledge Hub & Instructions</h2>
            <p class="text-xs text-slate-500 mt-1">Browse instructional articles, system glossaries, and developer guides.</p>
          </div>
          <span class="text-xs font-semibold text-slate-400">{{ docStore.documents.length }} Documents</span>
        </div>

        <!-- Document Cards Grid -->
        <div v-if="docStore.documents.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="doc in docStore.documents"
            :key="doc.id"
            @click="router.push(`/doc/${doc.slug}`)"
            class="p-5 rounded-2xl bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 hover:border-brand-500/50 dark:hover:border-brand-500/50 shadow-xs hover:shadow-md transition-all cursor-pointer flex flex-col justify-between group"
          >
            <div class="space-y-2">
              <div class="flex items-center justify-between text-xs text-slate-400">
                <span v-if="doc.category" class="px-2 py-0.5 rounded-md bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 font-medium">
                  {{ doc.category.name }}
                </span>
                <span class="text-[11px]">{{ formatDate(doc.updated_at) }}</span>
              </div>
              <h3 class="text-base font-bold text-slate-900 dark:text-white group-hover:text-brand-500 transition-colors">
                {{ doc.title }}
              </h3>
              <p class="text-xs text-slate-500 dark:text-slate-400 line-clamp-3 leading-relaxed">
                {{ doc.snippet || 'No preview available.' }}
              </p>
            </div>

            <div v-if="doc.tags && doc.tags.length > 0" class="flex flex-wrap gap-1 mt-4 pt-3 border-t border-slate-100 dark:border-slate-800/60">
              <span v-for="t in doc.tags" :key="t.id" class="text-[10px] text-slate-400">#{{ t.name }}</span>
            </div>
          </div>
        </div>

        <EmptyState v-else title="No Documents Found" description="Try clearing filters or search queries to view available guides.">
          <router-link
            to="/editor/new"
            class="px-4 py-2 rounded-xl bg-brand-500 text-white text-xs font-bold shadow-md hover:bg-brand-600 transition-all"
          >
            Create First Document
          </router-link>
        </EmptyState>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { watch, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import AppLayout from '@/layouts/AppLayout.vue';
import MarkdownRenderer from '@/components/MarkdownRenderer.vue';
import TableOfContents from '@/components/TableOfContents.vue';
import SkeletonLoader from '@/components/SkeletonLoader.vue';
import EmptyState from '@/components/EmptyState.vue';
import { useDocStore } from '@/stores/docStore';
import Icon from '@/components/Icon.vue';

const route = useRoute();
const router = useRouter();
const docStore = useDocStore();

function formatDate(dateStr: string) {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
}

function loadContent() {
  const slug = route.params.slug as string;
  if (slug) {
    docStore.fetchDocumentBySlug(slug);
  } else {
    docStore.activeDoc = null;
    docStore.fetchDocuments();
  }
}

watch(() => route.params.slug, loadContent);

onMounted(loadContent);
</script>
