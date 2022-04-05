CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "email" VARCHAR NOT NULL,
  "username" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "is_admin" BOOLEAN NOT NULL
);

CREATE TABLE "products" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "description" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "image_url" VARCHAR NOT NULL,
  "image_id" VARCHAR NOT NULL,
  "price" DECIMAL NOT NULL
);

CREATE TABLE "categories" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL
);

CREATE TABLE "sizes" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR NOT NULL
);

CREATE TABLE "category_product_map" (
  "id" SERIAL PRIMARY KEY,
  "product_id" INT NOT NULL,
  "category_id" INT NOT NULL
);

CREATE TABLE "size_product_map" (
  "id" SERIAL PRIMARY KEY,
  "product_id" INT NOT NULL,
  "size_id" INT NOT NULL
);

CREATE TABLE "orders" (
  "id" SERIAL PRIMARY KEY,
  "amount" INT NOT NULL,
  "user_id" INT NOT NULL,
  "status" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "address" VARCHAR NOT NULL,
  "prodcut_id" INT NOT NULL
);

CREATE TABLE "admins" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INT NOT NULL
);

ALTER TABLE "category_product_map" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id")  ON DELETE CASCADE;
ALTER TABLE "category_product_map" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE CASCADE;

ALTER TABLE "size_product_map" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE;
ALTER TABLE "size_product_map" ADD FOREIGN KEY ("size_id") REFERENCES "sizes" ("id") ON DELETE CASCADE;

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "orders" ADD FOREIGN KEY ("prodcut_id") REFERENCES "products" ("id") ON DELETE CASCADE;

ALTER TABLE "admins" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

-- CREATE INDEX ON "prodcuts" ("price");
-- CREATE INDEX ON "users" ("is_admin");
-- CREATE INDEX ON "category_product_map" ("category_id");
-- CREATE INDEX ON "size_product_map" ("size_id");
-- CREATE INDEX ON "orders" ("user_id");
-- CREATE INDEX ON "orders" ("product_id");