-- db/init.sql

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user'
);

-- Brands table
CREATE TABLE IF NOT EXISTS brands (
    id SERIAL PRIMARY KEY,
    brand_name VARCHAR(255) NOT NULL,
    image TEXT,
    status INT DEFAULT 1,
    is_feature BOOLEAN DEFAULT FALSE,
    is_publish BOOLEAN DEFAULT FALSE,
    is_special BOOLEAN DEFAULT FALSE,
    is_approved_by_admin BOOLEAN DEFAULT FALSE,
    is_visible_to_guest BOOLEAN DEFAULT TRUE,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50),
    category_name VARCHAR(255) NOT NULL,
    category_img TEXT,
    image TEXT,
    category_visibility INT DEFAULT 1,
    is_special INT DEFAULT 0,
    is_featured INT DEFAULT 0,
    is_approved BOOLEAN DEFAULT FALSE,
    is_published BOOLEAN DEFAULT FALSE,
    position INT,
    price_visibility INT DEFAULT 0,
    status INT DEFAULT 1,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Subcategories table
CREATE TABLE IF NOT EXISTS subcategories (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    subcategory_name VARCHAR(255) NOT NULL,
    image TEXT,
    status INT DEFAULT 1
);

-- Attributes table
CREATE TABLE IF NOT EXISTS attributes (
    id SERIAL PRIMARY KEY,
    attribute_name VARCHAR(100) NOT NULL,
    status INT DEFAULT 1
);

-- Attribute Values table
CREATE TABLE IF NOT EXISTS attribute_values (
    id SERIAL PRIMARY KEY,
    attribute_id INT NOT NULL REFERENCES attributes(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL,
    status INT DEFAULT 1
);

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    brand_id INT REFERENCES brands(id),
    category_id INT REFERENCES categories(id),
    sub_category_id INT REFERENCES subcategories(id),
    p_code VARCHAR(255),
    weight VARCHAR(255),
    product_name TEXT NOT NULL,
    product_code VARCHAR(255) NOT NULL,
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
    is_multi INT DEFAULT 0,
    serial_number INT,
    tax0 NUMERIC,
    tax1 NUMERIC,
    hsn_code VARCHAR(255),
    is_saleable INT DEFAULT 1,
    is_barcode INT DEFAULT 0,
    is_expirable INT DEFAULT 0,
    is_warranty INT DEFAULT 0,
    is_serviceable INT DEFAULT 0,
    is_variation INT DEFAULT 0,
    status INT DEFAULT 1,
    UOM_id INT,
    rating INT DEFAULT 0,
    current_stock INT DEFAULT 0
);

-- Product Attributes junction table
CREATE TABLE IF NOT EXISTS product_attributes (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    attribute_id INT NOT NULL REFERENCES attributes(id),
    attribute_value_id INT NOT NULL REFERENCES attribute_values(id),
    UNIQUE (product_id, attribute_id, attribute_value_id)
);

-- Variation Products table
CREATE TABLE IF NOT EXISTS variation_products (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(255) NOT NULL,
    sale_price NUMERIC,
    default_sell_price NUMERIC,
    discount NUMERIC,
    image TEXT,
    current_stock INT DEFAULT 0,
    UNIQUE (product_id, sku)
);

-- Sample data
INSERT INTO users (username, password, role) 
VALUES 
    ('admin', '$2a$10$YOUR_HASH_HERE', 'admin'), -- Replace with bcrypt hash for 'admin123'
    ('user1', '$2a$10$YOUR_HASH_HERE', 'user'); -- Replace with bcrypt hash for 'user123'

INSERT INTO brands (brand_name, image) 
VALUES 
    ('Nike', 'brands/nike_logo.jpg'),
    ('Adidas', 'brands/adidas_logo.jpg');

INSERT INTO categories (id, category_name, is_special, price_visibility, image) 
VALUES 
    (33, 'GLASS', 0, 0, 'categories/glass_icon.jpg');

INSERT INTO subcategories (category_id, subcategory_name, image) 
VALUES 
    (33, 'Sunglasses', 'subcategories/sunglasses.jpg');

INSERT INTO attributes (attribute_name) 
VALUES 
    ('Color'), 
    ('Size'), 
    ('Fabric');

INSERT INTO attribute_values (attribute_id, value) 
VALUES 
    ((SELECT id FROM attributes WHERE attribute_name = 'Color'), 'Red'),
    ((SELECT id FROM attributes WHERE attribute_name = 'Color'), 'Blue'),
    ((SELECT id FROM attributes WHERE attribute_name = 'Size'), 'Medium'),
    ((SELECT id FROM attributes WHERE attribute_name = 'Size'), 'Large'),
    ((SELECT id FROM attributes WHERE attribute_name = 'Fabric'), 'Cotton'),
    ((SELECT id FROM attributes WHERE attribute_name = 'Fabric'), 'Polyester');

-- Indexes for foreign keys in products table
CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_sub_category_id ON products(sub_category_id);