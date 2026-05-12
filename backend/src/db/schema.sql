CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL CHECK (LENGTH(username) >= 2) CHECK (LENGTH(username) <= 50),
    password_hash TEXT        NOT NULL,

    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER IF NOT EXISTS update_users_updated_at
    AFTER UPDATE
    ON users
    FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = new.id;
END;

CREATE TABLE IF NOT EXISTS projects (
    id          TEXT PRIMARY KEY,
    slug        TEXT UNIQUE NOT NULL CHECK (LENGTH(slug) >= 2) CHECK (LENGTH(slug) <= 50),
    name        TEXT        NOT NULL CHECK (LENGTH(name) >= 2) CHECK (LENGTH(name) <= 100),
    description TEXT        NOT NULL CHECK (LENGTH(description) <= 1000),

    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER IF NOT EXISTS update_projects_updated_at
    AFTER UPDATE
    ON projects
    FOR EACH ROW
BEGIN
    UPDATE projects SET updated_at = CURRENT_TIMESTAMP WHERE id = new.id;
END;

CREATE TABLE IF NOT EXISTS user_projects (
    user_id    TEXT NOT NULL,
    project_id TEXT NOT NULL,
    role       TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'viewer')),

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, project_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);
CREATE TRIGGER IF NOT EXISTS update_user_projects_updated_at
    AFTER UPDATE
    ON user_projects
    FOR EACH ROW
BEGIN
    UPDATE user_projects SET updated_at = CURRENT_TIMESTAMP WHERE user_id = new.user_id AND project_id = new.project_id;
END;

CREATE INDEX IF NOT EXISTS idx_user_projects_project_id
    ON user_projects (project_id);

CREATE TABLE IF NOT EXISTS monitors (
    id                       TEXT PRIMARY KEY,
    slug                     TEXT NOT NULL CHECK (LENGTH(slug) >= 2) CHECK (LENGTH(slug) <= 50),
    project_id               TEXT NOT NULL,                            -- UUID
    name                     TEXT NOT NULL CHECK (LENGTH(name) >= 2) CHECK (LENGTH(name) <= 100),
    description              TEXT NOT NULL CHECK (LENGTH(description) <= 1000),
    interval                 INT  NOT NULL CHECK (interval > 0),       -- in seconds
    kind                     TEXT NOT NULL,
    result_retention_seconds INT  NOT NULL CHECK (result_retention_seconds > 0),
    config                   TEXT NOT NULL CHECK (JSON_VALID(config)), -- JSON string

    created_at               DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at               DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (project_id, slug),
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);
CREATE TRIGGER IF NOT EXISTS update_monitors_updated_at
    AFTER UPDATE
    ON monitors
    FOR EACH ROW
BEGIN
    UPDATE monitors SET updated_at = CURRENT_TIMESTAMP WHERE id = new.id;
END;

CREATE INDEX IF NOT EXISTS idx_monitors_project_id ON monitors (project_id);

CREATE TABLE IF NOT EXISTS monitor_results (
    id                    TEXT PRIMARY KEY,
    monitor_id            TEXT    NOT NULL,
    is_success            BOOLEAN NOT NULL,
    is_manually_triggered BOOLEAN NOT NULL,
    duration_ms           INT     NOT NULL CHECK (duration_ms >= 0),

    error_details         TEXT CHECK (error_details IS NULL OR JSON_VALID(error_details)), -- JSON string

    details               TEXT    NOT NULL,

    created_at            DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (monitor_id) REFERENCES monitors (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_monitor_results_monitor_id_created ON monitor_results (monitor_id, created_at DESC);