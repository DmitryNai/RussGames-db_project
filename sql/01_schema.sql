CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL CHECK (status IN ('active','banned','inactive')),
    country TEXT,
    profile JSONB,
    CONSTRAINT username_length CHECK (char_length(username) BETWEEN 3 AND 64)
);

CREATE TABLE developers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    country TEXT,
    website TEXT,
    contact_email TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    metadata JSONB
);

CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id UUID NOT NULL REFERENCES developers(id) ON DELETE CASCADE ON UPDATE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    genre TEXT,
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    release_date DATE,
    avg_rating NUMERIC(3,2) DEFAULT 0,
    sales_count BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    metadata JSONB
);

CREATE TABLE game_licenses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE ON UPDATE CASCADE,
    key TEXT NOT NULL UNIQUE,
    assigned_to_user UUID REFERENCES users(id) ON DELETE SET NULL,
    assigned_at TIMESTAMPTZ,
    state TEXT NOT NULL CHECK (state IN ('available','assigned','revoked','consumed')) DEFAULT 'available',
    notes TEXT
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    amount NUMERIC(12,2) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'RUB',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    provider TEXT,
    status TEXT NOT NULL CHECK (status IN ('pending','completed','failed','refunded')),
    details JSONB
);

CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE RESTRICT,
    transaction_id UUID REFERENCES transactions(id) ON DELETE SET NULL,
    price_paid NUMERIC(10,2) NOT NULL,
    purchased_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    method TEXT,
    license_id UUID REFERENCES game_licenses(id) ON DELETE SET NULL,
    CONSTRAINT unique_user_game_purchase UNIQUE(user_id, game_id, transaction_id)
);


CREATE TABLE library (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    license_id UUID REFERENCES game_licenses(id) ON DELETE SET NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT unique_library_entry UNIQUE(user_id, game_id)
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 10),
    title TEXT,
    body TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    helpful_count INTEGER DEFAULT 0,
    CONSTRAINT unique_user_review UNIQUE(user_id, game_id)
);

CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    reward_points INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT unique_game_ach_code UNIQUE(game_id, code)
);

CREATE TABLE user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    achieved_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    progress JSONB,
    CONSTRAINT unique_user_achievement UNIQUE(user_id, achievement_id)
);

CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    table_name TEXT NOT NULL,
    operation CHAR(1) NOT NULL CHECK (operation IN ('I','U','D')),
    row_id UUID,
    performed_by UUID,
    performed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    old_data JSONB,
    new_data JSONB,
    query TEXT
);

CREATE TABLE batch_errors (
    id BIGSERIAL PRIMARY KEY,
    batch_name TEXT NOT NULL,
    row_data JSONB,
    error TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_games_dev ON games(developer_id);
CREATE INDEX idx_purchases_user ON purchases(user_id);
CREATE INDEX idx_reviews_game ON reviews(game_id);
CREATE INDEX idx_library_user ON library(user_id);
CREATE INDEX idx_transactions_user ON transactions(user_id);