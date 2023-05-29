DROP FUNCTION IF EXISTS public.search_packages CASCADE;
DROP FUNCTION IF EXISTS public.search_packages_no_keyword CASCADE;
DROP FUNCTION IF EXISTS public.get_package_results CASCADE;
DROP FUNCTION IF EXISTS public.search_and_get_results CASCADE;
DROP FUNCTION IF EXISTS public.browse_packages CASCADE;
DROP TYPE IF EXISTS public.search_result CASCADE;
---
---
CREATE TYPE search_result AS (
    package_id INTEGER,
    name TEXT,
    description TEXT,
    version TEXT,
    last_updated TIMESTAMP,
    downloads BIGINT,
    stars BIGINT,
    keywords TEXT[]
);
---
CREATE OR REPLACE FUNCTION browse_packages(_limit INTEGER DEFAULT 100, _offset INTEGER DEFAULT 0)
RETURNS SETOF search_result
AS $$
DECLARE
    _package_ids INTEGER[];
BEGIN
    SELECT array_agg(id) INTO _package_ids
    FROM (
        SELECT id
        FROM packages
        ORDER BY name
        LIMIT _limit OFFSET _offset
    ) AS p;
    
    RETURN QUERY SELECT * FROM get_package_results(_package_ids);
END;
$$ 
LANGUAGE plpgsql;
---
-- todo: split the `-` on names so they become two keywords
CREATE OR REPLACE FUNCTION search_and_get_results(_search TEXT, _limit INTEGER DEFAULT 50, _offset INTEGER DEFAULT 0)
RETURNS SETOF search_result
AS $$
DECLARE
    _package_ids INTEGER[];
BEGIN
    SELECT array_agg(id) INTO _package_ids
    FROM search_packages(_search, _limit, _offset);
    
    RETURN QUERY 
    SELECT * FROM get_package_results(_package_ids);
END;
$$ 
LANGUAGE plpgsql;




---
CREATE OR REPLACE FUNCTION search_packages(_search TEXT, _limit INTEGER DEFAULT 50, _offset INTEGER DEFAULT 0)
RETURNS TABLE (
    id INTEGER,
    rank REAL
) 
AS $$
BEGIN
    RETURN QUERY 
    WITH keyword_search AS (
        SELECT 
            p.id,
            to_tsvector(string_agg(k.keyword, ' ')) as keyword_tsv
        FROM 
            packages AS p
        LEFT JOIN
            package_keywords pk ON p.id = pk.package_id
        LEFT JOIN
            keywords k ON pk.keyword_id = k.id
        GROUP BY
            p.id
    )
    SELECT 
        p.id,
        ts_rank(p.TSV, to_tsquery(replace(trim(_search), ' ', ' | '))) + 
        ts_rank(ks.keyword_tsv, to_tsquery(replace(trim(_search), ' ', ' | '))) AS rank
    FROM 
        packages AS p
    JOIN
        keyword_search ks ON p.id = ks.id
    WHERE 
        p.TSV @@ to_tsquery(replace(trim(_search), ' ', ' | ')) OR
        ks.keyword_tsv @@ to_tsquery(replace(trim(_search), ' ', ' | '))
    ORDER BY 
        rank DESC
    LIMIT 
        _limit
    OFFSET 
        _offset;
END; $$
LANGUAGE plpgsql;


-- SELECT * FROM search_packages('http server async');
-- _offset => 20 to specify out of order named items
CREATE OR REPLACE FUNCTION search_packages_no_keyword(_search TEXT, _limit INTEGER DEFAULT 50, _offset INTEGER DEFAULT 0)
RETURNS TABLE (
    id INTEGER,
    rank REAL
) 
AS $$
BEGIN
    RETURN QUERY 
    SELECT 
        p.id,
        ts_rank(p.TSV, to_tsquery(replace(trim(_search), ' ', ' | '))) AS rank
    FROM 
        packages AS p
    WHERE 
        p.TSV @@ to_tsquery(replace(trim(_search), ' ', ' | '))
    ORDER BY 
        rank DESC
    LIMIT 
        _limit
    OFFSET 
        _offset;
END; $$
LANGUAGE plpgsql;
--------------------------------------------------------------------------------------------------------------------------
-- Hydrate Search
--------------------------------------------------------------------------------------------------------------------------
-- select * from public.get_package_results(ARRAY[1,2])
CREATE OR REPLACE FUNCTION get_package_results(_package_ids INT[])
RETURNS SETOF search_result
AS $$
BEGIN
    RETURN QUERY 
    SELECT 
        p.id AS package_id,
        p.name,
        p.description,
        v.version,
        v.created_at AS last_updated,
        (SELECT SUM(downloads) FROM versions WHERE package_id = p.id) AS all_downloads,
        count(distinct s.user_id) AS stars,
        array_agg(distinct k.keyword) AS keywords
    FROM 
        packages AS p
    JOIN 
        package_keywords AS pk ON p.id = pk.package_id
    JOIN 
        keywords AS k ON pk.keyword_id = k.id
    LEFT JOIN -- constrain to newest version id (serial)
        versions AS v ON p.id = v.package_id AND v.id = (SELECT MAX(id) FROM versions WHERE package_id = p.id)
    LEFT JOIN 
        stars AS s ON p.id = s.package_id
    WHERE 
        p.id = ANY(_package_ids)
    GROUP BY
        p.id,
        v.version,
        v.created_at;
END;
$$ 
LANGUAGE plpgsql;



-----


-- TODO: on pagination do i pass the id's back to client and they ask for results?
CREATE OR REPLACE FUNCTION search_package_count(_search TEXT)
RETURNS INTEGER 
AS $$
DECLARE
    _total INTEGER;
BEGIN
    SELECT COUNT(id) INTO _total
    FROM search_packages(_search, NULL, NULL);
    
    RETURN _total;
END;
$$ 
LANGUAGE plpgsql;
