package handlers

import (
	"ITLAFINAL/domain/models"
	"ITLAFINAL/domain/ports"
	"ITLAFINAL/domain/usecases/orderUseCases"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type fakeServiceTypeRepository struct {
	created *models.ServiceType
}

func (f *fakeServiceTypeRepository) Create(serviceType *models.ServiceType) error {
	f.created = serviceType
	return nil
}

func (f *fakeServiceTypeRepository) FindAll() ([]*models.ServiceType, error) {
	return nil, nil
}

func (f *fakeServiceTypeRepository) FindByName(name string) (*models.ServiceType, error) {
	return nil, nil
}

func (f *fakeServiceTypeRepository) Update(serviceType *models.ServiceType) error {
	return nil
}

func (f *fakeServiceTypeRepository) Delete(name string) error {
	return nil
}

var _ ports.ServiceTypeRepository = (*fakeServiceTypeRepository)(nil)

func TestCreateReturnsSuccessMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &fakeServiceTypeRepository{}
	createUseCase := orderUseCases.NewCreateServiceTypeUseCase(repo)
	handler := NewServiceTypeHandler(nil, createUseCase)

	payload := map[string]any{
		"name":             "Express",
		"description":      "Quick delivery",
		"base_price":       100.0,
		"price_per_weight": 2.5,
		"price_per_piece":  3.0,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/service-types", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	handler.Create(ctx)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid json response: %v", err)
	}

	if response["message"] != "service type created successfully" {
		t.Fatalf("expected success message, got %v", response)
	}

	if repo.created == nil {
		t.Fatal("expected service type to be created")
	}
	if repo.created.Name != "Express" {
		t.Fatalf("expected created service type name to be Express, got %s", repo.created.Name)
	}
	if repo.created.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be set")
	}
	if time.Since(repo.created.CreatedAt) < 0 {
		t.Fatal("created_at should not be in the future")
	}
}
