# mvc_assignment

A high-performance full-stack game server built in Go, utilizing Postgres for storage, Docker for containerization, and secure JWT authentication.

Local setup 

### 1. Clone the repository & Navigate inside
```
git clone https://github.com/gaurangivishwakarma-ui/mvc_assignment
cd mvc_assignment
```

### 2.Create a .env file in the root directory and paste these settings:
```
DB_USER=root
DB_PASSWORD=secret
DB_HOST=localhost
DB_PORT=5432
DB_NAME=game_db
DB_SSLMODE=disable
SERVER_PORT=8080
JWT_SECRET=yoursupersecuresecretkeychangeitlater
```
### 3. Run setup script 

```
chmod +x setup.sh
./setup.sh
```
### 4. Start the server

`go run cmd/main.go`

### For API Testing 
Once the server is running on port 8080, test end-to-end flow using curl:

#### Register
```
curl -X POST http://localhost:8080/api/register \
     -H "Content-Type: application/json" \
     -d '{"username": "VillageKing", "password": "supersecretpassword"}'
```
#### Login (Get Token)
```
curl -X POST http://localhost:8080/api/login \
     -H "Content-Type: application/json" \
     -d '{"username": "VillageKing", "password": "supersecretpassword"}'
```
