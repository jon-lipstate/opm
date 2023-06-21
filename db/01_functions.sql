-----------------------------------------------------------------------------
DROP FUNCTION IF EXISTS public.upsert_package CASCADE;
CREATE OR REPLACE FUNCTION upsert_package(
    _owner_id INTEGER,
    _host_name TEXT,
    _owner_name TEXT,
    _repo_name TEXT, 
    _description TEXT, 
    _url TEXT,
    _state package_state
)
RETURNS INTEGER
LANGUAGE plpgsql AS $$
DECLARE 
    _package_id INTEGER;
BEGIN
    INSERT INTO packages(owner_id, host_name, owner_name, repo_name, description, url, state)
    VALUES (_owner_id, _host_name, _owner_name, _repo_name, _description, _url, _state)
    ON CONFLICT (host_name, owner_name, repo_name) 
    DO UPDATE SET 
        description = EXCLUDED.description, 
        --url = EXCLUDED.url, 
        state = EXCLUDED.state
    RETURNING id INTO _package_id;

    RETURN _package_id;
END;
$$;
-----------------------------------------------------------------------------
DROP FUNCTION IF EXISTS public.insert_version CASCADE;
CREATE OR REPLACE FUNCTION insert_version(
    _package_id INTEGER,
    _version TEXT,
    _license TEXT,
    _size_kb INTEGER,
    _published_by INTEGER,
    _compiler TEXT,
    _commit_hash TEXT,
    _readme TEXT
)
RETURNS INTEGER
LANGUAGE plpgsql AS $$
DECLARE 
    _version_id INTEGER;
BEGIN
    INSERT INTO versions(package_id, version, license, size_kb, published_by, compiler, commit_hash, readme)
    VALUES (_package_id, _version, _license, _size_kb, _published_by, _compiler, _commit_hash, _readme)
    ON CONFLICT (package_id, version)
    DO NOTHING
    RETURNING id INTO _version_id;

    IF _version_id IS NULL THEN
        RAISE EXCEPTION 'Version `%` already exists for package id %', _version, _package_id;
    END IF;

    RETURN _version_id;
END;
$$;
-- NOT in use atm:
DROP FUNCTION IF EXISTS public.upsert_version CASCADE;
CREATE OR REPLACE FUNCTION upsert_version(
    _package_id INTEGER,
    _version TEXT,
    _license TEXT,
    _size_kb INTEGER,
    _published_by INTEGER,
    _compiler TEXT,
    _commit_hash TEXT,
    _readme TEXT
)
RETURNS INTEGER
LANGUAGE plpgsql AS $$
DECLARE 
    _version_id INTEGER;
BEGIN
    INSERT INTO versions(package_id, version, license, size_kb, published_by, compiler, commit_hash, readme)
    VALUES (_package_id, _version, _license, _size_kb, _published_by, _compiler, _commit_hash, _readme)
    ON CONFLICT (package_id, version)
    DO UPDATE SET
        license = EXCLUDED.license,
        size_kb = EXCLUDED.size_kb,
        published_by = EXCLUDED.published_by,
        compiler = EXCLUDED.compiler,
        commit_hash = EXCLUDED.commit_hash,
        readme = EXCLUDED.readme
    RETURNING id INTO _version_id;

    RETURN _version_id;
END;
$$;

-----------------------------------------------------------------------------
DROP FUNCTION IF EXISTS public.upsert_keywords CASCADE;
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

-----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS public.upsert_full_package CASCADE;
CREATE OR REPLACE PROCEDURE upsert_full_package(
    _owner_id INTEGER,
    _host_name TEXT,
    _owner_name TEXT,
    _repo_name TEXT,
    _description TEXT,
    _readme TEXT,
    _url TEXT,
    _version TEXT,
    _license TEXT,
    _size_kb INTEGER,
    _odin_compiler TEXT,
    _commit_hash TEXT,
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
    _package_id := upsert_package(_owner_id, _host_name, _owner_name, _repo_name, _description, _url, 'active'::package_state);

    -- Insert into versions table or update an existing one
    _version_id := insert_version(_package_id, _version, _license, _size_kb, _owner_id, _odin_compiler, _commit_hash,_readme);

    -- Insert keywords and their relation to the package or update existing ones
    PERFORM upsert_keywords(_package_id, _keywords);

    -- Insert dependencies of the package
    FOREACH _dependency_id IN ARRAY _dependency_ids
    LOOP
        INSERT INTO package_dependencies(version_id, depends_on_id)
        VALUES (_version_id, _dependency_id)
        ON CONFLICT ( version_id, depends_on_id) DO NOTHING;
    END LOOP;
END;
$$;

-------------------------------------------------------------------------------------------
-------------------------------------------------------------------------------------------
-------------------------------------------------------------------------------------------
DROP FUNCTION IF EXISTS public.get_package_ids CASCADE;
CREATE OR REPLACE FUNCTION get_package_ids(pkgs JSONB[])
RETURNS TABLE(
    package_id INT
) AS $$
DECLARE
    pkg JSONB;
BEGIN
    FOREACH pkg IN ARRAY pkgs
    LOOP
        RETURN QUERY
        SELECT 
            p.id AS package_id
        FROM 
            packages p
        WHERE 
            p.host_name = (pkg->>'host_name') AND
            p.owner_name = (pkg->>'owner_name') AND
            p.repo_name = (pkg->>'repo_name');
    END LOOP;
END; $$ LANGUAGE plpgsql;

-----------------------------------------------------------------------------
DROP FUNCTION IF EXISTS public.get_version_ids CASCADE;
CREATE OR REPLACE FUNCTION get_version_ids(pkgs JSONB[])
RETURNS TABLE(
    version_id INT
) AS $$
DECLARE
    pkg JSONB;
BEGIN
    FOREACH pkg IN ARRAY pkgs
    LOOP
        RETURN QUERY
        SELECT 
            v.id AS version_id
        FROM 
            versions v
        WHERE 
            v.package_id = (pkg->>'id')::INTEGER AND
            v.version = (pkg->>'version');
    END LOOP;
END; $$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------
DROP FUNCTION IF EXISTS insert_token;
CREATE OR REPLACE FUNCTION insert_token(_user_id INTEGER, _name TEXT, _token TEXT)
RETURNS VOID AS $$
BEGIN
    INSERT INTO public.api_tokens(user_id, name, token_hash)
    VALUES (_user_id, _name, digest(_token, 'sha256'));
END;
$$ LANGUAGE plpgsql;

----------------------------------------------------------------------------------------
DROP FUNCTION IF EXISTS verify_token;
CREATE OR REPLACE FUNCTION verify_token(_token TEXT)
RETURNS INTEGER AS $$
DECLARE
    _user_id INTEGER;
BEGIN
    SELECT user_id INTO _user_id 
    FROM public.api_tokens 
    WHERE digest(_token, 'sha256') = token_hash
    AND revoked = false;

    IF _user_id IS NULL THEN
        RAISE EXCEPTION 'Invalid or revoked token';
    ELSE
        UPDATE public.api_tokens 
        SET last_touched = NOW() 
        WHERE digest(_token, 'sha256') = token_hash
        AND revoked = false;
    END IF;

    RETURN _user_id;
END;
$$ LANGUAGE plpgsql;
