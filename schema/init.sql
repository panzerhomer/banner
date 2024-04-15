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

CREATE OR REPLACE FUNCTION maintain_banner_version_rows() RETURNS TRIGGER AS $$
BEGIN
    IF (SELECT COUNT(*) FROM banner_version WHERE banner_id = NEW.banner_id) > 3 THEN
        DELETE FROM banner_version
        WHERE id IN (
            SELECT id FROM banner_version
            WHERE banner_id = NEW.banner_id
            ORDER BY created_at ASC
            LIMIT 1
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER banner_version_row_limit_trigger
AFTER INSERT ON banner_version
FOR EACH ROW
EXECUTE FUNCTION maintain_banner_version_rows();

