CREATE TABLE IF NOT EXISTS shelf_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    recommended_by UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'queued', -- 'queued', 'finished', 'dropped'
    rating SMALLINT CHECK (rating IS NULL OR (rating >= 1 AND rating <= 10)),
    take VARCHAR(280) DEFAULT '', -- Short, intentional one-line reflection
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_item_entry UNIQUE(user_id, item_id)
);

CREATE INDEX IF NOT EXISTS idx_shelf_user_status ON shelf_entries(user_id, status);
CREATE INDEX IF NOT EXISTS idx_shelf_finished_at ON shelf_entries(user_id, finished_at DESC);
