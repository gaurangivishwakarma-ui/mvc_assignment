# mvc_assignment

A high-performance full-stack game server built in Go, utilizing Postgres for storage, Docker for containerization, and secure JWT authentication.

Local setup 

Please ensure you have docker and go installed.

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
### 3. Run docker compose

```
docker compose up --build -d
```
to stop everything:
```
docker compose down -v
```

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
####  Dashboard 
```
curl -X GET http://localhost:8080/api/player/dashboard \
     -H "Authorization: Bearer <token>"
```
#### Buy Buildings
run this after seeding the database by running, `make seed`
```
curl -X POST http://localhost:8080/api/village/buildings \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"building_id": 101, "x_coords": 10, "y_coords": 10}'
```
#### Collect resources
```
curl -X POST http://localhost:8080/api/village/collect      -H "Authorization: Bearer <token>"
```
#### Upgrade building
```
curl -X POST http://localhost:8080/api/village/buildings/upgrade \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json" \
     -d '{"placement_id": "placementid"}'
```
#### View army catalog
```
curl -X GET http://localhost:8080/api/army/catalog \
     -H "Authorization: Bearer <token>"
```
