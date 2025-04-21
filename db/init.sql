-- db/init.sql

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL
);

INSERT INTO categories (id, category_name, is_special, price_visibility)
VALUES (33, 'GLASS', 0, 0);

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    category_name VARCHAR(255),
    is_special INT DEFAULT 0,
    price_visibility INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    brand_id INT,
    category_id INT,  -- <- removed REFERENCES categories(id)
    sub_category_id INT,
    p_code VARCHAR(255),
    weight VARCHAR(255),
    product_name TEXT NOT NULL,
    product_code VARCHAR(255),
    price NUMERIC,
    m_total_price NUMERIC,
    unit VARCHAR(50),
    discount NUMERIC,
    tax NUMERIC,
    tax_type VARCHAR(50),
    serial_no VARCHAR(255),
    product_vat NUMERIC,
    product_model VARCHAR(255),
    warranty VARCHAR(255),
    minimum_qty_alert INT,
    image TEXT,
    is_multi INT,
    serial_number INT,
    tax0 NUMERIC,
    tax1 NUMERIC,
    hsn_code VARCHAR(255),
    is_saleable INT,
    is_barcode INT,
    is_expirable INT,
    is_warranty INT,
    is_serviceable INT,
    is_variation INT,
    status INT,
    UOM_id INT,
    rating INT,
    current_stock INT
);


CREATE TABLE IF NOT EXISTS variation_products (
    id SERIAL PRIMARY KEY,
    variation_id INT,
    variation_details_id INT,
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(255),
    value TEXT,
    color VARCHAR(100),
    sale_price NUMERIC,
    default_sell_price NUMERIC,
    discount NUMERIC,
    image TEXT
);
