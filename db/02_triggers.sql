DROP FUNCTION IF EXISTS public.packages_search_trigger CASCADE;
DROP FUNCTION IF EXISTS public.update_version_insecurity CASCADE;
DROP FUNCTION IF EXISTS public.set_updated_at CASCADE;
DROP TRIGGER IF EXISTS tsvectorupdate ON public.packages;
DROP TRIGGER IF EXISTS security_issues_inserted ON public.security_issues;
DROP TRIGGER IF EXISTS check_reserved_name_trigger ON public.packages;

-------------
CREATE OR REPLACE FUNCTION packages_search_trigger() RETURNS trigger AS $$
begin
  new.TSV :=
    setweight(to_tsvector(coalesce(new.name,'')), 'A') ||
    setweight(to_tsvector(coalesce(new.description,'')), 'B');
  return new;
end
$$ LANGUAGE plpgsql;

CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
ON packages FOR EACH ROW EXECUTE PROCEDURE packages_search_trigger();
---------
CREATE OR REPLACE FUNCTION update_version_insecurity() 
RETURNS TRIGGER AS $$
BEGIN
    UPDATE versions 
    SET insecure = TRUE 
    WHERE id = NEW.version_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER security_issues_inserted
AFTER INSERT ON security_issues
FOR EACH ROW
EXECUTE PROCEDURE update_version_insecurity();

CREATE FUNCTION public.set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF (
        NEW IS DISTINCT FROM OLD AND
        NEW.updated_at IS NOT DISTINCT FROM OLD.updated_at
    ) THEN
        NEW.updated_at = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END
$$;

CREATE OR REPLACE FUNCTION public.verify_package_name() RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS(SELECT 1 FROM reserved_names WHERE name = NEW.name) THEN
        RAISE EXCEPTION 'Package name is reserved.';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_reserved_name_trigger 
BEFORE INSERT OR UPDATE ON packages
FOR EACH ROW EXECUTE FUNCTION verify_package_name();
