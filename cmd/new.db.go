package cmd

func CreateDBFile() string {
	return `-- Postgresql
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	
	CREATE TABLE "user" (
		id             UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
		name           TEXT           NOT NULL,
		email 				 TEXT 					NOT NULL,
		password 			 TEXT 					NOT NULL,
		created_at     TIMESTAMPTZ    NOT NULL DEFAULT now(),
		created_by     UUID           NOT NULL,
		updated_at     TIMESTAMPTZ,
		updated_by     UUID
	);
	`
}
