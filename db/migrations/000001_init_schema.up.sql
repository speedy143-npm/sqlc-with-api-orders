

CREATE TABLE "thread" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "thread" VARCHAR(36) NOT NULL,
  "created_at" TIMESTAMP DEFAULT now()
);

CREATE TABLE "message" (
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "thread" VARCHAR(36) NOT NULL,
  "sender" VARCHAR(100) NOT NULL,
  "content" TEXT NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
 FOREIGN KEY ("thread") REFERENCES "thread" ("id") ON DELETE CASCADE
);

