-- Stories de platform (League Basket) pueden apuntar a un país específico.
-- GLOBAL = visible en todos los países.
ALTER TABLE stories
ADD COLUMN IF NOT EXISTS target_country_code VARCHAR(10) NOT NULL DEFAULT 'GLOBAL';

CREATE INDEX IF NOT EXISTS idx_stories_target_country
ON stories (target_country_code);
