CREATE EXTENSION if not exists cube;
CREATE EXTENSION if not exists earthdistance;

-- 1. Drop dependent views
DROP VIEW IF EXISTS providers_v CASCADE;
DROP VIEW IF EXISTS providers_base_v CASCADE;

-- 2. Change column type
ALTER TABLE provider_metadata ALTER COLUMN location TYPE point USING location::point;

---- create above / drop below ----

-- note: this doesn't work perfectly.  Assuming original string format was "X,Y" this will end up being "(X,Y)" when we convert 
-- back to text

-- Change back to text (if needed)
ALTER TABLE provider_metadata ALTER COLUMN location TYPE text;

-- Drop extensions (optional, if not used elsewhere)
DROP EXTENSION IF EXISTS earthdistance;
DROP EXTENSION IF EXISTS cube;
