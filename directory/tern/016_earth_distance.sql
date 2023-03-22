CREATE EXTENSION if not exists cube;
CREATE EXTENSION if not exists earthdistance;

ALTER TABLE provider_metadata ALTER COLUMN location TYPE point USING location::point;

---- create above / drop below ----
-- note: this doesn't work perfectly.  Assuming original string format was "X,Y" this will end up being "(X,Y)" when we convert 
-- back to text

ALTER TABLE provider_metadata ALTER COLUMN location TYPE text;
drop EXTENSION earthdistance;
drop EXTENSION cube;
