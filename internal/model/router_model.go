package model

type CreateRouterRequest struct {
	Name     string `json:"name" validate:"required"`
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,numeric"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateRouterRequest struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RouterResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Status    string `json:"status"` // online, offline
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
