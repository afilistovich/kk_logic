-- Создаём пользовательские типы (ENUM'ы)
CREATE TYPE request_type AS ENUM ('material', 'work', 'status_update');
CREATE TYPE request_status AS ENUM ('new', 'in_progress', 'done', 'rejected', 'cancelled');
CREATE TYPE user_role AS ENUM ('foreman', 'owner');

-- Создаем таблицу пользователей
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'foreman',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Создаем таблицу заявок
CREATE TABLE requests (
    id BIGSERIAL PRIMARY KEY,
    object TEXT NOT NULL,
    type request_type NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    status request_status NOT NULL DEFAULT 'new',
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Создаём таблицу истории статусов
CREATE TABLE request_status_history (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES requests(id) ON DELETE CASCADE,
    old_status request_status NOT NULL,
    new_status request_status NOT NULL,
    changed_by BIGINT NOT NULL REFERENCES users(id),
    comment TEXT,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Индексы для ускорения поиска
CREATE INDEX idx_requests_status ON requests(status);
CREATE INDEX idx_requests_object ON requests(object);
CREATE INDEX idx_history_request_id ON request_status_history(request_id);


