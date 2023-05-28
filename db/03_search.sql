DROP FUNCTION IF EXISTS public.search_packages CASCADE;
DROP FUNCTION IF EXISTS public.search_packages_no_keyword CASCADE;
DROP FUNCTION IF EXISTS public.get_package_details CASCADE;
DROP FUNCTION IF EXISTS public.search_and_get_details CASCADE;
---

CREATE OR REPLACE FUNCTION search_and_get_details(_search TEXT, _limit INTEGER DEFAULT 50, _offset INTEGER DEFAULT 0)
RETURNS TABLE (
    package_id INTEGER,
    name TEXT,
    description TEXT,
    version TEXT,
    last_updated TIMESTAMP,
    downloads INTEGER,
    all_downloads BIGINT,
    stars BIGINT,
    keywords TEXT[]
)
AS $$
DECLARE
    _package_ids INTEGER[];
BEGIN
    SELECT array_agg(id) INTO _package_ids
    FROM search_packages(_search, _limit, _offset);
    
    RETURN QUERY 
    SELECT 
        p.id AS package_id,
        p.name,
        p.description,
        v.version,
        v.created_at AS last_updated,
        v.downloads AS downloads,
        sum(v_all.downloads) AS all_downloads,
        count(s.*) AS stars,
        array_agg(distinct k.keyword) AS keywords
    FROM 
        packages AS p
    JOIN 
        package_keywords AS pk ON p.id = pk.package_id
    JOIN 
        keywords AS k ON pk.keyword_id = k.id
    LEFT JOIN 
        versions AS v ON p.id = v.package_id
    LEFT JOIN 
        versions AS v_all ON p.id = v_all.package_id
    LEFT JOIN 
        stars AS s ON p.id = s.package_id
    WHERE 
        p.id = ANY(_package_ids)
    GROUP BY
        p.id,
        v.version,
        v.created_at,
        v.downloads;
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
-- select * from public.get_package_details(ARRAY[1,2])
CREATE OR REPLACE FUNCTION get_package_details(_package_ids INT[])
RETURNS TABLE (
    package_id INTEGER,
    name TEXT,
    description TEXT,
    version TEXT,
    last_updated TIMESTAMP,
    downloads INTEGER,
    all_downloads BIGINT,
    stars BIGINT,
    keywords TEXT[]
)
AS $$
BEGIN
    RETURN QUERY 
    SELECT 
        p.id AS package_id,
        p.name,
        p.description,
        v.version,
        v.created_at AS last_updated,
        v.downloads AS downloads,
        sum(v_all.downloads) AS all_downloads,
        count(s.*) AS stars,
        array_agg(distinct k.keyword) AS keywords
    FROM 
        packages AS p
    JOIN 
        package_keywords AS pk ON p.id = pk.package_id
    JOIN 
        keywords AS k ON pk.keyword_id = k.id
    LEFT JOIN 
        versions AS v ON p.id = v.package_id
    LEFT JOIN 
        versions AS v_all ON p.id = v_all.package_id
    LEFT JOIN 
        stars AS s ON p.id = s.package_id
    WHERE 
        p.id = ANY(_package_ids)
    GROUP BY
        p.id,
        v.version,
        v.created_at,
        v.downloads; -- TODO: if older version is updated, shows it, not the new one, dont think can happen?
END; $$
LANGUAGE plpgsql;


-----


-- TODO: on pagination do i pass the id's back to client and they ask for details?
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
