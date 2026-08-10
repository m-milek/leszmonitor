PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE users (
    id            TEXT PRIMARY KEY,
    username      TEXT UNIQUE NOT NULL CHECK (LENGTH(username) >= 2) CHECK (LENGTH(username) <= 50),
    password_hash TEXT        NOT NULL,

    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO users VALUES('c074f2be-b8e8-4576-b46d-26529d8ecd34','leszmak','$2a$10$/fTzp6jHlClqweKnHP.UYOFT8lw0QAx19W9vRS3fkk4B4lBXlNgmm','2026-06-17 19:54:28','2026-06-17 19:54:28');
CREATE TABLE projects (
    id          TEXT PRIMARY KEY,
    slug        TEXT UNIQUE NOT NULL CHECK (LENGTH(slug) >= 2) CHECK (LENGTH(slug) <= 50),
    name        TEXT        NOT NULL CHECK (LENGTH(name) >= 2) CHECK (LENGTH(name) <= 100),
    description TEXT        NOT NULL CHECK (LENGTH(description) <= 1000),

    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO projects VALUES('36f5b816-1adf-4e51-a615-e2e68424aee0','leszmaks-sandbox','leszmak''s Sandbox','Your default sandbox project','2026-06-17 19:54:28','2026-06-17 19:54:28');
CREATE TABLE user_projects (
    user_id    TEXT NOT NULL,
    project_id TEXT NOT NULL,
    role       TEXT NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'viewer')),

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, project_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);
INSERT INTO user_projects VALUES('c074f2be-b8e8-4576-b46d-26529d8ecd34','36f5b816-1adf-4e51-a615-e2e68424aee0','owner','2026-06-17 19:54:28','2026-06-17 19:54:28');
CREATE TABLE monitors (
    id                       TEXT PRIMARY KEY,
    slug                     TEXT NOT NULL CHECK (LENGTH(slug) >= 2) CHECK (LENGTH(slug) <= 50),
    project_id               TEXT NOT NULL,                            -- UUID
    name                     TEXT NOT NULL CHECK (LENGTH(name) >= 2) CHECK (LENGTH(name) <= 100),
    description              TEXT NOT NULL CHECK (LENGTH(description) <= 1000),
    interval                 INT  NOT NULL CHECK (interval > 0),       -- in seconds
    kind                     TEXT NOT NULL,
    result_retention_seconds INT  NOT NULL CHECK (result_retention_seconds > 0),
    run_state                TEXT NOT NULL,
    config                   TEXT NOT NULL CHECK (JSON_VALID(config)), -- JSON string

    created_at               DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at               DATETIME DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (project_id, slug),
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);
INSERT INTO monitors VALUES('1ab9e321-42fe-48cb-9031-7acd38e7218f','http-monitor','36f5b816-1adf-4e51-a615-e2e68424aee0','HTTP Monitor','',20,'http',43200,'active','{"method":"GET","url":"https://example.com","headers":{},"body":"","saveResponseBody":false,"saveResponseHeaders":false,"expectedStatusCodes":[200],"expectedHeaders":{}}','2026-06-17 19:54:55','2026-06-17 19:54:55');
INSERT INTO monitors VALUES('bf526d03-9c24-4614-b5e8-f4b27dc20d3e','tcp-monitor','36f5b816-1adf-4e51-a615-e2e68424aee0','TCP Monitor','',20,'tcp',43200,'active','{"host":"example.com","port":443,"protocol":"tcp","timeout":5000,"retryCount":3}','2026-06-17 19:55:10','2026-06-17 19:55:10');
INSERT INTO monitors VALUES('b42f08c7-15e7-4783-a903-525347a19d6b','dns-monitor','36f5b816-1adf-4e51-a615-e2e68424aee0','DNS Monitor','',20,'dns',43200,'active','{"hostname":"example.com","recordType":"A","dnsServer":"1.1.1.1","expectedRecordValues":[]}','2026-06-17 19:56:33','2026-06-17 19:56:33');
CREATE TABLE monitor_results (
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
INSERT INTO monitor_results VALUES('9ae2721b-702b-4b61-b7ea-2ed0563b451b','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,185,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:55:15+02:00');
INSERT INTO monitor_results VALUES('3bb8cba4-d3df-49a7-9bab-7b8d582ab259','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,74,NULL,X'7b227472696573223a312c226c6174656e63794d73223a37347d','2026-06-17T21:55:30+02:00');
INSERT INTO monitor_results VALUES('3b8678cf-5dce-473f-9147-4bebd625735d','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,37,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:55:35+02:00');
INSERT INTO monitor_results VALUES('86d886b6-4896-4e97-95bc-0bbef1b184ed','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,32,NULL,X'7b227472696573223a312c226c6174656e63794d73223a33327d','2026-06-17T21:55:50+02:00');
INSERT INTO monitor_results VALUES('cbc75e9c-2dac-4f24-ad0d-818cc16ffb9f','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,64,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:55:55+02:00');
INSERT INTO monitor_results VALUES('d05730e2-820d-4780-84be-3193f9e85a96','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,28,NULL,X'7b227472696573223a312c226c6174656e63794d73223a32387d','2026-06-17T21:56:10+02:00');
INSERT INTO monitor_results VALUES('ad78c299-b377-4554-b005-c11b45c9624e','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,99,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:56:15+02:00');
INSERT INTO monitor_results VALUES('b8a59b0b-3532-4dd2-b29f-6cb0745e4d01','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,44,NULL,X'7b227472696573223a312c226c6174656e63794d73223a34347d','2026-06-17T21:56:30+02:00');
INSERT INTO monitor_results VALUES('1eb7a8bf-e9be-4d14-b8ba-60062e56c747','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,60,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:56:35+02:00');
INSERT INTO monitor_results VALUES('da1925b4-e4f0-41c5-b5d6-8271c12cf20c','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,28,NULL,X'7b227472696573223a312c226c6174656e63794d73223a32387d','2026-06-17T21:56:50+02:00');
INSERT INTO monitor_results VALUES('868f69a4-0e1c-45d7-8d70-8d64c7e1501b','b42f08c7-15e7-4783-a903-525347a19d6b','up',0,41,NULL,X'7b227265736f6c7665645265636f726473223a5b223137322e36362e3134372e323433222c223130342e32302e32332e313534225d7d','2026-06-17T21:56:53+02:00');
INSERT INTO monitor_results VALUES('cc6c134b-7778-4aab-a802-1f27b4b7ef27','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,65,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:56:55+02:00');
INSERT INTO monitor_results VALUES('2c6d2734-c6f6-4121-be3d-8934f9c461da','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,24,NULL,X'7b227472696573223a312c226c6174656e63794d73223a32347d','2026-06-17T21:57:10+02:00');
INSERT INTO monitor_results VALUES('d0b7fec7-13dd-4aab-a696-8646875cae83','b42f08c7-15e7-4783-a903-525347a19d6b','up',0,65,NULL,X'7b227265736f6c7665645265636f726473223a5b223130342e32302e32332e313534222c223137322e36362e3134372e323433225d7d','2026-06-17T21:57:13+02:00');
INSERT INTO monitor_results VALUES('21b553a5-7d8e-47dd-a6be-0f0bfea18b40','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,86,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:57:15+02:00');
INSERT INTO monitor_results VALUES('2dd21ce8-8c58-41e6-a238-1954296f7e08','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,37,NULL,X'7b227472696573223a312c226c6174656e63794d73223a33377d','2026-06-17T21:57:30+02:00');
INSERT INTO monitor_results VALUES('ed254bde-798b-42bd-9027-e2950311a26c','b42f08c7-15e7-4783-a903-525347a19d6b','up',0,84,NULL,X'7b227265736f6c7665645265636f726473223a5b223130342e32302e32332e313534222c223137322e36362e3134372e323433225d7d','2026-06-17T21:57:33+02:00');
INSERT INTO monitor_results VALUES('1f2d7a84-eeb7-43bd-b9fc-0c96c6b6ca89','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,46,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:57:35+02:00');
INSERT INTO monitor_results VALUES('e0cc7304-d7b6-444e-99ec-67d57041f996','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,43,NULL,X'7b227472696573223a312c226c6174656e63794d73223a34337d','2026-06-17T21:57:50+02:00');
INSERT INTO monitor_results VALUES('4640d639-791f-442c-80f3-5452deff922b','b42f08c7-15e7-4783-a903-525347a19d6b','up',0,37,NULL,X'7b227265736f6c7665645265636f726473223a5b223130342e32302e32332e313534222c223137322e36362e3134372e323433225d7d','2026-06-17T21:57:53+02:00');
INSERT INTO monitor_results VALUES('6f3153b0-b31a-41ec-9459-b72d40dc1b5f','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,28,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:57:55+02:00');
INSERT INTO monitor_results VALUES('05823b9c-7297-433c-869d-3cce77fb7e54','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,94,NULL,X'7b227472696573223a312c226c6174656e63794d73223a39347d','2026-06-17T21:58:10+02:00');
INSERT INTO monitor_results VALUES('e5ca8590-e94a-4024-b96a-5e6889855a49','b42f08c7-15e7-4783-a903-525347a19d6b','up',0,25,NULL,X'7b227265736f6c7665645265636f726473223a5b223137322e36362e3134372e323433222c223130342e32302e32332e313534225d7d','2026-06-17T21:58:13+02:00');
INSERT INTO monitor_results VALUES('a4ba3a9d-40a1-4f89-a4ad-7ef96fdf902f','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,64,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:58:15+02:00');
INSERT INTO monitor_results VALUES('9dac8efa-50f5-49f7-bed1-d77729f4388c','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,60,NULL,X'7b227472696573223a312c226c6174656e63794d73223a36307d','2026-06-17T21:58:30+02:00');
INSERT INTO monitor_results VALUES('4073bc4f-b9b0-44d2-a968-e7db0174b8a4','b42f08c7-15e7-4783-a903-525347a19d6b','up',0,30,NULL,X'7b227265736f6c7665645265636f726473223a5b223130342e32302e32332e313534222c223137322e36362e3134372e323433225d7d','2026-06-17T21:58:33+02:00');
INSERT INTO monitor_results VALUES('84393434-9ab7-42dd-b0dc-88a438006397','1ab9e321-42fe-48cb-9031-7acd38e7218f','up',0,84,NULL,X'7b22737461747573436f6465223a3230302c22636f6e74656e744c656e677468223a2d312c2270726f746f223a22485454502f322e30227d','2026-06-17T21:58:35+02:00');
INSERT INTO monitor_results VALUES('0b4ab3f3-81dc-4498-8552-5f0f23844d0e','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','up',0,67,NULL,X'7b227472696573223a312c226c6174656e63794d73223a36377d','2026-06-17T21:58:50+02:00');
INSERT INTO monitor_results VALUES('ccf635b2-fb1b-4cf9-8b5d-9e792c7b9b40','b42f08c7-15e7-4783-a903-525347a19d6b','up',0,33,NULL,X'7b227265736f6c7665645265636f726473223a5b223137322e36362e3134372e323433222c223130342e32302e32332e313534225d7d','2026-06-17T21:58:53+02:00');
CREATE TABLE audit_logs (
    id          TEXT PRIMARY KEY,
    username    TEXT,
    project_id  TEXT,
    resource_id TEXT,
    action      TEXT    NOT NULL,
    is_success  BOOLEAN NOT NULL,
    summary     TEXT CHECK (summary IS NULL OR LENGTH(summary) <= 1000),
    before      TEXT CHECK (before IS NULL OR JSON_VALID(before)), -- JSON string
    after       TEXT CHECK (after IS NULL OR JSON_VALID(after)),   -- JSON string
    trace_id    TEXT,

    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO audit_logs VALUES('10cd482b-e365-4f75-adb3-397ebc25ebbc','leszmak','36f5b816-1adf-4e51-a615-e2e68424aee0','1ab9e321-42fe-48cb-9031-7acd38e7218f','monitor.create',1,'Monitor with ID 1ab9e321-42fe-48cb-9031-7acd38e7218f created',NULL,'{"id":"1ab9e321-42fe-48cb-9031-7acd38e7218f","slug":"http-monitor","projectId":"36f5b816-1adf-4e51-a615-e2e68424aee0","name":"HTTP Monitor","description":"","interval":20,"type":"http","probeConfig":"{\"method\":\"GET\",\"url\":\"https://example.com\",\"headers\":{},\"body\":\"\",\"saveResponseBody\":false,\"saveResponseHeaders\":false,\"expectedStatusCodes\":[200],\"expectedHeaders\":{}}","resultRetentionSeconds":43200,"state":"active","createdAt":"2026-06-17T19:54:55Z","updatedAt":"2026-06-17T19:54:55Z"}',NULL,'2026-06-17 19:54:55.530801 +0000 UTC');
INSERT INTO audit_logs VALUES('c4875554-d2aa-439c-ad6d-b482aba90390','leszmak','36f5b816-1adf-4e51-a615-e2e68424aee0','bf526d03-9c24-4614-b5e8-f4b27dc20d3e','monitor.create',1,'Monitor with ID bf526d03-9c24-4614-b5e8-f4b27dc20d3e created',NULL,'{"id":"bf526d03-9c24-4614-b5e8-f4b27dc20d3e","slug":"tcp-monitor","projectId":"36f5b816-1adf-4e51-a615-e2e68424aee0","name":"TCP Monitor","description":"","interval":20,"type":"tcp","probeConfig":"{\"host\":\"example.com\",\"port\":443,\"protocol\":\"tcp\",\"timeout\":5000,\"retryCount\":3}","resultRetentionSeconds":43200,"state":"active","createdAt":"2026-06-17T19:55:10Z","updatedAt":"2026-06-17T19:55:10Z"}',NULL,'2026-06-17 19:55:10.328263 +0000 UTC');
INSERT INTO audit_logs VALUES('2f17077b-428a-4bde-b43e-540f5301002f','leszmak','36f5b816-1adf-4e51-a615-e2e68424aee0','b42f08c7-15e7-4783-a903-525347a19d6b','monitor.create',1,'Monitor with ID b42f08c7-15e7-4783-a903-525347a19d6b created',NULL,'{"id":"b42f08c7-15e7-4783-a903-525347a19d6b","slug":"dns-monitor","projectId":"36f5b816-1adf-4e51-a615-e2e68424aee0","name":"DNS Monitor","description":"","interval":20,"type":"dns","probeConfig":"{\"hostname\":\"example.com\",\"recordType\":\"A\",\"dnsServer\":\"1.1.1.1\",\"expectedRecordValues\":[]}","resultRetentionSeconds":43200,"state":"active","createdAt":"2026-06-17T19:56:33Z","updatedAt":"2026-06-17T19:56:33Z"}',NULL,'2026-06-17 19:56:33.433434 +0000 UTC');
CREATE TRIGGER update_users_updated_at
    AFTER UPDATE
    ON users
    FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = new.id;
END;
CREATE TRIGGER update_projects_updated_at
    AFTER UPDATE
    ON projects
    FOR EACH ROW
BEGIN
    UPDATE projects SET updated_at = CURRENT_TIMESTAMP WHERE id = new.id;
END;
CREATE TRIGGER update_user_projects_updated_at
    AFTER UPDATE
    ON user_projects
    FOR EACH ROW
BEGIN
    UPDATE user_projects SET updated_at = CURRENT_TIMESTAMP WHERE user_id = new.user_id AND project_id = new.project_id;
END;
CREATE TRIGGER update_monitors_updated_at
    AFTER UPDATE
    ON monitors
    FOR EACH ROW
BEGIN
    UPDATE monitors SET updated_at = CURRENT_TIMESTAMP WHERE id = new.id;
END;
CREATE INDEX idx_user_projects_project_id
    ON user_projects (project_id);
CREATE INDEX idx_monitors_project_id ON monitors (project_id);
CREATE INDEX idx_monitor_results_monitor_id_created ON monitor_results (monitor_id, created_at DESC);
CREATE INDEX idx_audit_logs_project_id_created ON audit_logs (project_id, created_at DESC);
COMMIT;
