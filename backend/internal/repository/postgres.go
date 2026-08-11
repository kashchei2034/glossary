package repository

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"kb-backend/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool

	// In-memory fallback fields
	mu         sync.RWMutex
	categories map[uuid.UUID]*model.Category
	tags       map[uuid.UUID]*model.Tag
	documents  map[uuid.UUID]*model.Document
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	repo := &Repository{
		pool:       pool,
		categories: make(map[uuid.UUID]*model.Category),
		tags:       make(map[uuid.UUID]*model.Tag),
		documents:  make(map[uuid.UUID]*model.Document),
	}

	if pool == nil {
		repo.seedInMemoryData()
	}

	return repo
}

func (r *Repository) IsInMemory() bool {
	return r.pool == nil
}

func slugify(text string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	slug := reg.ReplaceAllString(strings.ToLower(text), "-")
	return strings.Trim(slug, "-")
}

func (r *Repository) seedInMemoryData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Seed Categories
	cat1ID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	cat2ID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	cat3ID := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	cat4ID := uuid.MustParse("44444444-4444-4444-4444-444444444444")

	r.categories[cat1ID] = &model.Category{ID: cat1ID, Name: "Getting Started", Slug: "getting-started", Icon: "book-open", SortOrder: 1, CreatedAt: time.Now()}
	r.categories[cat2ID] = &model.Category{ID: cat2ID, Name: "Developer Guides", Slug: "developer-guides", Icon: "code", SortOrder: 2, CreatedAt: time.Now()}
	r.categories[cat3ID] = &model.Category{ID: cat3ID, Name: "Backend & DB", Slug: "backend-db", Icon: "database", ParentID: &cat2ID, SortOrder: 1, CreatedAt: time.Now()}
	r.categories[cat4ID] = &model.Category{ID: cat4ID, Name: "Frontend & Vue", Slug: "frontend-vue", Icon: "layout", ParentID: &cat2ID, SortOrder: 2, CreatedAt: time.Now()}

	// Seed Tags
	tag1ID := uuid.MustParse("a1111111-1111-1111-1111-111111111111")
	tag2ID := uuid.MustParse("a2222222-2222-2222-2222-222222222222")
	tag3ID := uuid.MustParse("a3333333-3333-3333-3333-333333333333")
	tag4ID := uuid.MustParse("a4444444-4444-4444-4444-444444444444")

	t1 := model.Tag{ID: tag1ID, Name: "Go", Slug: "go"}
	t2 := model.Tag{ID: tag2ID, Name: "Vue.js", Slug: "vuejs"}
	t3 := model.Tag{ID: tag3ID, Name: "PostgreSQL", Slug: "postgresql"}
	t4 := model.Tag{ID: tag4ID, Name: "Architecture", Slug: "architecture"}

	r.tags[tag1ID] = &t1
	r.tags[tag2ID] = &t2
	r.tags[tag3ID] = &t3
	r.tags[tag4ID] = &t4

	// Seed Documents
	doc1ID := uuid.MustParse("d1111111-1111-1111-1111-111111111111")
	doc2ID := uuid.MustParse("d2222222-2222-2222-2222-222222222222")
	doc3ID := uuid.MustParse("d3333333-3333-3333-3333-333333333333")

	r.documents[doc1ID] = &model.Document{
		ID:          doc1ID,
		Title:       "Welcome to Knowledge Base",
		Slug:        "welcome-to-knowledge-base",
		ContentMd:   "# Welcome to Knowledge Base Hub\n\nThis instructions hub and personal knowledge base is designed for quick lookup, clear architectural documentation, and Markdown-first editing.\n\n## Features Overview\n\n- **Markdown-First Reader**: Crisp distraction-free reading mode with automatically generated Table of Contents.\n- **Fast Search**: Built with full-text query matching.\n- **Admin Workspace**: Built-in side-by-side Markdown editor with live preview and metadata drawer.\n- **Unbreakable Layout**: Responsive CSS grid & flex layout.\n\n> [!TIP]\n> Press `Ctrl + K` or `Cmd + K` anytime to open global command palette!\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, Knowledge Base!\")\n}\n```\n",
		CategoryID:  &cat1ID,
		Category:    r.categories[cat1ID],
		Tags:        []model.Tag{t4},
		IsPublished: true,
		Priority:    1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	r.documents[doc2ID] = &model.Document{
		ID:          doc2ID,
		Title:       "Go Clean Architecture & API Design",
		Slug:        "go-clean-architecture",
		ContentMd:   "# Go Clean Architecture & REST Guidelines\n\nThis document details the architectural conventions used in our Go backend service.\n\n## Directory Layout\n\n```\n/cmd/server/       # Application entrypoints\n/internal/\n  /handler/        # HTTP handlers\n  /repository/     # Database queries\n  /model/          # Domain models\n```\n\n> [!IMPORTANT]\n> Always sanitize input before querying database or search indexes.\n",
		CategoryID:  &cat3ID,
		Category:    r.categories[cat3ID],
		Tags:        []model.Tag{t1, t3},
		IsPublished: true,
		Priority:    2,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	r.documents[doc3ID] = &model.Document{
		ID:          doc3ID,
		Title:       "Vue 3 Composition API & State Management",
		Slug:        "vue-3-composition-api",
		ContentMd:   "# Vue 3 Composition API Guidelines\n\nWe adopt Vue 3 `<script setup>` syntax with TypeScript and Pinia for reactive store management.\n\n```typescript\nimport { defineStore } from 'pinia';\nimport { ref } from 'vue';\n\nexport const useDocStore = defineStore('docStore', () => {\n  const docs = ref([]);\n  return { docs };\n});\n```\n",
		CategoryID:  &cat4ID,
		Category:    r.categories[cat4ID],
		Tags:        []model.Tag{t2},
		IsPublished: true,
		Priority:    3,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// GetCategoryTree fetches categories and structures them into a parent-child tree
func (r *Repository) GetCategoryTree(ctx context.Context) ([]model.Category, error) {
	if r.pool == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()

		allCats := make([]model.Category, 0, len(r.categories))
		catMap := make(map[uuid.UUID]*model.Category)

		// Calculate doc counts
		docCounts := make(map[uuid.UUID]int)
		for _, doc := range r.documents {
			if doc.CategoryID != nil {
				docCounts[*doc.CategoryID]++
			}
		}

		for _, cat := range r.categories {
			cCopy := *cat
			cCopy.DocCount = docCounts[cat.ID]
			cCopy.Children = []model.Category{}
			allCats = append(allCats, cCopy)
		}

		for i := range allCats {
			catMap[allCats[i].ID] = &allCats[i]
		}

		var rootCats []model.Category
		for i := range allCats {
			cat := catMap[allCats[i].ID]
			if cat.ParentID != nil {
				if parent, exists := catMap[*cat.ParentID]; exists {
					parent.Children = append(parent.Children, *cat)
				} else {
					rootCats = append(rootCats, *cat)
				}
			} else {
				rootCats = append(rootCats, *cat)
			}
		}

		sort.Slice(rootCats, func(i, j int) bool {
			return rootCats[i].SortOrder < rootCats[j].SortOrder
		})

		return rootCats, nil
	}

	query := `
		SELECT 
			c.id, c.name, c.slug, COALESCE(c.icon, 'folder'), c.parent_id, c.sort_order, c.created_at,
			COUNT(d.id) as doc_count
		FROM categories c
		LEFT JOIN documents d ON d.category_id = c.id
		GROUP BY c.id, c.name, c.slug, c.icon, c.parent_id, c.sort_order, c.created_at
		ORDER BY c.sort_order ASC, c.name ASC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var allCats []model.Category
	for rows.Next() {
		var cat model.Category
		err := rows.Scan(
			&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.ParentID, &cat.SortOrder, &cat.CreatedAt, &cat.DocCount,
		)
		if err != nil {
			return nil, err
		}
		cat.Children = []model.Category{}
		allCats = append(allCats, cat)
	}

	catMap := make(map[uuid.UUID]*model.Category)
	for i := range allCats {
		catMap[allCats[i].ID] = &allCats[i]
	}

	var rootCats []model.Category
	for i := range allCats {
		cat := catMap[allCats[i].ID]
		if cat.ParentID != nil {
			if parent, exists := catMap[*cat.ParentID]; exists {
				parent.Children = append(parent.Children, *cat)
			} else {
				rootCats = append(rootCats, *cat)
			}
		} else {
			rootCats = append(rootCats, *cat)
		}
	}

	return rootCats, nil
}

// CreateCategory adds a new category
func (r *Repository) CreateCategory(ctx context.Context, payload model.CreateCategoryPayload) (*model.Category, error) {
	if payload.Slug == "" {
		payload.Slug = slugify(payload.Name)
	}
	if payload.Icon == "" {
		payload.Icon = "folder"
	}

	if r.pool == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		catID := uuid.New()
		cat := &model.Category{
			ID:        catID,
			Name:      payload.Name,
			Slug:      payload.Slug,
			Icon:      payload.Icon,
			ParentID:  payload.ParentID,
			SortOrder: payload.SortOrder,
			CreatedAt: time.Now(),
			Children:  []model.Category{},
		}
		r.categories[catID] = cat
		return cat, nil
	}

	catID := uuid.New()
	now := time.Now()
	query := `
		INSERT INTO categories (id, name, slug, icon, parent_id, sort_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, slug, icon, parent_id, sort_order, created_at
	`

	var cat model.Category
	err := r.pool.QueryRow(ctx, query, catID, payload.Name, payload.Slug, payload.Icon, payload.ParentID, payload.SortOrder, now).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.ParentID, &cat.SortOrder, &cat.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	cat.Children = []model.Category{}
	return &cat, nil
}

// GetTags retrieves all tags
func (r *Repository) GetTags(ctx context.Context) ([]model.Tag, error) {
	if r.pool == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()

		tags := make([]model.Tag, 0, len(r.tags))
		for _, t := range r.tags {
			tags = append(tags, *t)
		}
		sort.Slice(tags, func(i, j int) bool {
			return tags[i].Name < tags[j].Name
		})
		return tags, nil
	}

	query := `SELECT id, name, slug FROM tags ORDER BY name ASC`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags: %w", err)
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// CreateTag creates a new tag
func (r *Repository) CreateTag(ctx context.Context, payload model.CreateTagPayload) (*model.Tag, error) {
	if payload.Slug == "" {
		payload.Slug = slugify(payload.Name)
	}

	if r.pool == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		for _, t := range r.tags {
			if t.Slug == payload.Slug {
				t.Name = payload.Name
				return t, nil
			}
		}

		tagID := uuid.New()
		tag := &model.Tag{ID: tagID, Name: payload.Name, Slug: payload.Slug}
		r.tags[tagID] = tag
		return tag, nil
	}

	tagID := uuid.New()
	query := `
		INSERT INTO tags (id, name, slug)
		VALUES ($1, $2, $3)
		ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, name, slug
	`
	var tag model.Tag
	err := r.pool.QueryRow(ctx, query, tagID, payload.Name, payload.Slug).Scan(&tag.ID, &tag.Name, &tag.Slug)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	return &tag, nil
}

// GetDocuments queries list of documents
func (r *Repository) GetDocuments(ctx context.Context, categorySlug, tagSlug, search string, publishedOnly bool) ([]model.DocumentSummary, error) {
	if r.pool == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()

		var results []model.DocumentSummary
		for _, doc := range r.documents {
			if publishedOnly && !doc.IsPublished {
				continue
			}

			if categorySlug != "" && (doc.Category == nil || doc.Category.Slug != categorySlug) {
				continue
			}

			if tagSlug != "" {
				hasTag := false
				for _, t := range doc.Tags {
					if t.Slug == tagSlug {
						hasTag = true
						break
					}
				}
				if !hasTag {
					continue
				}
			}

			if search != "" {
				sLower := strings.ToLower(search)
				tLower := strings.ToLower(doc.Title)
				cLower := strings.ToLower(doc.ContentMd)
				if !strings.Contains(tLower, sLower) && !strings.Contains(cLower, sLower) {
					continue
				}
			}

			snippet := doc.ContentMd
			if len(snippet) > 200 {
				snippet = snippet[:200]
			}

			results = append(results, model.DocumentSummary{
				ID:          doc.ID,
				Title:       doc.Title,
				Slug:        doc.Slug,
				CategoryID:  doc.CategoryID,
				Category:    doc.Category,
				Tags:        doc.Tags,
				IsPublished: doc.IsPublished,
				Priority:    doc.Priority,
				CreatedAt:   doc.CreatedAt,
				UpdatedAt:   doc.UpdatedAt,
				Snippet:     snippet,
			})
		}

		sort.Slice(results, func(i, j int) bool {
			return results[i].Priority < results[j].Priority
		})

		return results, nil
	}

	var conditions []string
	var args []interface{}
	argIdx := 1

	if publishedOnly {
		conditions = append(conditions, fmt.Sprintf("d.is_published = $%d", argIdx))
		args = append(args, true)
		argIdx++
	}

	if categorySlug != "" {
		conditions = append(conditions, fmt.Sprintf("c.slug = $%d", argIdx))
		args = append(args, categorySlug)
		argIdx++
	}

	if tagSlug != "" {
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM document_tags dt JOIN tags t ON dt.tag_id = t.id WHERE dt.document_id = d.id AND t.slug = $%d)", argIdx))
		args = append(args, tagSlug)
		argIdx++
	}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(d.title ILIKE $%d OR d.content_md ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			d.id, d.title, d.slug, d.category_id, d.is_published, d.priority, d.created_at, d.updated_at,
			c.id, c.name, c.slug, COALESCE(c.icon, 'folder'),
			SUBSTRING(d.content_md FROM 1 FOR 200) as snippet
		FROM documents d
		LEFT JOIN categories c ON d.category_id = c.id
		%s
		ORDER BY d.priority ASC, d.updated_at DESC
	`, whereClause)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	var docs []model.DocumentSummary
	for rows.Next() {
		var doc model.DocumentSummary
		var catID, cID *uuid.UUID
		var cName, cSlug, cIcon *string

		err := rows.Scan(
			&doc.ID, &doc.Title, &doc.Slug, &catID, &doc.IsPublished, &doc.Priority, &doc.CreatedAt, &doc.UpdatedAt,
			&cID, &cName, &cSlug, &cIcon, &doc.Snippet,
		)
		if err != nil {
			return nil, err
		}

		doc.CategoryID = catID
		if cID != nil {
			doc.Category = &model.Category{
				ID:   *cID,
				Name: *cName,
				Slug: *cSlug,
				Icon: *cIcon,
			}
		}

		tags, _ := r.getTagsForDocument(ctx, doc.ID)
		doc.Tags = tags
		docs = append(docs, doc)
	}

	return docs, nil
}

// GetDocumentByIDOrSlug fetches a single document
func (r *Repository) GetDocumentByIDOrSlug(ctx context.Context, identifier string) (*model.Document, error) {
	if r.pool == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()

		for _, doc := range r.documents {
			if doc.ID.String() == identifier || doc.Slug == identifier {
				return doc, nil
			}
		}
		return nil, nil
	}

	docUUID, parseErr := uuid.Parse(identifier)

	var query string
	var arg interface{}
	if parseErr == nil {
		query = `
			SELECT d.id, d.title, d.slug, d.content_md, d.category_id, d.is_published, d.priority, d.created_at, d.updated_at,
			       c.id, c.name, c.slug, COALESCE(c.icon, 'folder')
			FROM documents d
			LEFT JOIN categories c ON d.category_id = c.id
			WHERE d.id = $1
		`
		arg = docUUID
	} else {
		query = `
			SELECT d.id, d.title, d.slug, d.content_md, d.category_id, d.is_published, d.priority, d.created_at, d.updated_at,
			       c.id, c.name, c.slug, COALESCE(c.icon, 'folder')
			FROM documents d
			LEFT JOIN categories c ON d.category_id = c.id
			WHERE d.slug = $1
		`
		arg = identifier
	}

	var doc model.Document
	var catID, cID *uuid.UUID
	var cName, cSlug, cIcon *string

	err := r.pool.QueryRow(ctx, query, arg).Scan(
		&doc.ID, &doc.Title, &doc.Slug, &doc.ContentMd, &catID, &doc.IsPublished, &doc.Priority, &doc.CreatedAt, &doc.UpdatedAt,
		&cID, &cName, &cSlug, &cIcon,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch document: %w", err)
	}

	doc.CategoryID = catID
	if cID != nil {
		doc.Category = &model.Category{
			ID:   *cID,
			Name: *cName,
			Slug: *cSlug,
			Icon: *cIcon,
		}
	}

	tags, err := r.getTagsForDocument(ctx, doc.ID)
	if err == nil {
		doc.Tags = tags
	} else {
		doc.Tags = []model.Tag{}
	}

	return &doc, nil
}

// CreateDocument inserts a new document
func (r *Repository) CreateDocument(ctx context.Context, payload model.CreateDocumentPayload) (*model.Document, error) {
	if payload.Slug == "" {
		payload.Slug = slugify(payload.Title)
	}

	if r.pool == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		docID := uuid.New()
		now := time.Now()

		var cat *model.Category
		if payload.CategoryID != nil {
			cat = r.categories[*payload.CategoryID]
		}

		tags := make([]model.Tag, 0, len(payload.TagIDs))
		for _, tID := range payload.TagIDs {
			if t, exists := r.tags[tID]; exists {
				tags = append(tags, *t)
			}
		}

		doc := &model.Document{
			ID:          docID,
			Title:       payload.Title,
			Slug:        payload.Slug,
			ContentMd:   payload.ContentMd,
			CategoryID:  payload.CategoryID,
			Category:    cat,
			Tags:        tags,
			IsPublished: payload.IsPublished,
			Priority:    payload.Priority,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		r.documents[docID] = doc
		return doc, nil
	}

	docID := uuid.New()
	now := time.Now()

	query := `
		INSERT INTO documents (id, title, slug, content_md, category_id, is_published, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, title, slug, content_md, category_id, is_published, priority, created_at, updated_at
	`

	var doc model.Document
	err := r.pool.QueryRow(ctx, query, docID, payload.Title, payload.Slug, payload.ContentMd, payload.CategoryID, payload.IsPublished, payload.Priority, now, now).Scan(
		&doc.ID, &doc.Title, &doc.Slug, &doc.ContentMd, &doc.CategoryID, &doc.IsPublished, &doc.Priority, &doc.CreatedAt, &doc.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}

	if len(payload.TagIDs) > 0 {
		for _, tagID := range payload.TagIDs {
			_, _ = r.pool.Exec(ctx, "INSERT INTO document_tags (document_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", doc.ID, tagID)
		}
	}

	return r.GetDocumentByIDOrSlug(ctx, doc.ID.String())
}

// UpdateDocument updates an existing document
func (r *Repository) UpdateDocument(ctx context.Context, id uuid.UUID, payload model.UpdateDocumentPayload) (*model.Document, error) {
	if payload.Slug == "" {
		payload.Slug = slugify(payload.Title)
	}

	if r.pool == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		doc, exists := r.documents[id]
		if !exists {
			return nil, fmt.Errorf("document not found")
		}

		doc.Title = payload.Title
		doc.Slug = payload.Slug
		doc.ContentMd = payload.ContentMd
		doc.CategoryID = payload.CategoryID
		if payload.CategoryID != nil {
			doc.Category = r.categories[*payload.CategoryID]
		} else {
			doc.Category = nil
		}

		tags := make([]model.Tag, 0, len(payload.TagIDs))
		for _, tID := range payload.TagIDs {
			if t, ok := r.tags[tID]; ok {
				tags = append(tags, *t)
			}
		}
		doc.Tags = tags
		doc.IsPublished = payload.IsPublished
		doc.Priority = payload.Priority
		doc.UpdatedAt = time.Now()

		return doc, nil
	}

	query := `
		UPDATE documents
		SET title = $1, slug = $2, content_md = $3, category_id = $4, is_published = $5, priority = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`

	_, err := r.pool.Exec(ctx, query, payload.Title, payload.Slug, payload.ContentMd, payload.CategoryID, payload.IsPublished, payload.Priority, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	_, _ = r.pool.Exec(ctx, "DELETE FROM document_tags WHERE document_id = $1", id)
	if len(payload.TagIDs) > 0 {
		for _, tagID := range payload.TagIDs {
			_, _ = r.pool.Exec(ctx, "INSERT INTO document_tags (document_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", id, tagID)
		}
	}

	return r.GetDocumentByIDOrSlug(ctx, id.String())
}

// DeleteDocument removes a document
func (r *Repository) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	if r.pool == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.documents, id)
		return nil
	}

	_, err := r.pool.Exec(ctx, "DELETE FROM documents WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

// SearchDocuments performs search
func (r *Repository) SearchDocuments(ctx context.Context, queryStr string) ([]model.SearchResult, error) {
	if strings.TrimSpace(queryStr) == "" {
		return []model.SearchResult{}, nil
	}

	if r.pool == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()

		qLower := strings.ToLower(queryStr)
		var results []model.SearchResult
		for _, doc := range r.documents {
			tLower := strings.ToLower(doc.Title)
			cLower := strings.ToLower(doc.ContentMd)
			if strings.Contains(tLower, qLower) || strings.Contains(cLower, qLower) {
				snippet := doc.ContentMd
				if len(snippet) > 150 {
					snippet = snippet[:150]
				}

				tagNames := make([]string, 0, len(doc.Tags))
				for _, t := range doc.Tags {
					tagNames = append(tagNames, t.Name)
				}

				catName := ""
				if doc.Category != nil {
					catName = doc.Category.Name
				}

				results = append(results, model.SearchResult{
					ID:       doc.ID,
					Title:    doc.Title,
					Slug:     doc.Slug,
					Snippet:  snippet,
					Rank:     1.0,
					Category: catName,
					Tags:     tagNames,
				})
			}
		}
		return results, nil
	}

	words := strings.Fields(queryStr)
	tsQueryStr := strings.Join(words, " & ")

	query := `
		SELECT 
			d.id, d.title, d.slug,
			ts_headline('english', d.content_md, to_tsquery('english', $1), 'StartSel=<mark>, StopSel=</mark>, MaxWords=35, MinWords=15') as snippet,
			ts_rank_cd(d.search_vector, to_tsquery('english', $1)) as rank,
			COALESCE(c.name, '') as category_name
		FROM documents d
		LEFT JOIN categories c ON d.category_id = c.id
		WHERE d.search_vector @@ to_tsquery('english', $1) OR d.title ILIKE $2
		ORDER BY rank DESC, d.updated_at DESC
		LIMIT 20
	`

	ilikeArg := "%" + queryStr + "%"
	rows, err := r.pool.Query(ctx, query, tsQueryStr, ilikeArg)
	if err != nil {
		fallbackQuery := `
			SELECT d.id, d.title, d.slug, SUBSTRING(d.content_md FROM 1 FOR 150) as snippet, 1.0 as rank, COALESCE(c.name, '')
			FROM documents d
			LEFT JOIN categories c ON d.category_id = c.id
			WHERE d.title ILIKE $1 OR d.content_md ILIKE $1
			LIMIT 20
		`
		var fbErr error
		rows, fbErr = r.pool.Query(ctx, fallbackQuery, ilikeArg)
		if fbErr != nil {
			return nil, fmt.Errorf("search query failed: %w", fbErr)
		}
	}
	defer rows.Close()

	var results []model.SearchResult
	for rows.Next() {
		var res model.SearchResult
		if err := rows.Scan(&res.ID, &res.Title, &res.Slug, &res.Snippet, &res.Rank, &res.Category); err != nil {
			return nil, err
		}
		tags, _ := r.getTagsForDocument(ctx, res.ID)
		for _, t := range tags {
			res.Tags = append(res.Tags, t.Name)
		}
		results = append(results, res)
	}

	return results, nil
}

// BatchUpload inserts multiple documents
func (r *Repository) BatchUpload(ctx context.Context, items []model.BatchUploadItem) (int, error) {
	count := 0
	for _, item := range items {
		if strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.ContentMd) == "" {
			continue
		}

		var catID *uuid.UUID
		if item.Category != "" {
			catSlug := slugify(item.Category)
			newCat, err := r.CreateCategory(ctx, model.CreateCategoryPayload{Name: item.Category, Slug: catSlug})
			if err == nil {
				catID = &newCat.ID
			}
		}

		tagIDs := []uuid.UUID{}
		for _, tagName := range item.Tags {
			tSlug := slugify(tagName)
			newTag, err := r.CreateTag(ctx, model.CreateTagPayload{Name: tagName, Slug: tSlug})
			if err == nil {
				tagIDs = append(tagIDs, newTag.ID)
			}
		}

		_, err := r.CreateDocument(ctx, model.CreateDocumentPayload{
			Title:       item.Title,
			Slug:        slugify(item.Title),
			ContentMd:   item.ContentMd,
			CategoryID:  catID,
			TagIDs:      tagIDs,
			IsPublished: true,
			Priority:    0,
		})
		if err == nil {
			count++
		}
	}
	return count, nil
}

func (r *Repository) getTagsForDocument(ctx context.Context, docID uuid.UUID) ([]model.Tag, error) {
	query := `
		SELECT t.id, t.name, t.slug
		FROM tags t
		JOIN document_tags dt ON dt.tag_id = t.id
		WHERE dt.document_id = $1
		ORDER BY t.name ASC
	`
	rows, err := r.pool.Query(ctx, query, docID)
	if err != nil {
		return []model.Tag{}, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Slug); err == nil {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}
