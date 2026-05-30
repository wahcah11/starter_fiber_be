package crud_user_role_test

import (
	"errors"
	"testing"

	"starter-wahcah-be/internal/models"
	"starter-wahcah-be/internal/modules/user_role/crud_user_role"
	"starter-wahcah-be/internal/modules/user_role/crud_user_role/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// =========================================
// GetAll
// =========================================

func TestService_GetAll_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.UserRole{
		{UserID: 1, RoleID: 1},
		{UserID: 1, RoleID: 2},
		{UserID: 2, RoleID: 1},
	}, int64(3), nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetAll(1, 10, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "User roles retrieved successfully", res.Message)
	assert.Equal(t, 3, len(res.Data))
	assert.Equal(t, int64(3), res.Pagination.Total)
	assert.Equal(t, 1, res.Pagination.Page)
	assert.Equal(t, 10, res.Pagination.Limit)
	assert.Equal(t, 1, res.Pagination.TotalPage)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_FilterByUserID(t *testing.T) {
	mockRepo := new(mocks.Repository)

	userID := uint(1)
	mockRepo.On("FindAll", 1, 10, &userID).Return([]models.UserRole{
		{UserID: 1, RoleID: 1},
		{UserID: 1, RoleID: 2},
	}, int64(2), nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetAll(1, 10, &userID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(res.Data))
	assert.Equal(t, int64(2), res.Pagination.Total)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_Empty(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.UserRole{}, int64(0), nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(res.Data))
	assert.Equal(t, int64(0), res.Pagination.Total)
	assert.Equal(t, 0, res.Pagination.TotalPage)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_InvalidPage_DefaultsToOne(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.UserRole{}, int64(0), nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetAll(0, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.Pagination.Page)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_InvalidLimit_DefaultsToTen(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.UserRole{}, int64(0), nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetAll(1, 0, nil)

	assert.NoError(t, err)
	assert.Equal(t, 10, res.Pagination.Limit)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_LimitExceedsMax_DefaultsToTen(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return([]models.UserRole{}, int64(0), nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetAll(1, 999, nil)

	assert.NoError(t, err)
	assert.Equal(t, 10, res.Pagination.Limit)
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return(nil, int64(0), errors.New("db error"))

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetAll(1, 10, nil)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "failed to retrieve user roles", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_GetAll_PaginationCalculation(t *testing.T) {
	mockRepo := new(mocks.Repository)

	// Total 25 data dengan limit 10 = 3 halaman
	userRoles := make([]models.UserRole, 10)
	mockRepo.On("FindAll", 1, 10, (*uint)(nil)).Return(userRoles, int64(25), nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetAll(1, 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(25), res.Pagination.Total)
	assert.Equal(t, 3, res.Pagination.TotalPage)
	mockRepo.AssertExpectations(t)
}

// =========================================
// GetByID
// =========================================

func TestService_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(1)).Return(
		&models.UserRole{UserID: 1, RoleID: 2}, nil,
	)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "User role retrieved successfully", res.Message)
	assert.Equal(t, uint(1), res.Data.UserID)
	assert.Equal(t, uint(2), res.Data.RoleID)
	mockRepo.AssertExpectations(t)
}

func TestService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.GetByID(99)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "user role not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// =========================================
// Create
// =========================================

func TestService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByUserIDAndRoleID", uint(1), uint(3)).
		Return(nil, errors.New("record not found"))

	mockRepo.On("Create", mock.MatchedBy(func(ur *models.UserRole) bool {
		return ur.UserID == 1 && ur.RoleID == 3
	})).Run(func(args mock.Arguments) {
		ur := args.Get(0).(*models.UserRole)
		ur.ID = 1
	}).Return(nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.Create(crud_user_role.CreateUserRoleRequest{
		UserID: 1,
		RoleID: 3,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Role assigned to user successfully", res.Message)
	assert.Equal(t, uint(1), res.Data.UserID)
	assert.Equal(t, uint(3), res.Data.RoleID)
	mockRepo.AssertExpectations(t)
}

func TestService_Create_AlreadyAssigned(t *testing.T) {
	mockRepo := new(mocks.Repository)

	existing := &models.UserRole{UserID: 1, RoleID: 1}
	existing.ID = 1
	mockRepo.On("FindByUserIDAndRoleID", uint(1), uint(1)).Return(existing, nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.Create(crud_user_role.CreateUserRoleRequest{
		UserID: 1,
		RoleID: 1,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "role already assigned to this user", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Create_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByUserIDAndRoleID", uint(1), uint(3)).
		Return(nil, errors.New("record not found"))

	mockRepo.On("Create", mock.MatchedBy(func(ur *models.UserRole) bool {
		return ur.UserID == 1 && ur.RoleID == 3
	})).Return(errors.New("db error"))

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	res, err := svc.Create(crud_user_role.CreateUserRoleRequest{
		UserID: 1,
		RoleID: 3,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "failed to assign role to user", err.Error())
	mockRepo.AssertExpectations(t)
}

// =========================================
// Delete
// =========================================

func TestService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(1)).Return(
		&models.UserRole{UserID: 1, RoleID: 1}, nil,
	)
	mockRepo.On("Delete", uint(1)).Return(nil)

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	err := svc.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_NotFound(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	err := svc.Delete(99)

	assert.Error(t, err)
	assert.Equal(t, "user role not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestService_Delete_DBError(t *testing.T) {
	mockRepo := new(mocks.Repository)

	mockRepo.On("FindByID", uint(1)).Return(
		&models.UserRole{UserID: 1, RoleID: 1}, nil,
	)
	mockRepo.On("Delete", uint(1)).Return(errors.New("db error"))

	svc := crud_user_role.NewCrudUserRoleService(mockRepo)
	err := svc.Delete(1)

	assert.Error(t, err)
	assert.Equal(t, "failed to remove role from user", err.Error())
	mockRepo.AssertExpectations(t)
}
