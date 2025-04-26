-- Dropping the ideas table under the schema 'kanaka'
DROP TABLE IF EXISTS kanaka.ideas;

-- Dropping the comments table under the schema 'kanaka'
DROP TABLE IF EXISTS kanaka.comments;

-- Dropping the schema 'kanaka'
-- Be cautious: If there are any objects (tables, views, etc.) within the schema, it will be dropped too.
DROP SCHEMA IF EXISTS kanaka CASCADE;

