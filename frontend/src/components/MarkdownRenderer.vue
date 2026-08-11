<template>
  <div
    ref="containerRef"
    class="prose prose-slate dark:prose-invert max-w-prose w-full break-words leading-relaxed"
    v-html="parsedHtml"
  ></div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUpdated, nextTick } from 'vue';
import { Marked } from 'marked';
import hljs from 'highlight.js';

const props = defineProps<{
  content: string;
}>();

const containerRef = ref<HTMLElement | null>(null);

function slugifyHeader(text: string): string {
  return text
    .toLowerCase()
    .replace(/<[^>]*>/g, '')
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-');
}

// Configure Marked with custom renderer
const marked = new Marked({
  renderer: {
    heading(text: string, level: number) {
      const id = slugifyHeader(text);
      return `<h${level} id="${id}" class="group scroll-mt-20 relative">
        <a href="#${id}" class="no-underline text-inherit group-hover:text-brand-500 transition-colors">${text}</a>
      </h${level}>`;
    },
    code(code: string, language: string | undefined) {
      const validLang = language && hljs.getLanguage(language) ? language : 'plaintext';
      const highlighted = hljs.highlight(code, { language: validLang }).value;
      return `<div class="relative group my-6">
        <button class="copy-code-btn absolute top-3 right-3 px-2.5 py-1 text-xs font-mono rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 border border-slate-700/60 opacity-0 group-hover:opacity-100 transition-all duration-200 flex items-center gap-1.5 z-10" data-code="${encodeURIComponent(code)}">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 002-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/></svg>
          <span>Copy</span>
        </button>
        <pre><code class="hljs language-${validLang}">${highlighted}</code></pre>
      </div>`;
    },
    blockquote(quote: string) {
      // Check for Callout pattern: > [!NOTE], > [!TIP], > [!WARNING], > [!IMPORTANT], > [!CAUTION]
      const calloutMatch = quote.match(/^<p>\[\!(NOTE|TIP|WARNING|IMPORTANT|CAUTION)\]\s*/i);
      if (calloutMatch) {
        const type = calloutMatch[1].toLowerCase();
        const cleanContent = quote.replace(/^<p>\[\!(NOTE|TIP|WARNING|IMPORTANT|CAUTION)\]\s*/i, '<p>');
        return `<div class="callout-box callout-${type}">${cleanContent}</div>`;
      }
      return `<blockquote>${quote}</blockquote>`;
    },
    table(header: string, body: string) {
      return `<div class="overflow-x-auto my-6 rounded-xl border border-slate-200 dark:border-slate-800">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-200">${header}</thead>
          <tbody class="divide-y divide-slate-200 dark:divide-slate-800">${body}</tbody>
        </table>
      </div>`;
    },
  },
});

const parsedHtml = computed(() => {
  if (!props.content) return '';
  return marked.parse(props.content) as string;
});

function bindCopyButtons() {
  if (!containerRef.value) return;
  const buttons = containerRef.value.querySelectorAll('.copy-code-btn');
  buttons.forEach((btn) => {
    btn.addEventListener('click', (e) => {
      e.preventDefault();
      const target = e.currentTarget as HTMLButtonElement;
      const code = decodeURIComponent(target.getAttribute('data-code') || '');
      navigator.clipboard.writeText(code).then(() => {
        const label = target.querySelector('span');
        if (label) label.textContent = 'Copied!';
        target.classList.add('text-emerald-400');
        setTimeout(() => {
          if (label) label.textContent = 'Copy';
          target.classList.remove('text-emerald-400');
        }, 2000);
      });
    });
  });
}

onMounted(() => {
  nextTick(bindCopyButtons);
});

onUpdated(() => {
  nextTick(bindCopyButtons);
});
</script>
