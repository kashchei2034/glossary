<template>
  <nav v-if="tocItems.length > 0" class="space-y-1 text-sm">
    <div class="flex items-center gap-2 mb-3 text-xs font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-500">
      <Icon name="listTree" class="w-3.5 h-3.5" />
      <span>On This Page</span>
    </div>
    <ul class="space-y-1 border-l border-slate-200 dark:border-slate-800 font-medium">
      <li
        v-for="item in tocItems"
        :key="item.id"
        :style="{ paddingLeft: `${(item.level - 1) * 0.75 + 0.75}rem` }"
      >
        <a
          :href="`#${item.id}`"
          @click.prevent="scrollToHeading(item.id)"
          class="block py-1 text-xs truncate transition-colors border-l-2 -ml-[1px] pl-2"
          :class="
            activeId === item.id
              ? 'border-brand-500 text-brand-600 dark:text-brand-400 font-semibold'
              : 'border-transparent text-slate-500 hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-200'
          "
        >
          {{ item.text }}
        </a>
      </li>
    </ul>
  </nav>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import type { TocItem } from '@/types';
import Icon from '@/components/Icon.vue';

const props = defineProps<{
  content: string;
}>();

const activeId = ref<string>('');

const tocItems = computed<TocItem[]>(() => {
  if (!props.content) return [];
  const headingRegex = /^(#{1,3})\s+(.+)$/gm;
  const items: TocItem[] = [];
  let match;

  while ((match = headingRegex.exec(props.content)) !== null) {
    const level = match[1].length;
    const rawText = match[2].trim();
    // Clean markdown links/formatting from heading text
    const text = rawText.replace(/\[([^\]]+)\]\([^\)]+\)/g, '$1').replace(/[`*_*]/g, '');
    const id = text
      .toLowerCase()
      .replace(/[^\w\s-]/g, '')
      .replace(/\s+/g, '-');
    items.push({ id, text, level });
  }

  return items;
});

function scrollToHeading(id: string) {
  const el = document.getElementById(id);
  if (el) {
    activeId.value = id;
    el.scrollIntoView({ behavior: 'smooth', block: 'start' });
    history.replaceState(null, '', `#${id}`);
  }
}

let observer: IntersectionObserver | null = null;

function setupObserver() {
  if (observer) observer.disconnect();
  if (tocItems.value.length === 0) return;

  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          activeId.value = entry.target.id;
        }
      });
    },
    { rootMargin: '-80px 0px -40% 0px', threshold: 0.1 }
  );

  tocItems.value.forEach((item) => {
    const el = document.getElementById(item.id);
    if (el) observer?.observe(el);
  });
}

watch(
  () => props.content,
  () => {
    setTimeout(setupObserver, 200);
  }
);

onMounted(() => {
  setTimeout(setupObserver, 300);
});

onUnmounted(() => {
  if (observer) observer.disconnect();
});
</script>
