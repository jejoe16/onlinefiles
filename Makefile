postgres: 
	docker run --name digital_id_container -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=2XdncTB3NUDSA9dmd92f8nSerLSQMQE9GN -d postgres

createdatabase: 
	docker exec -it digital_id_container createdb --username=postgres --owner=postgres did_data

dropdatabase: 
	docker exec -it digital_id_container dropdb did_data --username=postgres

init:
	docker cp migration/db/init.sql digital_id_container:/init.sql
	docker exec -u postgres digital_id_container psql did_data postgres -f init.sql

.PHONY: createdatabase postgres createdb dropdb init