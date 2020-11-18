package product

//CreateRequest struct to manage product create request
type CreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	CategoryID  int    `json:"category_id" validate:"required,gt=0"`
}

// CreateResponse product details create response
type CreateResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	CategoryID  int    `json:"category_id"`
}

//UpdateRequest struct to represent the update request
type UpdateRequest struct {
	ProductID   int    `json:"product_id" validate:"required"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

// Variant to represent variant struct
type Variant struct {
	ID              int     `json:"variant_id"`
	Name            string  `json:"variant_name,omitempty"`
	MaxRetailPrice  float64 `json:"max_retail_price"`
	DiscountPrice float64 `json:"discount_price,omitempty"`
	Size            string  `json:"size,omitempty"`
	Color           string  `json:"color,omitempty"`
}

// ProductVariant to represent product struct with variants
type ProductVariant struct {
	ID          int       `json:"product_id"`
	Name        string    `json:"product_name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CategoryID  int       `json:"category_id"`
	Variants    []Variant `json:"variants"`
}

// ProductVariantRow to represent the product variant rows from DB
type ProductVariantRow struct {
	ProductID       int
	ProductName     string
	Description     string
	ImageURL        string
	CategoryID      int
	VariantID       int
	VariantName     string
	MRP  float64
	DiscountPrice float64
	VariantSize     string
	VariantColor   string
}
