CREATE TABLE IF NOT EXISTS banners(
    banner_id SERIAL PRIMARY KEY, 
    feature INT NOT NULL, 
    tags INT[] NOT NULL, 
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT unique_feature_tags UNIQUE (feature, tags)
);

CREATE TABLE IF NOT EXISTS banner_version(
    id SERIAL PRIMARY KEY, 
    banner_id INT REFERENCES banners(banner_id) ON DELETE CASCADE, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    banner_info JSONB
);

CREATE OR REPLACE FUNCTION check_banner_version_count()
RETURNS TRIGGER AS $$
DECLARE
    banner_count INT;
BEGIN
    SELECT count(*)
    INTO banner_count
    FROM banner_version
    WHERE banner_id = NEW.banner_id;

    IF banner_count >= 3 THEN
        RAISE EXCEPTION 'More than 3 rows with the same banner_id are not allowed in banner_version table';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER banner_version_count_trigger
BEFORE INSERT OR UPDATE ON banner_version
FOR EACH ROW
EXECUTE FUNCTION check_banner_version_count();

-- SELECT banners.banner_id, banner_version
  
  
-- INSERT INTO banners(feature, tags, is_active) VALUES (1, ARRAY[0], true);
  
-- INSERT INTO banner_version(banner_id, created_at, banner_info) VALUES (2, '2024-04-04', '{"a":"b"}');
-- INSERT INTO banner_version(banner_id, created_at, banner_info) VALUES (2, '2024-04-05', '{"a":"b1"}');
-- INSERT INTO banner_version(banner_id, created_at, banner_info) VALUES (2, '2024-04-06', '{"a":"b2"}');
  
-- SELECT * from banners;
-- select * from banner_version;
  
--   --Добавление баннера если его еще не существует
  
-- CREATE FUNCTION create_banner() RETURNS int AS '
--     INSERT INTO banners(feature, tag, is_active) VALUES (2, ARRAY[1], true)  RETURNING banner_id;
-- ' LANGUAGE SQL;  

-- INSERT INTO banner_version(banner_id, created_at, banner_info) VALUES (create_banner(), '2024-04-06', '{"a":"b2"}')

   
--  -- Добавление версии баннера если он уже есть
--  -- Для этого необохдимо проверить кол-во версий баннера и если их 3, то удалить посл
--  -- ПС этот дэлит работает нормально) удаляет только если их 3, если меньше то не удаляет
 
-- DELETE FROM banner_version 
-- where id = (SELECT id 
--             FROM banner_version 
--             WHERE banner_id = 5 AND created_at = (
--             SELECT MIN(created_at) FROM banner_version WHERE banner_id = 5 GROUP BY banner_id HAVING COUNT(*) = 3))
--   -- ну и собственно вставляем...

-- SELECT feature, tags FROM banner WHERE 

