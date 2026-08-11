package model

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a document category with potential nesting
type Category struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Icon      string     `json:"icon"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder int        `json:"sort_order"`
	CreatedAt time.Time  `json:"created_at"`
	Children  []Category `json:"children,omitempty"`
	DocCount  int        `json:"doc_count"`
}

// Tag represents a taxonomy label for documents
type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

// Document represents a knowledge base document
type Document struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	ContentMd   string     `json:"content_md"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	Category    *Category  `json:"category,omitempty"`
	Tags        []Tag      `json:"tags"`
	IsPublished bool       `json:"is_published"`
	Priority    int        `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DocumentSummary is a lightweight version for list views
type DocumentSummary struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	Category    *Category  `json:"category,omitempty"`
	Tags        []Tag      `json:"tags"`
	IsPublished bool       `json:"is_published"`
	Priority    int        `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Snippet     string     `json:"snippet,omitempty"`
}

// CreateDocumentPayload struct for incoming requests
type CreateDocumentPayload struct {
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	ContentMd   string      `json:"content_md"`
	CategoryID  *uuid.UUID  `json:"category_id"`
	TagIDs      []uuid.UUID `json:"tag_ids"`
	IsPublished bool        `json:"is_published"`
	Priority    int         `json:"priority"`
}

// UpdateDocumentPayload struct for updating documents
type UpdateDocumentPayload struct {
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	ContentMd   string      `json:"content_md"`
	CategoryID  *uuid.UUID  `json:"category_id"`
	TagIDs      []uuid.UUID `json:"tag_ids"`
	IsPublished bool        `json:"is_published"`
	Priority    int         `json:"priority"`
}

// CreateCategoryPayload struct
type CreateCategoryPayload struct {
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Icon      string     `json:"icon"`
	ParentID  *uuid.UUID `json:"parent_id"`
	SortOrder int        `json:"sort_order"`
}

// CreateTagPayload struct
type CreateTagPayload struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// SearchResult item
type SearchResult struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Snippet   string    `json:"snippet"`
	Rank      float64   `json:"rank"`
	Category  string    `json:"category,omitempty"`
	Tags      []string  `json:"tags"`
}

// BatchUploadItem for uploading multiple md files at once
type BatchUploadItem struct {
	Title      string   `json:"title"`
	ContentMd  string   `json:"content_md"`
	Category   string   `json:"category,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}
