package presenter

import "encoding/json"

type CustomerJsonResponse struct {
	ID        int    `json:"id" example:"1"`
	Name      string `json:"name" example:"John Doe"`
	Email     string `json:"email" example:"john.doe@email.com"`
	CPF       string `json:"cpf" example:"123.456.789-00"`
	CreatedAt string `json:"created_at" example:"2024-02-09T10:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-02-09T10:00:00Z"`
}

func (r CustomerJsonResponse) String() string {
	o, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(o)
}

type CustomerJsonPaginatedResponse struct {
	JsonPagination
	Customers []CustomerJsonResponse `json:"customers"`
}

func (r CustomerJsonPaginatedResponse) String() string {
	o, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(o)
}
