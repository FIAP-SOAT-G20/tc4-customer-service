package dto

type CreateCustomerInput struct {
	Name  string
	Email string
	CPF   string
}

type UpdateCustomerInput struct {
	ID    int
	Name  string
	Email string
}

type GetCustomerInput struct {
	ID int
}

type GetCustomerByCPFInput struct {
	CPF string
}

type DeleteCustomerInput struct {
	ID int
}

type ListCustomersInput struct {
	Name  string
	Page  int
	Limit int
}

type FindCustomerByCPFInput struct {
	CPF string
}
