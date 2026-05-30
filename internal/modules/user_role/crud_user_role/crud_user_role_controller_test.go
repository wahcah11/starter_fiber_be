package crud_user_role_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"starter-wahcah-be/internal/modules/user_role/crud_user_role"
	"starter-wahcah-be/internal/modules/user_role/crud_user_role/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// =========================================
// Helper
// =========================================

func setupUserRoleControllerApp(svc crud_user_role.Service) *fiber.App {
	app := fiber.New()
	ctrl := crud_user_role.NewCrudUserRoleController(svc)

	userRoles := app.Group("/user-roles")
	userRoles.Get("/", ctrl.GetAll)
	userRoles.Get("/:id", ctrl.GetByID)
	userRoles.Post("/", ctrl.Create)
	userRoles.Delete("/:id", ctrl.Delete)

	return app
}

// =========================================
// GetAll
// =========================================

func TestController_GetAll_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetAll", 1, 10, (*uint)(nil)).Return(&crud_user_role.UserRoleListResponse{
		Message: "User roles retrieved successfully",
		Data: []crud_user_role.UserRoleResponse{
			{ID: 1, UserID: 1, RoleID: 1},
			{ID: 2, UserID: 1, RoleID: 2},
			{ID: 3, UserID: 2, RoleID: 1},
		},
		Pagination: crud_user_role.PaginationMeta{
			Page: 1, Limit: 10, Total: 3, TotalPage: 1,
		},
	}, nil)

	app := setupUserRoleControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/user-roles?page=1&limit=10", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud_user_role.UserRoleListResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "User roles retrieved successfully", response.Message)
	assert.Equal(t, 3, len(response.Data))
	assert.Equal(t, int64(3), response.Pagination.Total)
	mockSvc.AssertExpectations(t)
}

func TestController_GetAll_FilterByUserID(t *testing.T) {
	mockSvc := new(mocks.Service)

	userID := uint(1)
	mockSvc.On("GetAll", 1, 10, &userID).Return(&crud_user_role.UserRoleListResponse{
		Message: "User roles retrieved successfully",
		Data: []crud_user_role.UserRoleResponse{
			{ID: 1, UserID: 1, RoleID: 1},
			{ID: 2, UserID: 1, RoleID: 2},
		},
		Pagination: crud_user_role.PaginationMeta{
			Page: 1, Limit: 10, Total: 2, TotalPage: 1,
		},
	}, nil)

	app := setupUserRoleControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/user-roles?page=1&limit=10&user_id=1", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud_user_role.UserRoleListResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, 2, len(response.Data))
	mockSvc.AssertExpectations(t)
}

func TestController_GetAll_InvalidUserID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupUserRoleControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/user-roles?user_id=abc", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Invalid user_id", response["error"])
}

func TestController_GetAll_DefaultPagination(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetAll", 1, 10, (*uint)(nil)).Return(&crud_user_role.UserRoleListResponse{
		Message:    "User roles retrieved successfully",
		Data:       []crud_user_role.UserRoleResponse{},
		Pagination: crud_user_role.PaginationMeta{Page: 1, Limit: 10, Total: 0, TotalPage: 0},
	}, nil)

	app := setupUserRoleControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/user-roles", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_GetAll_DBError(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetAll", 1, 10, (*uint)(nil)).Return(nil, assert.AnError)

	app := setupUserRoleControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/user-roles?page=1&limit=10", nil)

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

	mockSvc.On("GetByID", uint(1)).Return(&crud_user_role.UserRoleSingleResponse{
		Message: "User role retrieved successfully",
		Data:    crud_user_role.UserRoleResponse{ID: 1, UserID: 1, RoleID: 2},
	}, nil)

	app := setupUserRoleControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/user-roles/1", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response crud_user_role.UserRoleSingleResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "User role retrieved successfully", response.Message)
	assert.Equal(t, uint(1), response.Data.UserID)
	assert.Equal(t, uint(2), response.Data.RoleID)
	mockSvc.AssertExpectations(t)
}

func TestController_GetByID_NotFound(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("GetByID", uint(99)).Return(nil, assert.AnError)

	app := setupUserRoleControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/user-roles/99", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_GetByID_InvalidID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupUserRoleControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/user-roles/abc", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// =========================================
// Create
// =========================================

func TestController_Create_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Create", crud_user_role.CreateUserRoleRequest{
		UserID: 1,
		RoleID: 3,
	}).Return(&crud_user_role.UserRoleSingleResponse{
		Message: "Role assigned to user successfully",
		Data:    crud_user_role.UserRoleResponse{ID: 5, UserID: 1, RoleID: 3},
	}, nil)

	app := setupUserRoleControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"user_id": 1,
		"role_id": 3,
	})
	req := httptest.NewRequest(http.MethodPost, "/user-roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var response crud_user_role.UserRoleSingleResponse
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Role assigned to user successfully", response.Message)
	assert.Equal(t, uint(1), response.Data.UserID)
	assert.Equal(t, uint(3), response.Data.RoleID)
	mockSvc.AssertExpectations(t)
}

func TestController_Create_InvalidJSON(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupUserRoleControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/user-roles", bytes.NewBufferString("invalid-json"))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestController_Create_ValidationError(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupUserRoleControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{})
	req := httptest.NewRequest(http.MethodPost, "/user-roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.NotNil(t, response["validation"])
}

func TestController_Create_AlreadyAssigned(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Create", crud_user_role.CreateUserRoleRequest{
		UserID: 1,
		RoleID: 1,
	}).Return(nil, assert.AnError)

	app := setupUserRoleControllerApp(mockSvc)

	body, _ := json.Marshal(fiber.Map{
		"user_id": 1,
		"role_id": 1,
	})
	req := httptest.NewRequest(http.MethodPost, "/user-roles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

// =========================================
// Delete
// =========================================

func TestController_Delete_Success(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Delete", uint(1)).Return(nil)

	app := setupUserRoleControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/user-roles/1", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, "Role removed from user successfully", response["message"])
	mockSvc.AssertExpectations(t)
}

func TestController_Delete_NotFound(t *testing.T) {
	mockSvc := new(mocks.Service)

	mockSvc.On("Delete", uint(99)).Return(assert.AnError)

	app := setupUserRoleControllerApp(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/user-roles/99", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	mockSvc.AssertExpectations(t)
}

func TestController_Delete_InvalidID(t *testing.T) {
	mockSvc := new(mocks.Service)
	app := setupUserRoleControllerApp(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/user-roles/abc", nil)

	res, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}
