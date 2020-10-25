#/bin/sh
echo "Stopping previous containers"
docker rm -f $(docker ps -f "label=peteq-dev" -aq)

docker run -l peteq-dev -d -p 15672:15672 -p 5672:5672 rabbitmq:3-management
docker run -l peteq-dev -d -p 5432:5432 -e POSTGRES_DB=awesomedb -e POSTGRES_USER=amazinguser -e POSTGRES_PASSWORD=perfectpassword postgres:12
sleep 5
usql postgres://amazinguser:perfectpassword@localhost/awesomedb\?sslmode=disable -f ./hack/init_tables

pgweb --skip-open --url postgres://amazinguser:perfectpassword@localhost/awesomedb\?sslmode=disable