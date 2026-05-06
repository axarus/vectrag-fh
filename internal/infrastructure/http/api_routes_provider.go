package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/axarus/vectrag/internal/application"
)

type APIRoutesProvider struct{}

func (APIRoutesProvider) Register(mux *http.ServeMux) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	projectRoot, err := application.FindProjectRoot(wd)
	if err != nil {
		return err
	}

	modelsAPI, err := NewModelsAPI(projectRoot)
	if err != nil {
		return err
	}
	modelsAPI.Register(mux)

	return nil
}
