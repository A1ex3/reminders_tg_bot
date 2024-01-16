BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "users" (
	"user_id"	INTEGER NOT NULL UNIQUE,
	PRIMARY KEY("user_id")
);
CREATE TABLE IF NOT EXISTS "events" (
	"id"	INTEGER NOT NULL,
	"e_user_id"	INTEGER NOT NULL,
	"event_name"	TEXT NOT NULL,
	"start_time"	INTEGER NOT NULL,
	"notify_for"	INTEGER NOT NULL,
	FOREIGN KEY("e_user_id") REFERENCES "users"("user_id"),
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE UNIQUE INDEX IF NOT EXISTS "user_id_index_unique" ON "users" (
	"user_id"
);
COMMIT;
