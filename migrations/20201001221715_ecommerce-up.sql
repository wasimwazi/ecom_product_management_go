-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS tbl_category (
    category_id SERIAL,
    name VARCHAR(50) NOT NULL,
    parent_category_id INT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    PRIMARY KEY (category_id),
    FOREIGN KEY (parent_category_id) REFERENCES tbl_category(category_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS tbl_product (
    product_id SERIAL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(200),
    image_url VARCHAR(160),
    category_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    PRIMARY KEY (product_id),
    FOREIGN KEY (category_id) REFERENCES tbl_category(category_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS tbl_variant(
    variant_id SERIAL,
    name VARCHAR(50),
    max_retail_price FLOAT NOT NULL,
    discount_price FLOAT,
    size VARCHAR(10),
    color VARCHAR(15),
    product_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    PRIMARY KEY (variant_id),
    FOREIGN KEY (product_id) REFERENCES tbl_product(product_id) ON DELETE CASCADE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS tbl_variant;
DROP TABLE IF EXISTS tbl_product;
DROP TABLE IF EXISTS tbl_category;
