

-- name: CreateCustomer :one
INSERT INTO customer (name, phoneno, email)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCustomerByPhoneNo :many
SELECT phoneno FROM customer 
WHERE phoneno = $1;

-- name: GetCustomerById :many
SELECT * FROM customer 
WHERE id = $1
LIMIT 1;

-- name: CreateProduct :one
INSERT INTO product (name, price, stock)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetProductById :many
SELECT id FROM product 
WHERE id = $1;


-- name: CreateOrder :one
INSERT INTO "order" (customer_id, order_status, order_date, total_price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderById :many
SELECT id FROM "order" 
WHERE id = $1;


-- name: CreateOrderItem :one
INSERT INTO "order_item" (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderItemById :many
SELECT id FROM "order_item" 
WHERE id = $1;

-- name: UpdateOrderById :many 
UPDATE "order" 
SET order_status = $2
WHERE id = $1
RETURNING *;


-- -- name: CreateMessage :one
-- INSERT INTO message (thread, sender, content)
-- VALUES ($1, $2, $3)
-- RETURNING *;

-- -- name: GetMessageByID :one
-- SELECT * FROM message
-- WHERE id = $1;

-- -- name: GetMessagesByThread :many
-- SELECT * FROM message
-- WHERE thread = $1
-- ORDER BY created_at DESC;


-- -- creating a thread
-- -- name: CreateThread :one
-- INSERT INTO thread (thread)
-- VALUES ($1)
-- RETURNING *;