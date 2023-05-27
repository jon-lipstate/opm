DROP FUNCTION IF EXISTS public.packages_search_trigger CASCADE;

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
