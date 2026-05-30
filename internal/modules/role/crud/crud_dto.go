package crud

// =========================================
// Request
// =========================================

type CreateRoleRequest struct {
	Name           string `json:"name"            validate:"required,max=100"`
	SystemFunction string `json:"system_function" validate:"required,max=100"`
	IsDefault      bool   `json:"is_default"`
}

type UpdateRoleRequest struct {
	Name           string `json:"name"            validate:"required,max=100"`
	SystemFunction string `json:"system_function" validate:"required,max=100"`
	IsDefault      bool   `json:"is_default"`
}

// =========================================
// Response
// =========================================

type RoleResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	SystemFunction string `json:"system_function"`
	IsDefault      bool   `json:"is_default"`
}

// PaginationMeta berisi informasi halaman
type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

// RoleListResponse digunakan untuk response list dengan pagination
type RoleListResponse struct {
	Message    string         `json:"message"`
	Data       []RoleResponse `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// RoleSingleResponse digunakan untuk response single data
type RoleSingleResponse struct {
	Message string       `json:"message"`
	Data    RoleResponse `json:"data"`
}
