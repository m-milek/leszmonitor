CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$ BEGIN
    CREATE TYPE project_role AS ENUM ('owner', 'admin', 'member', 'viewer');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();

CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON projects
    FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

CREATE TABLE IF NOT EXISTS user_projects (
    user_id    UUID         NOT NULL,
    project_id UUID         NOT NULL,
    role       project_role NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, project_id),

    UNIQUE (user_id, project_id),
    FOREIGN KEY (user_id)    REFERENCES users(id)    ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON user_projects
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON permissions
    FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

CREATE TABLE IF NOT EXISTS monitors (
    id          UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    slug  VARCHAR(50)  NOT NULL,
    project_id  UUID         NOT NULL,
    name        VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    interval    INT          NOT NULL,
    kind        VARCHAR(50)  NOT NULL,
    config      JSON         NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (project_id, slug),

    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON monitors
    FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();