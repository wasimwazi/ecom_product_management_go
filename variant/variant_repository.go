package variant

import "database/sql"

//RepoInterface for DB operations
type RepoInterface interface {
	CreateVariant(*CreateRequest) (*CreateResponse, error)
	CheckProductExists(int) (bool, error)
	IsVariantIDExists(int) (bool, error)
	UpdateVariant(*UpdateRequest) error
	DeleteVariant(int) error
	ListVariant(*GetRequest) ([]Variant, error)
}

//NewRepo returns repository interface 
func NewRepo(db *sql.DB) RepoInterface {
	return &Repo{
		DB: db,
	}
}