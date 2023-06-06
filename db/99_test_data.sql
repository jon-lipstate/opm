ALTER SEQUENCE public.api_tokens_id_seq RESTART WITH 1;
ALTER SEQUENCE public.packages_id_seq RESTART WITH 1;
ALTER SEQUENCE public.users_id_seq RESTART WITH 1;
ALTER SEQUENCE public.versions_id_seq RESTART WITH 1;
ALTER SEQUENCE public.keywords_id_seq RESTART WITH 1;

INSERT INTO scopes(name) VALUES
('publish_own'),('update_own'),('delete_own'),('publish_any'),
('update_any'),('delete_any'),('comment'),('vote'),
('timeout_user'),('ban_user'),('update_user');

-- Insert 3 users
ALTER SEQUENCE public.api_tokens_id_seq RESTART WITH 1;
ALTER SEQUENCE public.packages_id_seq RESTART WITH 1;
ALTER SEQUENCE public.users_id_seq RESTART WITH 1;
ALTER SEQUENCE public.versions_id_seq RESTART WITH 1;
ALTER SEQUENCE public.keywords_id_seq RESTART WITH 1;
INSERT INTO public.users (gh_login, gh_access_token, gh_avatar, gh_id) 

VALUES ('jon', 'token1', 'avatar1', 1), 
       ('odie', 'token2', 'avatar2', 2), 
       ('freyja', 'token3', 'avatar3', 3);

-- Insert 2 packages by user1
CALL upsert_full_package(
    1, -- published_by
    'github.com', -- host_name
    'jon', -- owner_name
    'http-server', -- repo_name
    'a cool http/1.1 server', 
    'readme1', 
    'https://www.youtube.com/watch?v=xvFZjo5PgG0', 
    '1.0.0', -- version
    'MIT',
    1024,
    'dev-2023-05', -- compiler
    'abcd1234',
    ARRAY['fancy', 'pants'], -- keywords
    ARRAY[]::INTEGER[] -- no dependencies
);

CALL upsert_full_package(
    2, -- published_by
    'github.com', -- host_name
    'odie', -- owner_name
    'async-runtime', -- repo_name
    'it does stuff', 
    'readme2', 
    'https://vkontech.com/exploring-the-async-await-state-machine-main-workflow-and-state-transitions/', 
    '1.2.3',
    'BSD 3-Clause',
    2048,
    'dev-2023-05',
    'abcd1234',
    ARRAY['sweat', 'pants'], -- keywords
    ARRAY[]::INTEGER[] -- no dependencies
);

CALL upsert_full_package(
    2, -- published_by
    'github.com', -- host_name
    'odie', -- owner_name
    'i-am-a-teapot', -- repo_name
    '418', 
    'peekaboo', 
    'https://www.webfx.com/web-development/glossary/http-status-codes/what-is-a-418-status-code/', 
    '99.1.99',
    'LPG',
    1234,
    'dev-2023-05',
    'abcd1234',
    ARRAY['no', 'pants'], -- keywords
    ARRAY[1,2]::INTEGER[]  -- deps
);
