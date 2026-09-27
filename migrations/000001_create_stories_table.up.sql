CREATE TABLE IF NOT EXISTS publishers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(30) NOT NULL, -- 'league', 'team', 'user'
    entity_id UUID NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    avatar_url TEXT NULL,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_publisher_entity UNIQUE (entity_type, entity_id)
);

CREATE TABLE IF NOT EXISTS stories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    publisher_id UUID NOT NULL REFERENCES publishers(id) ON DELETE CASCADE,
    media_url TEXT NOT NULL,
    media_type VARCHAR(10) NOT NULL DEFAULT 'photo',
    cloudinary_public_id VARCHAR(255) NULL,
    duration_seconds INT DEFAULT 5,
    expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '24 hours'),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID NOT NULL,
    name_created_by VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS story_views (
    story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
    user_id UUID NOT NULL, 
    viewed_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (story_id, user_id)
);

CREATE INDEX idx_stories_expires ON stories (expires_at);
CREATE INDEX idx_story_views_user ON story_views (user_id);