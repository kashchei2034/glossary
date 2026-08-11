<template>
  <div class="fixed bottom-5 right-5 z-50 flex flex-col gap-2 max-w-sm w-full pointer-events-none">
    <TransitionGroup name="toast">
      <div
        v-for="t in toasts"
        :key="t.id"
        class="pointer-events-auto p-4 rounded-xl shadow-lg border flex items-start gap-3 transition-all duration-300"
        :class="{
          'bg-emerald-50 dark:bg-emerald-950/90 border-emerald-300 dark:border-emerald-700 text-emerald-900 dark:text-emerald-100': t.type === 'success',
          'bg-rose-50 dark:bg-rose-950/90 border-rose-300 dark:border-rose-700 text-rose-900 dark:text-rose-100': t.type === 'error',
          'bg-sky-50 dark:bg-sky-950/90 border-sky-300 dark:border-sky-700 text-sky-900 dark:text-sky-100': t.type === 'info',
        }"
      >
        <Icon :name="getIconName(t.type)" class="w-5 h-5 mt-0.5 shrink-0" />
        <div class="flex-1 min-w-0">
          <h4 class="text-sm font-semibold leading-snug">{{ t.title }}</h4>
          <p v-if="t.message" class="text-xs opacity-90 mt-0.5 break-words">{{ t.message }}</p>
        </div>
        <button
          @click="docStore.removeToast(t.id)"
          class="p-1 rounded-lg hover:bg-black/5 dark:hover:bg-white/10 opacity-60 hover:opacity-100 transition-opacity"
        >
          <Icon name="x" class="w-4 h-4" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useDocStore } from '@/stores/docStore';
import Icon from '@/components/Icon.vue';

const docStore = useDocStore();
const toasts = computed(() => docStore.toasts);

function getIconName(type: string) {
  if (type === 'success') return 'check';
  if (type === 'error') return 'alert';
  return 'info';
}
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-enter-from {
  opacity: 0;
  transform: translateY(1rem) scale(0.95);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(2rem);
}
</style>
