-- favorites table to track product likes
CREATE TABLE favorites (
    id BIGSERIAL PRIMARY KEY,  -- Use BIGSERIAL instead of BIGINT UNSIGNED AUTO_INCREMENT
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,  -- Use TIMESTAMPTZ for datetime
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT unique_user_product UNIQUE (user_id, product_id)  -- Define the constraint separately
);
