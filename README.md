# Mercado Libre
API that allows CRUD operations on users and items.
Language: Golang
Database: MySQL

## How to run in a local environment
1. Run a local database with docker. Please, follow these steps: 
    1.1 Create volume: **docker volume create dbdata**.
    1.2 Create MySQL container: **docker run -dp 3306:3306 --name mysql-db -e MYSQL_ROOT_PASSWORD=secret --mount src=dbdata,dst=/var/lib/mysql mysql:5.7**

2. Create database: create a local database and name it **mercadolibre**.

3. Create **.env** file: you must create a file in the root directory and name it .env. It have to contain the values ​​for your local database.
Help: you can make a copy of the **.env.template** file and rename.

4. From the command line in the root directory, run: **go mod tidy**

5. From the command line in the root directory, run the **go run** command to build and run the project: **go run ./cmd/api**


## Executing test

From the command line in the root directory, run: **go test -cover ./...**