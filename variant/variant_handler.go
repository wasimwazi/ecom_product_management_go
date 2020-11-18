package variant

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

//HandlerInterface for variant management
type HandlerInterface interface {
	CreateVariant(http.ResponseWriter, *http.Request)
	UpdateVariant(http.ResponseWriter, *http.Request)
	DeleteVariant(http.ResponseWriter, *http.Request)
	GetVariant(http.ResponseWriter, *http.Request)
	ListVariant(http.ResponseWriter, *http.Request)
}

//Handler struct for variant management
type Handler struct {
	cs ServiceInterface
}

//NewHTTPHandler to handle variant requests
func NewHTTPHandler(db *sql.DB) HandlerInterface {
	return &Handler{
		cs: NewService(db),
	}
}

//CreateVariant function to handle variant post request
func (h *Handler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /variant POST API")
	var request CreateRequest
	req := json.NewDecoder(r.Body)
	err := req.Decode(&request)
	if err != nil {
		log.Println("Error : Decode error(CreateVariant) -", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		log.Println("Error : Validation error(CreateVariant) -", err.Error())
		utils.Fail(w, 400, errors.New("Error validating request").Error())
		return
	}
	variant, err := h.cs.CreateVariant(&request)
	if err != nil {
		if err.Error() == utils.ProductIDNotExist {
			log.Println("Error : Product doen't exist error(CreateVariant) -", err.Error())
			utils.Fail(w, 400, err.Error())
			return
		}
		log.Println("Error : Variant creation error(CreateVariant) -", err.Error())
		utils.Fail(w, 500, err.Error())
		return
	}
	log.Println("App : Variant created successfully, Variant ID = ", variant.ID)
	utils.Send(w, 200, variant)
}

//UpdateVariant to handle the variant post request
func (h *Handler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /variant PATCH API")
	var request UpdateRequest
	req := json.NewDecoder(r.Body)
	err := req.Decode(&request)
	if err != nil {
		log.Println("Error : Decode error(UpdateVariant) -", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(&request)
	if err != nil {
		log.Println("Error : Validation error (UpdateVariant) -", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	err = h.cs.UpdateVariant(&request)
	if err != nil {
		log.Println("Error : (UpdateVariant) -", err.Error())
		if err.Error() == utils.InvalidVariantID {
			utils.Fail(w, 400, err.Error())
			return
		}
		if err.Error() == utils.NothingToUpdateInVariant {
			utils.Fail(w, 400, err.Error())
			return
		}
		utils.Fail(w, 500, err.Error())
		return
	}
	message := utils.Message{
		Message: fmt.Sprintf("Variant updated successfully, variant id = %d", request.VariantID),
	}
	log.Println("App : Variant updated successfully, variant id -", request.VariantID)
	utils.Send(w, 200, &message)
}

//DeleteVariant to handle the variant delete request
func (h *Handler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /delete/{variant_id} DELETE API")
	variantID, err := strconv.Atoi(chi.URLParam(r, "variant_id"))
	if err != nil {
		log.Println("Error :", utils.InvalidParameterError, " (DeleteVariant)")
		utils.Fail(w, 400, fmt.Errorf("%s %s", utils.InvalidParameterError, err.Error()).Error())
		return
	}
	err = h.cs.DeleteVariant(variantID)
	if err != nil {
		log.Println("Error : error while deleting variant (DeleteVariant)")
		if err.Error() == utils.VariantIDNotExist {
			utils.Fail(w, 400, err.Error())
			return
		}
		utils.Fail(w, 500, err.Error())
		return
	}
	message := utils.Message{
		Message: fmt.Sprintf("Variant deleted successfully, variant id = %d", variantID),
	}
	log.Println(message.Message)
	utils.Send(w, 200, &message)
}

// GetVariant to handle variant get request
func (h *Handler) GetVariant(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /product/{product_id}/variant/{variant_id} GET API")
	request, err := validateRequest(r)
	if err != nil {
		log.Println("Error : request validation error (GetVariant)", err.Error())
		utils.Fail(w, 400, err.Error())
		return
	}
	variants, err := h.cs.ListVariant(request)
	if err != nil {
		log.Println("Error : error while fetching variant details(GetVariant)", err.Error())
		if err.Error() == utils.NoDataFoundError {
			utils.Fail(w, 400, err.Error())
			return
		}
		utils.Fail(w, 500, err.Error())
		return
	}
	log.Println("App : variant details fetched successfully, variant id =", request.VariantID)
	utils.Send(w, 200, variants[0])
	return
}

// ListVariant to handle variant get request
func (h *Handler) ListVariant(w http.ResponseWriter, r *http.Request) {
	log.Println("App : /product/{product_id}/variant GET API")
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		log.Println("Error : request validation error (ListVariant)", err.Error())
		utils.Fail(w, 400, utils.InvalidProductID)
		return
	}
	request := &GetRequest{
		ProductID: productID,
	}
	variants, err := h.cs.ListVariant(request)
	if err != nil {
		log.Println("Error : error fetching variants(ListVariant)", err.Error())
		if err.Error() == utils.ProductIDNotExist {
			utils.Fail(w, 400, err.Error())
			return
		}
		utils.Fail(w, 500, err.Error())
		return
	}
	log.Println("App : List of variants in the product fetched successfully, product_id =", request.ProductID)
	utils.Send(w, 200, variants)
	return
}

func validateRequest(r *http.Request) (*GetRequest, error) {
	productID, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		return nil, errors.New(utils.InvalidProductID)
	}
	variantID, err := strconv.Atoi(chi.URLParam(r, "variant_id"))
	if err != nil {
		return nil, errors.New(utils.InvalidVariantID)
	}
	request := GetRequest{
		ProductID: productID,
		VariantID: variantID,
	}
	return &request, nil
}
