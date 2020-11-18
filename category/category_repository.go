package category

import "database/sql"

//RepoInterface for DB operations
type RepoInterface interface {
	CheckCategoryNameExists(string) (bool, error)
	CreateCategory(*CreateRequest) (*CreateResponse, error)
	IsCategoryIDExists(int) (bool, error)
	UpdateCategory(*UpdateRequest) error
	DeleteCategory(int) error
	IsSubCategoryExist(int) (bool, error)
	IsProductExist(int) (bool, error)
	GetProductVariantForEachCategory([]int) ([]Product, error)
	GetCategories() (*[]Category, error)
}

//NewRepo returns repository interface
func NewRepo(db *sql.DB) RepoInterface {
	return &Repo{
		DB: db,
	}
}
