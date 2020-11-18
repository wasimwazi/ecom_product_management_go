package utils

const (
	//DecodeError to represent error while decoding
	DecodeError = "Error Decoding"

	//CategoryExistsError to show category exists
	CategoryExistsError = "Category name already exists"

	//InvalidCategoryID to show invalid category error
	InvalidCategoryID = "Invalid category ID"

	//NothingToUpdateInCategory to show when nothing to update in a category update request
	NothingToUpdateInCategory = "Nothing to update in category"

	//InvalidParameterError invalid parameter error
	InvalidParameterError = "Invalid request parameter"

	//ProductExistsError to show product exists
	ProductExistsError = "Product name already exists"

	//InvalidProductID to show invalid product error
	InvalidProductID = "Invalid product ID"

	//NothingToUpdateInProduct to show when nothing to update in a product update request
	NothingToUpdateInProduct = "Nothing to update in product"

	//ProductIDNotExist to show if product id doesn't exist
	ProductIDNotExist = "Product ID doesn't exist"

	//InvalidVariantID to show invalid variant error
	InvalidVariantID = "Invalid variant ID"

	//NothingToUpdateInVariant to show when nothing to update in a variant update request
	NothingToUpdateInVariant = "Nothing to update in variant"

	//CategoryNOTExistsError to show category doesn't exist error
	CategoryNOTExistsError = "Category doesn't exist"

	//SubCategoryExists to show sub category exists for the given category
	SubCategoryExists = "Category can't be deleted since sub category exists for the given category"

	//ProductExistCategoryError to show product exists under the given category
	ProductExistCategoryError = "Category can't be deleted since products exist under this category"

	//VariantExistInProductError to show variant exists under the given product
	VariantExistInProductError = "Products can't be deleted since variants exist under this product"

	//VariantIDNotExist to show variant doesn't exist error
	VariantIDNotExist = "Given variant doens't exist"

	//NoDataFoundError to show no data found in DB
	NoDataFoundError = "No data found"
)
