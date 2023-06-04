DROP FUNCTION IF EXISTS public.insert_user CASCADE;
DROP FUNCTION IF EXISTS public.upsert_user CASCADE;

--maybe?
-- CREATE EXTENSION IF NOT EXISTS pg_stat_statements WITH SCHEMA public;
-- CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;
CREATE OR REPLACE FUNCTION public.upsert_user(
    _gh_id INT,
    _new_login TEXT,
    _new_access_token TEXT,
    _new_avatar TEXT,
    _new_created_at TIMESTAMP WITHOUT TIME ZONE
) RETURNS json
AS $$
DECLARE
    _result json;
BEGIN
    INSERT INTO public.users (
        gh_id,
        gh_login,
        gh_access_token,
        gh_avatar,
        gh_created_at
    ) 
    VALUES (
        _gh_id,
        _new_login,
        _new_access_token,
        _new_avatar,
        _new_created_at
    )
    ON CONFLICT (gh_id) DO 
    UPDATE SET 
        gh_login = excluded.gh_login,
        gh_access_token = excluded.gh_access_token,
        gh_avatar = excluded.gh_avatar,
        gh_created_at = excluded.gh_created_at
    WHERE 
        public.users.gh_id = _gh_id;

    SELECT 
        json_build_object(
            'id', id,
            'banned', banned,
            'ban_reason', ban_reason,
            'ban_timeout', ban_timeout
        ) INTO _result
    FROM
        public.users
    WHERE 
        gh_id = _gh_id;
    RETURN _result;
END;
$$
LANGUAGE plpgsql;


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
