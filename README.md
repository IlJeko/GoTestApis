# Test Project with Go language
This is a little Go project to test REST apis

## Project Structure
The project use Postgres as a Database manager, using localhost settings  
It has a Migrations folder to store database migrations  
main.go is the starting point  
database.go and service.go are for database configuration  
jwt.go is for JWT auth configuration and methods  

This project is also using Gin Gonic library for routing and go playground validation for validating data  

## Setup
To start this project, you need to download it and open a terminal where you saved it.  
To connect your Postgres database, in main.go, fill the connection string with your data  
```
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Eros2724@1"
)
```

To start the project just type the following command in the terminal
```
go run .
```

To create a migration file:
```
migrate create -ext sql -dir migrations -seq [your migration name]
```

To apply Up migration:
```
migrate -path migrations -database "postgres://[your username]:[your password]@localhost:5432/?sslmode=disable" up
```

To apply Down migration:
```
migrate -path migrations -database "postgres://[your username]:[your password]@localhost:5432/?sslmode=disable" down
```
