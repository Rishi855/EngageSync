CREATE TABLE kanaka.projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    admin_id VARCHAR(255)
);

CREATE TABLE kanaka.project_members (
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES kanaka.projects(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL, -- ✅ Change from VARCHAR(255) to INTEGER
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role SMALLINT,
    technology SMALLINT,
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);
