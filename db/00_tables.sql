DROP TABLE IF EXISTS public.package_dependencies CASCADE;
DROP TABLE IF EXISTS public.bookmarks CASCADE;
DROP TABLE IF EXISTS public.package_keywords CASCADE;
DROP TABLE IF EXISTS public.keywords CASCADE;
DROP TABLE IF EXISTS public.versions CASCADE;
DROP TABLE IF EXISTS public.users CASCADE;
DROP TYPE IF EXISTS public.package_state CASCADE;
DROP TYPE IF EXISTS public.severity CASCADE;
DROP TABLE IF EXISTS public.packages CASCADE;
DROP TABLE IF EXISTS public.actions CASCADE;
DROP TABLE IF EXISTS public.api_tokens CASCADE;
DROP TABLE IF EXISTS public.create_limits CASCADE;
DROP TABLE IF EXISTS public.package_authors CASCADE;
DROP TABLE IF EXISTS public.reserved_names CASCADE;
DROP TABLE IF EXISTS public.scopes CASCADE;
DROP TABLE IF EXISTS public.user_scopes CASCADE;
DROP TABLE IF EXISTS public.security_issues CASCADE;
DROP TABLE IF EXISTS public.background_jobs CASCADE;

DROP EXTENSION IF EXISTS pgcrypto;
--
CREATE TYPE package_state AS ENUM ('unpublished','active', 'archived','moderated','deleted');
CREATE TYPE severity AS ENUM ('low','medium', 'high', 'critical');

--

CREATE TABLE IF NOT EXISTS public.users (
	id SERIAL PRIMARY KEY,
	gh_login TEXT NOT NULL,
    gh_access_token TEXT NOT NULL,
	gh_avatar TEXT,
    gh_id INT NOT NULL UNIQUE, -- stable id, login can change
	gh_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- todo: move to own table?
    banned BOOLEAN NOT NULL DEFAULT false,
    ban_reason TEXT,
    ban_timeout TIMESTAMP WITH TIME ZONE 
);
CREATE TABLE IF NOT EXISTS public.packages (
    id SERIAL PRIMARY KEY,
    host_name TEXT NOT NULL, -- host e.g. 'github.com', 'sr.ht'
    owner_name TEXT NOT NULL, -- owner name in relation to the repo e.g. 'jon-lipstate', '~mjolnir'
    repo_name TEXT NOT NULL, -- repository name from url
    description TEXT,
    url TEXT, -- url or git, both?
    state package_state DEFAULT 'active',
    owner_id INTEGER REFERENCES users(id), -- github user-id
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    TSV tsvector,
    UNIQUE(host_name, owner_name, repo_name)
);

UPDATE packages SET TSV = to_tsvector('english', repo_name || ' ' || description);
CREATE INDEX packages_tsv_idx ON packages USING gin(tsv);


CREATE TABLE IF NOT EXISTS public.package_authors (
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES users(id),
	is_admin BOOLEAN DEFAULT false,
    PRIMARY KEY(package_id, author_id)
);

CREATE TABLE IF NOT EXISTS public.versions (
    id SERIAL PRIMARY KEY,
    package_id INTEGER REFERENCES packages(id) ON DELETE CASCADE,
    version TEXT NOT NULL,
    readme TEXT,
    commit_hash TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    license TEXT NOT NULL,
    size_kb INTEGER,
    published_by INTEGER REFERENCES users(id),
    insecure BOOLEAN DEFAULT false, 
    compiler TEXT NOT NULL,
    downloads INTEGER DEFAULT 0,
    UNIQUE(version,package_id)
);

-- idk if this is a keeper, only really care to point users to more info
CREATE TABLE security_issues (
    id SERIAL PRIMARY KEY,
    version_id INTEGER REFERENCES versions(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    link TEXT, -- http-url
    level severity NOT NULL,
    reporter_id INTEGER REFERENCES users(id), -- must login or anon ok??
	reported_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    pending BOOLEAN DEFAULT false, -- if left pending for x-days triggers a confirm?
    UNIQUE(version_id, name)
);
CREATE TABLE IF NOT EXISTS public.keywords (
    id SERIAL PRIMARY KEY,
    keyword TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS public.package_keywords (
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    keyword_id INTEGER NOT NULL REFERENCES keywords(id) ON DELETE CASCADE, -- TBD if this cacsade should be here? 
    PRIMARY KEY(package_id, keyword_id)
);

CREATE TABLE package_dependencies (
    version_id INTEGER REFERENCES versions(id) ON DELETE CASCADE,
    depends_on_id INTEGER REFERENCES versions(id) ON DELETE CASCADE,
    UNIQUE(version_id, depends_on_id)
);

CREATE TABLE IF NOT EXISTS public.bookmarks (
    user_id INTEGER NOT NULL REFERENCES users(id),
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, package_id)
);

-- meta-tables
--https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file#naming-conventions
--Odin Reserved Names
--https://github.com/LDNOOBW/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words/blob/master/en
--https://www.cs.cmu.edu/~biglou/resources/bad-words.txt
CREATE TABLE IF NOT EXISTS public.reserved_names (
    name TEXT NOT NULL
);

-- throttle package creation to prevent spamming
CREATE TABLE IF NOT EXISTS public.create_limits (
	-- TODO PK
    user_id INTEGER NOT NULL REFERENCES users(id),
	actions INTEGER NOT NULL,
    last_refill timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- action logs
CREATE TABLE IF NOT EXISTS public.actions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    table_name TEXT NOT NULL,
    row_id INTEGER NOT NULL,
    action TEXT NOT NULL,
    comment TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- user-cli auth tokens
CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;
CREATE TABLE IF NOT EXISTS public.api_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    token_hash bytea NOT NULL,
    created_at timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL, 
    last_touched timestamp WITH time zone, -- last used | edited
	revoked BOOLEAN DEFAULT false
);
-- INSERT INTO api_tokens (user_id, name, token_hash, scopes)
-- VALUES (_user_id, _name, digest(_token, 'sha256'), _scopes);

-- SELECT * FROM api_tokens 
-- WHERE user_id = _user_id 
-- AND token_hash = digest(_token, 'sha256');

CREATE TABLE public.background_jobs (
    id bigint NOT NULL,
    job_type text NOT NULL,
    data jsonb DEFAULT '{}'::jsonb NOT NULL,
    retries integer DEFAULT 0 NOT NULL,
    last_retry timestamp WITH time zone DEFAULT '1970-01-01 00:00:00'::timestamp WITH time zone NOT NULL,
    created_at timestamp WITH time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS public.scopes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS public.user_scopes (
    user_id INTEGER NOT NULL REFERENCES users(id),
    scope_id INTEGER NOT NULL REFERENCES scopes(id) ON DELETE CASCADE,
    granted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, scope_id)
);
