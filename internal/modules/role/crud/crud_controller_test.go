package crud_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"starter-wahcah-be/internal/modules/role/crud"
	"starter-wahcah-be/internal/modules/role/crud/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// =========================================
// Helper
// =========================================

func setupCrudControllerApp(svc crud.Service) *fiber.App {
	app := fiber.New()
	ctrl := crud.NewCrudController(svc)

	roles := app.Group("/roles")
	roles.Get("/", ctrl.GetAll)
	roles.Get("/:id", ctrl.GetByID)
	roles.Post("/", ctrl.Create)
	roles.Put("/:id", ctrl.Update)
	roles.Delete("/:id", ctrl.Delete)

	return app
}

// =========================================
// GetAll
// =========================================

func TestController_GetAll_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetAll", 1, 10).Return(&crud.RoleListResponse{
		Message: "Roles retrieved successfully",
		Data: []crud.RoleResponse{
			{ID: 1, Name: "superadmin", SystemFunction: "full_access", IsDefault: false},
			{ID: 2, Name: "member", SystemFunction: "basic_access", IsDefault: true},
		},
		Pagination: crud.PaginationMeta{
			Page: 1, Limit: 10, Total: 2, TotalPage: 1,
		},
	}, nil)

	app := setupCrudControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/roles?page=1&limit=10", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud.RoleListResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Roles retrieved successfully", response.Message)
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, int64(2), response.Pagination.Total)
	mockSvc.AssertExpectations(t)
}

func TestController_GetAll_DefaultPagination(t *testing.T) {
	mockSvc := new(mocks.Service)

	// Tanpa query param, default page=1 limit=10
	mockSvc.On("GetAll", 1, 10).Return(&crud.RoleListResponse{
		Message:    "Roles retrieved successfully",
		Data:       []crud.RoleResponse{},
		Pagination: crud.PaginationMeta{Page: 1, Limit: 10, Total: 0, TotalPage: 0},
	}, nil)

	app := setupCrudControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/roles", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_GetAll_DBError(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetAll", 1, 10).Return(nil, assert.AnError)

	app := setupCrudControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/roles?page=1&limit=10", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

// =========================================
// GetByID
// =========================================

func TestController_GetByID_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetByID", uint(1)).Return(&crud.RoleSingleResponse{
		Message: "Role retrieved successfully",
		Data:    crud.RoleResponse{ID: 1, Name: "superadmin", SystemFunction: "full_access"},
	}, nil)

	app := setupCrudControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/roles/1", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud.RoleSingleResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Role retrieved successfully", response.Message)
	assert.Equal(t, "superadmin", response.Data.Name)
	mockSvc.AssertExpectations(t)
}

func TestController_GetByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetByID", uint(99)).Return(nil, assert.AnError)

	app := setupCrudControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/roles/99", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_GetByID_InvalidID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupCrudControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/roles/abc", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// =========================================
// Create
// =========================================

func TestController_Create_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Create", crud.CreateRoleRequest{
		Name:           "editor",
		SystemFunction: "edit_access",
		IsDefault:      false,
	}).Return(&crud.RoleSingleResponse{
		Message: "Role created successfully",
		Data:    crud.RoleResponse{ID: 3, Name: "editor", SystemFunction: "edit_access"},
	}, nil)

	app := setupCrudControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"name":            "editor",
		"system_function": "edit_access",
		"is_default":      false,
	})
	req := httptest.NewRequest(http.MethodPost, "/roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var response crud.RoleSingleResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Role created successfully", response.Message)
	assert.Equal(t, "editor", response.Data.Name)
	mockSvc.AssertExpectations(t)
}

func TestController_Create_InvalidJSON(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupCrudControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/roles", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestController_Create_ValidationError(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupCrudControllerApp(mockSvc)

	// Name kosong, system_function kosong
	body, _ := json.Marshal(fiber.Map{})
	req := httptest.NewRequest(http.MethodPost, "/roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotNil(t, response["validation"])
}

func TestController_Create_NameAlreadyExists(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Create", crud.CreateRoleRequest{
		Name:           "member",
		SystemFunction: "basic_access",
		IsDefault:      false,
	}).Return(nil, assert.AnError)

	app := setupCrudControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"name":            "member",
		"system_function": "basic_access",
		"is_default":      false,
	})
	req := httptest.NewRequest(http.MethodPost, "/roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

// =========================================
// Update
// =========================================

func TestController_Update_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Update", uint(1), crud.UpdateRoleRequest{
		Name:           "editor-updated",
		SystemFunction: "edit_access",
		IsDefault:      false,
	}).Return(&crud.RoleSingleResponse{
		Message: "Role updated successfully",
		Data:    crud.RoleResponse{ID: 1, Name: "editor-updated", SystemFunction: "edit_access"},
	}, nil)

	app := setupCrudControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"name":            "editor-updated",
		"system_function": "edit_access",
		"is_default":      false,
	})
	req := httptest.NewRequest(http.MethodPut, "/roles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud.RoleSingleResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Role updated successfully", response.Message)
	assert.Equal(t, "editor-updated", response.Data.Name)
	mockSvc.AssertExpectations(t)
}

func TestController_Update_NotFound(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Update", uint(99), crud.UpdateRoleRequest{
		Name:           "editor",
		SystemFunction: "edit_access",
		IsDefault:      false,
	}).Return(nil, assert.AnError)

	app := setupCrudControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"name":            "editor",
		"system_function": "edit_access",
		"is_default":      false,
	})
	req := httptest.NewRequest(http.MethodPut, "/roles/99", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_Update_InvalidID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupCrudControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"name":            "editor",
		"system_function": "edit_access",
	})
	req := httptest.NewRequest(http.MethodPut, "/roles/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestController_Update_ValidationError(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupCrudControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{})
	req := httptest.NewRequest(http.MethodPut, "/roles/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotNil(t, response["validation"])
}

// =========================================
// Delete
// =========================================

func TestController_Delete_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Delete", uint(1)).Return(nil)

	app := setupCrudControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/roles/1", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Role deleted successfully", response["message"])
	mockSvc.AssertExpectations(t)
}

func TestController_Delete_NotFound(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Delete", uint(99)).Return(assert.AnError)

	app := setupCrudControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/roles/99", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_Delete_InvalidID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupCrudControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/roles/abc", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
