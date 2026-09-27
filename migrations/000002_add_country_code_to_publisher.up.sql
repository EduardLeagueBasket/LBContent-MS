ALTER TABLE publishers 
ADD COLUMN country_code VARCHAR(10) NOT NULL DEFAULT 'GLOBAL';

CREATE INDEX idx_publishers_country ON publishers (country_code);