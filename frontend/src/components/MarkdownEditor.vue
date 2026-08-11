<template>
  <div class="flex flex-col h-full bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl overflow-hidden shadow-sm">
    <!-- Toolbar -->
    <div class="flex flex-wrap items-center justify-between gap-2 p-3 bg-slate-50 dark:bg-slate-950 border-b border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300">
      <div class="flex flex-wrap items-center gap-1 text-xs font-medium">
        <!-- Text Styles -->
        <button @click="insertText('# ')" title="Heading 1" class="btn-toolbar"><span class="font-bold text-xs">H1</span></button>
        <button @click="insertText('## ')" title="Heading 2" class="btn-toolbar"><span class="font-bold text-xs">H2</span></button>
        <button @click="insertText('### ')" title="Heading 3" class="btn-toolbar"><span class="font-bold text-xs">H3</span></button>
        <div class="w-px h-4 bg-slate-300 dark:bg-slate-700 mx-1"></div>
        <button @click="wrapText('**', '**')" title="Bold" class="btn-toolbar"><Icon name="bold" class="w-4 h-4" /></button>
        <button @click="wrapText('*', '*')" title="Italic" class="btn-toolbar"><Icon name="italic" class="w-4 h-4" /></button>
        <button @click="wrapText('~~', '~~')" title="Strikethrough" class="btn-toolbar"><Icon name="strikethrough" class="w-4 h-4" /></button>
        <div class="w-px h-4 bg-slate-300 dark:bg-slate-700 mx-1"></div>

        <!-- Lists & Quote -->
        <button @click="insertText('- ')" title="Bullet List" class="btn-toolbar"><Icon name="list" class="w-4 h-4" /></button>
        <button @click="insertText('1. ')" title="Numbered List" class="btn-toolbar"><Icon name="listOrdered" class="w-4 h-4" /></button>
        <button @click="insertText('> ')" title="Quote" class="btn-toolbar"><Icon name="quote" class="w-4 h-4" /></button>
        <button @click="wrapText('`', '`')" title="Inline Code" class="btn-toolbar"><Icon name="code" class="w-4 h-4" /></button>
        <button @click="insertText('\n```go\n// code block\n```\n')" title="Code Block" class="btn-toolbar"><Icon name="code" class="w-4 h-4" /></button>

        <div class="w-px h-4 bg-slate-300 dark:bg-slate-700 mx-1"></div>

        <!-- Callouts & Table -->
        <div class="relative group">
          <button title="Insert Callout Box" class="btn-toolbar flex items-center gap-1">
            <Icon name="alert" class="w-4 h-4 text-brand-500" />
            <Icon name="chevronDown" class="w-3 h-3" />
          </button>
          <div class="absolute left-0 top-full mt-1 hidden group-hover:flex flex-col bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl shadow-xl p-1 z-30 w-36">
            <button @click="insertCallout('NOTE')" class="px-2 py-1.5 text-xs text-left hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg text-blue-600 dark:text-blue-400 font-semibold">Note</button>
            <button @click="insertCallout('TIP')" class="px-2 py-1.5 text-xs text-left hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg text-emerald-600 dark:text-emerald-400 font-semibold">Tip</button>
            <button @click="insertCallout('WARNING')" class="px-2 py-1.5 text-xs text-left hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg text-amber-600 dark:text-amber-400 font-semibold">Warning</button>
            <button @click="insertCallout('IMPORTANT')" class="px-2 py-1.5 text-xs text-left hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg text-purple-600 dark:text-purple-400 font-semibold">Important</button>
          </div>
        </div>

        <button @click="insertTable" title="Insert Table" class="btn-toolbar"><Icon name="table" class="w-4 h-4" /></button>
        <button @click="wrapText('[', '](https://)')" title="Insert Link" class="btn-toolbar"><Icon name="link" class="w-4 h-4" /></button>
        <button @click="wrapText('![Alt text](', ')')" title="Insert Image" class="btn-toolbar"><Icon name="image" class="w-4 h-4" /></button>
      </div>

      <!-- View Mode Selector -->
      <div class="flex items-center gap-1 bg-slate-200 dark:bg-slate-800 p-1 rounded-xl">
        <button
          @click="viewMode = 'edit'"
          class="px-2.5 py-1 text-xs rounded-lg font-medium transition-all"
          :class="viewMode === 'edit' ? 'bg-white dark:bg-slate-700 shadow text-slate-900 dark:text-white' : 'text-slate-500 hover:text-slate-800 dark:text-slate-400'"
        >
          Edit
        </button>
        <button
          @click="viewMode = 'split'"
          class="px-2.5 py-1 text-xs rounded-lg font-medium transition-all hidden sm:block"
          :class="viewMode === 'split' ? 'bg-white dark:bg-slate-700 shadow text-slate-900 dark:text-white' : 'text-slate-500 hover:text-slate-800 dark:text-slate-400'"
        >
          Split
        </button>
        <button
          @click="viewMode = 'preview'"
          class="px-2.5 py-1 text-xs rounded-lg font-medium transition-all"
          :class="viewMode === 'preview' ? 'bg-white dark:bg-slate-700 shadow text-slate-900 dark:text-white' : 'text-slate-500 hover:text-slate-800 dark:text-slate-400'"
        >
          Preview
        </button>
      </div>
    </div>

    <!-- Workspace Body -->
    <div class="flex-1 grid grid-cols-1 min-h-0 overflow-hidden" :class="{ 'sm:grid-cols-2': viewMode === 'split' }">
      <!-- Editor Area -->
      <div v-show="viewMode !== 'preview'" class="flex flex-col h-full min-h-0 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800">
        <textarea
          ref="textareaRef"
          :value="modelValue"
          @input="onInput"
          placeholder="Write your Markdown content here..."
          class="flex-1 w-full p-4 font-mono text-sm bg-transparent border-0 resize-none focus:outline-none focus:ring-0 text-slate-900 dark:text-slate-100 placeholder-slate-400 dark:placeholder-slate-600 leading-relaxed overflow-y-auto"
        ></textarea>
        <!-- Footer Stats -->
        <div class="px-4 py-2 bg-slate-50 dark:bg-slate-950 border-t border-slate-200 dark:border-slate-800 text-xs text-slate-400 dark:text-slate-500 flex justify-between">
          <span>{{ wordCount }} words | {{ charCount }} characters</span>
          <span>Markdown Mode</span>
        </div>
      </div>

      <!-- Preview Area -->
      <div v-show="viewMode !== 'edit'" class="h-full overflow-y-auto p-6 bg-slate-50/50 dark:bg-slate-900/50">
        <MarkdownRenderer :content="modelValue" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import MarkdownRenderer from '@/components/MarkdownRenderer.vue';
import Icon from '@/components/Icon.vue';

const props = defineProps<{
  modelValue: string;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', val: string): void;
}>();

const textareaRef = ref<HTMLTextAreaElement | null>(null);
const viewMode = ref<'edit' | 'split' | 'preview'>('split');

const wordCount = computed(() => {
  if (!props.modelValue.trim()) return 0;
  return props.modelValue.trim().split(/\s+/).length;
});

const charCount = computed(() => props.modelValue.length);

function onInput(e: Event) {
  const target = e.target as HTMLTextAreaElement;
  emit('update:modelValue', target.value);
}

function insertText(text: string) {
  if (!textareaRef.value) return;
  const textarea = textareaRef.value;
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  const current = props.modelValue;
  const updated = current.substring(0, start) + text + current.substring(end);
  emit('update:modelValue', updated);

  setTimeout(() => {
    textarea.focus();
    textarea.setSelectionRange(start + text.length, start + text.length);
  }, 10);
}

function wrapText(prefix: string, suffix: string) {
  if (!textareaRef.value) return;
  const textarea = textareaRef.value;
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  const current = props.modelValue;
  const selected = current.substring(start, end) || 'text';
  const updated = current.substring(0, start) + prefix + selected + suffix + current.substring(end);
  emit('update:modelValue', updated);

  setTimeout(() => {
    textarea.focus();
    textarea.setSelectionRange(start + prefix.length, end + prefix.length);
  }, 10);
}

function insertCallout(type: string) {
  insertText(`\n> [!${type}]\n> Write callout content here...\n\n`);
}

function insertTable() {
  insertText(`\n| Column 1 | Column 2 | Column 3 |\n| --- | --- | --- |\n| Item 1 | Detail 1 | Value 1 |\n| Item 2 | Detail 2 | Value 2 |\n\n`);
}
</script>

<style scoped>
.btn-toolbar {
  @apply p-1.5 rounded-lg hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors;
}
</style>
