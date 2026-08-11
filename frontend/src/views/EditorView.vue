<template>
  <AppLayout>
    <div class="flex-1 flex flex-col h-full min-h-0 min-w-0 p-3 sm:p-4 space-y-3">
      <!-- Top Action Bar -->
      <div class="flex flex-wrap items-center justify-between gap-2 border-b border-slate-200 dark:border-slate-800 pb-3 shrink-0">
        <div class="flex items-center gap-2.5">
          <button @click="router.back()" class="p-2 rounded-xl border border-slate-200 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-500 transition-colors">
            <Icon name="arrowLeft" class="w-4 h-4" />
          </button>
          <div>
            <h2 class="text-base sm:text-lg font-bold text-slate-900 dark:text-white leading-snug">
              {{ isEditMode ? 'Edit Document' : 'Create New Document' }}
            </h2>
            <p class="text-[11px] text-slate-400">Markdown Editor Workspace</p>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <!-- Slide-Over Metadata Settings Button -->
          <button
            @click="isMetadataDrawerOpen = true"
            class="px-3 py-1.5 sm:py-2 rounded-xl text-xs font-semibold border border-slate-200 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 transition-colors flex items-center gap-1.5 shadow-xs"
          >
            <Icon name="sliders" class="w-4 h-4 text-brand-500" />
            <span class="hidden xs:inline">Document Settings</span>
            <span class="xs:hidden">Settings</span>
          </button>

          <!-- Save Document CTA -->
          <button
            @click="saveDocument"
            :disabled="saving"
            class="px-4 sm:px-5 py-1.5 sm:py-2 rounded-xl bg-brand-600 hover:bg-brand-700 text-white text-xs font-bold shadow-md shadow-brand-500/20 transition-all flex items-center gap-1.5 disabled:opacity-50"
          >
            <Icon name="save" class="w-4 h-4" />
            <span>{{ saving ? 'Saving...' : 'Save' }}</span>
          </button>
        </div>
      </div>

      <!-- Main Full-Width Editor Area (100% width & 100% height) -->
      <div class="flex-1 flex flex-col min-h-0 min-w-0 space-y-2.5">
        <!-- Title Input -->
        <input
          v-model="title"
          type="text"
          placeholder="Document Title (e.g., How to Configure REST API CORS)"
          class="w-full px-4 py-2.5 text-base sm:text-lg font-bold bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-xl focus:outline-none focus:ring-2 focus:ring-brand-500 text-slate-900 dark:text-white placeholder-slate-400 shrink-0 shadow-xs"
        />

        <!-- Live Markdown Workspace Container (100% Flex-1 height) -->
        <div class="flex-1 min-h-0 min-w-0">
          <MarkdownEditor v-model="contentMd" />
        </div>
      </div>

      <!-- Document Settings Slide-Over Drawer Modal -->
      <Teleport to="body">
        <Transition name="slide-over">
          <div
            v-if="isMetadataDrawerOpen"
            class="fixed inset-0 z-50 flex justify-end bg-slate-950/60 backdrop-blur-xs"
            @click.self="isMetadataDrawerOpen = false"
          >
            <div class="w-full max-w-sm h-full bg-white dark:bg-slate-900 border-l border-slate-200 dark:border-slate-800 p-5 space-y-5 overflow-y-auto shadow-2xl flex flex-col justify-between">
              <div class="space-y-5">
                <div class="flex items-center justify-between border-b border-slate-200 dark:border-slate-800 pb-3">
                  <h3 class="text-sm font-bold text-slate-900 dark:text-white flex items-center gap-2">
                    <Icon name="sliders" class="w-4 h-4 text-brand-500" />
                    <span>Document Settings</span>
                  </h3>
                  <button
                    @click="isMetadataDrawerOpen = false"
                    class="p-1 rounded-lg text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
                  >
                    <Icon name="x" class="w-5 h-5" />
                  </button>
                </div>

                <!-- Category Select -->
                <div class="space-y-1.5">
                  <label class="text-xs font-semibold text-slate-500">Category</label>
                  <select
                    v-model="selectedCategoryId"
                    class="w-full px-3 py-2 text-xs bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-900 dark:text-white focus:ring-2 focus:ring-brand-500"
                  >
                    <option :value="null">Uncategorized</option>
                    <option v-for="c in categoryStore.categories" :key="c.id" :value="c.id">{{ c.name }}</option>
                  </select>
                </div>

                <!-- Tags Multi Select -->
                <div class="space-y-1.5">
                  <label class="text-xs font-semibold text-slate-500">Tags</label>
                  <div class="flex flex-wrap gap-1 border border-slate-200 dark:border-slate-800 p-2 rounded-xl bg-slate-50/50 dark:bg-slate-800/40">
                    <button
                      v-for="t in categoryStore.tags"
                      :key="t.id"
                      @click="toggleTag(t.id)"
                      type="button"
                      class="px-2 py-1 rounded-lg text-[11px] font-medium transition-all"
                      :class="selectedTagIds.includes(t.id) ? 'bg-brand-500 text-white' : 'bg-slate-200 dark:bg-slate-700 text-slate-700 dark:text-slate-300'"
                    >
                      #{{ t.name }}
                    </button>
                  </div>
                </div>

                <!-- URL Slug -->
                <div class="space-y-1.5">
                  <label class="text-xs font-semibold text-slate-500">URL Slug</label>
                  <input
                    v-model="slug"
                    type="text"
                    placeholder="auto-generated-slug"
                    class="w-full px-3 py-2 text-xs font-mono bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-900 dark:text-white"
                  />
                </div>

                <!-- Priority Order -->
                <div class="space-y-1.5">
                  <label class="text-xs font-semibold text-slate-500">Sort Priority (Order)</label>
                  <input
                    v-model.number="priority"
                    type="number"
                    class="w-full px-3 py-2 text-xs bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-slate-900 dark:text-white"
                  />
                </div>

                <!-- Published Switch -->
                <div class="flex items-center justify-between pt-3 border-t border-slate-200 dark:border-slate-800">
                  <span class="text-xs font-semibold text-slate-700 dark:text-slate-300">Publish Document</span>
                  <button
                    @click="isPublished = !isPublished"
                    type="button"
                    class="w-11 h-6 rounded-full p-1 transition-colors duration-200 flex items-center"
                    :class="isPublished ? 'bg-brand-500 justify-end' : 'bg-slate-300 dark:bg-slate-700 justify-start'"
                  >
                    <div class="w-4 h-4 rounded-full bg-white shadow-md"></div>
                  </button>
                </div>
              </div>

              <!-- Drawer Footer Apply Button -->
              <button
                @click="isMetadataDrawerOpen = false"
                class="w-full py-2.5 rounded-xl bg-brand-500 hover:bg-brand-600 text-white text-xs font-bold shadow-md transition-all"
              >
                Apply Settings
              </button>
            </div>
          </div>
        </Transition>
      </Teleport>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import AppLayout from '@/layouts/AppLayout.vue';
import MarkdownEditor from '@/components/MarkdownEditor.vue';
import Icon from '@/components/Icon.vue';
import { useDocStore } from '@/stores/docStore';
import { useCategoryStore } from '@/stores/categoryStore';

const route = useRoute();
const router = useRouter();
const docStore = useDocStore();
const categoryStore = useCategoryStore();

const isEditMode = computed(() => !!route.params.id && route.params.id !== 'new');

const title = ref('');
const contentMd = ref('# Sample Title\n\nWrite content here...');
const selectedCategoryId = ref<string | null>(null);
const selectedTagIds = ref<string[]>([]);
const slug = ref('');
const priority = ref(0);
const isPublished = ref(true);
const saving = ref(false);
const isMetadataDrawerOpen = ref(false);

function toggleTag(tagId: string) {
  if (selectedTagIds.value.includes(tagId)) {
    selectedTagIds.value = selectedTagIds.value.filter((id) => id !== tagId);
  } else {
    selectedTagIds.value.push(tagId);
  }
}

async function loadExistingDoc() {
  if (isEditMode.value) {
    const docId = route.params.id as string;
    await docStore.fetchDocumentBySlug(docId);
    if (docStore.activeDoc) {
      const doc = docStore.activeDoc;
      title.value = doc.title;
      contentMd.value = doc.content_md;
      selectedCategoryId.value = doc.category_id || null;
      selectedTagIds.value = doc.tags ? doc.tags.map((t) => t.id) : [];
      slug.value = doc.slug;
      priority.value = doc.priority || 0;
      isPublished.value = doc.is_published;
    }
  }
}

async function saveDocument() {
  if (!title.value.trim()) {
    docStore.addToast('error', 'Validation Error', 'Document Title cannot be empty.');
    return;
  }

  saving.value = true;
  const payload = {
    title: title.value,
    slug: slug.value,
    content_md: contentMd.value,
    category_id: selectedCategoryId.value,
    tag_ids: selectedTagIds.value,
    is_published: isPublished.value,
    priority: priority.value,
  };

  let savedDoc = null;
  if (isEditMode.value) {
    savedDoc = await docStore.updateDocument(route.params.id as string, payload);
  } else {
    savedDoc = await docStore.createDocument(payload);
  }

  saving.value = false;
  if (savedDoc) {
    router.push(`/doc/${savedDoc.slug}`);
  }
}

onMounted(() => {
  categoryStore.fetchCategories();
  categoryStore.fetchTags();
  loadExistingDoc();
});
</script>

<style scoped>
.slide-over-enter-active,
.slide-over-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.slide-over-enter-from,
.slide-over-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
