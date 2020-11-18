package category

import (
	"database/sql"
	"ecommerce/utils"
	"errors"
	"fmt"
	"log"
	"strings"
)

//Repo is the DB repo struct
type Repo struct {
	DB *sql.DB
}

//CheckCategoryNameExists function to check if the category with the given name already exists
func (repo *Repo) CheckCategoryNameExists(name string) (bool, error) {
	var categoryExists bool
	var count int
	query := `
		SELECT
			count(*)
		FROM
			tbl_category
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
		categoryExists = true
	}
	return categoryExists, nil
}

func getNullInt32(value int) sql.NullInt32 {
	if value == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: int32(value),
		Valid: true,
	}
}

//CreateCategory - DB function to create category
func (repo *Repo) CreateCategory(request *CreateRequest) (*CreateResponse, error) {
	var createResponse CreateResponse
	var parentID sql.NullInt32
	query := `
		INSERT INTO 
			tbl_category (name, parent_category_id, created_at, updated_at)
		VALUES
			($1, $2, NOW(), NOW())
		RETURNING
			category_id, name, parent_category_id
	`
	row := repo.DB.QueryRow(query, request.Name, getNullInt32(request.ParentID))
	err := row.Scan(&createResponse.ID, &createResponse.Name, &parentID)
	if err != nil {
		return nil, err
	}
	if parentID.Valid {
		createResponse.ParentID = int(parentID.Int32)
	}
	return &createResponse, nil
}

//IsCategoryIDExists function to check if the category ID exists or not
func (repo *Repo) IsCategoryIDExists(id int) (bool, error) {
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
	err := repo.DB.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

//UpdateCategory to update a category
func (repo *Repo) UpdateCategory(request *UpdateRequest) error {
	var slice []string
	if len(request.Name) > 0 && request.Name != utils.EmptyString {
		slice = append(slice, fmt.Sprintf(" name = '%s' ", request.Name))
	}
	if request.ParentID != 0 {
		slice = append(slice, fmt.Sprintf(" parent_category_id = %d ", request.ParentID))
	}
	slice = append(slice, fmt.Sprintf(" updated_at = NOW() "))
	updateQuery := strings.Join(slice, ", ")
	mainQuery := `
		UPDATE
			tbl_category
		SET
			%s
		WHERE
			category_id = $1
		AND 
			deleted_at IS NULL
	`
	query := fmt.Sprintf(mainQuery, updateQuery)
	result, err := repo.DB.Exec(query, request.CategoryID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New(utils.InvalidCategoryID)
	}
	return err
}

//IsSubCategoryExist to check if there are any subcategories for the given category
func (repo *Repo) IsSubCategoryExist(categoryID int) (bool, error) {
	var count int
	query := `
		SELECT 
			count(*)
		FROM
			tbl_category
		WHERE
			parent_category_id = $1
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

// IsProductExist to check if there are any products under the given category
func (repo *Repo) IsProductExist(categoryID int) (bool, error) {
	var count int
	query := `
		SELECT 
			count(*)
		FROM
			tbl_product
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

//DeleteCategory to check if the given category exists
func (repo *Repo) DeleteCategory(categoryID int) error {
	query := `
		UPDATE 
			tbl_category
		SET
			deleted_at = NOW()
		WHERE
			category_id = $1
		AND
			deleted_at IS NULL
	`
	result, err := repo.DB.Exec(query, categoryID)
	if err != nil {
		return err
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return errors.New(utils.InvalidCategoryID)
	}
	return nil
}

// GetProductVariantForEachCategory DB function to get the product and variants for each category
func (repo *Repo) GetProductVariantForEachCategory(categoryIDs []int) ([]Product, error) {
	var params []string
	for i := range categoryIDs {
		params = append(params, fmt.Sprintf("$%d", i+1))
	}
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
		ON 
			p.product_id = v.product_id
		WHERE
			p.product_id IN (%s)
		AND
			p.deleted_at IS NULL
		AND
			v.deleted_at IS NULL
		ORDER BY 
			product_id ASC,
			variant_id ASC
	`
	mainQuery := fmt.Sprintf(query, strings.Join(params, ", "))
	categoryIDInterface := make([]interface{}, len(categoryIDs))
	for i, v := range categoryIDs {
		categoryIDInterface[i] = v
	}
	rows, err := repo.DB.Query(mainQuery, categoryIDInterface...)
	if err != nil {
		return nil, err
	}
	var productVariantList []ProductVariantRow
	for rows.Next() {
		var prodVar ProductVariantRow
		err = rows.Scan(&prodVar.ProductID, &prodVar.ProductName, &description, &imageURL,
			&prodVar.CategoryID, &variantID, &variantName, &maxRetailPrice, &discountPrice,
			&size, &color)
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
			if maxRetailPrice.Valid {
				prodVar.MRP = float64(maxRetailPrice.Float64)
			}
			if discountPrice.Valid {
				prodVar.DiscountPrice = float64(discountPrice.Float64)
			}
			if size.Valid {
				prodVar.Size = size.String
			}
			if color.Valid {
				prodVar.Color = color.String
			}
		}
		log.Println("Heree", prodVar)
		productVariantList = append(productVariantList, prodVar)
	}
	var productList []Product
	if len(productVariantList) > 0 {
		for _, row := range productVariantList {
			var product Product
			var variant Variant
			var variants []Variant
			var isAppended bool
			for k, v := range productList {
				if v.ProductID == row.ProductID {
					variant.VariantID = row.VariantID
					variant.Name = row.VariantName
					variant.MRP = row.MRP
					variant.DiscountPrice = row.DiscountPrice
					variant.Size = row.Size
					variant.Color = row.Color
					v.Variants = append(v.Variants, variant)
					productList[k] = v
					isAppended = true
					break
				}
			}
			if isAppended {
				continue
			}
			if row.VariantID != 0 {
				variant.VariantID = row.VariantID
				variant.Name = row.VariantName
				variant.MRP = row.MRP
				variant.DiscountPrice = row.DiscountPrice
				variant.Size = row.Size
				variant.Color = row.Color
				variants = append(variants, variant)
			}
			product.ProductID = row.ProductID
			product.Name = row.ProductName
			product.Description = row.Description
			product.ImageURL = row.ImageURL
			product.CategoryID = row.CategoryID
			product.Variants = variants
			log.Println("Heree======", product)
			productList = append(productList, product)
		}
	}
	return productList, nil
}

// GetCategories to get the categories from DB
func (repo *Repo) GetCategories() (*[]Category, error) {
	var categories []Category
	var parentID sql.NullInt32
	query := `
		SELECT 
			category_id, name, parent_category_id
		FROM 
			tbl_category
		WHERE
			deleted_at IS NULL
		ORDER BY
			category_id ASC
	`
	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &parentID)
		if err != nil {
			return nil, err
		}
		if parentID.Valid {
			category.ParentID = int(parentID.Int32)
		}
		categories = append(categories, category)
	}
	return &categories, nil
}
