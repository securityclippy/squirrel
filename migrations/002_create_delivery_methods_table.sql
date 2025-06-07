-- +goose Up
CREATE TABLE delivery_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    config JSONB,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

INSERT INTO delivery_methods (name, config) VALUES 
('email', '{"smtp_host": "", "smtp_port": 587, "username": "", "password": ""}'),
('sms', '{"api_key": "", "service": "twilio"}'),
('slack', '{"webhook_url": ""}'),
('discord', '{"webhook_url": ""}');

-- +goose Down
DROP TABLE delivery_methods;