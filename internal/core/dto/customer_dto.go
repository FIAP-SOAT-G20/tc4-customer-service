package dto

type GetCustomerInput struct {
	CPF string
}

type CustomerDatasourceFilter struct {
	CPF *string
}
