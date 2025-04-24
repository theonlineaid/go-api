-- db/init.sql

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role VARCHAR(20) NOT NULL, -- e.g., user, admin, manager
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    image TEXT,
    address_line1 VARCHAR(100),
    city VARCHAR(50),
    country VARCHAR(50),
    postal_code VARCHAR(20),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Login sessions table
CREATE TABLE IF NOT EXISTS login_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    browser VARCHAR(100),
    os VARCHAR(100),
    device VARCHAR(100),
    ip_address VARCHAR(45), -- Supports IPv4/IPv6
    login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Shipping addresses table
CREATE TABLE IF NOT EXISTS shipping_addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    address_line1 VARCHAR(100) NOT NULL,
    city VARCHAR(50) NOT NULL,
    country VARCHAR(50) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('home', 'office', 'other')),
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Billing addresses table
CREATE TABLE IF NOT EXISTS billing_addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    address_line1 VARCHAR(100) NOT NULL,
    city VARCHAR(50) NOT NULL,
    country VARCHAR(50) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('home', 'office', 'other')),
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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