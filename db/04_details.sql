DROP FUNCTION IF EXISTS public.get_all_dependencies CASCADE;
DROP FUNCTION IF EXISTS public.get_insecure_versions CASCADE;
DROP FUNCTION IF EXISTS public.get_all_dependency_licenses CASCADE;
DROP FUNCTION IF EXISTS public.get_version_details CASCADE;
DROP FUNCTION IF EXISTS public.get_dependencies_flat CASCADE;
DROP FUNCTION IF EXISTS public.get_package_details CASCADE;

--------
CREATE OR REPLACE FUNCTION get_all_dependencies(_version_id INTEGER)
RETURNS TABLE(
    package_id INTEGER,
    version_id INTEGER
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE dependency_tree AS (
        SELECT 
            pd.dependency_package_id AS package_id,
            pd.dependency_version_id AS version_id
        FROM 
            package_dependencies AS pd
        WHERE 
            pd.version_id = _version_id
        UNION 
        SELECT 
            pd.dependency_package_id,
            pd.dependency_version_id
        FROM 
            package_dependencies AS pd
        INNER JOIN 
            dependency_tree AS dt ON pd.version_id = dt.version_id
    )
    SELECT DISTINCT
        dt.package_id,
        dt.version_id
    FROM 
        dependency_tree AS dt;
END;
$$;
--------

CREATE OR REPLACE FUNCTION get_insecure_versions()
RETURNS TABLE(
    package_id INTEGER,
    version_id INTEGER
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    WITH all_dependencies AS (
        SELECT
            v.package_id,
            v.id AS version_id,
            (get_all_dependencies(v.id)).package_id AS dependency_package_id,
            (get_all_dependencies(v.id)).version_id AS dependency_version_id
        FROM 
            versions AS v
    )
    SELECT DISTINCT
        ad.package_id,
        ad.version_id
    FROM 
        all_dependencies AS ad
    INNER JOIN 
        versions AS v ON ad.dependency_version_id = v.id
    WHERE 
        v.insecure = true
    UNION 
    SELECT 
        v.package_id,
        v.id AS version_id
    FROM
        versions v
    WHERE 
        v.insecure = true;
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
            (get_all_dependencies(version_id)).*
    )
    SELECT 
        v.license,
        array_agg(distinct p.name)
    FROM 
        all_dependencies AS ad
    INNER JOIN 
        versions AS v ON ad.version_id = v.id
    INNER JOIN
        packages AS p ON ad.package_id = p.id
    GROUP BY 
        v.license;
END;
$$;

-----------------------------------------------------------

CREATE OR REPLACE FUNCTION get_dependencies_flat(_version_id INTEGER)
RETURNS TABLE(
    owner TEXT,
    slug TEXT,
    package_name TEXT,
    version TEXT,
    license TEXT,
    last_updated TIMESTAMP,
    archived BOOLEAN,
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
        u.gh_login AS owner,
        p.slug,
        p.name AS package_name,
        v.version,
        v.license,
        v.created_at AS last_updated,
        p.archived,
        v.insecure
    FROM 
        packages p
    INNER JOIN
        versions v ON p.id = v.package_id
    INNER JOIN
        users u ON p.owner = u.id
    WHERE 
        v.id != _version_id AND v.id = ANY(SELECT version_id FROM all_dependencies);
END;
$$;

-------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION get_version_details(_package_id INTEGER)
RETURNS TABLE(
    id INTEGER,
    version TEXT,
    insecure BOOLEAN,
    createdAt TIMESTAMP,
    size_kb INTEGER,
    dependency_count BIGINT,
    compiler TEXT,
    license TEXT,
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
        (SELECT COUNT(*) 
         FROM package_dependencies pd
         WHERE v.id = pd.version_id) AS dependency_count,
        v.compiler,
        v.license,
        EXISTS(
            SELECT 1
            FROM get_insecure_versions() giv
            WHERE giv.package_id = v.package_id AND giv.version_id = v.id
        ) AS has_insecure_dependency
    FROM 
        versions AS v
    WHERE 
        v.package_id = _package_id;
END;
$$
LANGUAGE plpgsql;

-----------------------------------------------------------------
CREATE OR REPLACE FUNCTION get_package_details(_package_id INTEGER)
RETURNS TABLE(
    id INTEGER,
    name TEXT,
    slug TEXT,
    description TEXT,
    archived BOOLEAN,
    keywords TEXT[],
    stars BIGINT,
    repository TEXT,
    readme TEXT,
    owner TEXT,
    authors TEXT[]
)
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.name,
        p.slug,
        p.description,
        p.archived,
        (SELECT array_agg(k.keyword) FROM package_keywords pk INNER JOIN keywords k ON pk.keyword_id = k.id WHERE pk.package_id = _package_id) AS keywords,
        (SELECT COUNT(*) FROM stars s WHERE s.package_id = _package_id) AS stars,
        p.repository,
        p.readme,
        u.gh_login AS owner,
        ARRAY(
            SELECT a.gh_login
            FROM package_authors pa
            INNER JOIN users a ON pa.author_id = a.id
            WHERE pa.package_id = _package_id
        ) AS authors
    FROM 
        packages p
    INNER JOIN
        users u ON p.owner = u.id
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
        package_dependencies pd ON p.id = pd.package_id
    WHERE 
        pd.dependency_package_id = _package_id;
END;
$$
LANGUAGE plpgsql;
