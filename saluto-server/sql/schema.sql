CREATE TABLE post (
  id BIGSERIAL PRIMARY KEY,
  displayName text NOT NULL,
  body text NOT NULL,
  timePosted timestamp NOT NULL
);