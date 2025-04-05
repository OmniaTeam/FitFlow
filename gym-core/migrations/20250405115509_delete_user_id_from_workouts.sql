-- Удаляем индекс, связанный с user_id
DROP INDEX IF EXISTS idx_workout_user_id;

-- Удаляем внешний ключ, связанный с user_id
ALTER TABLE workouts DROP CONSTRAINT IF EXISTS workout_user_id_fkey;

-- Удаляем столбец user_id из таблицы workouts
ALTER TABLE workouts DROP COLUMN user_id;