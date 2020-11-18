package variant

import (
	"database/sql"
	"ecommerce/utils"
	"errors"
	"fmt"
	"strings"
)

//Repo is the DB repository struct
type Repo struct {
	DB *sql.DB
}

//CheckProductExists function to check if the product with the given ID exists
func (repo *Repo) CheckProductExists(productID int) (bool, error) {
	var count int
	query := `
		SELECT
			count(*)
		FROM
			tbl_product
		WHERE
			product_id = $1
		AND 
			deleted_at IS NULL
	`
	err := repo.DB.QueryRow(query, productID).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func getNullFloat64(value float64) sql.NullFloat64 {
	if value == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: value,
		Valid:   true,
	}
}

//CreateVariant to create a variant in DB
func (repo *Repo) CreateVariant(request *CreateRequest) (*CreateResponse, error) {
	var createResponse CreateResponse
	var name, size, color sql.NullString
	var discountPrice sql.NullFloat64
	query := `
		INSERT INTO 
			tbl_variant (name, max_retail_price, discount_price, size, color, product_id, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING
			variant_id, name, max_retail_price, discount_price, size, color, product_id
	`
	row := repo.DB.QueryRow(query, request.Name, request.MRP, getNullFloat64(request.DiscountPrice), request.Size, request.Color, request.ProductID)
	err := row.Scan(&createResponse.ID, &name, &createResponse.MRP, &discountPrice, &size, &color, &createResponse.ProductID)
	if err != nil {
		return nil, err
	}
	if name.Valid {
		createResponse.Name = name.String
	}
	if size.Valid {
		createResponse.Size = size.String
	}
	if color.Valid {
		createResponse.Color = color.String
	}
	if discountPrice.Valid {
		createResponse.DiscountPrice = discountPrice.Float64
	}
	return &createResponse, nil
}

//IsVariantIDExists function to check if the variant ID already exists
func (repo *Repo) IsVariantIDExists(id int) (bool, error) {
	var count int
	query := `
		SELECT
			count(*)
		FROM
			tbl_variant
		WHERE
			variant_id = $1
		AND 
			deleted_at IS NULL
	`
	err := repo.DB.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

//UpdateVariant to update a variant
func (repo *Repo) UpdateVariant(request *UpdateRequest) error {
	var slice []string
	if len(request.Name) > 0 && request.Name != utils.EmptyString {
		slice = append(slice, fmt.Sprintf(" name = '%s' ", request.Name))
	}
	if len(request.Size) > 0 && request.Size != utils.EmptyString {
		slice = append(slice, fmt.Sprintf(" size = '%s' ", request.Size))
	}
	if len(request.Color) > 0 && request.Color != utils.EmptyString {
		slice = append(slice, fmt.Sprintf(" color = '%s' ", request.Color))
	}
	if request.DiscountPrice != 0 {
		slice = append(slice, fmt.Sprintf(" discount_price = %f ", request.DiscountPrice))
	}
	if request.MRP != 0 {
		slice = append(slice, fmt.Sprintf(" max_retail_price = %f ", request.MRP))
	}
	slice = append(slice, fmt.Sprintf(" updated_at = NOW() "))
	updateQuery := strings.Join(slice, ", ")
	mainQuery := `
		UPDATE
			tbl_variant
		SET
			%s
		WHERE
			variant_id = $1
		AND 
			deleted_at IS NULL
	`
	query := fmt.Sprintf(mainQuery, updateQuery)
	result, err := repo.DB.Exec(query, request.VariantID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New(utils.InvalidVariantID)
	}
	return err
}

// DeleteVariant to delete the variant from DB
func (repo *Repo) DeleteVariant(variantID int) error {
	query := `
		UPDATE
			tbl_variant
		SET
			deleted_at = NOW()
		WHERE
			variant_id = $1
		AND 
			deleted_at IS NULL
	`
	result, err := repo.DB.Exec(query, variantID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New(utils.InvalidVariantID)
	}
	return nil
}

// ListVariant is the DB function to list variants
func (repo *Repo) ListVariant(request *GetRequest) ([]Variant, error) {
	var variants []Variant
	var name, size, color sql.NullString
	var discountPrice sql.NullFloat64
	var subQuery string
	if request.VariantID != 0 {
		subQuery = fmt.Sprintf(" AND variant_id = %d ", request.VariantID)
	}
	query := `
		SELECT
			variant_id, name, max_retail_price, discount_price, size, color
		FROM
			tbl_variant
		WHERE
			product_id = $1
			%s
		AND 
			deleted_at IS NULL
	`
	mainQuery := fmt.Sprintf(query, subQuery)
	rows, err := repo.DB.Query(mainQuery, request.ProductID)
	if err != nil {
		return nil, err
	}
	var variant Variant
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&variant.ID, &name, &variant.MRP, &discountPrice, &size, &color)
		if err != nil {
			return nil, err
		}
		if name.Valid {
			variant.Name = name.String
		}
		if size.Valid {
			variant.Size = size.String
		}
		if color.Valid {
			variant.Color = color.String
		}
		if discountPrice.Valid {
			variant.DiscountPrice = discountPrice.Float64
		}
		variant.ProductID = request.ProductID
		variants = append(variants, variant)
	}
	if len(variants) <= 0 {
		return nil, errors.New(utils.NoDataFoundError)
	}
	return variants, nil
}
