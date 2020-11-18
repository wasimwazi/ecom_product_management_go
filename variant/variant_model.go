package variant

//CreateRequest struct to manage variant create request
type CreateRequest struct {
	Name          string  `json:"name"`
	MRP           float64 `json:"max_retail_price" validate:"required"`
	DiscountPrice float64 `json:"discount_price"`
	Size          string  `json:"size"`
	Color         string  `json:"color"`
	ProductID     int     `json:"product_id" validate:"required,gt=0"`
}

// CreateResponse variant details create response
type CreateResponse struct {
	ID            int     `json:"id"`
	Name          string  `json:"name,omitempty"`
	MRP           float64 `json:"max_retail_price"`
	DiscountPrice float64 `json:"discount_price,omitempty"`
	Size          string  `json:"size,omitempty"`
	Color         string  `json:"color,omitempty"`
	ProductID     int     `json:"product_id"`
}

//UpdateRequest struct to represent the variant update request
type UpdateRequest struct {
	VariantID     int     `json:"variant_id" validate:"required"`
	Name          string  `json:"name"`
	MRP           float64 `json:"max_retail_price"`
	DiscountPrice float64 `json:"discount_price"`
	Size          string  `json:"size"`
	Color         string  `json:"color"`
}

// GetRequest to represent get variant request
type GetRequest struct {
	ProductID int
	VariantID int
}

// Variant to represent variant struct
type Variant struct {
	ID              int     `json:"variant_id"`
	Name            string  `json:"name,omitempty"`
	MRP             float64 `json:"max_retail_price"`
	DiscountPrice float64 `json:"discount_price,omitempty"`
	Size            string  `json:"size,omitempty"`
	Color          string  `json:"color,omitempty"`
	ProductID       int     `json:"product_id"`
}
