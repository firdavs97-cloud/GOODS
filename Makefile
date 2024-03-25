all: run-queries

run-queries:
	#sudo apt install postgresql-client-common
	#sudo apt install postgresql-client
	psql -h localhost -U myuser -W -d mydatabase -f sql/migrate.sql

