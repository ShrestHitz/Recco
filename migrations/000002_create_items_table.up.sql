CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(32) NOT NULL, -- 'movie', 'youtube', 'book', 'anime', 'podcast'
    title VARCHAR(255) NOT NULL,
    creator VARCHAR(255) NOT NULL DEFAULT '',
    cover_url TEXT DEFAULT '',
    external_url TEXT NOT NULL,
    source_id VARCHAR(255) DEFAULT '', -- e.g. 'yt:dQw4w9WgXcQ', 'tmdb:550'
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_items_type ON items(type);
CREATE INDEX IF NOT EXISTS idx_items_source_id ON items(source_id);
CREATE INDEX IF NOT EXISTS idx_items_external_url ON items(external_url);
CREATE INDEX IF NOT EXISTS idx_items_metadata ON items USING GIN(metadata);
