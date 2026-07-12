package model

type UserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"username"`
	RoleID      int    `json:"role_id"`
	CompanyName string `json:"company_name"`
	Email       string `json:"email"`
	Token       string `json:"token"`
	RegionId    string `json:"region_id"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type UserResponse1 struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Token     string `json:"token,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UpdatedAt int64  `json:"updated_at,omitempty"`
}

type ForUserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"username"`
	RoleID      int    `json:"role_id"`
	CompanyName string `json:"company_name"`
	Email       string `json:"email"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type VerifyUserRequest struct {
	Token string `json:"required,omitempty"`
}

type RegisterUserRequest struct {
	ID          string `json:"id" validate:"required,max=100"`
	Password    string `json:"password" validate:"required,max=100"`
	Name        string `json:"name" validate:"required,max=100"`
	Email       string `json:"email" validate:"required,max=100"`
	CompanyName string `json:"company_name" validate:"required,max=100"`
	RegionId    string `json:"region_id" validate:"required,max=50"`
}

type UpdateUserRequest struct {
	ID          string `json:"-" validate:"required,max=100"`
	Password    string `json:"password" validate:"max=100"`
	Name        string `json:"username" validate:"max=100"`
	Email       string `json:"email" validate:"max=100"`
	CompanyName string `json:"company_name" validate:"max=100"`
	RegionId    string `json:"region_id" validate:"max=50"`
}

type UpdateUserAllRequest struct {
	ID          string `json:"-" validate:"required,max=100"`
	Password    string `json:"password" validate:"max=100"`
	Name        string `json:"username" validate:"max=100"`
	Email       string `json:"email" validate:"max=100"`
	CompanyName string `json:"company_name" validate:"max=250"`
	RegionId    string `json:"region_id" validate:"max=50"`
}

type UserSearchRequest struct {
	ID          string `json:"id" validate:"max=100"`
	Name        string `json:"name" validate:"max=100"`
	Email       string `json:"email" validate:"max=100"`
	CompanyName string `json:"company_name" validate:"max=100"`
	RegionId    string `json:"region_id" validate:"max=50"`
	Page        int    `json:"page" validate:"min=1"`
	Size        int    `json:"size" validate:"min=1,max=100"`
}

type LoginUserRequest struct {
	ID       string `json:"id" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}
type LogoutUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetUserID struct {
	ID string `json:"-" validate:"required,max=100"`
}

type RequestUpdateForAdmin struct {
	ID     string `json:"-" validate:"required,max=100"`
	RoleId int    `json:"role-id" validate:"required,max=10"`
}

type FindUserIDRequest struct {
	ID       string `json:"id" validate:"max=100"`
	RegionId string `json:"region_id" validate:"max=50"`
}

type DeleteUsersIDRequest struct {
	ID string `json:"-" validate:"required,max=100"`
}
