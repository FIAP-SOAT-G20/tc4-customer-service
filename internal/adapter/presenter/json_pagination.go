package presenter

type JsonPagination struct {
	Total int64 `json:"total" example:"100"`
	Page  int   `json:"page" example:"1"`
	Limit int   `json:"limit" example:"10"`
}
