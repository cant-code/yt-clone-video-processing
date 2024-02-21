CREATE TABLE IF NOT EXISTS JOB_STATUS (
    vid BIGINT NOT NULL,
    quality VARCHAR(4) NOT NULL,
    success BOOLEAN NOT NULL,
    error TEXT,
    created_at TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (vid, quality)
);

