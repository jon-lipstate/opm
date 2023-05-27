-- this file is just for scratch queries

-- get all packages and collect keywords
SELECT 
    p.id,
    p.name,
    p.description,
    string_agg(k.keyword, ', ') AS keywords
FROM 
    packages AS p
LEFT JOIN
    package_keywords pk ON p.id = pk.package_id
LEFT JOIN
    keywords k ON pk.keyword_id = k.id
GROUP BY
    p.id;
