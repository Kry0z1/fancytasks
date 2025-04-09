CREATE TABLE IF NOT EXISTS users(
    username VARCHAR(128) PRIMARY KEY,
	hashed_password VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS base_tasks(
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    description VARCHAR(512) NOT NULL,
    done BOOLEAN NOT NULL,
    owner VARCHAR(128) NOT NULL,

    FOREIGN KEY (owner) REFERENCES users(username)
);

CREATE TABLE IF NOT EXISTS events(
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    description VARCHAR(512) NOT NULL,
    done BOOLEAN NOT NULL,
    owner VARCHAR(128) NOT NULL,
    starts_at TIME NOT NULL,
    ends_at TIME NOT NULL,

    FOREIGN KEY (owner) REFERENCES users(username)
);

CREATE TABLE IF NOT EXISTS tasks_with_deadline(
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    description VARCHAR(512) NOT NULL,
    done BOOLEAN NOT NULL,
    owner VARCHAR(128) NOT NULL,
    deadline TIME NOT NULL,

    FOREIGN KEY (owner) REFERENCES users(username)
);

CREATE TABLE IF NOT EXISTS repeating_tasks(
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    description VARCHAR(512) NOT NULL,
    done BOOLEAN NOT NULL,
    owner VARCHAR(128) NOT NULL,
    starts_at TIME,
    ends_at TIME,
    period INTERVAL NOT NULL,
    loop INTEGER NOT NULL,
    excepts INTEGER[] NOT NULL,

    FOREIGN KEY (owner) REFERENCES users(username)
);