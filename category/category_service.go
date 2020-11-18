package category

import (
	"database/sql"
	"ecommerce/utils"
	"errors"
	"log"
)

//Visited to mark the category visited while looping through categories
var Visited = make(map[int]bool)

//ServiceInterface is category service interface
type ServiceInterface interface {
	CreateCategory(*CreateRequest) (*CreateResponse, error)
	UpdateCategory(*UpdateRequest) error
	DeleteCategory(int) error
	ListCategory() (*[]CategoryList, error)
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

//CreateCategory service function to create category
func (service *Service) CreateCategory(req *CreateRequest) (*CreateResponse, error) {
	categoryExists, err := service.repo.CheckCategoryNameExists(req.Name)
	if err != nil {
		return nil, err
	}
	if categoryExists {
		return nil, errors.New(utils.CategoryExistsError)
	}
	category, err := service.repo.CreateCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}

//UpdateCategory to update the category
func (service *Service) UpdateCategory(request *UpdateRequest) error {
	isExist, err := service.repo.IsCategoryIDExists(request.CategoryID)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New(utils.InvalidCategoryID)
	}
	if len(request.Name) <= 0 && request.ParentID == 0 {
		return errors.New(utils.NothingToUpdateInCategory)
	}
	categoryExists, err := service.repo.CheckCategoryNameExists(request.Name)
	if err != nil {
		return err
	}
	if categoryExists {
		return errors.New(utils.CategoryExistsError)
	}
	return service.repo.UpdateCategory(request)
}

//DeleteCategory to delete a category
func (service Service) DeleteCategory(categoryID int) error {
	isCategoryExist, err := service.repo.IsCategoryIDExists(categoryID)
	if err != nil {
		return err
	}
	if !isCategoryExist {
		return errors.New(utils.CategoryNOTExistsError)
	}
	isSubCategoryExists, err := service.repo.IsSubCategoryExist(categoryID)
	if err != nil {
		return err
	}
	if isSubCategoryExists {
		return errors.New(utils.SubCategoryExists)
	}
	isExistProduct, err := service.repo.IsProductExist(categoryID)
	if err != nil {
		return err
	}
	if isExistProduct {
		return errors.New(utils.ProductExistCategoryError)
	}
	return service.repo.DeleteCategory(categoryID)
}

//ListCategory lists all the categories and its child elements
func (service *Service) ListCategory() (*[]CategoryList, error) {
	Visited = map[int]bool{}
	//Get the details of existing categories
	categoryDetails, err := service.repo.GetCategories()
	log.Println(categoryDetails)
	if err != nil {
		return nil, err
	}
	var categoryIDs []int                   //to store all the category IDs
	var mainCategories []int                //to store main categories which doesn't have a child
	categoryChildMap := make(map[int][]int) //to store category and its child relation
	categoryNameMap := make(map[int]string) //to map category and its name
	//generating categoryChildMap, categoryNameMap, categoryIDs and mainCategories
	for _, v := range *categoryDetails {
		if v.ParentID == 0 {
			mainCategories = append(mainCategories, v.ID)
			categoryNameMap[v.ID] = v.Name
			categoryIDs = append(categoryIDs, v.ID)
			continue
		}
		categoryChildMap[v.ParentID] = append(categoryChildMap[v.ParentID], v.ID)
		categoryNameMap[v.ID] = v.Name
		categoryIDs = append(categoryIDs, v.ID)
	}
	//To get all the product and its variants
	productVariantForCategory, err := service.repo.GetProductVariantForEachCategory(categoryIDs)
	if err != nil {
		return nil, err
	}
	log.Println(productVariantForCategory)
	//generating category and its associated products mapping
	categoryProductMap := make(map[int][]Product)
	for _, v := range productVariantForCategory {
		categoryProductMap[v.CategoryID] = append(categoryProductMap[v.CategoryID], v)
	}
	log.Println(categoryProductMap)
	var categoryList []CategoryList //Final result category listing
	for _, categoryID := range mainCategories {

		catList := formatCategory(categoryID, categoryProductMap, categoryChildMap, categoryNameMap)
		if catList.CategoryID != 0 {
			categoryList = append(categoryList, catList)
		}
	}
	return &categoryList, nil
}

//To format the categories and its sub categories
func formatCategory(categoryID int, categoryProductMap map[int][]Product, categoryChildMap map[int][]int, categoryNameMap map[int]string) CategoryList {
	//if already visited, return null for the category
	if Visited[categoryID] {
		return CategoryList{}
	}
	Visited[categoryID] = true
	var catList CategoryList
	catList.Products = categoryProductMap[categoryID]
	catList.Name = categoryNameMap[categoryID]
	catList.CategoryID = categoryID
	for _, childID := range categoryChildMap[categoryID] {
		cList := formatCategory(childID, categoryProductMap, categoryChildMap, categoryNameMap)
		catList.Categories = append(catList.Categories, cList)
	}
	return catList
}
