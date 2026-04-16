CREATE TABLE IF NOT EXISTS subscriptions (
    id            SERIAL PRIMARY KEY,
    service_name  VARCHAR(255) NOT NULL,
    price         INTEGER NOT NULL,
    user_id       VARCHAR(255) NOT NULL,
    start_date    DATE NOT NULL,           -- храним как дату (01 число месяца)
    end_date      DATE NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX idx_subscriptions_start_date ON subscriptions(start_date);