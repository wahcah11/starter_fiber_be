package crud_user_role

// =========================================
// Request
// =========================================

type CreateUserRoleRequest struct {
	UserID uint `json:"user_id" validate:"required"`
	RoleID uint `json:"role_id" validate:"required"`
}

// =========================================
// Response
// =========================================

type UserRoleResponse struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
}

type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type UserRoleListResponse struct {
	Message    string             `json:"message"`
	Data       []UserRoleResponse `json:"data"`
	Pagination PaginationMeta     `json:"pagination"`
}

type UserRoleSingleResponse struct {
	Message string           `json:"message"`
	Data    UserRoleResponse `json:"data"`
}
