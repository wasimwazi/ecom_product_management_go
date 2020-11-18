package product

import "database/sql"

//RepoInterface for DB operations
type RepoInterface interface {
	CheckProductNameExists(string) (bool, error)
	CreateProduct(*CreateRequest) (*CreateResponse, error)
	CheckCategoryExists(int) (bool, error)
	IsProductIDExists(int) (bool, error)
	UpdateProduct(*UpdateRequest) error
	DeleteProduct(int) error
	GetProduct(int) ([]ProductVariantRow, error)
}

//NewRepo returns repository interface
func NewRepo(db *sql.DB) RepoInterface {
	return &Repo{
		DB: db,
	}
}
