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

// MockDB is a mock implementation of the DBLayer defined in handlers.go.
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

// Table-driven tests for RegisterUser
func TestRegisterUser(t *testing.T) {
	jwtSecret := []byte("test-secret")

	tests := []struct {
		name      string
		existing  bool
		createErr error
		wantCode  int
	}{
		{name: "successful registration", existing: false, createErr: nil, wantCode: http.StatusCreated},
		{name: "user already exists", existing: true, createErr: nil, wantCode: http.StatusConflict},
		{name: "db error on create", existing: false, createErr: errors.New("db fail"), wantCode: http.StatusInternalServerError},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockDB := &MockDB{
				GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
					if tc.existing {
						return &models.User{ID: 2, Email: email}, nil
					}
					return nil, errors.New("not found")
				},
				CreateUserFunc: func(ctx context.Context, user *models.User) (int64, error) {
					if tc.createErr != nil {
						return 0, tc.createErr
					}
					return 1, nil
				},
			}

			handler := New(mockDB, jwtSecret)

			creds := models.Credentials{Email: "test@example.com", Password: "password123"}
			body, _ := json.Marshal(creds)
			req, err := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.RegisterUser(rr, req)

			if rr.Code != tc.wantCode {
				t.Fatalf("[%s] expected status %d got %d; body: %s", tc.name, tc.wantCode, rr.Code, rr.Body.String())
			}

			if tc.wantCode == http.StatusCreated {
				var created models.User
				if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
					t.Fatalf("[%s] failed to decode response: %v", tc.name, err)
				}
				if created.Email != creds.Email {
					t.Fatalf("[%s] expected email %q got %q", tc.name, creds.Email, created.Email)
				}
				if created.ID != 1 {
					t.Fatalf("[%s] expected id %d got %d", tc.name, 1, created.ID)
				}
			}
		})
	}
}

// Table-driven tests for LoginUser
func TestLoginUser(t *testing.T) {
	jwtSecret := []byte("test-secret")
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		userPresent    bool
		passwordHash   []byte // hash to return as stored password
		credPassword   string
		wantStatusCode int
	}{
		{name: "successful login", userPresent: true, passwordHash: hashedPassword, credPassword: password, wantStatusCode: http.StatusOK},
		{name: "invalid credentials", userPresent: true, passwordHash: func() []byte { h, _ := bcrypt.GenerateFromPassword([]byte("otherpass"), bcrypt.DefaultCost); return h }(), credPassword: "wrongpassword", wantStatusCode: http.StatusUnauthorized},
		{name: "user not found", userPresent: false, passwordHash: nil, credPassword: "irrelevant", wantStatusCode: http.StatusUnauthorized},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockDB := &MockDB{
				GetUserByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
					if !tc.userPresent {
						return nil, errors.New("not found")
					}
					return &models.User{
						ID:       1,
						Email:    "test@example.com",
						Password: string(tc.passwordHash),
						IsAdmin:  false,
					}, nil
				},
			}

			handler := New(mockDB, jwtSecret)

			creds := models.Credentials{Email: "test@example.com", Password: tc.credPassword}
			body, _ := json.Marshal(creds)
			req, err := http.NewRequest("POST", "/api/users/login", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.LoginUser(rr, req)

			if rr.Code != tc.wantStatusCode {
				t.Fatalf("[%s] expected status %d got %d; body: %s", tc.name, tc.wantStatusCode, rr.Code, rr.Body.String())
			}

			if tc.wantStatusCode == http.StatusOK {
				var resp map[string]string
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("[%s] failed to decode response: %v", tc.name, err)
				}
				if _, ok := resp["token"]; !ok {
					t.Fatalf("[%s] expected token in response, got: %v", tc.name, resp)
				}
			}
		})
	}
}

// Table-driven tests for GetUserProfile
func TestGetUserProfile(t *testing.T) {
	jwtSecret := []byte("test-secret")

	tests := []struct {
		name        string
		userID      string
		mockUser    *models.User
		mockErr     error
		wantCode    int
		expectEmail string
	}{
		{name: "success", userID: "42", mockUser: &models.User{ID: 42, Email: "profile@example.com", IsAdmin: false}, mockErr: nil, wantCode: http.StatusOK, expectEmail: "profile@example.com"},
		{name: "user not found", userID: "99", mockUser: nil, mockErr: errors.New("not found"), wantCode: http.StatusNotFound, expectEmail: ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockDB := &MockDB{
				GetUserByIDFunc: func(ctx context.Context, id int64) (*models.User, error) {
					if tc.mockErr != nil {
						return nil, tc.mockErr
					}
					return tc.mockUser, nil
				},
			}

			handler := New(mockDB, jwtSecret)

			req, err := http.NewRequest("GET", "/api/users/me", nil)
			if err != nil {
				t.Fatal(err)
			}
			// The GetUserProfile handler expects a string userID in the context.
			ctx := context.WithValue(req.Context(), "userID", tc.userID)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			handler.GetUserProfile(rr, req)

			if rr.Code != tc.wantCode {
				t.Fatalf("[%s] expected status %d got %d; body: %s", tc.name, tc.wantCode, rr.Code, rr.Body.String())
			}

			if tc.wantCode == http.StatusOK {
				var returned models.User
				if err := json.NewDecoder(rr.Body).Decode(&returned); err != nil {
					t.Fatalf("[%s] failed to decode response: %v", tc.name, err)
				}
				if returned.ID != tc.mockUser.ID || returned.Email != tc.mockUser.Email {
					t.Fatalf("[%s] returned user mismatch: got %+v want %+v", tc.name, returned, tc.mockUser)
				}
			}
		})
	}
}
