-- Check if the user has liked the product
-- db: products
SELECT EXISTS (
    SELECT 1
    FROM favorites
    WHERE user_id = $1 AND product_id = $2
);

-- Unlike a product (remove a like)
-- db: products
DELETE FROM favorites
WHERE user_id = $1 AND product_id = $2;

-- Retrieve liked products with pagination
-- db: products
SELECT p.*
FROM products p
INNER JOIN (
    SELECT product_id
    FROM favorites
    WHERE user_id = $1
    ORDER BY created_at 
    LIMIT $2
    OFFSET ($3 - 1) * $2
) AS liked_products ON p.id = liked_products.product_id;

-- Retrieve the total count of liked products for a user
-- db: products
SELECT COUNT(*)
FROM favorites
WHERE user_id = $1;

-- Add a product like (insert a record into favorites)
-- db: products
INSERT INTO favorites (user_id, product_id, created_at)
VALUES ($1, $2, NOW())
RETURNING id;
