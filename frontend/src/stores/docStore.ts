import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { DocumentSummary, DocumentDetail, SearchResult, ToastMessage } from '@/types';

export const useDocStore = defineStore('docStore', () => {
  const documents = ref<DocumentSummary[]>([]);
  const activeDoc = ref<DocumentDetail | null>(null);
  const searchResults = ref<SearchResult[]>([]);
  const loading = ref<boolean>(false);
  const docLoading = ref<boolean>(false);
  const searchLoading = ref<boolean>(false);
  const isSearchOpen = ref<boolean>(false);
  const isSidebarOpen = ref<boolean>(false);
  const toasts = ref<ToastMessage[]>([]);

  function addToast(type: 'success' | 'error' | 'info', title: string, message?: string) {
    const id = Math.random().toString(36).substring(2, 9);
    toasts.value.push({ id, type, title, message });
    setTimeout(() => {
      removeToast(id);
    }, 4000);
  }

  function removeToast(id: string) {
    toasts.value = toasts.value.filter((t) => t.id !== id);
  }

  async function fetchDocuments(categorySlug?: string, tagSlug?: string, searchQuery?: string) {
    loading.value = true;
    try {
      const params = new URLSearchParams();
      if (categorySlug) params.append('category', categorySlug);
      if (tagSlug) params.append('tag', tagSlug);
      if (searchQuery) params.append('q', searchQuery);

      const res = await fetch(`/api/documents?${params.toString()}`);
      if (res.ok) {
        documents.value = await res.json();
      }
    } catch (e) {
      console.error('Failed to fetch documents:', e);
      addToast('error', 'Network Error', 'Unable to load documents from backend server.');
    } finally {
      loading.value = false;
    }
  }

  async function fetchDocumentBySlug(identifier: string) {
    docLoading.value = true;
    try {
      const res = await fetch(`/api/documents/${identifier}`);
      if (res.ok) {
        activeDoc.value = await res.json();
      } else {
        activeDoc.value = null;
      }
    } catch (e) {
      console.error('Failed to fetch document:', e);
      activeDoc.value = null;
      addToast('error', 'Error', 'Failed to retrieve requested document.');
    } finally {
      docLoading.value = false;
    }
  }

  async function searchDocs(query: string) {
    if (!query.trim()) {
      searchResults.value = [];
      return;
    }
    searchLoading.value = true;
    try {
      const res = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
      if (res.ok) {
        searchResults.value = await res.json();
      }
    } catch (e) {
      console.error('Search failed:', e);
    } finally {
      searchLoading.value = false;
    }
  }

  async function createDocument(payload: any) {
    try {
      const res = await fetch('/api/documents', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      if (res.ok) {
        const doc = await res.json();
        addToast('success', 'Document Created', `"${doc.title}" has been saved.`);
        return doc;
      } else {
        const err = await res.json();
        addToast('error', 'Failed to Create Document', err.error || 'Unknown error');
      }
    } catch (e) {
      addToast('error', 'Network Error', 'Could not reach server.');
    }
    return null;
  }

  async function updateDocument(id: string, payload: any) {
    try {
      const res = await fetch(`/api/documents/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });
      if (res.ok) {
        const doc = await res.json();
        activeDoc.value = doc;
        addToast('success', 'Document Saved', `Changes to "${doc.title}" were updated.`);
        return doc;
      } else {
        const err = await res.json();
        addToast('error', 'Failed to Update Document', err.error || 'Unknown error');
      }
    } catch (e) {
      addToast('error', 'Network Error', 'Could not reach server.');
    }
    return null;
  }

  async function deleteDocument(id: string) {
    try {
      const res = await fetch(`/api/documents/${id}`, { method: 'DELETE' });
      if (res.ok) {
        documents.value = documents.value.filter((d) => d.id !== id);
        if (activeDoc.value?.id === id) {
          activeDoc.value = null;
        }
        addToast('success', 'Document Deleted', 'The document was permanently removed.');
        return true;
      }
    } catch (e) {
      addToast('error', 'Failed to Delete Document');
    }
    return false;
  }

  async function batchUpload(items: any[]) {
    try {
      const res = await fetch('/api/documents/batch', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(items),
      });
      if (res.ok) {
        const data = await res.json();
        addToast('success', 'Batch Upload Complete', `Successfully imported ${data.uploaded_count} Markdown files.`);
        await fetchDocuments();
        return true;
      }
    } catch (e) {
      addToast('error', 'Batch Upload Failed');
    }
    return false;
  }

  return {
    documents,
    activeDoc,
    searchResults,
    loading,
    docLoading,
    searchLoading,
    isSearchOpen,
    isSidebarOpen,
    toasts,
    addToast,
    removeToast,
    fetchDocuments,
    fetchDocumentBySlug,
    searchDocs,
    createDocument,
    updateDocument,
    deleteDocument,
    batchUpload,
  };
});
