

-- Customers Table (Optional)
CREATE TABLE "customer" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "name" VARCHAR(255) NOT NULL,
  "phoneno" VARCHAR(255) NOT NULL UNIQUE,
  "email" VARCHAR(255) NOT NULL UNIQUE,
  "created_at" TIMESTAMP DEFAULT now()
);

-- Products Table
-- CREATE TABLE "product" (
--   "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
--   "name" VARCHAR(255) NOT NULL,
--   "price" DECIMAL(10, 2) NOT NULL,
--   "stock" INT NOT NULL DEFAULT 0,
--   "created_at" TIMESTAMP DEFAULT now()
-- );

-- Orders Table
CREATE TABLE "order" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "customer_id" VARCHAR(36) NOT NULL,
  "order_status" VARCHAR(20) DEFAULT 'PENDING',
  "order_date" TIMESTAMP DEFAULT now(),
  "total_price" DECIMAL(10, 2) NOT NULL,
  CONSTRAINT fk_customer FOREIGN KEY ("customer_id") REFERENCES "customer" ("id") ON DELETE CASCADE
);

-- Order Items Table
-- CREATE TABLE "order_item" (
--   "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
--   "order_id" VARCHAR(36) NOT NULL,
--   "product_id" VARCHAR(36) NOT NULL,
--   "quantity" INT NOT NULL CHECK (quantity > 0),
--   "price" DECIMAL(10, 2) NOT NULL,
--   "created_at" TIMESTAMP DEFAULT now(),
--   CONSTRAINT fk_order FOREIGN KEY ("order_id") REFERENCES "order" ("id") ON DELETE CASCADE,
--   CONSTRAINT fk_product FOREIGN KEY ("product_id") REFERENCES "product" ("id") ON DELETE CASCADE
-- );

