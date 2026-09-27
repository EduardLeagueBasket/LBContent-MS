DROP INDEX IF EXISTS idx_articles_target_country;

ALTER TABLE articles
DROP COLUMN IF EXISTS target_country_code;
