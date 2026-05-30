package crud_permission_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"starter-wahcah-be/internal/modules/permission/crud_permission"
	"starter-wahcah-be/internal/modules/permission/crud_permission/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// =========================================
// Helper
// =========================================

func setupPermissionControllerApp(svc crud_permission.Service) *fiber.App {
	app := fiber.New()
	ctrl := crud_permission.NewCrudPermissionController(svc)

	permissions := app.Group("/permissions")
	permissions.Get("/", ctrl.GetAll)
	permissions.Get("/:id", ctrl.GetByID)
	permissions.Post("/", ctrl.Create)
	permissions.Put("/:id", ctrl.Update)
	permissions.Delete("/:id", ctrl.Delete)

	return app
}

// =========================================
// GetAll
// =========================================

func TestController_GetAll_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetAll", 1, 10, (*uint)(nil)).Return(&crud_permission.PermissionListResponse{
		Message: "Permissions retrieved successfully",
		Data: []crud_permission.PermissionResponse{
			{ID: 1, RoleID: 1, Name: "role:read"},
			{ID: 2, RoleID: 1, Name: "role:create"},
			{ID: 3, RoleID: 2, Name: "user:read"},
		},
		Pagination: crud_permission.PaginationMeta{
			Page: 1, Limit: 10, Total: 3, TotalPage: 1,
		},
	}, nil)

	app := setupPermissionControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/permissions?page=1&limit=10", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud_permission.PermissionListResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Permissions retrieved successfully", response.Message)
	assert.Equal(t, 3, len(response.Data))
	assert.Equal(t, int64(3), response.Pagination.Total)
	mockSvc.AssertExpectations(t)
}

func TestController_GetAll_FilterByRoleID(t *testing.T) {
	mockSvc := new(mocks.Service)

	roleID := uint(1)
	mockSvc.On("GetAll", 1, 10, &roleID).Return(&crud_permission.PermissionListResponse{
		Message: "Permissions retrieved successfully",
		Data: []crud_permission.PermissionResponse{
			{ID: 1, RoleID: 1, Name: "role:read"},
			{ID: 2, RoleID: 1, Name: "role:create"},
		},
		Pagination: crud_permission.PaginationMeta{
			Page: 1, Limit: 10, Total: 2, TotalPage: 1,
		},
	}, nil)

	app := setupPermissionControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/permissions?page=1&limit=10&role_id=1", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud_permission.PermissionListResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, 2, len(response.Data))
	mockSvc.AssertExpectations(t)
}

func TestController_GetAll_InvalidRoleID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupPermissionControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/permissions?role_id=abc", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Invalid role_id", response["error"])
}

func TestController_GetAll_DefaultPagination(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetAll", 1, 10, (*uint)(nil)).Return(&crud_permission.PermissionListResponse{
		Message:    "Permissions retrieved successfully",
		Data:       []crud_permission.PermissionResponse{},
		Pagination: crud_permission.PaginationMeta{Page: 1, Limit: 10, Total: 0, TotalPage: 0},
	}, nil)

	app := setupPermissionControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/permissions", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_GetAll_DBError(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetAll", 1, 10, (*uint)(nil)).Return(nil, assert.AnError)

	app := setupPermissionControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/permissions?page=1&limit=10", nil)

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

	mockSvc.On("GetByID", uint(1)).Return(&crud_permission.PermissionSingleResponse{
		Message: "Permission retrieved successfully",
		Data:    crud_permission.PermissionResponse{ID: 1, RoleID: 1, Name: "role:read"},
	}, nil)

	app := setupPermissionControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/permissions/1", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud_permission.PermissionSingleResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Permission retrieved successfully", response.Message)
	assert.Equal(t, "role:read", response.Data.Name)
	mockSvc.AssertExpectations(t)
}

func TestController_GetByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetByID", uint(99)).Return(nil, assert.AnError)

	app := setupPermissionControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/permissions/99", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_GetByID_InvalidID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupPermissionControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/permissions/abc", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// =========================================
// Create
// =========================================

func TestController_Create_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Create", crud_permission.CreatePermissionRequest{
		RoleID: 1,
		Name:   "role:delete",
	}).Return(&crud_permission.PermissionSingleResponse{
		Message: "Permission created successfully",
		Data:    crud_permission.PermissionResponse{ID: 5, RoleID: 1, Name: "role:delete"},
	}, nil)

	app := setupPermissionControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"role_id": 1,
		"name":    "role:delete",
	})
	req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var response crud_permission.PermissionSingleResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Permission created successfully", response.Message)
	assert.Equal(t, "role:delete", response.Data.Name)
	assert.Equal(t, uint(1), response.Data.RoleID)
	mockSvc.AssertExpectations(t)
}

func TestController_Create_InvalidJSON(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupPermissionControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestController_Create_ValidationError(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupPermissionControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{})
	req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotNil(t, response["validation"])
}

func TestController_Create_AlreadyExists(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Create", crud_permission.CreatePermissionRequest{
		RoleID: 1,
		Name:   "role:read",
	}).Return(nil, assert.AnError)

	app := setupPermissionControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"role_id": 1,
		"name":    "role:read",
	})
	req := httptest.NewRequest(http.MethodPost, "/permissions", bytes.NewBuffer(body))
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

	mockSvc.On("Update", uint(1), crud_permission.UpdatePermissionRequest{
		RoleID: 1,
		Name:   "role:read-updated",
	}).Return(&crud_permission.PermissionSingleResponse{
		Message: "Permission updated successfully",
		Data:    crud_permission.PermissionResponse{ID: 1, RoleID: 1, Name: "role:read-updated"},
	}, nil)

	app := setupPermissionControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"role_id": 1,
		"name":    "role:read-updated",
	})
	req := httptest.NewRequest(http.MethodPut, "/permissions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud_permission.PermissionSingleResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Permission updated successfully", response.Message)
	assert.Equal(t, "role:read-updated", response.Data.Name)
	mockSvc.AssertExpectations(t)
}

func TestController_Update_NotFound(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Update", uint(99), crud_permission.UpdatePermissionRequest{
		RoleID: 1,
		Name:   "role:read",
	}).Return(nil, assert.AnError)

	app := setupPermissionControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"role_id": 1,
		"name":    "role:read",
	})
	req := httptest.NewRequest(http.MethodPut, "/permissions/99", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_Update_InvalidID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupPermissionControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"role_id": 1,
		"name":    "role:read",
	})
	req := httptest.NewRequest(http.MethodPut, "/permissions/abc", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestController_Update_ValidationError(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupPermissionControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{})
	req := httptest.NewRequest(http.MethodPut, "/permissions/1", bytes.NewBuffer(body))
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

	app := setupPermissionControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/permissions/1", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Permission deleted successfully", response["message"])
	mockSvc.AssertExpectations(t)
}

func TestController_Delete_NotFound(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Delete", uint(99)).Return(assert.AnError)

	app := setupPermissionControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/permissions/99", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_Delete_InvalidID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupPermissionControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/permissions/abc", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
