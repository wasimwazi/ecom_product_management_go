package product

import (
	"database/sql"
	"ecommerce/utils"
	"errors"
)

//ServiceInterface is product service interface
type ServiceInterface interface {
	CreateProduct(*CreateRequest) (*CreateResponse, error)
	UpdateProduct(*UpdateRequest) error
	DeleteProduct(int) error
	GetProduct(int) (*ProductVariant, error)
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

//CreateProduct service function to create a product
func (service *Service) CreateProduct(request *CreateRequest) (*CreateResponse, error) {
	categoryExists, err := service.repo.CheckCategoryExists(request.CategoryID)
	if err != nil {
		return nil, err
	}
	if !categoryExists {
		return nil, errors.New(utils.CategoryNOTExistsError)
	}
	productExists, err := service.repo.CheckProductNameExists(request.Name)
	if err != nil {
		return nil, err
	}
	if productExists {
		return nil, errors.New(utils.ProductExistsError)
	}
	product, err := service.repo.CreateProduct(request)
	if err != nil {
		return nil, err
	}
	return product, nil
}

//UpdateProduct to update the product
func (service *Service) UpdateProduct(request *UpdateRequest) error {
	isExist, err := service.repo.IsProductIDExists(request.ProductID)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New(utils.ProductIDNotExist)
	}
	if len(request.Name) <= 0 && len(request.Description) <= 0 && len(request.ImageURL) <= 0 {
		return errors.New(utils.NothingToUpdateInCategory)
	}
	productExists, err := service.repo.CheckProductNameExists(request.Name)
	if err != nil {
		return err
	}
	if productExists {
		return errors.New(utils.ProductExistsError)
	}
	return service.repo.UpdateProduct(request)
}

//DeleteProduct to delete the given product
func (service *Service) DeleteProduct(productID int) error {
	isProductExist, err := service.repo.IsProductIDExists(productID)
	if err != nil {
		return err
	}
	if !isProductExist {
		return errors.New(utils.ProductIDNotExist)
	}
	return service.repo.DeleteProduct(productID)
}

// GetProduct  to get a product
func (service *Service) GetProduct(productID int) (*ProductVariant, error) {
	isExist, err := service.repo.IsProductIDExists(productID)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.New(utils.ProductIDNotExist)
	}
	productDetails, err := service.repo.GetProduct(productID)
	if err != nil {
		return nil, err
	}
	var (
		product  ProductVariant
		variant  Variant
		variants []Variant
	)
	for _, row := range productDetails {
		if row.VariantID != 0 {
			variant.ID = row.VariantID
			variant.Name = row.VariantName
			variant.MaxRetailPrice = row.MRP
			variant.DiscountPrice = row.DiscountPrice
			variant.Size = row.VariantSize
			variant.Color = row.VariantColor
			variants = append(variants, variant)
		}
	}
	product.ID = productDetails[0].ProductID
	product.Name = productDetails[0].ProductName
	product.Description = productDetails[0].Description
	product.ImageURL = productDetails[0].ImageURL
	product.CategoryID = productDetails[0].CategoryID
	product.Variants = variants
	return &product, nil
}
