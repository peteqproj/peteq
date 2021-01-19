#/bin/sh

curl -X POST -H "Content-Type: application/json" --data "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" http://localhost:8080/c/user/register