CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE stock_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stock_id UUID NOT NULL,
    movement_type TEXT NOT NULL,
    quantity INT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL,
    CONSTRAINT fk_stock_movement_stock FOREIGN KEY (stock_id) REFERENCES stocks(id)
);

CREATE INDEX idx_stock_movements_stock_id ON stock_movements(stock_id);