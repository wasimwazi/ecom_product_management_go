package category

//CreateRequest to represent the category post request
type CreateRequest struct {
	Name     string `json:"name" validate:"required"`
	ParentID int    `json:"parent_id" validate:"omitempty,gt=0"`
}

//CreateResponse to represent the category post request
type CreateResponse struct {
	ID       int    `json:"category_id"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
}

//UpdateRequest to represent category update request
type UpdateRequest struct {
	CategoryID int    `json:"category_id" validate:"required"`
	Name       string `json:"name"`
	ParentID   int    `json:"parent_id" validate:"omitempty, gt=0"`
}

//Variant to represent variant struct
type Variant struct {
	VariantID     int     `json:"variant_id"`
	Name          string  `json:"name,omitempty"`
	MRP           float64 `json:"max_retail_price"`
	DiscountPrice float64 `json:"discount_price,omitempty"`
	Size          string  `json:"size,omitempty"`
	Color         string  `json:"color,omitempty"`
}

//Product to represent product struct
type Product struct {
	ProductID   int       `json:"product_id"`
	Name        string    `json:"product_name"`
	Description string    `json:"description,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	CategoryID  int       `json:"category_id"`
	Variants    []Variant `json:"variants,omitempty"`
}

//CategoryList to represent the category listing struct
type CategoryList struct {
	CategoryID int            `json:"category_id"`
	Name       string         `json:"category_name"`
	Products   []Product      `json:"products"`
	Categories []CategoryList `json:"categories"`
}

//CategoryRelationship represent the category sub category relationships
type CategoryRelationship struct {
	CategoryID  int
	Ancestor    []int
	level       int
	FirstParent int
}

// ProductVariantRow to represent product variant combination row
type ProductVariantRow struct {
	ProductID     int
	ProductName   string
	Description   string
	ImageURL      string
	CategoryID    int
	VariantID     int
	VariantName   string
	MRP           float64
	DiscountPrice float64
	Size          string
	Color         string
}

//Category to represent categories
type Category struct {
	ID       int
	Name     string
	ParentID int
}
