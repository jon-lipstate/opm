CREATE OR REPLACE FUNCTION search_packages(search_term TEXT, limit INT, offset INT)
RETURNS TABLE (id INTEGER, rank NUMERIC) 
LANGUAGE plpgsql AS $$
BEGIN
  RETURN QUERY 
  SELECT 
    p.id,
    ts_rank(p.TSV, plainto_tsquery(search_term)) AS rank
  FROM 
    packages AS p
  WHERE 
    p.TSV @@ plainto_tsquery(search_term)
  ORDER BY 
    rank DESC
  LIMIT 
    limit
  OFFSET 
    offset;
END; 
$$;
