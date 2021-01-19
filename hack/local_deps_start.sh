#/bin/sh

# sleep 5
usql postgres://amazinguser:perfectpassword@10.152.183.234:5432/awesomedb\?sslmode=disable -f ./hack/init_tables
