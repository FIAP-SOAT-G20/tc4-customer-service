package datasource

import (
	"context"
	"errors"

	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/domain/entity"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/dto"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/core/port"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/infrastructure/database"
	"gorm.io/gorm"
)

type customerDataSource struct {
	db *gorm.DB
}

func NewCustomerDataSource(db database.Database) port.CustomerDataSource {
	return &customerDataSource{db.DB}
}

func (ds *customerDataSource) FindOne(ctx context.Context, filter dto.CustomerDatasourceFilter) (*entity.Customer, error) {
	var customer entity.Customer
	if err := ds.db.WithContext(ctx).First(&customer, filter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}
