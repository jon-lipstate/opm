DROP FUNCTION IF EXISTS public.search_packages CASCADE;
DROP FUNCTION IF EXISTS public.search_packages_no_keyword CASCADE;
CREATE OR REPLACE FUNCTION search_packages(_search TEXT, _limit INTEGER, _offset INTEGER)
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


-- SELECT * FROM search_packages('http server async',999,0);
CREATE OR REPLACE FUNCTION search_packages_no_keyword(_search TEXT, _limit INTEGER, _offset INTEGER)
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
