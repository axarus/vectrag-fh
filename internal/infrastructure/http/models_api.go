package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/axarus/vectrag/internal/application"
	"github.com/axarus/vectrag/internal/domain"
	"github.com/axarus/vectrag/internal/infrastructure/filestore"
	"github.com/google/uuid"
)

type ModelsAPI struct {
	mu         sync.Mutex
	modelsDir  string
	modelSvc   *application.ModelService
	enableCORS bool
}

type CreateModelRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Status      string             `json:"status"`
	Fields      []CreateFieldInput `json:"fields"`
}

type CreateFieldInput struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Unique      bool   `json:"unique,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Status      string `json:"status"`
}

type UpdateModelRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Status      string             `json:"status"`
	Fields      []UpdateFieldInput `json:"fields"`
}

type UpdateFieldInput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Unique      bool   `json:"unique,omitempty"`
	Required    bool   `json:"required,omitempty"`
	Status      string `json:"status"`
}

func NewModelsAPI(projectRoot string) (*ModelsAPI, error) {
	cfg, err := application.LoadProjectConfig(projectRoot)
	if err != nil {
		return nil, err
	}
	modelsDir, err := application.ResolveModelsDir(projectRoot, cfg)
	if err != nil {
		return nil, err
	}

	repo, err := filestore.NewYamlRepository(modelsDir)
	if err != nil {
		return nil, err
	}

	return &ModelsAPI{
		modelsDir:  modelsDir,
		modelSvc:   application.NewModelService(repo),
		enableCORS: cfg.Development.EnableCORS,
	}, nil
}

func (api *ModelsAPI) Register(mux *http.ServeMux) {
	mux.Handle("/api/models", api)
	mux.Handle("/api/models/", api)
}

func (api *ModelsAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if api.enableCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "api/models" {
		switch r.Method {
		case http.MethodGet:
			api.handleList(w, r)
			return
		case http.MethodPost:
			api.handleCreate(w, r)
			return
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	}

	if strings.HasPrefix(path, "api/models/") {
		slug := strings.TrimPrefix(path, "api/models/")
		slug = strings.Trim(slug, "/")
		if slug == "" {
			writeError(w, http.StatusNotFound, "not found")
			return
		}

		switch r.Method {
		case http.MethodGet:
			api.handleGet(w, r, slug)
			return
		case http.MethodPut:
			api.handleUpdate(w, r, slug)
			return
		case http.MethodDelete:
			api.handleDelete(w, r, slug)
			return
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
	}

	writeError(w, http.StatusNotFound, "not found")
}

func (api *ModelsAPI) handleList(w http.ResponseWriter, r *http.Request) {
	api.mu.Lock()
	defer api.mu.Unlock()
	
	models, err := api.modelSvc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, models)
}

func (api *ModelsAPI) handleGet(w http.ResponseWriter, r *http.Request, slug string) {
	api.mu.Lock()
	defer api.mu.Unlock()

	model, err := api.modelSvc.Get(slug)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, model)
}

func (api *ModelsAPI) handleCreate(w http.ResponseWriter, r *http.Request) {
	api.mu.Lock()
	defer api.mu.Unlock()

	var req CreateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	slug, err := slugify(req.Name)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := os.Stat(filepath.Join(api.modelsDir, slug+".yaml")); err == nil {
		writeError(w, http.StatusConflict, "model already exists")
		return
	}

	modelID := newID()

	fields := make([]domain.Field, len(req.Fields))
	for i, f := range req.Fields {
		fieldID := newID()

		fields[i] = domain.Field{
			ID:          fieldID,
			Name:        f.Name,
			Type:        domain.FieldType(f.Type),
			Description: f.Description,
			Unique:      f.Unique,
			Required:    f.Required,
			Status:      domain.Status(f.Status),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
	}

	model := domain.Model{
		ID:            modelID,
		Name:          req.Name,
		Slug:          slug,
		Description:   req.Description,
		Fields:        fields,
		Status:        domain.Status(req.Status),
		SchemaVersion: 1,
	}

	if err := domain.ValidateModel(model); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := api.modelSvc.Create(model); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, model)
}

func (api *ModelsAPI) handleUpdate(w http.ResponseWriter, r *http.Request, slug string) {
	api.mu.Lock()
	defer api.mu.Unlock()

	existing, err := api.modelSvc.Get(slug)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	existingFieldByID := make(map[string]domain.Field, len(existing.Fields))
	for _, f := range existing.Fields {
		existingFieldByID[f.ID] = f
	}

	var req UpdateModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	fields := make([]domain.Field, len(req.Fields))
	for i, f := range req.Fields {
		fieldID := strings.TrimSpace(f.ID)
		if fieldID == "" {
			fieldID = newID()
		}

		createdAt := time.Time{}
		if prev, ok := existingFieldByID[fieldID]; ok {
			createdAt = prev.CreatedAt
		}
		if createdAt.IsZero() {
			createdAt = time.Now().UTC()
		}

		fields[i] = domain.Field{
			ID:          fieldID,
			Name:        f.Name,
			Type:        domain.FieldType(f.Type),
			Description: f.Description,
			Unique:      f.Unique,
			Required:    f.Required,
			Status:      domain.Status(f.Status),
			CreatedAt:   createdAt,
			UpdatedAt:   time.Now().UTC(),
		}
	}

	updated := domain.Model{
		ID:            existing.ID,
		Name:          req.Name,
		Slug:          existing.Slug,
		Description:   req.Description,
		Fields:        fields,
		Status:        domain.Status(req.Status),
		SchemaVersion: existing.SchemaVersion,
	}

	if err := domain.ValidateModel(updated); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := api.modelSvc.Update(updated); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (api *ModelsAPI) handleDelete(w http.ResponseWriter, r *http.Request, slug string) {
	api.mu.Lock()
	defer api.mu.Unlock()

	if err := api.modelSvc.Delete(slug); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": true})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": msg})
}

var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(name string) (string, error) {
	s := strings.TrimSpace(name)
	if s == "" {
		return "", fmt.Errorf("name cannot be empty")
	}

	s = strings.ToLower(s)
	s = slugRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "", fmt.Errorf("name results in empty slug")
	}

	m, _ := regexp.MatchString(`^[a-z0-9-]+$`, s)
	if !m {
		return "", fmt.Errorf("invalid slug")
	}

	return s, nil
}

func newID() string {
	id := uuid.New().String()
	return id 
}
