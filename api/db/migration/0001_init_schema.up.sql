CREATE TABLE IF NOT EXISTS "users"
(
    id       SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "todo"
(
    id       SERIAL PRIMARY KEY,
    body     TEXT                                            NOT NULL,
    end_date DATE                                            NOT NULL,
    user_id  INT REFERENCES "users" ("id") ON DELETE CASCADE NOT NULL,
    done     BIT -- 1 True, 0 False
);
