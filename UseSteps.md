### migration steps

- install migration package
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

- create migration file
migrate create -ext sql -dir db/migrations/ -seq __FILENAME__

- up migration
go run db/migration.go up                                                                    

- down migration
go run db/migration.go down


PS D:\VS code\WebSocket\EngageSync> migrate create -ext sql -dir migrations -seq adding_kanaka_schema

PS D:\VS code\WebSocket\EngageSync> migrate -database "postgres://postgres:root@localhost:5432/engagesyncdb?sslmode=disable" -path migrations up