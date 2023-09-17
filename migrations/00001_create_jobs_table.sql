CREATE TABLE jobs (
    id VARCHAR(80) PRIMARY KEY,
    name TEXT NOT NULL,
    data TEXT,
    run_at TIMESTAMPTZ,
    execute_at TIMESTAMPTZ,
    status VARCHAR(80) NOT NULL,
    ttl INT NOT NULL,
    times INT NOT NULL,
    executed_times INT,
    note TEXT,
    level INT NOT NULL,
    type VARCHAR(200) NOT NULL,
    logs TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);