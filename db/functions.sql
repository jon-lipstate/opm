CREATE EXTENSION pldbgapi;

DROP FUNCTION IF EXISTS public.insert_user CASCADE;
DROP FUNCTION IF EXISTS public.insert_new_package CASCADE;
DROP PROCEDURE IF EXISTS public.create_new_package CASCADE;

-- CALL insert_user('gh_login_value', 'gh_access_token_value', 'gh_avatar_value', gh_id_value, 'gh_email_value');
CREATE OR REPLACE PROCEDURE insert_user(
    _gh_login VARCHAR,
    _gh_access_token VARCHAR,
    _gh_avatar TEXT,
    _gh_id INT,
    _gh_email VARCHAR
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


-- SELECT insert_package(' Name', ' Description', ' README', ' Repository', false);
CREATE OR REPLACE FUNCTION insert_new_package(
    _name VARCHAR, 
    _description TEXT, 
    _readme TEXT, 
    _repository TEXT, 
    _owner INTEGER
)
RETURNS INTEGER
LANGUAGE plpgsql AS $$
DECLARE 
    _package_id INTEGER;
    existing_count INTEGER;
BEGIN
    -- Check if package with same name and owner exists
    SELECT COUNT(*) INTO existing_count
    FROM public.packages
    WHERE name = _name AND owner = _owner;
    IF existing_count > 0 THEN
        RAISE EXCEPTION 'Package with same name already exists for this owner';
    END IF;

    INSERT INTO packages(name, description, readme, repository, archived, owner)
    VALUES (_name, _description, _readme, _repository, false, _owner) 
    RETURNING id INTO _package_id;

    RETURN _package_id;
END;
$$;

-- SELECT insert_version(1, '1.0.0', 'MIT', 500, 1, false, '1.23.0', 'abc123');
CREATE OR REPLACE FUNCTION insert_version(
    _package_id INTEGER, 
    _version VARCHAR(255), 
    _license VARCHAR(255),
    _size_kb INTEGER,
    _published_by INTEGER,
    _insecure BOOLEAN,
    _odin_compiler VARCHAR(255),
    _checksum CHAR(64)
) 
RETURNS VOID AS $$
BEGIN
    INSERT INTO public.versions (package_id, version, license, size_kb, published_by, insecure, odin_compiler, checksum)
    VALUES (_package_id, _version, _license, _size_kb, _published_by, _insecure, _odin_compiler, _checksum);
END;
$$ LANGUAGE plpgsql;

--
-- 
CREATE OR REPLACE PROCEDURE create_new_package(
    _name VARCHAR, 
    _description TEXT, 
    _readme TEXT, 
    _repository TEXT, 
    _version VARCHAR,
    _license VARCHAR,
    _size_kb INTEGER,
    _published_by INTEGER,
    _odin_compiler VARCHAR,
    _checksum CHAR(64),
    _keywords TEXT[],
    _dependencies INTEGER[] -- this is PKs of pk-dep table
)
LANGUAGE plpgsql
AS $$
DECLARE 
    _package_id INTEGER;
    _version_id INTEGER;
    _keyword_id INTEGER;
    _keyword VARCHAR;
	_dependency INTEGER;
    existing_count INTEGER;
BEGIN
    -- Call the function to insert a new package
    _package_id := insert_new_package(_name, _description, _readme, _repository, _published_by);
    
    
    -- Insert into versions table
        -- Check if Same package/version exists
        SELECT COUNT(*) INTO existing_count
        FROM public.versions
        WHERE version = _version AND package_id = _package_id;
        IF existing_count > 0 THEN
            RAISE EXCEPTION 'Same Version for this package detected';
        END IF;

    INSERT INTO versions(package_id, version, license, size_kb, published_by, insecure, odin_compiler, checksum)
    VALUES (_package_id, _version, _license, _size_kb, _published_by, false, _odin_compiler, _checksum)
    RETURNING id INTO _version_id;

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

    -- Insert dependencies of the package
    FOREACH _dependency IN ARRAY _dependencies
    LOOP
        INSERT INTO package_dependencies(package_id, version_id)
        VALUES (_package_id, _dependency)
        ON CONFLICT (package_id, version_id) DO NOTHING;
    END LOOP;
    
END;
$$;
