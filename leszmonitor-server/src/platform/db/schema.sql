CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL CHECK (LENGTH(username) >= 2) CHECK (LENGTH(username) <= 50),
    role          TEXT        NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
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

CREATE TABLE IF NOT EXISTS monitors (
    id                       TEXT PRIMARY KEY,
    slug                     TEXT NOT NULL CHECK (LENGTH(slug) >= 2) CHECK (LENGTH(slug) <= 50),
    name                     TEXT NOT NULL CHECK (LENGTH(name) >= 2) CHECK (LENGTH(name) <= 100),
    description              TEXT NOT NULL CHECK (LENGTH(description) <= 1000),
    interval                 INT  NOT NULL CHECK (interval > 0),       -- in seconds
    owner_id                 TEXT NOT NULL,                            -- user who created the monitor
    kind                     TEXT NOT NULL,
    result_retention_seconds INT  NOT NULL CHECK (result_retention_seconds > 0),
    run_state                TEXT NOT NULL,
    config                   TEXT NOT NULL CHECK (JSON_VALID(config)), -- JSON string

    created_at               DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at               DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (slug)
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
    status                TEXT    NOT NULL,
    is_manually_triggered BOOLEAN NOT NULL,
    duration_ms           INT     NOT NULL CHECK (duration_ms >= 0),

    error_details         TEXT CHECK (error_details IS NULL OR JSON_VALID(error_details)), -- JSON string

    details               TEXT    NOT NULL,

    created_at            DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (monitor_id) REFERENCES monitors (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_monitor_results_monitor_id_created ON monitor_results (monitor_id, created_at DESC);

CREATE TABLE IF NOT EXISTS audit_logs (
    id          TEXT PRIMARY KEY,
    username    TEXT,
    resource_id TEXT,
    action      TEXT    NOT NULL,
    is_success  BOOLEAN NOT NULL,
    summary     TEXT CHECK (summary IS NULL OR LENGTH(summary) <= 1000),
    before      TEXT CHECK (before IS NULL OR JSON_VALID(before)), -- JSON string
    after       TEXT CHECK (after IS NULL OR JSON_VALID(after)),   -- JSON string
    trace_id    TEXT,

    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS monitor_status_changes (
    id              TEXT PRIMARY KEY,
    monitor_id      TEXT NOT NULL,
    caused_by_id    TEXT NOT NULL, -- UUID of the monitor result that caused the status change
    previous_status TEXT NOT NULL,
    next_status     TEXT NOT NULL,

    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (monitor_id) REFERENCES monitors (id) ON DELETE CASCADE,
    FOREIGN KEY (caused_by_id) REFERENCES monitor_results (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS tags (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    color_hex   TEXT NOT NULL,

    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS update_tags_updated_at
    AFTER UPDATE
    ON tags
    FOR EACH ROW
BEGIN
    UPDATE tags SET updated_at = CURRENT_TIMESTAMP WHERE id = new.id;
END;

CREATE TABLE IF NOT EXISTS monitor_tags (
    monitor_id TEXT NOT NULL,
    tag_id     TEXT NOT NULL,

    PRIMARY KEY (monitor_id, tag_id),
    FOREIGN KEY (monitor_id) REFERENCES monitors (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_monitor_tags_tag_id ON monitor_tags (tag_id);
