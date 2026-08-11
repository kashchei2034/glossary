import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { Category, Tag } from '@/types';

export const useCategoryStore = defineStore('categoryStore', () => {
  const categories = ref<Category[]>([]);
  const tags = ref<Tag[]>([]);
  const loading = ref<boolean>(false);
  const selectedCategorySlug = ref<string>('');
  const selectedTagSlug = ref<string>('');

  async function fetchCategories() {
    loading.value = true;
    try {
      const res = await fetch('/api/categories');
      if (res.ok) {
        categories.value = await res.json();
      }
    } catch (e) {
      console.error('Failed to fetch categories:', e);
    } finally {
      loading.value = false;
    }
  }

  async function fetchTags() {
    try {
      const res = await fetch('/api/tags');
      if (res.ok) {
        tags.value = await res.json();
      }
    } catch (e) {
      console.error('Failed to fetch tags:', e);
    }
  }

  async function createCategory(name: string, parentId?: string, icon?: string) {
    try {
      const res = await fetch('/api/categories', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, parent_id: parentId || null, icon: icon || 'folder' }),
      });
      if (res.ok) {
        await fetchCategories();
        return true;
      }
    } catch (e) {
      console.error('Failed to create category:', e);
    }
    return false;
  }

  async function createTag(name: string) {
    try {
      const res = await fetch('/api/tags', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name }),
      });
      if (res.ok) {
        await fetchTags();
        return true;
      }
    } catch (e) {
      console.error('Failed to create tag:', e);
    }
    return false;
  }

  function setCategory(slug: string) {
    selectedCategorySlug.value = slug;
  }

  function setTag(slug: string) {
    selectedTagSlug.value = slug;
  }

  return {
    categories,
    tags,
    loading,
    selectedCategorySlug,
    selectedTagSlug,
    fetchCategories,
    fetchTags,
    createCategory,
    createTag,
    setCategory,
    setTag,
  };
});
