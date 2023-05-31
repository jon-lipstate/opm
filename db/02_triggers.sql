DROP FUNCTION IF EXISTS public.packages_search_trigger CASCADE;
DROP FUNCTION IF EXISTS public.update_version_insecurity CASCADE;
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
