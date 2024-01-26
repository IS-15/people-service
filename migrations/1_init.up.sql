CREATE TABLE IF NOT EXISTS people(
    id SERIAL NOT NULL,
    name varchar(255) NOT NULL,
    surname varchar(255) NOT NULL,
    patronymic varchar(255),
    age integer,
    gender varchar(8),
    nationality varchar(8),
    PRIMARY KEY(id)
);
CREATE UNIQUE INDEX people_un ON "people" USING btree ("name");