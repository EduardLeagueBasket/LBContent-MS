DROP INDEX IF EXISTS idx_publishers_country;

ALTER TABLE publishers 
DROP COLUMN IF EXISTS country_code;