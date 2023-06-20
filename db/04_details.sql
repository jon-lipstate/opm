DROP FUNCTION IF EXISTS public.get_all_dependencies CASCADE;
DROP FUNCTION IF EXISTS public.get_all_dependency_licenses CASCADE;
DROP FUNCTION IF EXISTS public.get_package_details CASCADE;

--------
CREATE OR REPLACE FUNCTION get_all_dependencies(_version_id INTEGER)
RETURNS TABLE(depends_on_id INTEGER)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE dependency_tree AS (
        SELECT 
            pd.depends_on_id AS depends_on_id
        FROM 
            package_dependencies AS pd
        WHERE 
            pd.version_id = _version_id
        UNION 
        SELECT 
            pd.depends_on_id
        FROM 
            package_dependencies AS pd
        INNER JOIN 
            dependency_tree AS dt ON pd.version_id = dt.depends_on_id
    )
    SELECT DISTINCT
        dt.depends_on_id
    FROM 
        dependency_tree AS dt;
END;
$$;

-----------------------------------------------------
CREATE OR REPLACE FUNCTION get_all_dependency_licenses(version_id INTEGER)
RETURNS TABLE(
    license TEXT,
    packages TEXT[]
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    WITH all_dependencies AS (
        SELECT 
            get_all_dependencies(version_id) as version_id
    )
    SELECT 
        v.license,
        array_agg(distinct CONCAT(p.host_name, '/', p.owner_name, '/', p.repo_name))
    FROM 
        all_dependencies AS ad
    INNER JOIN 
        versions AS v ON ad.version_id = v.id
    INNER JOIN
        packages AS p ON v.package_id = p.id
    GROUP BY 
        v.license;
END;
$$;



-----------------------------------------------------------
DROP FUNCTION IF EXISTS public.get_dependencies_flat CASCADE;
CREATE OR REPLACE FUNCTION get_dependencies_flat(_version_id INTEGER)
RETURNS TABLE(
    host_name TEXT,
    owner_name TEXT,
    repo_name TEXT,
    version TEXT,
    license TEXT,
    last_updated TIMESTAMP WITH TIME ZONE,
    state package_state,
    insecure BOOLEAN
)
LANGUAGE plpgsql
AS $$
DECLARE
    _package_id INTEGER;
BEGIN
    SELECT v.package_id
    INTO _package_id
    FROM versions v
    WHERE v.id = _version_id;


    RETURN QUERY
    WITH all_dependencies AS (
        SELECT * FROM get_all_dependencies(_version_id)
    )
    SELECT 
        p.host_name,
        p.owner_name,
        p.repo_name,
        v.version,
        v.license,
        v.created_at AS last_updated,
        p.state,
        v.insecure
    FROM 
        packages p
    INNER JOIN
        versions v ON p.id = v.package_id
    WHERE 
        v.id != _version_id AND v.id = ANY(SELECT depends_on_id FROM all_dependencies);
END;
$$;


-------------------------------------------------------------------------------------------
-- NOTE: reamde will get heavy when many versions exist
-- TODO: refactor to be lighter?
DROP FUNCTION IF EXISTS public.get_version_details CASCADE;
CREATE OR REPLACE FUNCTION get_version_details(_package_id INTEGER)
RETURNS TABLE(
    id INTEGER,
    version TEXT,
    insecure BOOLEAN,
    createdAt TIMESTAMP WITH TIME ZONE,
    size_kb INTEGER,
    dependency_count INTEGER,
    compiler TEXT,
    license TEXT,
    readme TEXT,
    commit_hash TEXT,
    has_insecure_dependency BOOLEAN
)
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        v.id,
        v.version,
        v.insecure,
        v.created_at,
        v.size_kb,
        (SELECT COUNT(*)::INTEGER 
         FROM package_dependencies pd
         WHERE v.id = pd.version_id) AS dependency_count,
        v.compiler,
        v.license,
        v.readme,
        v.commit_hash,
        EXISTS(
            SELECT 1
            FROM package_dependencies pd
            INNER JOIN versions dep_v ON pd.depends_on_id = dep_v.id
            WHERE pd.version_id = v.id AND dep_v.insecure = true
        ) AS has_insecure_dependency
    FROM 
        versions AS v
    WHERE 
        v.package_id = _package_id
    ORDER BY
        v.created_at DESC;
END;
$$
LANGUAGE plpgsql;

-----------------------------------------------------------------
DROP FUNCTION IF EXISTS public.get_package_details CASCADE;
CREATE OR REPLACE FUNCTION get_package_details(_package_id INTEGER)
RETURNS TABLE(
    id INTEGER,
    host_name TEXT,
    owner_name TEXT,
    repo_name TEXT,
    description TEXT,
    state package_state,
    keywords TEXT[],
    bookmarks INTEGER,
    url TEXT,
    authors TEXT[]
)
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.host_name,
        p.owner_name,
        p.repo_name,
        p.description,
        p.state,
        (SELECT array_agg(k.keyword) FROM package_keywords pk INNER JOIN keywords k ON pk.keyword_id = k.id WHERE pk.package_id = _package_id) AS keywords,
        (SELECT COUNT(*)::INTEGER FROM bookmarks bm WHERE bm.package_id = _package_id) AS bookmarks,
        p.url,
        ARRAY(
            SELECT a.gh_login
            FROM package_authors pa
            INNER JOIN users a ON pa.author_id = a.id
            WHERE pa.package_id = _package_id
        ) AS authors
    FROM 
        packages p
    WHERE 
        p.id = _package_id;
END;
$$
LANGUAGE plpgsql;


-----------------------------------------------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION get_dependent_packages(_package_id INTEGER)
RETURNS TABLE(
    id INTEGER,
    name TEXT
)
AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT 
        p.id,
        p.name
    FROM 
        packages p
    INNER JOIN
        versions v ON p.id = v.package_id
    INNER JOIN
        package_dependencies pd ON v.id = pd.version_id
    WHERE 
        pd.depends_on_id IN (SELECT id FROM versions WHERE package_id = _package_id);
END;
$$
LANGUAGE plpgsql;
