-- Phase 1: Oura Ring データ保存テーブル群
-- gen_random_uuid() は pgcrypto 拡張が必要
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS readiness_records (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date        DATE NOT NULL UNIQUE,
    score       INTEGER,
    hrv_balance INTEGER,
    raw_json    JSONB,
    fetched_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sleep_records (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date          DATE NOT NULL UNIQUE,
    score         INTEGER,
    total_minutes INTEGER,
    efficiency    INTEGER,
    wake_time     TIMESTAMPTZ,
    raw_json      JSONB
);

-- Phase 2: HRV 生データ（poller が insert、件数が多いため COPY バッチ推奨）
CREATE TABLE IF NOT EXISTS ibi_records (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recorded_at TIMESTAMPTZ NOT NULL,
    interval_ms FLOAT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_ibi_records_recorded_at ON ibi_records (recorded_at);

CREATE TABLE IF NOT EXISTS daily_summaries (
    date                DATE PRIMARY KEY,
    condition_score     INTEGER,
    focus_peak_start    TIMESTAMPTZ,
    focus_peak_end      TIMESTAMPTZ,
    recommend_bedtime   TIMESTAMPTZ,
    sleep_debt_minutes  INTEGER,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);
