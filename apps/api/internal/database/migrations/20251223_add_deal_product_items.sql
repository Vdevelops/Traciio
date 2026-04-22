-- Migration: Add Deal Product Items
-- Date: 2025-12-23
-- Description: Adds deal_product_items table to attach product line items to deals/opportunities

CREATE TABLE IF NOT EXISTS deal_product_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deal_id UUID NOT NULL REFERENCES deals(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,

    -- Snapshot fields for historical accuracy
    product_name VARCHAR(200) NOT NULL,
    product_sku VARCHAR(100) NOT NULL,

    -- Monetary values in smallest currency unit (sen)
    unit_price BIGINT NOT NULL DEFAULT 0,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity >= 1),
    discount_amount BIGINT NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
    subtotal BIGINT NOT NULL DEFAULT 0,

    notes TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_deal_product_items_deal_id
    ON deal_product_items(deal_id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_deal_product_items_product_id
    ON deal_product_items(product_id) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uq_deal_product_items_deal_product
    ON deal_product_items(deal_id, product_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE deal_product_items IS 'Product line items attached to deals/opportunities';
COMMENT ON COLUMN deal_product_items.unit_price IS 'Unit price in smallest currency unit (sen) at time of adding';
COMMENT ON COLUMN deal_product_items.discount_amount IS 'Discount amount in smallest currency unit (sen) for the line item';
COMMENT ON COLUMN deal_product_items.subtotal IS 'Computed subtotal: unit_price * quantity - discount_amount';
