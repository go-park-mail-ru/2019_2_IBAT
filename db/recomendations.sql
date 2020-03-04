DROP TRIGGER expire_recommendations_trigger ON recommendations;
DROP FUNCTION expire_recommendations;
DROP TABLE recommendations;

CREATE TABLE recommendations (
    timestamp timestamp NOT NULL DEFAULT NOW(),
    person_id uuid NOT NULL,
    tag_id uuid NOT NULL,
    -- name TEXT NOT NULL
    PRIMARY KEY(person_id, tag_id)
);

CREATE UNIQUE INDEX recom_idx ON recommendations (timestamp);

CREATE FUNCTION expire_recommendations() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  DELETE FROM recommendations WHERE timestamp < NOW() - INTERVAL '7 days';--'1 minute';
  RETURN NEW;
END;
$$;

CREATE TRIGGER expire_recommendations_trigger
    AFTER INSERT ON recommendations
    EXECUTE PROCEDURE expire_recommendations();