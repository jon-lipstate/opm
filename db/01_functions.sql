DROP FUNCTION IF EXISTS public.insert_user CASCADE;
DROP FUNCTION IF EXISTS public.insert_new_package CASCADE;
DROP PROCEDURE IF EXISTS public.create_new_package CASCADE;

--maybe?
-- CREATE EXTENSION IF NOT EXISTS pg_stat_statements WITH SCHEMA public;
-- CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


-- CALL insert_user('gh_login_value', 'gh_access_token_value', 'gh_avatar_value', gh_id_value, 'gh_email_value');
CREATE OR REPLACE PROCEDURE insert_user(
    _gh_login TEXT,
    _gh_access_token TEXT,
    _gh_avatar TEXT,
    _gh_id INT,
    _gh_email TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO public.users(
        gh_login,
        gh_access_token,
        gh_avatar,
        gh_id,
        gh_email
    ) VALUES (
        _gh_login,
        _gh_access_token,
        _gh_avatar,
        _gh_id,
        _gh_email
    );
END;
$$;

-- INSERT INTO users (gh_login, gh_access_token, gh_avatar, gh_id, gh_created_at) 
-- VALUES ('new_login', 'new_token', 'new_avatar', new_id, 'new_created_at') 
-- ON CONFLICT (gh_login) DO UPDATE 
-- SET gh_access_token = excluded.gh_access_token,
--     gh_avatar = excluded.gh_avatar,
--     gh_id = excluded.gh_id,
--     gh_created_at = excluded.gh_created_at;

CREATE OR REPLACE FUNCTION insert_new_package(
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
    RETURNING id INTO _package_id;

    RETURN _package_id;
END;
$$;

CREATE OR REPLACE FUNCTION insert_new_version(
    _package_id INTEGER,
    _version TEXT,
    _license TEXT,
    _size_kb INTEGER,
    _published_by INTEGER,
    _odin_compiler TEXT
) RETURNS INTEGER AS $$
DECLARE
    _version_id INTEGER;
    existing_count INTEGER;
BEGIN
    -- Check if Same package/version exists
    SELECT COUNT(*) INTO existing_count
    FROM public.versions
    WHERE version = _version AND package_id = _package_id;

    IF existing_count > 0 THEN
        RAISE EXCEPTION 'Same Version for this package detected';
    END IF;

    INSERT INTO versions(package_id, version, license, size_kb, published_by, compiler )
    VALUES (_package_id, _version, _license, _size_kb, _published_by, _odin_compiler )
    RETURNING id INTO _version_id;

    RETURN _version_id;
END;
$$ LANGUAGE plpgsql;
-----
CREATE OR REPLACE FUNCTION insert_keywords(
    _package_id INTEGER,
    _keywords TEXT[]
) RETURNS VOID AS $$
DECLARE
    _keyword_id INTEGER;
    _keyword TEXT;
BEGIN
    -- Insert keywords and their relation to the package
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
CREATE OR REPLACE PROCEDURE create_new_package(
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
    _dependency_ids INTEGER[] -- this is PKs of pk-dep table
)
LANGUAGE plpgsql
AS $$
DECLARE 
    _package_id INTEGER;
    _version_id INTEGER;
    _keyword_id INTEGER;
    _keyword TEXT;
	_dependency INTEGER;
    existing_count INTEGER;
BEGIN
    -- Call the function to insert a new package -- TODO:: SET TO UNPUBLISHED
    _package_id := insert_new_package(_owner_id, _name, _slug, _description, _readme, _url, 'active'::package_state);

    -- Insert into versions table
    _version_id := insert_new_version(_package_id, _version, _license, _size_kb, _owner_id, _odin_compiler);

       -- Insert keywords and their relation to the package
    PERFORM insert_keywords(_package_id, _keywords);


    -- Insert dependencies of the package
    FOREACH _dependency IN ARRAY _dependency_ids
    LOOP
        INSERT INTO package_dependencies(package_id, version_id, dependency_package_id, dependency_version_id)
        VALUES (_package_id, _version_id, _dependency, _dependency)
        ON CONFLICT (package_id, version_id, dependency_package_id, dependency_version_id) DO NOTHING;
    END LOOP;
    
END;
$$;
