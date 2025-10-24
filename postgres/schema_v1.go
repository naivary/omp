package postgres

func provisionStatementV1() *provisionStatement {
	return &provisionStatement{
		version: 1,
		sql: `
		CREATE OR REPLACE FUNCTION pseudo_encrypt(value bigint) returns bigint AS $$
		DECLARE
		l1 int;
		l2 int;
		r1 int;
		r2 int;
		i int:=0;
		BEGIN
		 l1:= (value >> 16) & 65535;
		 r1:= value & 65535;
		 WHILE i < 3 LOOP
		   l2 := r1;
		   r2 := l1 # ((((1366 * r1 + 150889) % 714025) / 714025.0) * 32767)::int;
		   l1 := l2;
		   r1 := r2;
		   i := i + 1;
		 END LOOP;
		 return ((r1 << 16) + l1);
		END;
		$$ LANGUAGE plpgsql strict immutable;

		CREATE SEQUENCE IF NOT EXISTS id_seq START 1;

		CREATE TABLE IF NOT EXISTS club(
			id bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('id_seq')),
			name text NOT NULL,
			location text NOT NULL
		);

		CREATE TABLE IF NOT EXISTS team(
			id bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('id_seq')),
			name text NOT NULL UNIQUE,
			league text NOT NULL,
			-- foreign keys
			club_id bigint REFERENCES club(id)
		);

		CREATE TYPE foot AS ENUM ('right', 'left', 'both');
		CREATE TABLE IF NOT EXISTS player_profile(
			id bigint PRIMARY KEY DEFAULT pseudo_encrypt(nextval('id_seq')),
			email text NOT NULL UNIQUE,
			first_name text NOT NULL,
			last_name text NOT NULL,
			jersey_number int,
			position text NOT NULL,
			strong_foot foot NOT NULL,
			-- foreign keys
			team_id bigint REFERENCES team(id)
		);

		CREATE TABLE omp_metadata(
			key text PRIMARY KEY,
			value text NOT NULL
		);
		`,
	}
}
