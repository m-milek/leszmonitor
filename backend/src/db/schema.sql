CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    username      VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255)       NOT NULL,

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
    slug        VARCHAR(50) UNIQUE NOT NULL,
    name        VARCHAR(100)       NOT NULL,
    description VARCHAR(1000)      NOT NULL,

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

CREATE TABLE IF NOT EXISTS permissions (
    id          TEXT PRIMARY KEY,
    slug        VARCHAR(50)   NOT NULL,
    name        VARCHAR(100)  NOT NULL,
    description VARCHAR(1000) NOT NULL,

    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER IF NOT EXISTS update_permissions_updated_at
    AFTER UPDATE
    ON permissions
    FOR EACH ROW
BEGIN
    UPDATE permissions SET updated_at = CURRENT_TIMESTAMP WHERE id = new.id;
END;

CREATE TABLE IF NOT EXISTS monitors (
    id          TEXT PRIMARY KEY,
    slug        VARCHAR(50)   NOT NULL,
    project_id  TEXT          NOT NULL,
    name        VARCHAR(100)  NOT NULL,
    description VARCHAR(1000) NOT NULL,
    interval    INT           NOT NULL,
    kind        VARCHAR(50)   NOT NULL,
    config      TEXT          NOT NULL, -- JSON string

    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP,

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

CREATE TABLE IF NOT EXISTS monitor_results (
    id                    TEXT PRIMARY KEY,
    monitor_id            TEXT    NOT NULL,
    is_success            BOOLEAN NOT NULL,
    is_manually_triggered BOOLEAN NOT NULL,
    duration_ms           INT     NOT NULL,

    error_message         TEXT,
    error_details         TEXT, -- JSON string

    details               TEXT    NOT NULL,

    created_at            DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (monitor_id) REFERENCES monitors (id) ON DELETE CASCADE
);
