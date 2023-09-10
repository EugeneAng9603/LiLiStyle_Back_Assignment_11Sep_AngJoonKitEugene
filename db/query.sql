-- Check if the user has liked the product
-- $1 and $2 are placeholders for query parameters
-- Returns a boolean (true if liked, false otherwise)
-- db: products
SELECT EXISTS (
    SELECT 1
    FROM favorites
    WHERE user_id = $1 AND product_id = $2
);

-- Unlike a product (remove a like)
-- $1 and $2 are placeholders for query parameters
-- Deletes the record from the favorites table
-- Returns the number of rows affected (usually 1 if the like was successfully removed)
-- db: products
DELETE FROM favorites
WHERE user_id = $1 AND product_id = $2;

-- Retrieve liked products with pagination
-- $1, $2, and $3 are placeholders for query parameters
-- Returns a list of liked products and the total count
-- Implement pagination and join with the products table to get product details
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
-- $1 is a placeholder for the query parameter
-- Returns the total count of liked products
-- db: products
SELECT COUNT(*)
FROM favorites
WHERE user_id = $1;

-- Add a product like (insert a record into favorites)
-- $1 and $2 are placeholders for query parameters
-- Inserts a record into the favorites table
-- db: products
INSERT INTO favorites (user_id, product_id, created_at)
VALUES ($1, $2, NOW())
RETURNING id;
