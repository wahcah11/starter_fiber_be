package crud_permission

// =========================================
// Request
// =========================================

type CreatePermissionRequest struct {
	RoleID uint   `json:"role_id" validate:"required"`
	Name   string `json:"name"    validate:"required,max=100"`
}

type UpdatePermissionRequest struct {
	RoleID uint   `json:"role_id" validate:"required"`
	Name   string `json:"name"    validate:"required,max=100"`
}

// =========================================
// Response
// =========================================

type PermissionResponse struct {
	ID     uint   `json:"id"`
	RoleID uint   `json:"role_id"`
	Name   string `json:"name"`
}

type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type PermissionListResponse struct {
	Message    string               `json:"message"`
	Data       []PermissionResponse `json:"data"`
	Pagination PaginationMeta       `json:"pagination"`
}

type PermissionSingleResponse struct {
	Message string             `json:"message"`
	Data    PermissionResponse `json:"data"`
}
