createmigration:
	migrate create -ext sql -dir migrations -seq init_table  

migrate:
	migrate -source file://migrations -database "mysql://root:root@tcp(localhost:3306)/courses" -verbose up

migratedown:
	migrate -source file://migrations -database "mysql://root:root@tcp(localhost:3306)/courses" -verbose down

.PHONY: migrate