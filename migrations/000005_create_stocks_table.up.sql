CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE stocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_stocks_product FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX idx_stocks_product_id ON stocks(product_id);