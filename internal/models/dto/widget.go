package dto

// dto/ holds API request/response shapes only.

type CreateWidgetReq struct {
	Name  string `json:"name" binding:"required"`
	Units int    `json:"units" binding:"required,gt=0"`
}

type WidgetResponse struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	Price     float64 `json:"price"`
	CreatedAt int64   `json:"created_at"`
}
