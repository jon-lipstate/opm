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



-- Insert 3 users
INSERT INTO public.users (gh_login, gh_access_token, gh_avatar, gh_id, gh_email) 
VALUES ('jon', 'token1', 'avatar1', 1, 'jon@email.com'), 
       ('odie', 'token2', 'avatar2', 2, 'odie@email.com'), 
       ('freyja', 'token3', 'avatar3', 3, 'freyja@email.com');

-- Insert 2 packages by user1
CALL create_new_package(
    'http server', 
    'a cool http/1.1 server', 
    'readme1', 
    'https://repository1', 
    '1.0.0', -- version
    'MIT',
    1024,
    1, -- published_by
    'dev-2023-05', -- compiler
    'checksum1',
    ARRAY['fancy', 'pants'], -- keywords
    ARRAY[]::INTEGER[] -- no dependencies
);

CALL create_new_package(
    'async runtime', 
    'it does stuff', 
    'readme2', 
    'https://repository2', 
    '1.2.3',
    'BSD 3-Clause',
    2048,
    2, -- published_by
    'dev-2023-05',
    'checksum2',
    ARRAY['sweat', 'pants'], -- keywords
    ARRAY[]::INTEGER[] -- no dependencies
);

CALL create_new_package(
    'i am a tthnteapot', 
    '415', 
    'dont find me', 
    'http://nothere', 
    '99.99.99',
    'GPL',
    99999,
    2, -- published_by
    'dev-2023-05',
    'checksum2',
    ARRAY['no', 'pants'], -- keywords
    ARRAY[1,2]::INTEGER[]  -- deps
);
