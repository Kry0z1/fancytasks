CREATE TABLE IF NOT EXISTS users(
    username VARCHAR(128) PRIMARY KEY,
	HashedPassword VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS base_tasks(
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    done BOOLEAN NOT NULL,
    owner VARCHAR(128) NOT NULL,

    FOREIGN KEY (owner) REFERENCES users(username)
);

CREATE TABLE IF NOT EXISTS events(
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    done BOOLEAN NOT NULL,
    owner VARCHAR(128) NOT NULL,
    starts_at TIME NOT NULL,
    ends_at TIME NOT NULL,

    FOREIGN KEY (owner) REFERENCES users(username)
);

CREATE TABLE IF NOT EXISTS tasks_with_deadline(
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    done BOOLEAN NOT NULL,
    owner VARCHAR(128) NOT NULL,
    deadline TIME NOT NULL,

    FOREIGN KEY (owner) REFERENCES users(username)
);

CREATE TABLE IF NOT EXISTS repeating_task(
    id SERIAL PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    done BOOLEAN NOT NULL,
    owner VARCHAR(128) NOT NULL,
    starts_at TIME,
    ends_at TIME,
    period INTERVAL NOT NULL,
    loop INTEGER NOT NULL,

    FOREIGN KEY (owner) REFERENCES users(username)
);

CREATE TABLE IF NOT EXISTS excepts (
    task_id INTEGER PRIMARY KEY,
    value INTEGER NOT NULL,

    FOREIGN KEY (task_id) REFERENCES repeating_task(id)
);