-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    icon VARCHAR(100) DEFAULT 'folder',
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tags table
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE
);

-- Documents table
CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    content_md TEXT NOT NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    is_published BOOLEAN DEFAULT true,
    priority INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    search_vector TSVECTOR
);

-- Document Tags junction table
CREATE TABLE IF NOT EXISTS document_tags (
    document_id UUID REFERENCES documents(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (document_id, tag_id)
);

-- Indexes for performance & search
CREATE INDEX IF NOT EXISTS idx_documents_category_id ON documents(category_id);
CREATE INDEX IF NOT EXISTS idx_documents_slug ON documents(slug);
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);

-- Full-text search index
CREATE INDEX IF NOT EXISTS idx_documents_search_vector ON documents USING GIN(search_vector);

-- Function and trigger to update search_vector automatically
CREATE OR REPLACE FUNCTION documents_search_vector_trigger() RETURNS trigger AS $$
BEGIN
  NEW.search_vector :=
    setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(NEW.content_md, '')), 'B');
  NEW.updated_at := CURRENT_TIMESTAMP;
  RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_documents_search_vector ON documents;
CREATE TRIGGER trg_documents_search_vector
BEFORE INSERT OR UPDATE ON documents
FOR EACH ROW EXECUTE FUNCTION documents_search_vector_trigger();

-- SEED DATA
INSERT INTO categories (id, name, slug, icon, parent_id, sort_order) VALUES
  ('11111111-1111-1111-1111-111111111111', 'Getting Started', 'getting-started', 'book-open', NULL, 1),
  ('22222222-2222-2222-2222-222222222222', 'Developer Guides', 'developer-guides', 'code', NULL, 2),
  ('33333333-3333-3333-3333-333333333333', 'Backend & DB', 'backend-db', 'database', '22222222-2222-2222-2222-222222222222', 1),
  ('44444444-4444-4444-4444-444444444444', 'Frontend & Vue', 'frontend-vue', 'layout', '22222222-2222-2222-2222-222222222222', 2),
  ('55555555-5555-5555-5555-555555555555', 'DevOps & Deployment', 'devops', 'server', NULL, 3)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO tags (id, name, slug) VALUES
  ('a1111111-1111-1111-1111-111111111111', 'Go', 'go'),
  ('a2222222-2222-2222-2222-222222222222', 'Vue.js', 'vuejs'),
  ('a3333333-3333-3333-3333-333333333333', 'PostgreSQL', 'postgresql'),
  ('a4444444-4444-4444-4444-444444444444', 'Docker', 'docker'),
  ('a5555555-5555-5555-5555-555555555555', 'Architecture', 'architecture')
ON CONFLICT (slug) DO NOTHING;

INSERT INTO documents (id, title, slug, content_md, category_id, is_published, priority) VALUES
(
  'd1111111-1111-1111-1111-111111111111',
  'Welcome to Knowledge Base',
  'welcome-to-knowledge-base',
  '# Welcome to Knowledge Base Hub

This instructions hub and personal knowledge base is designed for quick lookup, clear architectural documentation, and Markdown-first editing.

## Features Overview

- **Markdown-First Reader**: Crisp distraction-free reading mode with automatically generated Table of Contents.
- **Fast Full-Text Search**: Built on PostgreSQL `TSVECTOR` and trigram matching.
- **Admin Workspace**: Built-in side-by-side Markdown editor with live preview and metadata drawer.
- **Unbreakable Layout**: Responsive CSS grid & flex layout designed for desktop, tablet, and mobile.

> [!TIP]
> Press `Ctrl + K` or `Cmd + K` anytime to open the global command palette and search through documents!

### Code Block Syntax Highlighting

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Knowledge Base!")
}
```

### Callout Boxes Support

> [!NOTE]
> Knowledge documents support custom callout notes, warnings, and code references.

---

### Quick Glossary Table

| Term | Definition | Primary Tech |
| --- | --- | --- |
| **REST API** | High-performance GoChi REST service | Go |
| **Full-Text Index** | TSVECTOR GIN index for text search | PostgreSQL |
| **Pinia Store** | Reactive Vue 3 state management | TypeScript |
',
  '11111111-1111-1111-1111-111111111111',
  true,
  1
),
(
  'd2222222-2222-2222-2222-222222222222',
  'Go Clean Architecture & API Design',
  'go-clean-architecture',
  '# Go Clean Architecture & REST Guidelines

This document details the architectural conventions used in our Go backend service.

## Directory Layout

```
/cmd/server/       # Application entrypoints
/internal/
  /handler/        # HTTP handlers & request parsing
  /repository/     # Database queries & transactions
  /model/          # Core domain models & DTOs
```

## PostgreSQL Full-Text Search Implementation

We utilize PostgreSQL `tsvector` and `tsquery` for instant search relevance ranking.

```sql
SELECT id, title, slug, ts_rank(search_vector, query) AS rank
FROM documents, to_tsquery(''english'', ''go & architecture'') query
WHERE search_vector @@ query
ORDER BY rank DESC;
```

> [!IMPORTANT]
> Always sanitize user input before passing it into `to_tsquery` or use `phraseto_tsquery`.
',
  '33333333-3333-3333-3333-333333333333',
  true,
  2
),
(
  'd3333333-3333-3333-3333-333333333333',
  'Vue 3 Composition API & State Management',
  'vue-3-composition-api',
  '# Vue 3 Composition API Guidelines

We adopt Vue 3 `<script setup>` syntax with TypeScript and Pinia for reactive store management.

## Store Pattern Example

```typescript
import { defineStore } from ''pinia'';
import { ref, computed } from ''vue'';

export const useDocStore = defineStore(''docStore'', () => {
  const documents = ref([]);
  const activeDoc = ref(null);

  async function fetchDocs() {
    const res = await fetch(''/api/documents'');
    documents.value = await res.json();
  }

  return { documents, activeDoc, fetchDocs };
});
```

## Dynamic Table of Contents Parsing

The table of contents is automatically extracted from Markdown heading tags (`h1`, `h2`, `h3`) using a DOM observer or AST parser.
',
  '44444444-4444-4444-4444-444444444444',
  true,
  3
)
ON CONFLICT (slug) DO NOTHING;

INSERT INTO document_tags (document_id, tag_id) VALUES
  ('d1111111-1111-1111-1111-111111111111', 'a5555555-5555-5555-5555-555555555555'),
  ('d2222222-2222-2222-2222-222222222222', 'a1111111-1111-1111-1111-111111111111'),
  ('d2222222-2222-2222-2222-222222222222', 'a3333333-3333-3333-3333-333333333333'),
  ('d3333333-3333-3333-3333-333333333333', 'a2222222-2222-2222-2222-222222222222')
ON CONFLICT DO NOTHING;
