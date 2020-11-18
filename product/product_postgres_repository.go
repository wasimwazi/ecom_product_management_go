package product

import (
	"database/sql"
	"ecommerce/utils"
	"errors"
	"fmt"
	"strings"
)

//Repo is the DB repo struct
type Repo struct {
	DB *sql.DB
}

//CheckCategoryExists function to check if the given category exist in our DB
func (repo *Repo) CheckCategoryExists(categoryID int) (bool, error) {
	var count int
	query := `
		SELECT
			count(*)
		FROM
			tbl_category
		WHERE
			category_id = $1
		AND 
			deleted_at IS NULL
	`
	err := repo.DB.QueryRow(query, categoryID).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

//CheckProductNameExists function to check if the product with the given name already exists
func (repo *Repo) CheckProductNameExists(name string) (bool, error) {
	var count int
	query := `
		SELECT
			count(*)
		FROM
			tbl_product
		WHERE
			name = $1
		AND 
			deleted_at IS NULL
	`
	err := repo.DB.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

//CreateProduct - DB function to create product
func (repo *Repo) CreateProduct(request *CreateRequest) (*CreateResponse, error) {
	var createResponse CreateResponse
	var description, imageURL sql.NullString
	query := `
		INSERT INTO 
			tbl_product (name, description, image_url, category_id, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, NOW(), NOW())
		RETURNING
			product_id, name, description, image_url, category_id
	`
	row := repo.DB.QueryRow(query, request.Name, request.Description, request.ImageURL, request.CategoryID)
	err := row.Scan(&createResponse.ID, &createResponse.Name, &description, &imageURL, &createResponse.CategoryID)
	if err != nil {
		return nil, err
	}
	if description.Valid {
		createResponse.Description = description.String
	}
	if imageURL.Valid {
		createResponse.ImageURL = imageURL.String
	}
	return &createResponse, nil
}

//IsProductIDExists function to check if the product ID already exists
func (repo *Repo) IsProductIDExists(id int) (bool, error) {
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
	err := repo.DB.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

//UpdateProduct to update a category
func (repo *Repo) UpdateProduct(request *UpdateRequest) error {
	var slice []string
	if len(request.Name) > 0 && request.Name != utils.EmptyString {
		slice = append(slice, fmt.Sprintf(" name = '%s' ", request.Name))
	}
	if len(request.Description) > 0 && request.Description != utils.EmptyString {
		slice = append(slice, fmt.Sprintf(" description = '%s' ", request.Description))
	}
	if len(request.ImageURL) > 0 && request.ImageURL != utils.EmptyString {
		slice = append(slice, fmt.Sprintf(" image_url = '%s' ", request.ImageURL))
	}
	slice = append(slice, fmt.Sprintf(" updated_at = NOW() "))
	updateQuery := strings.Join(slice, ", ")
	mainQuery := `
		UPDATE
			tbl_product
		SET
			%s
		WHERE
			product_id = $1
		AND 
			deleted_at IS NULL
	`
	query := fmt.Sprintf(mainQuery, updateQuery)
	result, err := repo.DB.Exec(query, request.ProductID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New(utils.InvalidProductID)
	}
	return err
}

// DeleteProduct function to remove a product from DB
func (repo *Repo) DeleteProduct(productID int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return err
	}
	query := `
		UPDATE
			tbl_product
		SET
			deleted_at = NOW()
		WHERE
			product_id = $1
		AND 
			deleted_at IS NULL
	`
	result, err := repo.DB.Exec(query, productID)
	if err != nil {
		tx.Rollback()
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if affectedRows == 0 {
		tx.Rollback()
		return errors.New(utils.InvalidProductID)
	}
	query = `
		UPDATE
			tbl_variant
		SET
			deleted_at = NOW()
		WHERE
			product_id = $1
		AND 
			deleted_at IS NULL
	`
	_, err = repo.DB.Exec(query, productID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
	}
	return nil
}

// GetProduct : Postgres function to get a product
func (repo *Repo) GetProduct(productID int) ([]ProductVariantRow, error) {
	var productVariantList []ProductVariantRow
	var description, imageURL, variantName, size, color sql.NullString
	var maxRetailPrice, discountPrice sql.NullFloat64
	var variantID sql.NullInt32
	query := `
		SELECT
			p.product_id, p.name AS product_name, p.description, p.image_url, p.category_id,
			v.variant_id, v.name AS variant_name, v.max_retail_price, v.discount_price,
			v.size, v.color
		FROM
			tbl_product p
			LEFT JOIN
				tbl_variant v
			ON p.product_id = v.product_id
		WHERE
			p.product_id = $1
			AND p.deleted_at IS NULL
			AND v.deleted_at IS NULL
	`
	rows, err := repo.DB.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var prodVar ProductVariantRow
	for rows.Next() {
		err := rows.Scan(&prodVar.ProductID, &prodVar.ProductName, &description, &imageURL, &prodVar.CategoryID,
			&variantID, &variantName, &maxRetailPrice, &discountPrice, &size, &color)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			prodVar.Description = description.String
		}
		if imageURL.Valid {
			prodVar.ImageURL = imageURL.String
		}
		if variantID.Valid {
			prodVar.VariantID = int(variantID.Int32)
			if variantName.Valid {
				prodVar.VariantName = variantName.String
			}
			if size.Valid {
				prodVar.VariantSize = size.String
			}
			if color.Valid {
				prodVar.VariantColor = color.String
			}
			if maxRetailPrice.Valid {
				prodVar.MRP = maxRetailPrice.Float64
			}
			if discountPrice.Valid {
				prodVar.DiscountPrice = discountPrice.Float64
			}
		}
		productVariantList = append(productVariantList, prodVar)
	}
	return productVariantList, nil
}

