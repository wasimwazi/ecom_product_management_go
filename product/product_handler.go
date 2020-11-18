package product

import (
	"database/sql"
	"ecommerce/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gopkg.in/go-playground/validator.v9"
)

//HandlerInterface for product management
type HandlerInterface interface {
	CreateProduct(http.ResponseWriter, *http.Request)
	UpdateProduct(http.ResponseWriter, *http.Request)
	DeleteProduct(http.ResponseWriter, *http.Request)
	GetProduct(http.ResponseWriter, *http.Request)
}

//Handler struct for product management
type Handler struct {
	cs ServiceInterface
}

//NewHTTPHandler to handle product requests
func NewHTTPHandler(db *sql.DB) HandlerInterface {
	return &Handler{
		cs: NewService(db),
	}
}

//CreateProduct function to handle product post request
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /product POST API")
	var request CreateRequest
	req := json.NewDecoder(r.Body)
	err := req.Decode(&request)
	if err != nil {
		log.Println("Error : Decode error(CreateProduct) -", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		log.Println("Error : Validation error(CreateProduct) -", err.Error())
		utils.Fail(w, 400, errors.New("Error validating request").Error())
		return
	}
	product, err := h.cs.CreateProduct(&request)
	if err != nil {
		if err.Error() == utils.ProductExistsError {
			log.Println("Error : Product exists error(CreateProduct) -", err.Error())
			utils.Fail(w, 200, err.Error())
			return
		}
		log.Println("Error : Product creation error(CreateProduct) -", err.Error())
		utils.Fail(w, 500, err.Error())
		return
	}
	log.Println("App : Product created successfully, Product ID = ", product.ID)
	utils.Send(w, 200, product)
}

//UpdateProduct to handle the product post request
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /product PATCH API")
	var request UpdateRequest
	req := json.NewDecoder(r.Body)
	err := req.Decode(&request)
	if err != nil {
		log.Println("Error : Decode error(UpdateProduct) -", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(&request)
	if err != nil {
		log.Println("Error : Validation error (UpdateProduct) -", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	err = h.cs.UpdateProduct(&request)
	if err != nil {
		log.Println("Error : (UpdateProduct) -", err.Error())
		if err.Error() == utils.ProductExistsError {
			utils.Fail(w, 200, err.Error())
			return
		}
		if err.Error() == utils.ProductIDNotExist {
			utils.Fail(w, 400, err.Error())
			return
		}
		if err.Error() == utils.NothingToUpdateInProduct {
			utils.Fail(w, 400, err.Error())
			return
		}
		utils.Fail(w, 500, err.Error())
		return
	}
	message := utils.Message{
		Message: fmt.Sprintf("Product updated successfully, product id = %d", request.ProductID),
	}
	log.Println("App : Product updated successfully, product id -", request.ProductID)
	utils.Send(w, 200, &message)
}

//DeleteProduct to handle the product delete request
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /delete/{product_id} DELETE API")
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		log.Println("Error :", utils.InvalidParameterError, " (DeleteProduct)")
		utils.Fail(w, 400, fmt.Errorf("%s %s", utils.InvalidParameterError, err.Error()).Error())
		return
	}
	err = h.cs.DeleteProduct(productID)
	if err != nil {
		log.Println("Error : error while deleting product (DeleteProduct)")
		if err.Error() == utils.ProductIDNotExist {
			utils.Fail(w, 400, err.Error())
			return
		}
		utils.Fail(w, 500, err.Error())
		return
	}
	message := utils.Message{
		Message: fmt.Sprintf("Product deleted successfully, product id = %d", productID),
	}
	log.Println(message.Message)
	utils.Send(w, 200, &message)
}

// GetProduct to handle the product get request
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /product/{product_id} GET API")
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		log.Println("Error : (GetProduct)", err.Error())
		utils.Fail(w, 400, errors.New(utils.InvalidProductID).Error())
		return
	}
	product, err := h.cs.GetProduct(productID)
	if err != nil {
		log.Println("Error : error fetching product details(GetProduct)", err.Error())
		if err.Error() == utils.ProductIDNotExist {
			utils.Fail(w, 400, err.Error())
			return
		}
		utils.Fail(w, 500, err.Error())
		return
	}
	log.Println("App : Product fetched successfully, product_id : ", productID)
	utils.Send(w, 200, product)
}
