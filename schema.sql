-- CREATE DATABASE dogstore;
-- \c dogstore;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- automated
CREATE TABLE breeds (
  breed_id INT GENERATED ALWAYS AS IDENTITY,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,
  PRIMARY KEY(breed_id)
);

CREATE FUNCTION breed_id_by_name(breed_name TEXT)
  RETURNS INT
  LANGUAGE PLPGSQL
AS
$$
DECLARE
  result_breed_id INT;
BEGIN
  SELECT breed_id
  INTO result_breed_id
  FROM breeds
  WHERE name like breed_name;

  RETURN result_breed_id;
END
$$;

INSERT INTO breeds(name, description)
VALUES
  ('German Shepherd', 'The German Shepherd or Alsatian is a German breed of working dog of medium to large size'),
  ('Bulldog', 'The Bulldog is a British breed of dog of mastiff type'),
  ('Labrador Retriever', 'The Labrador Retriever or simply Labrador is a British breed of retriever gun dog'),
  ('Golden Retriever', 'The Golden Retriever is a Scottish breed of retriever dog of medium size');

CREATE TABLE owners (
  owner_id UUID DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  surname TEXT NOT NULL,
  birthday DATE NOT NULL,
  PRIMARY KEY(owner_id)
);

CREATE FUNCTION owner_id_by_name(owner_name TEXT)
  RETURNS UUID
  LANGUAGE PLPGSQL
AS
$$
DECLARE
  result_owner_id UUID;
BEGIN
  SELECT owner_id
  INTO result_owner_id
  FROM owners
  WHERE name like owner_name
  LIMIT 1;

  RETURN result_owner_id;
END
$$;

INSERT INTO owners(name, surname, birthday)
VALUES
  ('Pedro', 'Gonzalez Gonzalez', '1996-11-02'),
  ('Ivan', 'Martinez Alberte', '1999-05-15'),
  ('Juan', 'Perez Gonzalez', '1990-02-12');

CREATE TABLE dogs (
  dog_id INT GENERATED ALWAYS AS IDENTITY,
  name TEXT NOT NULL,
  breed_id INT NOT NULL REFERENCES breeds(breed_id) ON UPDATE CASCADE ON DELETE CASCADE,
  owner_id UUID NOT NULL REFERENCES owners(owner_id) ON UPDATE CASCADE ON DELETE CASCADE,
  PRIMARY KEY(dog_id),
  UNIQUE(name, owner_id)
);

INSERT INTO dogs (name, breed_id, owner_id)
VALUES
  ('Max'   ,( SELECT breed_id_by_name('German Shepherd'    )), ( SELECT owner_id_by_name('Pedro' ))),
  ('Kobe'  ,( SELECT breed_id_by_name('Bulldog'            )), ( SELECT owner_id_by_name('Juan'  ))),
  ('Oscar' ,( SELECT breed_id_by_name('Labrador Retriever' )), ( SELECT owner_id_by_name('Ivan'  ))),
  ('Teddy' ,( SELECT breed_id_by_name('Golden Retriever'   )), ( SELECT owner_id_by_name('Ivan'  ))),
  ('Bailey',( SELECT breed_id_by_name('Bulldog'            )), ( SELECT owner_id_by_name('Juan'  ))),
  ('Chip'  ,( SELECT breed_id_by_name('German Shepherd'    )), ( SELECT owner_id_by_name('Pedro' )));

-- Some queries here
SELECT p.name, b.name, o.name, o.surname, o.birthday
FROM dogs AS p INNER JOIN owners AS o USING(owner_id)
  INNER JOIN breeds AS b USING(breed_id);

SELECT STRING_AGG(p.name || ' ' || b.name, ','), o.name
FROM dogs AS p INNER JOIN owners AS o USING(owner_id)
  INNER JOIN breeds AS b USING(breed_id)
GROUP BY o.name;

-- select name of the breed grouping dog name in an array
SELECT b.name, json_agg(p.name)
FROM breeds AS b INNER JOIN dogs AS p USING(breed_id)
GROUP BY b.name;
