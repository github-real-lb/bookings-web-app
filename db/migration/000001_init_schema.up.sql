CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar(255) NOT NULL,
  "last_name" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "password" varchar(255) NOT NULL,
  "access_level" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "reservations" (
  "id" bigserial PRIMARY KEY,
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
  "room" varchar(255),
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
  "name" varchar(255),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("first_name");

CREATE INDEX ON "users" ("last_name");

CREATE INDEX ON "users" ("first_name", "last_name");

CREATE INDEX ON "reservations" ("first_name");

CREATE INDEX ON "reservations" ("last_name");

CREATE INDEX ON "reservations" ("first_name", "last_name");

CREATE INDEX ON "reservations" ("start_date");

CREATE INDEX ON "reservations" ("end_date");

CREATE INDEX ON "reservations" ("room_id");

CREATE INDEX ON "rooms" ("room");

CREATE INDEX ON "room_restrictions" ("start_date");

CREATE INDEX ON "room_restrictions" ("end_date");

CREATE INDEX ON "room_restrictions" ("room_id");

CREATE INDEX ON "room_restrictions" ("reservation_id");

CREATE INDEX ON "restrictions" ("name");

ALTER TABLE "reservations" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "room_restrictions" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "room_restrictions" ADD FOREIGN KEY ("reservation_id") REFERENCES "reservations" ("id");

ALTER TABLE "room_restrictions" ADD FOREIGN KEY ("restrictions_id") REFERENCES "restrictions" ("id");
