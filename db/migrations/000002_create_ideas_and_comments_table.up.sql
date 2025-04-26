-- Creating the schema 'kanaka'
CREATE SCHEMA IF NOT EXISTS kanaka;

-- Creating the ideas table under the schema 'kanaka'
CREATE TABLE kanaka.ideas (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES public.users(id),  -- Referring to the users table in the 'public' schema
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    content TEXT NOT NULL,
    total_likes INTEGER DEFAULT 0,
    total_comments INTEGER DEFAULT 0
);

-- Creating the comments table under the schema 'kanaka'
CREATE TABLE kanaka.comments (
    id SERIAL PRIMARY KEY,
    idea_id INTEGER NOT NULL REFERENCES kanaka.ideas(id) ON DELETE CASCADE,  -- Referring to the ideas table in the 'kanaka' schema
    user_id INTEGER NOT NULL REFERENCES public.users(id),  -- Referring to the users table in the 'public' schema
    comment_text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
