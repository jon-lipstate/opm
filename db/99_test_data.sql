DELETE FROM public.package_dependencies CASCADE;
DELETE FROM public.stars CASCADE;
DELETE FROM public.package_keywords CASCADE;
DELETE FROM public.keywords CASCADE;
DELETE FROM public.versions CASCADE;
DELETE FROM public.users CASCADE;
DELETE FROM public.packages CASCADE;
DELETE FROM public.actions;
DELETE FROM public.api_tokens;
DELETE FROM public.create_limits;
DELETE FROM public.package_authors;
DELETE FROM public.reserved_names;

ALTER SEQUENCE public.package_dependencies_id_seq RESTART WITH 1;
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
INSERT INTO public.users (gh_login, gh_access_token, gh_avatar, gh_id) 
VALUES ('jon', 'token1', 'avatar1', 1), 
       ('odie', 'token2', 'avatar2', 2), 
       ('freyja', 'token3', 'avatar3', 3);

-- Insert 2 packages by user1
CALL create_new_package(
    1, -- published_by
    'http server', 
    'http-server', 
    'a cool http/1.1 server', 
    'readme1', 
    'https://www.youtube.com/watch?v=xvFZjo5PgG0', 
    '1.0.0', -- version
    'MIT',
    1024,
    'dev-2023-05', -- compiler
    ARRAY['fancy', 'pants'], -- keywords
    ARRAY[]::INTEGER[] -- no dependencies
);

CALL create_new_package(
    2, -- published_by
    'async runtime', 
    'async-runtime', 
    'it does stuff', 
    'readme2', 
    'https://vkontech.com/exploring-the-async-await-state-machine-main-workflow-and-state-transitions/', 
    '1.2.3',
    'BSD 3-Clause',
    2048,
    'dev-2023-05',
    ARRAY['sweat', 'pants'], -- keywords
    ARRAY[]::INTEGER[] -- no dependencies
);

CALL create_new_package(
    2, -- published_by
    'i am a teapot', 
    'i-am-a-teapot', 
    '418', 
    'dont find me', 
    'https://www.webfx.com/web-development/glossary/http-status-codes/what-is-a-418-status-code/', 
    '99.1.99',
    'GPL',
    1234,
    'dev-2023-05',
    ARRAY['no', 'pants'], -- keywords
    ARRAY[1,2]::INTEGER[]  -- deps
);
