DROP FUNCTION IF EXISTS public.insert_new_package CASCADE;
DROP FUNCTION IF EXISTS public.insert_new_version CASCADE;
DROP FUNCTION IF EXISTS public.insert_keywords CASCADE;
DROP PROCEDURE IF EXISTS public.create_new_package CASCADE;
DROP FUNCTION IF EXISTS public.get_package_ids CASCADE;
DROP FUNCTION IF EXISTS public.get_version_ids CASCADE;

CREATE OR REPLACE FUNCTION upsert_package(
    _owner_id INTEGER,
    _name TEXT, 
    _slug TEXT, 
    _description TEXT, 
    _readme TEXT, 
    _url TEXT,
    _state package_state
)
RETURNS INTEGER
LANGUAGE plpgsql AS $$
DECLARE 
    _package_id INTEGER;
BEGIN
    INSERT INTO packages(owner,name, slug, description, readme, url, state)
    VALUES (_owner_id, _name, _slug, _description, _readme, _url, _state)
    ON CONFLICT (owner, name) 
    DO UPDATE SET 
        --slug = EXCLUDED.slug,  --QUESTION: ALLOW RENAMES ???
        description = EXCLUDED.description, 
        readme = EXCLUDED.readme, 
        --url = EXCLUDED.url,  --QUESTION: ALLOW RENAMES ???
        state = EXCLUDED.state
    RETURNING id INTO _package_id;

    RETURN _package_id;
END;
$$;

CREATE OR REPLACE FUNCTION upsert_version(
    _package_id INTEGER,
    _version TEXT,
    _license TEXT,
    _size_kb INTEGER,
    _published_by INTEGER,
    _compiler TEXT
)
RETURNS INTEGER
LANGUAGE plpgsql AS $$
DECLARE 
    _version_id INTEGER;
BEGIN
    INSERT INTO versions(package_id, version, license, size_kb, published_by, compiler)
    VALUES (_package_id, _version, _license, _size_kb, _published_by, _compiler)
    ON CONFLICT (package_id, version)
    DO UPDATE SET
        license = EXCLUDED.license,
        size_kb = EXCLUDED.size_kb,
        published_by = EXCLUDED.published_by,
        compiler = EXCLUDED.compiler
    RETURNING id INTO _version_id;

    RETURN _version_id;
END;
$$;

-----
CREATE OR REPLACE FUNCTION upsert_keywords(
    _package_id INTEGER,
    _keywords TEXT[]
) RETURNS VOID AS $$
DECLARE
    _keyword_id INTEGER;
    _keyword TEXT;
BEGIN
    -- Delete any keywords from package_keywords that are not in the new _keywords array
    DELETE FROM package_keywords 
    WHERE package_id = _package_id 
    AND keyword_id NOT IN (
        SELECT id FROM keywords WHERE keyword = ANY(_keywords)
    );
    
    -- Insert or update keywords and their relation to the package
    FOREACH _keyword IN ARRAY _keywords
    LOOP
        INSERT INTO keywords(keyword) 
        VALUES (_keyword) 
        ON CONFLICT (keyword) DO NOTHING;
        
        SELECT id INTO _keyword_id FROM keywords WHERE keyword = _keyword;

        INSERT INTO package_keywords(package_id, keyword_id)
        VALUES (_package_id, _keyword_id)
        ON CONFLICT (package_id, keyword_id) DO NOTHING;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

--
-- 
CREATE OR REPLACE PROCEDURE upsert_full_package(
    _owner_id INTEGER,
    _name TEXT, 
    _slug TEXT, 
    _description TEXT, 
    _readme TEXT, 
    _url TEXT, 
    _version TEXT,
    _license TEXT,
    _size_kb INTEGER,
    _odin_compiler TEXT,
    _keywords TEXT[],
    _dependency_ids INTEGER[] -- this is PKs of versions table
)
LANGUAGE plpgsql
AS $$
DECLARE 
    _package_id INTEGER;
    _version_id INTEGER;
	_dependency_id INTEGER;
BEGIN
    -- Call the function to insert a new package or update an existing one
    _package_id := upsert_package(_owner_id, _name, _slug, _description, _readme, _url, 'active'::package_state);

    -- Insert into versions table or update an existing one
    _version_id := upsert_version(_package_id, _version, _license, _size_kb, _owner_id, _odin_compiler);

    -- Insert keywords and their relation to the package or update existing ones
    PERFORM upsert_keywords(_package_id, _keywords);

    -- Insert dependencies of the package
    FOREACH _dependency_id IN ARRAY _dependency_ids
    LOOP
        INSERT INTO package_dependencies(version_id, depends_on_version_id)
        VALUES (_version_id, _dependency_id)
        ON CONFLICT ( version_id, depends_on_version_id) DO NOTHING;
    END LOOP;    
END;
$$;




-------------------------------------------------------------------------------------------
-------------------------------------------------------------------------------------------
-------------------------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION get_package_ids(pkgs JSONB[])
RETURNS TABLE(
    package_slug TEXT,
    package_id INT
) AS $$
DECLARE
    pkg JSONB;
BEGIN
    FOREACH pkg IN ARRAY pkgs
    LOOP
        RETURN QUERY
        SELECT 
            p.slug AS package_slug,
            p.id AS package_id
        FROM 
            packages p
        WHERE 
            p.slug = (pkg->>'slug') AND
            (SELECT gh_login FROM users WHERE id = p.owner) = (pkg->>'name');
    END LOOP;
END; $$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_version_ids(pkgs JSONB[])
RETURNS TABLE(
    --package_slug TEXT,
    version_id INT
) AS $$
DECLARE
    pkg JSONB;
BEGIN
    FOREACH pkg IN ARRAY pkgs
    LOOP
        RETURN QUERY
        SELECT 
         --   p.slug AS package_slug,
            v.id AS version_id
        FROM 
            packages p
        JOIN 
            versions v ON v.package_id = p.id
        WHERE 
            v.version = (pkg->>'version');
    END LOOP;
END; $$ LANGUAGE plpgsql;