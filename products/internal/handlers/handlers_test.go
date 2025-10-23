package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.MixieMelts.products/internal/models"
)

// MockDB is a mock implementation of the DBLayer for testing purposes.
type MockDB struct {
	GetProductsFunc           func(ctx context.Context, limit int) ([]models.Product, error)
	CreateProductFunc         func(ctx context.Context, product *models.Product) (int64, error)
	GetSubscriptionBoxesFunc  func(ctx context.Context, limit int) ([]models.SubscriptionBox, error)
	CreateSubscriptionBoxFunc func(ctx context.Context, box *models.SubscriptionBox) (int64, error)
}

func (m *MockDB) GetProducts(ctx context.Context, limit int) ([]models.Product, error) {
	if m.GetProductsFunc != nil {
		return m.GetProductsFunc(ctx, limit)
	}
	return nil, errors.New("GetProductsFunc not implemented")
}

func (m *MockDB) CreateProduct(ctx context.Context, product *models.Product) (int64, error) {
	if m.CreateProductFunc != nil {
		return m.CreateProductFunc(ctx, product)
	}
	return 0, errors.New("CreateProductFunc not implemented")
}

func (m *MockDB) GetSubscriptionBoxes(ctx context.Context, limit int) ([]models.SubscriptionBox, error) {
	if m.GetSubscriptionBoxesFunc != nil {
		return m.GetSubscriptionBoxesFunc(ctx, limit)
	}
	return nil, errors.New("GetSubscriptionBoxesFunc not implemented")
}

func (m *MockDB) CreateSubscriptionBox(ctx context.Context, box *models.SubscriptionBox) (int64, error) {
	if m.CreateSubscriptionBoxFunc != nil {
		return m.CreateSubscriptionBoxFunc(ctx, box)
	}
	return 0, errors.New("CreateSubscriptionBoxFunc not implemented")
}

// Table-driven tests for GetProducts
func TestGetProducts(t *testing.T) {
	baseProducts := []models.Product{
		{ID: 1, Name: "Vanilla Wax Melts", Price: 10.99},
		{ID: 2, Name: "Lavender Candle", Price: 15.99},
	}

	tests := []struct {
		name         string
		limitQuery   string
		mockResponse []models.Product
		mockError    error
		wantStatus   int
		wantCount    int
	}{
		{
			name:         "ok no limit",
			limitQuery:   "",
			mockResponse: baseProducts,
			mockError:    nil,
			wantStatus:   http.StatusOK,
			wantCount:    len(baseProducts),
		},
		{
			name:         "ok with limit",
			limitQuery:   "?limit=1",
			mockResponse: baseProducts[:1],
			mockError:    nil,
			wantStatus:   http.StatusOK,
			wantCount:    1,
		},
		{
			name:         "database error",
			limitQuery:   "",
			mockResponse: nil,
			mockError:    errors.New("db failure"),
			wantStatus:   http.StatusInternalServerError,
			wantCount:    0,
		},
	}

	for _, tc := range tests {
		tc := tc // capture
		t.Run(tc.name, func(t *testing.T) {
			mockDB := &MockDB{
				GetProductsFunc: func(ctx context.Context, limit int) ([]models.Product, error) {
					if tc.mockError != nil {
						return nil, tc.mockError
					}
					return tc.mockResponse, nil
				},
			}

			handler := New(mockDB)

			req, err := http.NewRequest("GET", "/products"+tc.limitQuery, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.GetProducts(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("[%s] unexpected status: got %d want %d; body: %s", tc.name, rr.Code, tc.wantStatus, rr.Body.String())
			}

			if tc.wantStatus == http.StatusOK {
				var returned []models.Product
				if err := json.NewDecoder(rr.Body).Decode(&returned); err != nil {
					t.Fatalf("[%s] failed to decode response: %v", tc.name, err)
				}
				if len(returned) != tc.wantCount {
					t.Fatalf("[%s] expected %d products, got %d", tc.name, tc.wantCount, len(returned))
				}
			}
		})
	}
}

// Table-driven tests for CreateProduct
func TestCreateProduct(t *testing.T) {
	validProduct := models.Product{Name: "New Product", Price: 20.00}
	createdID := int64(123)

	tests := []struct {
		name          string
		payload       []byte
		mockCreateID  int64
		mockCreateErr error
		wantStatus    int
		wantCreatedID int64
	}{
		{
			name:          "successful creation",
			payload:       func() []byte { b, _ := json.Marshal(validProduct); return b }(),
			mockCreateID:  createdID,
			mockCreateErr: nil,
			wantStatus:    http.StatusCreated,
			wantCreatedID: createdID,
		},
		{
			name:          "invalid body",
			payload:       []byte("{invalid-json"),
			mockCreateID:  0,
			mockCreateErr: nil,
			wantStatus:    http.StatusBadRequest,
			wantCreatedID: 0,
		},
		{
			name:          "database error on creation",
			payload:       func() []byte { b, _ := json.Marshal(validProduct); return b }(),
			mockCreateID:  0,
			mockCreateErr: errors.New("db error"),
			wantStatus:    http.StatusInternalServerError,
			wantCreatedID: 0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockDB := &MockDB{
				CreateProductFunc: func(ctx context.Context, product *models.Product) (int64, error) {
					if tc.mockCreateErr != nil {
						return 0, tc.mockCreateErr
					}
					// basic validation to ensure payload was parsed
					if product.Name == "" {
						return 0, errors.New("invalid product")
					}
					return tc.mockCreateID, nil
				},
			}

			handler := New(mockDB)

			req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(tc.payload))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.CreateProduct(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("[%s] unexpected status: got %d want %d; body: %s", tc.name, rr.Code, tc.wantStatus, rr.Body.String())
			}

			if tc.wantStatus == http.StatusCreated {
				var created models.Product
				if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
					t.Fatalf("[%s] failed to decode response: %v", tc.name, err)
				}
				if created.ID != tc.wantCreatedID {
					t.Fatalf("[%s] expected created ID %d got %d", tc.name, tc.wantCreatedID, created.ID)
				}
			}
		})
	}
}

// Table-driven tests for subscription boxes
func TestSubscriptionBoxes(t *testing.T) {
	baseBoxes := []models.SubscriptionBox{
		{ID: 1, Name: "Monthly Surprise", Description: "curated", Price: 29.99},
		{ID: 2, Name: "Seasonal Box", Description: "seasonal", Price: 25.00},
	}

	t.Run("GetSubscriptionBoxes", func(t *testing.T) {
		tests := []struct {
			name       string
			query      string
			mockResp   []models.SubscriptionBox
			mockErr    error
			wantStatus int
			wantCount  int
		}{
			{name: "ok no limit", query: "", mockResp: baseBoxes, mockErr: nil, wantStatus: http.StatusOK, wantCount: len(baseBoxes)},
			{name: "ok with limit", query: "?limit=1", mockResp: baseBoxes[:1], mockErr: nil, wantStatus: http.StatusOK, wantCount: 1},
			{name: "db error", query: "", mockResp: nil, mockErr: errors.New("db fail"), wantStatus: http.StatusInternalServerError, wantCount: 0},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				mockDB := &MockDB{
					GetSubscriptionBoxesFunc: func(ctx context.Context, limit int) ([]models.SubscriptionBox, error) {
						if tc.mockErr != nil {
							return nil, tc.mockErr
						}
						return tc.mockResp, nil
					},
				}

				handler := New(mockDB)
				req, _ := http.NewRequest("GET", "/products/subscription-boxes"+tc.query, nil)
				rr := httptest.NewRecorder()
				handler.GetSubscriptionBoxes(rr, req)

				if rr.Code != tc.wantStatus {
					t.Fatalf("[%s] unexpected status: got %d want %d; body: %s", tc.name, rr.Code, tc.wantStatus, rr.Body.String())
				}

				if tc.wantStatus == http.StatusOK {
					var boxes []models.SubscriptionBox
					if err := json.NewDecoder(rr.Body).Decode(&boxes); err != nil {
						t.Fatalf("[%s] failed to decode response: %v", tc.name, err)
					}
					if len(boxes) != tc.wantCount {
						t.Fatalf("[%s] expected %d boxes got %d", tc.name, tc.wantCount, len(boxes))
					}
				}
			})
		}
	})

	t.Run("CreateSubscriptionBox", func(t *testing.T) {
		validBox := models.SubscriptionBox{Name: "Monthly Surprise", Description: "curated", Price: 29.99}
		createdID := int64(77)

		tests := []struct {
			name    string
			payload []byte
			mockID  int64
			mockErr error
			want    int
			wantID  int64
		}{
			{
				name:    "successful create",
				payload: func() []byte { b, _ := json.Marshal(validBox); return b }(),
				mockID:  createdID,
				mockErr: nil,
				want:    http.StatusCreated,
				wantID:  createdID,
			},
			{
				name:    "invalid body",
				payload: []byte("{invalid-json"),
				mockID:  0,
				mockErr: nil,
				want:    http.StatusBadRequest,
				wantID:  0,
			},
			{
				name:    "db error",
				payload: func() []byte { b, _ := json.Marshal(validBox); return b }(),
				mockID:  0,
				mockErr: errors.New("db error"),
				want:    http.StatusInternalServerError,
				wantID:  0,
			},
		}

		for _, tc := range tests {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				mockDB := &MockDB{
					CreateSubscriptionBoxFunc: func(ctx context.Context, box *models.SubscriptionBox) (int64, error) {
						if tc.mockErr != nil {
							return 0, tc.mockErr
						}
						if box.Name == "" {
							return 0, errors.New("invalid")
						}
						return tc.mockID, nil
					},
				}

				handler := New(mockDB)
				req, _ := http.NewRequest("POST", "/products/subscription-boxes", bytes.NewBuffer(tc.payload))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()
				handler.CreateSubscriptionBox(rr, req)

				if rr.Code != tc.want {
					t.Fatalf("[%s] unexpected status: got %d want %d; body: %s", tc.name, rr.Code, tc.want, rr.Body.String())
				}

				if tc.want == http.StatusCreated {
					var created models.SubscriptionBox
					if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
						t.Fatalf("[%s] failed to decode response: %v", tc.name, err)
					}
					if created.ID != tc.wantID {
						t.Fatalf("[%s] expected created id %d got %d", tc.name, tc.wantID, created.ID)
					}
				}
			})
		}
	})
}
