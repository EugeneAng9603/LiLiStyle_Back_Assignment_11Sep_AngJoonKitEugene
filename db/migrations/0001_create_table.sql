-- db/migrations/0001_create_tables.sql

-- products table
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    shop_id BIGINT,
    name VARCHAR(255) COLLATE "C", 
    description TEXT,
    thumbnail_url TEXT NOT NULL,
    origin_price BIGINT NOT NULL,
    discounted_price BIGINT NOT NULL,
    discounted_rate DOUBLE PRECISION,
    status VARCHAR(191) NOT NULL,
    in_stock BOOLEAN DEFAULT FALSE NOT NULL,
    is_preorder BOOLEAN DEFAULT FALSE NOT NULL,
    is_purchasable BOOLEAN DEFAULT FALSE NOT NULL,
    delivery_condition VARCHAR(255) NOT NULL,
    delivery_display TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,  
    name VARCHAR(50),
    email VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    status VARCHAR(191) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,  
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL, 
    deleted_at TIMESTAMPTZ DEFAULT NOW() NOT NULL 
);

-- favorites table to track product likes
CREATE TABLE favorites (
    id BIGSERIAL PRIMARY KEY,  
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL, 
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT unique_user_product UNIQUE (user_id, product_id)  
);
