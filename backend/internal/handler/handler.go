package handler

import (
	"encoding/json"
	"net/http"

	"kb-backend/internal/model"
	"kb-backend/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// GetDocuments handler
func (h *Handler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	categorySlug := r.URL.Query().Get("category")
	tagSlug := r.URL.Query().Get("tag")
	search := r.URL.Query().Get("q")
	publishedOnly := r.URL.Query().Get("published") != "false"

	docs, err := h.repo.GetDocuments(r.Context(), categorySlug, tagSlug, search, publishedOnly)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, docs)
}

// GetDocumentByIDOrSlug handler
func (h *Handler) GetDocumentByIDOrSlug(w http.ResponseWriter, r *http.Request) {
	identifier := chi.URLParam(r, "identifier")
	if identifier == "" {
		respondError(w, http.StatusBadRequest, "Document identifier required")
		return
	}

	doc, err := h.repo.GetDocumentByIDOrSlug(r.Context(), identifier)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if doc == nil {
		respondError(w, http.StatusNotFound, "Document not found")
		return
	}

	respondJSON(w, http.StatusOK, doc)
}

// CreateDocument handler
func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var payload model.CreateDocumentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if payload.Title == "" || payload.ContentMd == "" {
		respondError(w, http.StatusBadRequest, "Title and ContentMd are required")
		return
	}

	doc, err := h.repo.CreateDocument(r.Context(), payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, doc)
}

// UpdateDocument handler
func (h *Handler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Document ID")
		return
	}

	var payload model.UpdateDocumentPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	doc, err := h.repo.UpdateDocument(r.Context(), id, payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, doc)
}

// DeleteDocument handler
func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid Document ID")
		return
	}

	if err := h.repo.DeleteDocument(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Document deleted successfully"})
}

// BatchUpload handler
func (h *Handler) BatchUpload(w http.ResponseWriter, r *http.Request) {
	var items []model.BatchUploadItem
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid payload array")
		return
	}

	count, err := h.repo.BatchUpload(r.Context(), items)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Batch upload completed",
		"uploaded_count": count,
	})
}

// GetCategories handler
func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.repo.GetCategoryTree(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, cats)
}

// CreateCategory handler
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var payload model.CreateCategoryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if payload.Name == "" {
		respondError(w, http.StatusBadRequest, "Category name is required")
		return
	}

	cat, err := h.repo.CreateCategory(r.Context(), payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, cat)
}

// GetTags handler
func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.repo.GetTags(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, tags)
}

// CreateTag handler
func (h *Handler) CreateTag(w http.ResponseWriter, r *http.Request) {
	var payload model.CreateTagPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if payload.Name == "" {
		respondError(w, http.StatusBadRequest, "Tag name is required")
		return
	}

	tag, err := h.repo.CreateTag(r.Context(), payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, tag)
}

// SearchDocuments handler
func (h *Handler) SearchDocuments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	results, err := h.repo.SearchDocuments(r.Context(), query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, results)
}
