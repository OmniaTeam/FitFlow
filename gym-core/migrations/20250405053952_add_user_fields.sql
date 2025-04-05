-- Создаем перечисления для новых полей на русском языке
CREATE TYPE purpose_enum AS ENUM ('gain', 'slim', 'stable');
CREATE TYPE place_enum AS ENUM ('home', 'gym');
CREATE TYPE level_enum AS ENUM ('new', 'medium', 'pro');

-- Добавляем новые поля в таблицу users
ALTER TABLE users
    ADD COLUMN purpose purpose_enum,
    ADD COLUMN placement place_enum,
    ADD COLUMN level level_enum,
    ADD COLUMN training_count INTEGER,
    ADD COLUMN food_prompt TEXT;