package variant

import (
	"database/sql"
	"ecommerce/utils"
	"errors"
)

//ServiceInterface is variant service interface
type ServiceInterface interface {
	CreateVariant(*CreateRequest) (*CreateResponse, error)
	UpdateVariant(*UpdateRequest) error
	DeleteVariant(int) error
	ListVariant(*GetRequest) ([]Variant, error)
}

//Service struct for service functionalities
type Service struct {
	repo RepoInterface
}

//NewService :
func NewService(db *sql.DB) ServiceInterface {
	return &Service{
		repo: NewRepo(db),
	}
}

//CreateVariant service function to create a variant
func (service *Service) CreateVariant(request *CreateRequest) (*CreateResponse, error) {
	isValidProduct, err := service.repo.CheckProductExists(request.ProductID)
	if err != nil {
		return nil, err
	}
	if !isValidProduct {
		return nil, errors.New(utils.ProductIDNotExist)
	}
	variant, err := service.repo.CreateVariant(request)
	if err != nil {
		return nil, err
	}
	return variant, nil
}

//UpdateVariant to update the variant
func (service *Service) UpdateVariant(request *UpdateRequest) error {
	isExist, err := service.repo.IsVariantIDExists(request.VariantID)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New(utils.InvalidVariantID)
	}
	if len(request.Name) <= 0 && len(request.Size) <= 0 && len(request.Color) <= 0 && request.MRP == 0 && request.DiscountPrice == 0 {
		return errors.New(utils.NothingToUpdateInVariant)
	}
	return service.repo.UpdateVariant(request)
}

//DeleteVariant to delete the given variant
func (service *Service) DeleteVariant(variantID int) error {
	isVariantExist, err := service.repo.IsVariantIDExists(variantID)
	if err != nil {
		return err
	}
	if !isVariantExist {
		return errors.New(utils.VariantIDNotExist)
	}
	return service.repo.DeleteVariant(variantID)
}

// ListVariant : to list out all variants of a product
func (service *Service) ListVariant(request *GetRequest) ([]Variant, error) {
	isValid, err := service.repo.CheckProductExists(request.ProductID)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, errors.New(utils.ProductIDNotExist)
	}
	return service.repo.ListVariant(request)
}
