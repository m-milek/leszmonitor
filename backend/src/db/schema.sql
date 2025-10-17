CREATE OR REPLACE FUNCTION update_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();

CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON teams
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();

CREATE TABLE IF NOT EXISTS user_teams (
    user_id UUID NOT NULL,
    team_id UUID NOT NULL,
    role team_role NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, team_id),

    UNIQUE (user_id, team_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON user_teams
    FOR EACH ROW
    EXECUTE PROCEDURE update_timestamp();

CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id VARCHAR(50) UNIQUE NOT NULL,
    team_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (id, team_id),
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON groups
    FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON permissions
    FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

DO $$ BEGIN
    CREATE TYPE team_role AS ENUM ('owner', 'admin', 'member', 'viewer');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS monitors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id VARCHAR(50) NOT NULL,
    group_id UUID,
    team_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    interval INT NOT NULL,
    kind VARCHAR(50) NOT NULL,
    config JSON NOT NULL,

    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (id, group_id),

    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL
);
CREATE OR REPLACE TRIGGER set_timestamp
    BEFORE UPDATE ON monitors
    FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();