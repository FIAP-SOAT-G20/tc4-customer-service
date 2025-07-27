package datasource_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FIAP-SOAT-G20/tc4-customer-service/internal/core/domain/entity"
)

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) TestDynamoCreate() {
	tests := []struct {
		name     string
		customer *entity.Customer
		wantErr  bool
	}{
		{
			name: "should create customer successfully",
			customer: &entity.Customer{
				Name:      "John Doe",
				Email:     "john.doe@example.com",
				CPF:       "123.456.789-01",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "should create customer with special characters",
			customer: &entity.Customer{
				Name:      "José da Silva",
				Email:     "jose.silva@example.com",
				CPF:       "987.654.321-09",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			// Clear table before each test
			suite.clearTestTable()

			err := suite.dataSource.Create(suite.ctx, tt.customer)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.customer.ID, "Customer ID should be set after creation")

				// Verify customer was created by finding it
				found, err := suite.dataSource.FindByID(suite.ctx, tt.customer.ID)
				assert.NoError(t, err)
				assert.NotNil(t, found)
				assert.Equal(t, tt.customer.Name, found.Name)
				assert.Equal(t, tt.customer.Email, found.Email)
				assert.Equal(t, tt.customer.CPF, found.CPF)
			}
		})
	}
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) TestDynamoFindByID() {
	customer := &entity.Customer{
		Name:      "Jane Doe",
		Email:     "jane.doe@example.com",
		CPF:       "111.222.333-44",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.dataSource.Create(suite.ctx, customer)
	require.NoError(suite.T(), err)

	tests := []struct {
		name      string
		id        int
		wantFound bool
	}{
		{
			name:      "should find existing customer",
			id:        customer.ID,
			wantFound: true,
		},
		{
			name:      "should return nil for non-existent customer",
			id:        999999,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			result, err := suite.dataSource.FindByID(suite.ctx, tt.id)

			if tt.wantFound {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, customer.Name, result.Name)
				assert.Equal(t, customer.Email, result.Email)
				assert.Equal(t, customer.CPF, result.CPF)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, result)
			}
		})
	}
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) TestDynamoFindByCPF() {
	customer := &entity.Customer{
		Name:      "Bob Smith",
		Email:     "bob.smith@example.com",
		CPF:       "555.666.777-88",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.dataSource.Create(suite.ctx, customer)
	require.NoError(suite.T(), err)

	tests := []struct {
		name      string
		cpf       string
		wantFound bool
	}{
		{
			name:      "should find existing customer by CPF",
			cpf:       "555.666.777-88",
			wantFound: true,
		},
		{
			name:      "should return nil for non-existent CPF",
			cpf:       "999.888.777-66",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			result, err := suite.dataSource.FindByCPF(suite.ctx, tt.cpf)

			if tt.wantFound {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, customer.Name, result.Name)
				assert.Equal(t, customer.Email, result.Email)
				assert.Equal(t, tt.cpf, result.CPF)
			} else {
				assert.NoError(t, err)
				assert.Nil(t, result)
			}
		})
	}
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) TestDynamoFindAll() {
	customers := []*entity.Customer{
		{
			Name:      "Alice Johnson",
			Email:     "alice.johnson@example.com",
			CPF:       "111.111.111-11",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Alice Smith",
			Email:     "alice.smith@example.com",
			CPF:       "222.222.222-22",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Bob Johnson",
			Email:     "bob.johnson@example.com",
			CPF:       "333.333.333-33",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, customer := range customers {
		err := suite.dataSource.Create(suite.ctx, customer)
		require.NoError(suite.T(), err)
	}

	tests := []struct {
		name        string
		filters     map[string]interface{}
		page        int
		limit       int
		wantCount   int
		wantTotal   int64
		checkResult func(*testing.T, []*entity.Customer)
	}{
		{
			name:      "should find all customers without filters",
			filters:   map[string]interface{}{},
			page:      1,
			limit:     10,
			wantCount: 3,
			wantTotal: 3,
			checkResult: func(t *testing.T, result []*entity.Customer) {
				assert.Len(t, result, 3)
			},
		},
		{
			name:      "should filter customers by name",
			filters:   map[string]interface{}{"name": "Alice Johnson"},
			page:      1,
			limit:     10,
			wantCount: 1,
			wantTotal: 1,
			checkResult: func(t *testing.T, result []*entity.Customer) {
				assert.Len(t, result, 1)
				assert.Equal(t, "Alice Johnson", result[0].Name)
			},
		},
		{
			name:      "should return empty result for non-matching filter",
			filters:   map[string]interface{}{"name": "Non-existent"},
			page:      1,
			limit:     10,
			wantCount: 0,
			wantTotal: 0,
			checkResult: func(t *testing.T, result []*entity.Customer) {
				assert.Empty(t, result)
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			result, total, err := suite.dataSource.FindAll(suite.ctx, tt.filters, tt.page, tt.limit)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantTotal, total)
			assert.Len(t, result, tt.wantCount)

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) TestDynamoUpdate() {
	customer := &entity.Customer{
		Name:      "Original Name",
		Email:     "original@example.com",
		CPF:       "123.456.789-00",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.dataSource.Create(suite.ctx, customer)
	require.NoError(suite.T(), err)

	tests := []struct {
		name        string
		updateData  func(*entity.Customer)
		wantErr     bool
		checkResult func(*testing.T, *entity.Customer)
	}{
		{
			name: "should update customer successfully",
			updateData: func(c *entity.Customer) {
				c.Name = "Updated Name"
				c.Email = "updated@example.com"
				c.UpdatedAt = time.Now()
			},
			wantErr: false,
			checkResult: func(t *testing.T, updated *entity.Customer) {
				assert.Equal(t, "Updated Name", updated.Name)
				assert.Equal(t, "updated@example.com", updated.Email)
				assert.Equal(t, "123.456.789-00", updated.CPF)
			},
		},
		{
			name: "should fail to update with non-existent ID",
			updateData: func(c *entity.Customer) {
				c.ID = 999999
				c.Name = "Should not update"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			testCustomer := *customer
			tt.updateData(&testCustomer)

			err := suite.dataSource.Update(suite.ctx, &testCustomer)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tt.checkResult != nil {
					updated, err := suite.dataSource.FindByID(suite.ctx, customer.ID)
					assert.NoError(t, err)
					assert.NotNil(t, updated)
					tt.checkResult(t, updated)
				}
			}
		})
	}
}

func (suite *CustomerDynamoDataSourceIntegrationTestSuite) TestDynamoDelete() {
	customer := &entity.Customer{
		Name:      "To Delete",
		Email:     "delete@example.com",
		CPF:       "999.999.999-99",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.dataSource.Create(suite.ctx, customer)
	require.NoError(suite.T(), err)

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "should delete existing customer",
			id:      customer.ID,
			wantErr: false,
		},
		{
			name:    "should fail to delete non-existent customer",
			id:      999999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			err := suite.dataSource.Delete(suite.ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tt.id == customer.ID {
					found, err := suite.dataSource.FindByID(suite.ctx, customer.ID)
					assert.NoError(t, err)
					assert.Nil(t, found, "Customer should not exist after deletion")
				}
			}
		})
	}
}
