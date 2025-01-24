# Backend
learn back-end 

----------------------------------------------------------------
# installation
https://go.dev/dl/
-----chek install----
go mod tidy


----------------------------------------------------------------
# run app
go run main.go // [namefile.go]

----------------------------------------------------------------
# Postman Test
-rest_api-
curl -X GET http://localhost:8080/bentuk
curl -X DELETE http://localhost:8080/bentuk/bdfa5dd4-8b93-48cf-87aa-65a5a5e558a4

-api_db- 
create database "rest_api"
curl -X POST http://localhost:8080/bentuk \
-H "Content-Type: application/json" \
-d '{"jenis": "lingkaran", "parameter": 5}'

curl -X GET http://localhost:8080/bentuk
curl -X DELETE http://localhost:8080/bentuk/1

----------------------------------------------------------------
# Output: ApI 
![alt text]({549954DE-51B0-42AF-AA44-A7D5E2C8930B}.png)
![alt text]({9C699A3D-EDCC-4A7D-8446-0A9E69648C52}.png)

