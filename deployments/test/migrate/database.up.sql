-- CREATE DATABASE photoflux;

CREATE TABLE IF NOT EXISTS users
(
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    email text UNIQUE NOT NULL, -- ????? unique not working
    password text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS categories
(
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS photos
(
    id uuid NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    is_uploaded boolean,
    name text,
    user_id uuid,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT "FK_users" FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS photo_categories
(
    photo_id uuid NOT NULL,
    category_id uuid NOT NULL,
    CONSTRAINT "PK_photo_categories" PRIMARY KEY (photo_id, category_id),
    CONSTRAINT "FK_category" FOREIGN KEY (category_id)
        REFERENCES categories (id) MATCH SIMPLE
        ON DELETE CASCADE,
    CONSTRAINT "FK_photo" FOREIGN KEY (photo_id)
        REFERENCES photos (id) MATCH SIMPLE
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS stars
(
    user_id uuid NOT NULL,
    photo_id uuid NOT NULL,
    CONSTRAINT "PK_star" PRIMARY KEY (user_id, photo_id),
    CONSTRAINT "FK_photo" FOREIGN KEY (photo_id)
        REFERENCES photos (id) MATCH SIMPLE
        ON DELETE CASCADE,
    CONSTRAINT "FK_user" FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refresh_tokens
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    token_id uuid NOT NULL,
    user_id uuid NOT NULL,
    CONSTRAINT "FK_user" FOREIGN KEY (user_id)
        REFERENCES users (id) MATCH SIMPLE
        ON DELETE CASCADE
);

INSERT INTO categories (name, created_at, updated_at)
VALUES
    (
        'sunset',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'day',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'sunrise',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'night',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'portrait',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'landscape',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'animals',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'nature',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'black&white',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'lights',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    ),
    (
        'moody',
		'2023-05-24 19:10:25-07',
        '2023-05-24 19:10:25-07'
    );
