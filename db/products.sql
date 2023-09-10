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