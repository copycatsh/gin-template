package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-template/internal/handler"
	"gin-template/internal/model"
	"gin-template/internal/service"

	"gorm.io/gorm"
)

type mockUserRepository struct {
	users   []model.User
	findErr error
	saveErr error
}

func (m *mockUserRepository) FindAll() ([]model.User, error) {
	return m.users, m.findErr
}

func (m *mockUserRepository) FindByID(id string) (model.User, error) {
	for _, u := range m.users {
		if fmt.Sprint(u.ID) == id {
			return u, nil
		}
	}
	return model.User{}, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) Create(user *model.User) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	user.ID = uint(len(m.users) + 1)
	m.users = append(m.users, *user)
	return nil
}

func (m *mockUserRepository) Save(user *model.User) error {
	return m.saveErr
}

func setupRouter(mock *mockUserRepository) http.Handler {
	svc := service.NewUserService(mock)
	h := handler.NewUserHandler(svc)
	return handler.SetupRouter(h)
}

func TestGetRoot(t *testing.T) {
	router := setupRouter(&mockUserRepository{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetPing(t *testing.T) {
	router := setupRouter(&mockUserRepository{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetUsers_Success(t *testing.T) {
	mock := &mockUserRepository{
		users: []model.User{{ID: 1, Name: "Alice", Email: "a@b.com"}},
	}
	router := setupRouter(mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetUsers_DBError(t *testing.T) {
	mock := &mockUserRepository{findErr: errors.New("db down")}
	router := setupRouter(mock)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestCreateUser_Success(t *testing.T) {
	mock := &mockUserRepository{}
	router := setupRouter(mock)

	body, _ := json.Marshal(model.User{Name: "Bob", Email: "bob@test.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestCreateUser_BadJSON(t *testing.T) {
	router := setupRouter(&mockUserRepository{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateUser_DBError(t *testing.T) {
	mock := &mockUserRepository{saveErr: errors.New("db error")}
	router := setupRouter(mock)

	body, _ := json.Marshal(model.User{Name: "Bob", Email: "bob@test.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	mock := &mockUserRepository{
		users: []model.User{{ID: 1, Name: "Alice", Email: "a@b.com"}},
	}
	router := setupRouter(mock)

	body, _ := json.Marshal(model.User{Name: "Updated", Email: "new@b.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	mock := &mockUserRepository{}
	router := setupRouter(mock)

	body, _ := json.Marshal(model.User{Name: "X", Email: "x@b.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestUpdateUser_BadJSON(t *testing.T) {
	router := setupRouter(&mockUserRepository{
		users: []model.User{{ID: 1, Name: "Alice", Email: "a@b.com"}},
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateUser_SaveError(t *testing.T) {
	mock := &mockUserRepository{
		users:   []model.User{{ID: 1, Name: "Alice", Email: "a@b.com"}},
		saveErr: errors.New("save failed"),
	}
	router := setupRouter(mock)

	body, _ := json.Marshal(model.User{Name: "Updated", Email: "new@b.com"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
