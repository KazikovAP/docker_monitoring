-- Создание таблицы для хранения информации о пингах
CREATE TABLE IF NOT EXISTS pings (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(255) UNIQUE NOT NULL,
    ping_time BIGINT,
    last_success_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для ускорения поиска по IP-адресу
CREATE UNIQUE INDEX IF NOT EXISTS idx_ip_address ON pings (ip_address);

-- Триггер для автоматического обновления поля updated_at при изменении записи
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS set_updated_at ON pings;

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON pings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
