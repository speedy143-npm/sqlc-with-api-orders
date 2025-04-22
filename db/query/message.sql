

-- name: CreateCustomer :one
INSERT INTO customer (name, phoneno, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCustomerByPhoneNo :many
SELECT phoneno FROM customer 
WHERE phoneno = $1;

-- name: GetCustomerById :one
SELECT * FROM customer 
WHERE id = $1
LIMIT 1;

-- name: CreateProduct :one
INSERT INTO product (name, price, stock)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProductById :one
SELECT id FROM product 
WHERE id = $1;


-- name: CreateOrder :one
INSERT INTO "order" (customer_id, product_id, price, quantity, total_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetOrderById :one
SELECT id FROM "order" 
WHERE id = $1;


-- name: CreateOrderItem :one
INSERT INTO "order_item" (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderItemById :one
SELECT id FROM "order_item" 
WHERE id = $1;

-- name: UpdateOrderById :one 
UPDATE "order" 
SET 
order_status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateOrderTotalPriceById :one 
UPDATE "order" 
SET 
 total_price = quantity * price
WHERE id = $1
RETURNING *;


