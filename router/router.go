package router

import (
	"database/sql"
	"ecommerce/category"
	"ecommerce/product"
	"ecommerce/variant"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

//Router interface for routing
type Router interface {
	Setup() *chi.Mux
}

//ChiRouter chi struct
type ChiRouter struct {
	DB *sql.DB
}

//NewRouter returns a router struct
func NewRouter(db *sql.DB) Router {
	return &ChiRouter{
		DB: db,
	}
}

//Setup function to initialize the routing
func (router *ChiRouter) Setup() *chi.Mux {
	cr := chi.NewRouter()
	categoryHandler := category.NewHTTPHandler(router.DB)
	productHandler := product.NewHTTPHandler(router.DB)
	variantHandler := variant.NewHTTPHandler(router.DB)
	cr.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	cr.Post("/category", categoryHandler.CreateCategory)
	cr.Patch("/category", categoryHandler.UpdateCategory)
	cr.Get("/category", categoryHandler.ListCategory)
	cr.Delete("/category/{category_id}", categoryHandler.DeleteCategory)
	cr.Post("/product", productHandler.CreateProduct)
	cr.Patch("/product", productHandler.UpdateProduct)
	cr.Get("/product/{product_id}", productHandler.GetProduct)
	cr.Delete("/product/{product_id}", productHandler.DeleteProduct)
	cr.Post("/variant", variantHandler.CreateVariant)
	cr.Patch("/variant", variantHandler.UpdateVariant)
	cr.Get("/product/{product_id}/variant/{variant_id}", variantHandler.GetVariant)
	cr.Get("/product/{product_id}/variant", variantHandler.ListVariant)
	cr.Delete("/variant/{variant_id}", variantHandler.DeleteVariant)
	return cr
}
