DROP TABLE IF EXISTS public.package_dependencies CASCADE;
DROP TABLE IF EXISTS public.stars CASCADE;
DROP TABLE IF EXISTS public.package_keywords CASCADE;
DROP TABLE IF EXISTS public.keywords CASCADE;
DROP TABLE IF EXISTS public.versions CASCADE;
DROP TABLE IF EXISTS public.users CASCADE;
DROP TABLE IF EXISTS public.packages CASCADE;
DROP TABLE IF EXISTS public.actions CASCADE;
DROP TABLE IF EXISTS public.api_tokens CASCADE;
DROP TABLE IF EXISTS public.create_limits CASCADE;
DROP TABLE IF EXISTS public.package_authors CASCADE;
DROP TABLE IF EXISTS public.reserved_names CASCADE;


CREATE TABLE IF NOT EXISTS public.users (
	id SERIAL PRIMARY KEY,
	gh_login TEXT NOT NULL,
    gh_access_token TEXT NOT NULL,
	gh_avatar TEXT,
    gh_id INT NOT NULL, -- not certain if i need?
    gh_email TEXT NOT NULL, -- not 100% sure i need this? make nullable, only need if publisher?
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
	last_login TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS public.packages (
    id SERIAL PRIMARY KEY,
	owner INTEGER REFERENCES users(id),
	name TEXT NOT NULL,
	slug TEXT NOT NULL,
	readme_rendered_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL, -- internal use only i think?
	updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL, -- internal use only i think?
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
	description TEXT,
	readme TEXT,
	TSV tsvector,
	repository TEXT, -- url or git, both?
	archived BOOLEAN DEFAULT false,
    UNIQUE(name,owner)
);
UPDATE packages SET TSV = to_tsvector('english', name || ' ' || description);
CREATE INDEX packages_tsv_idx ON packages USING gin(tsv);


CREATE TABLE IF NOT EXISTS public.package_authors (
    package_id INTEGER NOT NULL REFERENCES packages(id),
    author_id INTEGER NOT NULL REFERENCES users(id),
	is_admin BOOLEAN DEFAULT false,
    PRIMARY KEY(package_id, author_id)
);

CREATE TABLE IF NOT EXISTS public.versions (
    id SERIAL PRIMARY KEY,
	package_id INTEGER REFERENCES packages(id),
	version TEXT NOT NULL,
	created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL,
	license TEXT NOT NULL,
	size_kb INTEGER,
	published_by INTEGER REFERENCES users(id),
	insecure BOOLEAN DEFAULT false, -- aka yank
	compiler TEXT NOT NULL,
	downloads INTEGER DEFAULT 0,
	checksum CHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS public.keywords (
    id SERIAL PRIMARY KEY,
    keyword TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS public.package_keywords (
    package_id INTEGER NOT NULL REFERENCES packages(id),
    keyword_id INTEGER NOT NULL REFERENCES keywords(id),
    PRIMARY KEY(package_id, keyword_id)
);

CREATE TABLE package_dependencies (
    id SERIAL PRIMARY KEY,
    package_id INTEGER REFERENCES packages(id),
    version_id INTEGER REFERENCES versions(id),
    dependency_package_id INTEGER REFERENCES packages(id),
    dependency_version_id INTEGER REFERENCES versions(id),
	-- dependency_type TEXT -- e.g., 'dev', 'prod', etc.
    UNIQUE(package_id, version_id, dependency_package_id, dependency_version_id)
);


CREATE TABLE IF NOT EXISTS public.stars (
    user_id INTEGER NOT NULL REFERENCES users(id),
    package_id INTEGER NOT NULL REFERENCES packages(id),
    PRIMARY KEY(user_id, package_id)
);
-- meta-tables
CREATE TABLE IF NOT EXISTS public.reserved_names (
    name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS public.create_limits (
	-- TODO PK
    user_id INTEGER NOT NULL REFERENCES users(id),
	actions INTEGER NOT NULL,
    last_refill timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS public.actions (
	-- TODO PK
    user_id INTEGER NOT NULL REFERENCES users(id),
    version_id INTEGER NOT NULL REFERENCES versions(id),
	action TEXT NOT NULL,
    time_of timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS public.api_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL, 
    last_used_at timestamp without time zone,
	revoked BOOLEAN DEFAULT false,
	scopes TEXT[]
);