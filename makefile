SHELL := /bin/bash

export PROJECT = ardan-starter-kit

# ==============================================================================
# Testing running system

# For testing a simple query on the system. Don't forget to `make seed` first.
# curl --user "admin@example.com:gophers" http://localhost:3000/v1/users/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2

# For testing load on the service.
# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/1/2
# zipkin: http://localhost:9411
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# Used to install expvarmon program for metrics dashboard.
# go install github.com/divan/expvarmon@latest

# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./sales-admin genkey

# ==============================================================================
# Mine

# Get general transaction info from database, inc. earliest and latest transaction IDs
# and number of transactions
get-txn-info-from-db:
	go run ./backend/app/algo-admin/main.go get-txn-info-from-db | go run backend/app/logfmt/main.go

get-blocks-count-from-db:
	go run ./backend/app/algo-admin/main.go get-blocks-count-from-db | go run backend/app/logfmt/main.go

try-get-txns-by-acct-from-db:
	go run backend/app/algo-admin/main.go get-txns-by-acct-from-db 2255PMXS65R54KKH5FQVV5UQZSAQCYL5U3OWQ2E5IZGOLK5XVTAVKNRPPQ | go run backend/app/logfmt/main.go
#	go run backend/app/algo-admin/main.go get-txns-by-acct-from-db 22NA4OMQB46PO5MD22EXW5JAYNWAPQBFBYMM6OJSLAJJ23ZKQA6MPRZKG4 | go run backend/app/logfmt/main.go

try-get-txns-by-acct-pagination-from-db:
	go run backend/app/algo-admin/main.go get-txns-by-acct-pagination-from-db 2255PMXS65R54KKH5FQVV5UQZSAQCYL5U3OWQ2E5IZGOLK5XVTAVKNRPPQ 100 1 asc | go run backend/app/logfmt/main.go

# https://github.com/ThiagoBarradas/woocommerce-docker/issues/2
start-wp:
	docker run --name woocommerce -p80:80 -d thiagobarradas/woocommerce:3.5.3-wp5.0.2-php7.2

# https://developer.algorand.org/docs/run-a-node/setup/install/
stop-algorand:
	sudo systemctl stop algorand

# https://www.cyberciti.biz/faq/star-stop-restart-apache2-webserver/
stop-apache:
	/etc/init.d/apache2 stop

kill-postgres:
	# https://askubuntu.com/questions/547434/how-to-nicely-stop-all-postgres-processes
	sudo pkill -u postgres


get-current-round:
	go run app/algo-admin/main.go get-current-round

# Run the cal-engine with all the defaults, except the private keys
it-rain:
# 	go run app/cal-engine/main.go --auth-keys-folder=./keys-dir
#	go run backend/app/algosearch/main.go --web-enable-sync=true | go run backend/app/logfmt/main.go
	go run backend/app/algosearch/main.go | go run backend/app/logfmt/main.go

restart-couch: stop-couch start-couch

start-couch: run-couch migrate-couch

migrate-couch:
	go run backend/app/algo-admin/main.go migrate

run-couch:
	# sudo rm -rf db-data
	# Create a folder called db-data
	mkdir -p db-data
	echo $(pwd)/db-data
	#docker run -e COUCHDB_USER=admin -e COUCHDB_PASSWORD=password -p 5984:5984 --name my-couchdb -v $(pwd)/db-data:/opt/couchdb/data -d couchdb
	docker run -e COUCHDB_USER=kevin -e COUCHDB_PASSWORD=makechesterproud! -p 5984:5984 --name algosearch-couchdb -v $(shell pwd)/db-data:/opt/couchdb/data -d couchdb
	# https://github.com/apache/couchdb-docker/issues/54
	# curl -X PUT http://127.0.0.1:5984/_users
	# https://guide.couchdb.org/draft/security.html
	# > HOST="http://anna:secret@127.0.0.1:5984"
    # > curl -X PUT $HOST/somedatabase
    # {"ok":true}
    # curl -X PUT 'http://kevin:makechesterproud!@127.0.0.1:5984/_users'

stop-couch:
	docker stop algosearch-couchdb
	docker rm algosearch-couchdb
	sudo rm -rf db-data

fuck-this:
	go run ./backend/app/algo-admin/main.go get-and-insert-blocks 11 10000

# ==============================================================================
# Monitoring

monitor:
	expvarmon -ports="http://localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 1.1

all: algosearch-backend metrics

algosearch-backend:
	docker build \
		-f zarf/docker/dockerfile.algosearch-backend \
		-t algosearch-backend-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

algosearch-metrics:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t metrics-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

algosearch-backend-latest:
	docker build \
		-f zarf/docker/dockerfile.algosearch-backend \
		-t algosearch-backend-amd64:latest \
		--build-arg BUILD_REF=latest \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# build the cal-engine image and push it to AWS ECR
build-cal-engine-for-m1:
	docker buildx build --platform linux/amd64 \
		-f zarf/docker/dockerfile.cal-engine \
		--push -t 938897780349.dkr.ecr.ap-southeast-1.amazonaws.com/cal-engine-amd64:$(VERSION) \
		--build-arg VCS_REF=1.0 \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

metrics:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t metrics-amd64:$(VERSION) \
		--build-arg VCS_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

metrics-latest:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t metrics-amd64:latest \
		--build-arg VCS_REF=latest \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# build the cal-engine image and push it to AWS ECR
build-metrics-for-m1:
	docker buildx build --platform linux/amd64 \
		-f zarf/docker/dockerfile.metrics \
		--push -t 938897780349.dkr.ecr.ap-southeast-1.amazonaws.com/cal-engine-metrics-amd64:$(VERSION) \
		--build-arg VCS_REF=1.1 \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within docker compose

up:
	docker-compose -f zarf/compose/compose.yaml -f zarf/compose/compose-config.yaml up --detach --remove-orphans

down:
	docker-compose -f zarf/compose/compose.yaml down --remove-orphans

logs:
	docker-compose -f zarf/compose/compose.yaml logs -f

# ==============================================================================
# Administration

migrate:
	go run app/sales-admin/main.go migrate

# ==============================================================================
# Running tests within the local computer

test:
	#go test ./... -count=1
	go test -v ./... -count=1
	staticcheck -checks=all ./...

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

docker-down-local:
	docker stop $(FILES)
	docker rm $(FILES)

docker-logs-local:
	docker logs -f $(FILES)

docker-down:
	docker rm -f $(shell docker ps -aq)

docker-clean:
	docker system prune -f

docker-delete-all-containers:
	docker rm -f $(docker ps -a -q)
