package dto

type CreateCustomerInput struct {
	Name  string
	Email string
	CPF   string
}

type UpdateCustomerInput struct {
	ID    string
	Name  string
	Email string
}

type GetCustomerInput struct {
	ID string
}

type GetCustomerByCPFInput struct {
	CPF string
}

type DeleteCustomerInput struct {
	ID string
}

type ListCustomersInput struct {
	Name  string
	Page  int
	Limit int
}

type FindCustomerByCPFInput struct {
	CPF string
}
