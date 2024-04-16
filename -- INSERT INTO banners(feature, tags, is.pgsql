-- INSERT INTO banners(feature, tags, is_active) VALUES (1, ARRAY[0, 1], true);
-- INSERT INTO banners(feature, tags, is_active) VALUES (1, ARRAY[0], true);
-- INSERT INTO banners(feature, tags, is_active) VALUES (1, ARRAY[1, 2], true);


-- select * from banners where feature = 2 and (SELECT UNNEST (ARRAY[0, 1]) INTERSECT (SELECT UNNEST (tags))) is not null;

-- INSERT INTO banner_version(banner_id, banner_info) VALUES (1, '{"333neeew!!":"b", "ab":"c"}');

-- select * from banner_version where banner_id = 1 limit 1 offset 2;

-- select * from banners where feature = 2 and tags = ARRAY[0,1];
-- select * from banners;

-- SELECT 
-- 		b.banner_id, 
-- 		b.feature, 
-- 		b.tags, 
-- 		b.is_active, 
-- 		bv.banner_info, 
-- 		bv.created_at, 
-- 		bv.updated_at 
-- 	FROM 
-- 		banners as b
-- 	JOIN banner_version as bv
-- 	ON 
-- 		b.banner_id = bv.banner_id 
-- 	WHERE 
-- 		b.feature = 2 AND tags = ARRAY[0,1] AND b.is_active is not false
--     ORDER BY
--         b.banner_id, bv.created_at;

-- select * from banners;
-- select * FROM banner_version order by banner_id;
-- delete from banners where banner_id = 1

-- select max(updated_at) FROM 
-- banner_version 
-- where banner_id = 1
-- GROUP by updated_at
-- order by updated_at desc
-- LIMIT 1;

-- SELECT id
--     FROM banner_version
--     WHERE banner_id = 1
--     GROUP BY id
--     ORDER BY max(updated_at) DESC
--     LIMIT 1;

-- UPDATE 
-- 	banners
-- SET 
-- 	feature = 1,
-- 	tags = $2,
-- 	is_active = $3
-- WHERE 
-- 	banner_id = $4