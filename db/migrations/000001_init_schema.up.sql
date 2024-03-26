CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar(255) NOT NULL,
  "last_name" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "password" varchar(255) NOT NULL,
  "access_level" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "reservations" (
  "id" bigserial PRIMARY KEY,
  "code" varchar(255) NOT NULL,
  "first_name" varchar(255) NOT NULL,
  "last_name" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "phone" varchar(255),
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  "room_id" bigint NOT NULL,
  "notes" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "rooms" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "room_restrictions" (
  "id" bigserial PRIMARY KEY,
  "start_date" date NOT NULL,
  "end_date" date NOT NULL,
  "room_id" bigint NOT NULL,
  "reservation_id" bigint,
  "restrictions_id" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "restrictions" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX ON "users" ("email");

CREATE UNIQUE INDEX ON "reservations" ("code", "last_name");

CREATE INDEX ON "reservations" ("start_date");

CREATE INDEX ON "reservations" ("end_date");

CREATE INDEX ON "reservations" ("room_id");

CREATE INDEX ON "room_restrictions" ("start_date", "end_date");

CREATE INDEX ON "room_restrictions" ("room_id");

CREATE INDEX ON "room_restrictions" ("reservation_id");

ALTER TABLE "reservations" ADD CONSTRAINT "fk_reservations_room_id" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "room_restrictions" ADD CONSTRAINT "fk_room_restrictions_room_id" FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "room_restrictions" ADD CONSTRAINT "fk_room_restrictions_reservation_id" FOREIGN KEY ("reservation_id") REFERENCES "reservations" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "room_restrictions" ADD CONSTRAINT "fk_room_restrictions_restrictions_id" FOREIGN KEY ("restrictions_id") REFERENCES "restrictions" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
