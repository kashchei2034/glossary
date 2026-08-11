export interface Tag {
  id: string;
  name: string;
  slug: string;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  icon: string;
  parent_id?: string | null;
  sort_order: number;
  created_at: string;
  children?: Category[];
  doc_count?: number;
}

export interface DocumentSummary {
  id: string;
  title: string;
  slug: string;
  category_id?: string | null;
  category?: Category | null;
  tags: Tag[];
  is_published: boolean;
  priority: number;
  created_at: string;
  updated_at: string;
  snippet?: string;
}

export interface DocumentDetail extends DocumentSummary {
  content_md: string;
}

export interface SearchResult {
  id: string;
  title: string;
  slug: string;
  snippet: string;
  rank: number;
  category?: string;
  tags?: string[];
}

export interface TocItem {
  id: string;
  text: string;
  level: number;
}

export interface ToastMessage {
  id: string;
  type: 'success' | 'error' | 'info';
  title: string;
  message?: string;
}
