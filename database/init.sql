-- Create auditable (parent) table for other tables and update function for updated_at column
CREATE TABLE public.auditable
(
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT NULL
);

CREATE
    OR REPLACE FUNCTION update_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE SCHEMA "auth";
CREATE SCHEMA "geo";
CREATE SCHEMA "census";

CREATE TYPE auth.role_type AS ENUM ('agent', 'administrator');
CREATE TYPE census.gender_type AS ENUM ('male', 'female');

-- Create tables users and user_auth in auth schema
CREATE TABLE auth.users
(
    id           uuid PRIMARY KEY        DEFAULT gen_random_uuid(),
    email        text           NOT NULL UNIQUE,
    first_name   text           NOT NULL DEFAULT '',
    last_name    text           NOT NULL DEFAULT '',
    role         auth.role_type NOT NULL DEFAULT 'agent',
    default_user boolean                 DEFAULT false
) INHERITS (public.auditable);

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON auth.users
    FOR EACH ROW
EXECUTE FUNCTION update_at();

CREATE TABLE auth.user_auth
(
    user_id              uuid PRIMARY KEY REFERENCES auth.users (id) ON DELETE CASCADE,
    password             text                     NOT NULL, -- Encrypted password (bcrypt)
    last_login           timestamp with time zone NOT NULL DEFAULT now(),
    last_password_change timestamp with time zone NOT NULL DEFAULT now()
) INHERITS (public.auditable);

CREATE TRIGGER update_user_auth_updated_at
    BEFORE UPDATE
    ON auth.user_auth
    FOR EACH ROW
EXECUTE FUNCTION update_at();


-- Create tables regions, cities, buildings, and addresses in geo schema
CREATE TABLE geo.regions
(
    id   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL UNIQUE
) INHERITS (public.auditable);

CREATE TRIGGER update_regions_updated_at
    BEFORE UPDATE
    ON geo.regions
    FOR EACH ROW
EXECUTE FUNCTION update_at();

CREATE TABLE geo.cities
(
    id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    region_id uuid NOT NULL REFERENCES geo.regions (id) ON DELETE CASCADE,
    name      text NOT NULL UNIQUE
) INHERITS (public.auditable);

CREATE TRIGGER update_cities_updated_at
    BEFORE UPDATE
    ON geo.cities
    FOR EACH ROW
EXECUTE FUNCTION update_at();

CREATE TABLE geo.buildings
(
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    city_id      uuid NOT NULL REFERENCES geo.cities (id) ON DELETE CASCADE,
    street       text NOT NULL,
    house_number text             DEFAULT NULL
) INHERITS (public.auditable);

CREATE TRIGGER update_buildings_updated_at
    BEFORE UPDATE
    ON geo.buildings
    FOR EACH ROW
EXECUTE FUNCTION update_at();

CREATE TABLE geo.addresses
(
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    building_id      uuid NOT NULL REFERENCES geo.buildings (id) ON DELETE CASCADE,
    apartment_number text             DEFAULT NULL
) INHERITS (public.auditable);

CREATE TRIGGER update_addresses_updated_at
    BEFORE UPDATE
    ON geo.addresses
    FOR EACH ROW
EXECUTE FUNCTION update_at();

-- Create table census in census schema
CREATE TABLE census.households
(
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    address_id      uuid    NOT NULL REFERENCES geo.addresses (id) ON DELETE CASCADE,
    enumerator_id   uuid    NOT NULL REFERENCES auth.users (id),

    total_residents integer NOT NULL, -- Total number of residents in the household
    dwelling_type   text,             -- Type of dwelling (e.g., apartment, house)
    building_year   text,             -- Year of building construction
    total_area      integer,          -- Total area of the dwelling
    living_area     integer,          -- Living area of the dwelling
    rooms_count     integer,          -- Number of rooms in the dwelling

    notes           text
) INHERITS (public.auditable);

CREATE TRIGGER update_households_updated_at
    BEFORE UPDATE
    ON census.households
    FOR EACH ROW
EXECUTE FUNCTION update_at();

CREATE TABLE census.persons
(
    id                    uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    household_id          uuid               NOT NULL REFERENCES census.households (id) ON DELETE CASCADE,

    gender                census.gender_type NOT NULL,
    birth_date            date               NOT NULL,
    citizenship           text,
    has_dual_citizenship  boolean,
    nationality           text,
    native_language       text,
    speaks_russian        boolean,
    other_languages       text[], -- Languages spoken by the person

    education_level       text,
    marital_status        text,
    children_count        INTEGER,
    relation_to_household text,

    place_of_birth        text,
    current_residence     text,

    income_sources        text[], -- Array: ["salary", "pension", "business", "other"]
    employment_status     text    -- "employed", "unemployed", "student", "retired", "other"
) INHERITS (public.auditable);

CREATE TRIGGER update_persons_updated_at
    BEFORE UPDATE
    ON census.persons
    FOR EACH ROW
EXECUTE FUNCTION update_at();