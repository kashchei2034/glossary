<template>
  <AppLayout>
    <div class="p-6 md:p-10 max-w-6xl mx-auto w-full space-y-6">
      <!-- Header Controls -->
      <div class="flex flex-wrap items-center justify-between gap-4 border-b border-slate-200 dark:border-slate-800 pb-4">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Document & Category Manager</h2>
          <p class="text-xs text-slate-500 mt-1">Manage documents, create categories, and perform batch Markdown imports.</p>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <button
            @click="isCategoryModalOpen = true"
            class="px-3 py-2 rounded-xl text-xs font-semibold border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 transition-all flex items-center gap-1.5"
          >
            <Icon name="folder-plus" class="w-4 h-4 text-brand-500" />
            <span>Add Category</span>
          </button>

          <button
            @click="isBatchModalOpen = true"
            class="px-3 py-2 rounded-xl text-xs font-semibold border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 transition-all flex items-center gap-1.5"
          >
            <Icon name="upload" class="w-4 h-4 text-brand-500" />
            <span>Mass Import .MD</span>
          </button>

          <router-link
            to="/editor/new"
            class="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md transition-all flex items-center gap-1.5"
          >
            <Icon name="plus" class="w-4 h-4" />
            <span>Create Document</span>
          </router-link>
        </div>
      </div>

      <!-- Filters & Search Bar -->
      <div class="flex flex-wrap items-center justify-between gap-3 bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800">
        <div class="flex-1 max-w-sm relative">
          <Icon name="search" class="w-4 h-4 absolute left-3 top-3 text-slate-400" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Filter documents..."
            class="w-full pl-9 pr-3 py-1.5 text-xs bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl focus:outline-none focus:ring-2 focus:ring-brand-500 text-slate-900 dark:text-white"
            @input="filterDocs"
          />
        </div>
        <div class="text-xs text-slate-500">
          Showing <b>{{ filteredDocs.length }}</b> documents
        </div>
      </div>

      <!-- Documents Table View -->
      <div class="bg-white dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 overflow-hidden shadow-xs">
        <div class="overflow-x-auto">
          <table class="w-full text-left text-xs">
            <thead class="bg-slate-50 dark:bg-slate-950 border-b border-slate-200 dark:border-slate-800 text-slate-500 dark:text-slate-400 font-semibold uppercase tracking-wider">
              <tr>
                <th class="p-4">Title & Slug</th>
                <th class="p-4">Category</th>
                <th class="p-4">Tags</th>
                <th class="p-4">Status</th>
                <th class="p-4">Last Updated</th>
                <th class="p-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-200 dark:divide-slate-800">
              <tr v-for="doc in filteredDocs" :key="doc.id" class="hover:bg-slate-50/50 dark:hover:bg-slate-800/50 transition-colors">
                <!-- Title & Slug -->
                <td class="p-4">
                  <router-link :to="`/doc/${doc.slug}`" class="font-bold text-sm text-slate-900 dark:text-white hover:text-brand-500 transition-colors block truncate max-w-xs">
                    {{ doc.title }}
                  </router-link>
                  <span class="text-[11px] text-slate-400 font-mono">{{ doc.slug }}</span>
                </td>

                <!-- Category -->
                <td class="p-4">
                  <span v-if="doc.category" class="px-2 py-1 rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 font-medium">
                    {{ doc.category.name }}
                  </span>
                  <span v-else class="text-slate-400 italic">Uncategorized</span>
                </td>

                <!-- Tags -->
                <td class="p-4">
                  <div v-if="doc.tags && doc.tags.length > 0" class="flex flex-wrap gap-1">
                    <span v-for="t in doc.tags" :key="t.id" class="px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-[10px] text-slate-500">
                      #{{ t.name }}
                    </span>
                  </div>
                  <span v-else class="text-slate-400">-</span>
                </td>

                <!-- Status -->
                <td class="p-4">
                  <span
                    class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider"
                    :class="doc.is_published ? 'bg-emerald-100 dark:bg-emerald-950 text-emerald-600 dark:text-emerald-400' : 'bg-amber-100 dark:bg-amber-950 text-amber-600 dark:text-amber-400'"
                  >
                    {{ doc.is_published ? 'Published' : 'Draft' }}
                  </span>
                </td>

                <!-- Date -->
                <td class="p-4 text-slate-500">
                  {{ formatDate(doc.updated_at) }}
                </td>

                <!-- Actions -->
                <td class="p-4 text-right">
                  <div class="flex items-center justify-end gap-1">
                    <router-link
                      :to="`/editor/${doc.id}`"
                      class="p-1.5 rounded-lg hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500 hover:text-slate-900 dark:hover:text-white"
                      title="Edit Document"
                    >
                      <Icon name="edit" class="w-4 h-4" />
                    </router-link>
                    <button
                      @click="confirmDelete(doc)"
                      class="p-1.5 rounded-lg hover:bg-rose-50 dark:hover:bg-rose-950 text-rose-500 transition-colors"
                      title="Delete Document"
                    >
                      <Icon name="trash" class="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Add Category Modal -->
      <Teleport to="body">
        <div v-if="isCategoryModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs">
          <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-6 w-full max-w-md space-y-4 shadow-2xl">
            <h3 class="text-lg font-bold text-slate-900 dark:text-white">Create New Category</h3>
            <div class="space-y-3">
              <div>
                <label class="text-xs font-semibold text-slate-500">Category Name</label>
                <input
                  v-model="newCategoryName"
                  type="text"
                  placeholder="e.g., API Specs"
                  class="w-full mt-1 px-3 py-2 text-sm bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl focus:ring-2 focus:ring-brand-500 text-slate-900 dark:text-white"
                />
              </div>
              <div>
                <label class="text-xs font-semibold text-slate-500">Parent Category (Optional)</label>
                <select
                  v-model="newCategoryParent"
                  class="w-full mt-1 px-3 py-2 text-sm bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-900 dark:text-white"
                >
                  <option value="">None (Top-Level Category)</option>
                  <option v-for="c in categoryStore.categories" :key="c.id" :value="c.id">{{ c.name }}</option>
                </select>
              </div>
            </div>
            <div class="flex justify-end gap-2 pt-2">
              <button @click="isCategoryModalOpen = false" class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800">
                Cancel
              </button>
              <button @click="submitCategory" class="px-4 py-2 rounded-xl text-xs font-bold bg-brand-500 hover:bg-brand-600 text-white shadow-md">
                Create Category
              </button>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Mass Import Markdown Modal -->
      <Teleport to="body">
        <div v-if="isBatchModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/60 backdrop-blur-xs">
          <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-6 w-full max-w-lg space-y-4 shadow-2xl">
            <h3 class="text-lg font-bold text-slate-900 dark:text-white">Mass Import .MD Files</h3>
            <p class="text-xs text-slate-500">Select one or multiple Markdown files from your computer to automatically ingest them into the knowledge base.</p>

            <div
              @dragover.prevent
              @drop.prevent="handleFileDrop"
              class="border-2 border-dashed border-slate-300 dark:border-slate-700 rounded-2xl p-8 text-center bg-slate-50/50 dark:bg-slate-950/50 hover:border-brand-500 transition-colors cursor-pointer"
              @click="$refs.fileInput.click()"
            >
              <Icon name="upload" class="w-10 h-10 text-brand-500 mx-auto mb-2" />
              <p class="text-sm font-semibold text-slate-700 dark:text-slate-300">Click or drop .md files here</p>
              <p class="text-xs text-slate-400 mt-1">{{ batchFiles.length }} files selected</p>
              <input ref="fileInput" type="file" multiple accept=".md,.markdown" class="hidden" @change="handleFileSelect" />
            </div>

            <div v-if="batchFiles.length > 0" class="max-h-32 overflow-y-auto space-y-1">
              <div v-for="(f, i) in batchFiles" :key="i" class="text-xs text-slate-600 dark:text-slate-400 flex justify-between bg-slate-100 dark:bg-slate-800 p-2 rounded-lg">
                <span class="truncate">{{ f.name }}</span>
                <span>{{ (f.size / 1024).toFixed(1) }} KB</span>
              </div>
            </div>

            <div class="flex justify-end gap-2 pt-2">
              <button @click="isBatchModalOpen = false" class="px-4 py-2 rounded-xl text-xs font-semibold text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800">
                Cancel
              </button>
              <button
                @click="submitBatchUpload"
                :disabled="batchFiles.length === 0"
                class="px-4 py-2 rounded-xl text-xs font-bold bg-brand-500 hover:bg-brand-600 disabled:opacity-50 text-white shadow-md"
              >
                Upload {{ batchFiles.length }} Files
              </button>
            </div>
          </div>
        </div>
      </Teleport>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import AppLayout from '@/layouts/AppLayout.vue';
import Icon from '@/components/Icon.vue';
import { useDocStore } from '@/stores/docStore';
import { useCategoryStore } from '@/stores/categoryStore';
import type { DocumentSummary } from '@/types';

const docStore = useDocStore();
const categoryStore = useCategoryStore();

const searchQuery = ref('');
const isCategoryModalOpen = ref(false);
const newCategoryName = ref('');
const newCategoryParent = ref('');

const isBatchModalOpen = ref(false);
const batchFiles = ref<File[]>([]);
const fileInput = ref<any>(null);

const filteredDocs = computed(() => {
  if (!searchQuery.value.trim()) return docStore.documents;
  const q = searchQuery.value.toLowerCase();
  return docStore.documents.filter(
    (d) => d.title.toLowerCase().includes(q) || d.slug.toLowerCase().includes(q)
  );
});

function formatDate(dStr: string) {
  if (!dStr) return '';
  return new Date(dStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
}

function filterDocs() {}

async function confirmDelete(doc: DocumentSummary) {
  if (confirm(`Are you sure you want to delete "${doc.title}"?`)) {
    await docStore.deleteDocument(doc.id);
  }
}

async function submitCategory() {
  if (!newCategoryName.value.trim()) return;
  const ok = await categoryStore.createCategory(newCategoryName.value, newCategoryParent.value || undefined);
  if (ok) {
    isCategoryModalOpen.value = false;
    newCategoryName.value = '';
    newCategoryParent.value = '';
  }
}

function handleFileSelect(e: Event) {
  const target = e.target as HTMLInputElement;
  if (target.files) {
    batchFiles.value = Array.from(target.files);
  }
}

function handleFileDrop(e: DragEvent) {
  if (e.dataTransfer?.files) {
    batchFiles.value = Array.from(e.dataTransfer.files).filter((f) => f.name.endsWith('.md') || f.name.endsWith('.markdown'));
  }
}

async function submitBatchUpload() {
  const items: any[] = [];
  for (const file of batchFiles.value) {
    const text = await file.text();
    const title = file.name.replace(/\.md$/i, '').replace(/[-_]/g, ' ');
    items.push({
      title: title.charAt(0).toUpperCase() + title.slice(1),
      content_md: text,
    });
  }

  const ok = await docStore.batchUpload(items);
  if (ok) {
    isBatchModalOpen.value = false;
    batchFiles.value = [];
  }
}

onMounted(() => {
  docStore.fetchDocuments('', '', '');
});
</script>
