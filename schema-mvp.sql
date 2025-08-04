-- Extensions
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- For fuzzy text search

-----------------------------------------------------------------------------------

-- Users table (OAuth authenticated)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    github_id VARCHAR(255) UNIQUE,
    discord_id VARCHAR(255) UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(255),
    slug VARCHAR(100) UNIQUE CHECK (slug ~ '^[a-z0-9-]+$'), -- URL Safe
    avatar_url TEXT,
    is_moderator BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    reputation INTEGER NOT NULL DEFAULT 0,
    reputation_rank VARCHAR(50) NOT NULL DEFAULT 'newbie',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    discord_verified BOOLEAN DEFAULT FALSE,
    github_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    CONSTRAINT users_has_oauth CHECK (github_id IS NOT NULL OR discord_id IS NOT NULL)
);
CREATE INDEX idx_users_slug ON users(slug);

-----------------------------------------------------------------------------------

CREATE TYPE package_type AS ENUM ('library', 'project');
CREATE TYPE package_status AS ENUM ('in_work', 'ready', 'archived', 'abandoned');

-- Packages table
CREATE TABLE packages (
    id SERIAL PRIMARY KEY,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slug VARCHAR(100) NOT NULL UNIQUE, -- URL safe
    display_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    type package_type NOT NULL DEFAULT 'project',
    status package_status NOT NULL DEFAULT 'in_work',
    repository_url TEXT NOT NULL,
    license VARCHAR(100),
    view_count BIGINT NOT NULL DEFAULT 0,
    bookmark_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT packages_slug_format CHECK (slug ~ '^[a-z0-9_-]+$'),
    search_vector tsvector
);
CREATE INDEX idx_packages_slug ON packages(slug);
CREATE INDEX idx_packages_search_vector ON packages USING GIN(search_vector);
CREATE INDEX idx_packages_view_count ON packages(view_count DESC);
CREATE INDEX idx_packages_bookmark_count ON packages(bookmark_count DESC);

-- Updated search vector function that includes tags
CREATE OR REPLACE FUNCTION update_package_search_vector() RETURNS trigger AS $$
DECLARE
    tag_text TEXT;
BEGIN
    -- Get all tag names for this package
    SELECT string_agg(t.name, ' ') INTO tag_text
    FROM tags t
    JOIN package_tags pt ON t.id = pt.tag_id
    WHERE pt.package_id = NEW.id;
    
    NEW.search_vector := 
        setweight(to_tsvector('english', COALESCE(NEW.slug, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.display_name, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(tag_text, '')), 'B');
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- Function to update search vector for a specific package (called when tags change)
CREATE OR REPLACE FUNCTION update_package_search_vector_by_id(p_package_id INTEGER) RETURNS VOID AS $$
DECLARE
    tag_text TEXT;
BEGIN
    -- Get all tag names for this package
    SELECT string_agg(t.name, ' ') INTO tag_text
    FROM tags t
    JOIN package_tags pt ON t.id = pt.tag_id
    WHERE pt.package_id = p_package_id;
    
    -- Update the search vector
    UPDATE packages 
    SET search_vector = 
        setweight(to_tsvector('english', COALESCE(slug, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(display_name, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(tag_text, '')), 'B')
    WHERE id = p_package_id;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER packages_search_vector_update
    BEFORE INSERT OR UPDATE OF slug, display_name, description
    ON packages
    FOR EACH ROW
    EXECUTE FUNCTION update_package_search_vector();

-----------------------------------------------------------------------------------

-- Package views tracking
CREATE TABLE package_views (
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    viewed_at DATE NOT NULL DEFAULT CURRENT_DATE
);

-- Create unique constraint that handles NULL user_id
CREATE UNIQUE INDEX idx_package_views_unique ON package_views(package_id, user_id, viewed_at) 
WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX idx_package_views_unique_anon ON package_views(package_id, viewed_at) 
WHERE user_id IS NULL;

-- Create index for view tracking
CREATE INDEX idx_package_views_package_date ON package_views(package_id, viewed_at DESC);

-- Stored procedure for tracking views
CREATE OR REPLACE FUNCTION track_package_view(p_package_id INTEGER, p_user_id INTEGER DEFAULT NULL) 
RETURNS VOID AS $$
DECLARE
    v_inserted BOOLEAN;
BEGIN
    -- Try to insert, if exists already, do nothing (user already viewed today)
    BEGIN
        INSERT INTO package_views (package_id, user_id, viewed_at)
        VALUES (p_package_id, p_user_id, CURRENT_DATE);
        v_inserted := TRUE;
    EXCEPTION WHEN unique_violation THEN
        -- View already exists for this user/package/date
        v_inserted := FALSE;
    END;
    
    -- Only increment count if we actually inserted a new view
    IF v_inserted THEN
        UPDATE packages 
        SET view_count = view_count + 1
        WHERE id = p_package_id;
    END IF;
END;
$$ LANGUAGE plpgsql;

-----------------------------------------------------------------------------------

-- Tags table
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    added_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    search_vector tsvector
);
CREATE INDEX idx_tags_search_vector ON tags USING GIN(search_vector);

CREATE OR REPLACE FUNCTION update_tag_search_vector() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', COALESCE(NEW.name, ''));
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

-- Trigger for tags
CREATE TRIGGER tags_search_vector_update
    BEFORE INSERT OR UPDATE OF name
    ON tags
    FOR EACH ROW
    EXECUTE FUNCTION update_tag_search_vector();

-- Package tags (many-to-many)
CREATE TABLE package_tags (
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    score INTEGER NOT NULL DEFAULT 0, -- driven by tag_votes
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (package_id, tag_id)
);

-- Tag votes table
CREATE TABLE tag_votes (
    id SERIAL PRIMARY KEY,
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vote_value INTEGER NOT NULL CHECK (vote_value BETWEEN -10 AND 10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(package_id, tag_id, user_id)
);

-- Bookmarks table (user favorites)
CREATE TABLE bookmarks (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, package_id)
);

-- Flags table for moderation
CREATE TABLE flags (
    id SERIAL PRIMARY KEY,
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason VARCHAR(50) NOT NULL,
    details TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, reviewed, resolved, dismissed
    resolved_by INTEGER REFERENCES users(id),
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-----------------------------------------------------------------------------------
-- Indexes for performance

CREATE INDEX idx_packages_author ON packages(author_id);
CREATE INDEX idx_packages_status ON packages(status);
CREATE INDEX idx_packages_type ON packages(type);
CREATE INDEX idx_packages_created_at ON packages(created_at DESC);
CREATE INDEX idx_packages_display_name_trgm ON packages USING GIN(display_name gin_trgm_ops);

CREATE INDEX idx_bookmarks_user ON bookmarks(user_id);
CREATE INDEX idx_bookmarks_package ON bookmarks(package_id);

CREATE INDEX idx_flags_package ON flags(package_id);
CREATE INDEX idx_flags_status ON flags(status);

CREATE INDEX idx_tag_votes_package_tag ON tag_votes(package_id, tag_id);
CREATE INDEX idx_tag_votes_user ON tag_votes(user_id);

-----------------------------------------------------------------------------------
-- Triggers

-- Update timestamp trigger
CREATE OR REPLACE FUNCTION update_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER packages_updated_at BEFORE UPDATE ON packages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER flags_updated_at BEFORE UPDATE ON flags
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tag_votes_updated_at BEFORE UPDATE ON tag_votes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Update tag usage count trigger
CREATE OR REPLACE FUNCTION update_tag_usage_count() RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE tags SET usage_count = usage_count + 1 WHERE id = NEW.tag_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE tags SET usage_count = usage_count - 1 WHERE id = OLD.tag_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER package_tags_usage_count
    AFTER INSERT OR DELETE ON package_tags
    FOR EACH ROW EXECUTE FUNCTION update_tag_usage_count();

-- Update bookmark count trigger
CREATE OR REPLACE FUNCTION update_bookmark_count() RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE packages SET bookmark_count = bookmark_count + 1 WHERE id = NEW.package_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE packages SET bookmark_count = bookmark_count - 1 WHERE id = OLD.package_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER bookmarks_count_update
    AFTER INSERT OR DELETE ON bookmarks
    FOR EACH ROW EXECUTE FUNCTION update_bookmark_count();

-- Trigger to update package search vector when tags change
CREATE OR REPLACE FUNCTION trigger_update_package_search_on_tag_change() RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        PERFORM update_package_search_vector_by_id(OLD.package_id);
        RETURN OLD;
    ELSE
        PERFORM update_package_search_vector_by_id(NEW.package_id);
        RETURN NEW;
    END IF;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER package_tags_search_update
    AFTER INSERT OR DELETE
    ON package_tags
    FOR EACH ROW
    EXECUTE FUNCTION trigger_update_package_search_on_tag_change();