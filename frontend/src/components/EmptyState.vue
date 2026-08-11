<template>
  <div class="flex flex-col items-center justify-center py-16 px-4 text-center">
    <div class="w-16 h-16 rounded-2xl bg-slate-100 dark:bg-slate-800 text-slate-400 dark:text-slate-500 flex items-center justify-center mb-4 shadow-sm">
      <Icon name="file" class="w-8 h-8" />
    </div>
    <h3 class="text-lg font-semibold text-slate-800 dark:text-slate-200 mb-1">{{ title }}</h3>
    <p class="text-sm text-slate-500 dark:text-slate-400 max-w-sm mb-6">{{ description }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
import Icon from '@/components/Icon.vue';

withDefaults(
  defineProps<{
    title?: string;
    description?: string;
  }>(),
  {
    title: 'No Documents Found',
    description: 'There are no documents matching your criteria yet.',
  }
);
</script>
