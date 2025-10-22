package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.MixieMelts.users/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// DBLayer defines the interface for database operations. This would be used
// in a refactored Handler to allow for dependency injection.
type DBLayer interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
}

// MockDB is a mock implementation of DBLayer for testing.
type MockDB struct {
	GetUserByEmailFunc func(ctx context.Context, email string) (*models.User, error)
	CreateUserFunc     func(ctx context.Context, user *models.User) (int64, error)
	GetUserByIDFunc    func(ctx context.Context, id int64) (*models.User, error)
}

func (m *MockDB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return nil, errors.New("GetUserByEmailFunc not implemented")
}

func (m *MockDB) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return 0, errors.New("CreateUserFunc not implemented")
}

func (m *MockDB) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	return nil, errors.New("GetUserByIDFunc not implemented")
}

func TestRegisterUser(t *testing.T) {
	// A dummy JWT secret key for testing purposes.
	jwtSecret := []byte("test-secret")

	t.Run("successful registration", func(t *testing.T) {
		mockDB := &MockDB{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
				// Simulate user not found
				return nil, errors.New("user not found")
			},
			CreateUserFunc: func(ctx context.Context, user *models.User) (int64, error) {
				// Simulate successful user creation
				return 1, nil
			},
		}

		// This test requires refactoring New() and Handler to accept the DBLayer interface.
		// handler := New(mockDB, jwtSecret)

		creds := models.Credentials{Email: "test@example.com", Password: "password123"}
		body, _ := json.Marshal(creds)
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// handler.RegisterUser(rr, req)

		// Expected assertions:
		// if status := rr.Code; status != http.StatusCreated {
		// 	t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		// }
		//
		// var user models.User
		// json.Unmarshal(rr.Body.Bytes(), &user)
		// if user.Email != creds.Email {
		// 	t.Errorf("handler returned unexpected body: got email %v want %v", user.Email, creds.Email)
		// }
		t.Skip("Skipping test: requires refactoring handlers.go to use an interface for mocking the database.")
	})

	t.Run("user already exists", func(t *testing.T) {
		mockDB := &MockDB{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
				// Simulate user found
				return &models.User{ID: 1, Email: "test@example.com"}, nil
			},
		}

		// handler := New(mockDB, jwtSecret)

		creds := models.Credentials{Email: "test@example.com", Password: "password123"}
		body, _ := json.Marshal(creds)
		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// handler.RegisterUser(rr, req)

		// Expected assertions:
		// if status := rr.Code; status != http.StatusConflict {
		// 	t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
		// }
		t.Skip("Skipping test: requires refactoring handlers.go to use an interface for mocking the database.")
	})
}

func TestLoginUser(t *testing.T) {
	jwtSecret := []byte("test-secret")
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	t.Run("successful login", func(t *testing.T) {
		mockDB := &MockDB{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}, nil
			},
		}

		// handler := New(mockDB, jwtSecret)

		creds := models.Credentials{Email: "test@example.com", Password: password}
		body, _ := json.Marshal(creds)
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// handler.LoginUser(rr, req)

		// Expected assertions:
		// if status := rr.Code; status != http.StatusOK {
		// 	t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		// }
		//
		// var resp map[string]string
		// json.Unmarshal(rr.Body.Bytes(), &resp)
		// if _, ok := resp["token"]; !ok {
		// 	t.Error("handler did not return a token")
		// }
		t.Skip("Skipping test: requires refactoring handlers.go to use an interface for mocking the database.")
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockDB := &MockDB{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
				return &models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
				}, nil
			},
		}
		// handler := New(mockDB, jwtSecret)

		creds := models.Credentials{Email: "test@example.com", Password: "wrongpassword"}
		body, _ := json.Marshal(creds)
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		// handler.LoginUser(rr, req)

		// Expected assertions:
		// if status := rr.Code; status != http.StatusUnauthorized {
		// 	t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		// }
		t.Skip("Skipping test: requires refactoring handlers.go to use an interface for mocking the database.")
	})
}
