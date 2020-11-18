package category

import (
	"database/sql"
	"ecommerce/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gopkg.in/go-playground/validator.v9"
)

//HandlerInterface for category management
type HandlerInterface interface {
	CreateCategory(http.ResponseWriter, *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
	ListCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
}

//Handler struct for category management
type Handler struct {
	cs ServiceInterface
}

//NewHTTPHandler to handle category requests
func NewHTTPHandler(db *sql.DB) HandlerInterface {
	return &Handler{
		cs: NewService(db),
	}
}

//CreateCategory to handle the category post request
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /category POST API")
	var request CreateRequest
	req := json.NewDecoder(r.Body)
	err := req.Decode(&request)
	if err != nil {
		log.Println("Error : Decode error(CreateCategory) -", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		log.Println("Error : Validation error(CreateCategory) -", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	category, err := h.cs.CreateCategory(&request)
	if err != nil {
		if err.Error() == utils.CategoryExistsError {
			log.Println("Error : Category exists error(CreateCategory) -", err.Error())
			utils.Fail(w, 200, err.Error())
			return
		}
		log.Println("Error : Create category error(CreateCategory) -", err.Error())
		utils.Fail(w, 500, err.Error())
		return
	}
	log.Println("App : Category created successfully, Category ID = ", category.ID)
	utils.Send(w, 200, category)
}

//UpdateCategory to handle the category post request
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /category PATCH API")
	var request UpdateRequest
	req := json.NewDecoder(r.Body)
	err := req.Decode(&request)
	if err != nil {
		log.Println("Error : Decode error(UpdateCategory) -", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	err = h.cs.UpdateCategory(&request)
	if err != nil {
		log.Println("Error : (UpdateCategory) -", err.Error())
		if err.Error() == utils.CategoryExistsError {
			utils.Fail(w, 200, err.Error())
			return
		}
		if err.Error() == utils.InvalidCategoryID {
			utils.Fail(w, 400, err.Error())
			return
		}
		if err.Error() == utils.NothingToUpdateInCategory {
			utils.Fail(w, 400, err.Error())
			return
		}
		utils.Fail(w, 500, err.Error())
		return
	}
	message := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Category updated successfully, category id = %d", request.CategoryID),
	}
	log.Println("App : Category updated successfully, category id -", request.CategoryID)
	utils.Send(w, 200, &message)
}

//DeleteCategory to list all the categories
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /category/{category_id} DELETE API")
	categoryID, err := strconv.Atoi(chi.URLParam(r, "category_id"))
	if err != nil {
		log.Println("Error :", utils.InvalidParameterError, " (DeleteCategory)")
		utils.Fail(w, 400, fmt.Errorf("%s %s", utils.InvalidParameterError, err.Error()).Error())
		return
	}
	err = h.cs.DeleteCategory(categoryID)
	if err != nil {
		log.Println("Error : error while deleting category (DeleteCategory)")
		if err.Error() == utils.CategoryNOTExistsError {
			utils.Fail(w, 400, err.Error())
			return
		}
		if err.Error() == utils.SubCategoryExists {
			utils.Fail(w, 500, err.Error())
			return
		}
		if err.Error() == utils.ProductExistCategoryError {
			utils.Fail(w, 500, err.Error())
			return
		}
	}
	message := utils.Message{
		Message: fmt.Sprintf("Category deleted successfully, category id = %d", categoryID),
	}
	log.Println(message.Message)
	utils.Send(w, 200, &message)
}

//ListCategory to list all the categories
func (h *Handler) ListCategory(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /category GET API")
	categoryList, err := h.cs.ListCategory()
	if err != nil {
		log.Println("Error : category listing error(ListCategory) -", err.Error())
		utils.Fail(w, 500, err.Error())
	}
	log.Println("App : Category listed successfully")
	utils.Send(w, 200, &categoryList)
}