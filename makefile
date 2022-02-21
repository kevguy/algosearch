SHELL := /bin/bash

export PROJECT = algosearch

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

# Access metrics directly (4000) or through the sidecar (3001)
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# expvarmon -ports=":3001" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# Used to install expvarmon program for metrics dashboard.
# go install github.com/divan/expvarmon@latest

# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./sales-admin genkey

# ==============================================================================
# Local Algod Utilities

start-node:
	sudo systemctl start algorand

stop-node:
	sudo systemctl stop algorand

node-status:
	goal node status -d /var/lib/algorand

# ==============================================================================
# Other Utilities


kill-postgres:
	# https://askubuntu.com/questions/547434/how-to-nicely-stop-all-postgres-processes
	sudo pkill -u postgres

# https://www.cyberciti.biz/faq/star-stop-restart-apache2-webserver/
stop-apache:
	/etc/init.d/apache2 stop

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

get-current-round:
	go run app/algo-admin/main.go get-current-round

# ==============================================================================
# Help

algosearch-backend-help:
	go run ./backend/app/algosearch/main.go --help

algosearch-metrics-help:
	go run ./backend/app/sidecar/metrics/main.go --help

# ==============================================================================
# Local

start-algosearch-backend:
	go run backend/app/algosearch/main.go \
		--web-enable-sync=true \
		| go run backend/app/logfmt/main.go

start-algosearch-metrics:
	go run ./backend/app/sidecar/metrics/main.go

# ==============================================================================
# Sandbox

start-sandbox-algosearch-backend:
	go run backend/app/algosearch/main.go \
		--web-enable-sync=true \
	 	--algorand-indexer-protocol=http \
	  	--algorand-indexer-addr=localhost:8980 \
		| go run backend/app/logfmt/main.go

start-sandbox-algosearch-metrics:
	go run ./backend/app/sidecar/metrics/main.go

migrate-couch-sandbox:
	go run backend/app/algo-admin/main.go migrate

migrate-couch-sandbox-2:
	# curl -X PUT "http://kevin:makechesterproud!@$89.39.110.254:5984/algo_beta?partitioned=false"
	go run backend/app/algo-admin/main.go \
		--couch-db-protocol=http \
		--couch-db-user=kevin \
		--couch-db-password=makechesterproud! \
		--couch-db-host=89.39.110.254:5984 \
		--couch-db-name=algo_beta \
		migrate

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

# ==============================================================================
# Monitoring

monitor:
	expvarmon -ports="http://localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 1.1

all: algosearch-backend algosearch-metrics algosearch-frontend

deploy-to-docker-hub: algosearch algosearch-latest algosearch-backend algosearch-backend-latest algosearch-frontend algosearch-frontend-latest algosearch-metrics algosearch-metrics-latest
	docker tag algosearch:latest kevguy/algosearch:latest
	docker tag algosearch:$(VERSION) kevguy/algosearch:$(VERSION)
	docker tag algosearch-backend:latest kevguy/algosearch-backend:latest
	docker tag algosearch-backend:$(VERSION) kevguy/algosearch-backend:$(VERSION)
	docker tag algosearch-metrics:latest kevguy/algosearch-metrics:latest
	docker tag algosearch-metrics:$(VERSION) kevguy/algosearch-metrics:$(VERSION)
	docker tag algosearch-frontend:latest kevguy/algosearch-frontend:latest
	docker tag algosearch-frontend:$(VERSION) kevguy/algosearch-frontend:$(VERSION)
	docker push kevguy/algosearch:latest
	docker push kevguy/algosearch:$(VERSION)
	docker push kevguy/algosearch-backend:latest
	docker push kevguy/algosearch-backend:$(VERSION)
	docker push kevguy/algosearch-metrics:latest
	docker push kevguy/algosearch-metrics:$(VERSION)
	docker push kevguy/algosearch-frontend:latest
	docker push kevguy/algosearch-frontend:$(VERSION)

algosearch:
	docker build \
		-f zarf/docker/dockerfile.all \
		-t algosearch:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

algosearch-latest:
	docker build \
		-f zarf/docker/dockerfile.all \
		-t algosearch:latest \
		--build-arg BUILD_REF=latest \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

algosearch-backend:
	docker build \
		-f zarf/docker/dockerfile.algosearch-backend \
		-t algosearch-backend:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

algosearch-backend-latest:
	docker build \
		-f zarf/docker/dockerfile.algosearch-backend \
		-t algosearch-backend:latest \
		--build-arg BUILD_REF=latest \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

algosearch-metrics:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t algosearch-metrics:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

algosearch-metrics-latest:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t algosearch-metrics:latest \
		--build-arg BUILD_REF=latest \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

algosearch-frontend:
	docker build \
		-f zarf/docker/dockerfile.algosearch-frontend \
		-t algosearch-frontend:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

algosearch-frontend-latest:
	docker build \
		-f zarf/docker/dockerfile.algosearch-frontend \
		-t algosearch-frontend:latest \
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

deploy-backend-to-heroku:
	# https://devcenter.heroku.com/articles/container-registry-and-runtime#dockerfile-commands-and-runtime
	# rmb to run `heroku login` and `heroku container:login` first
	docker tag algosearch-backend:latest kevguy/algosearch-backend:latest
	docker push kevguy/algosearch-backend:latest
	docker tag kevguy/algosearch-backend:latest registry.heroku.com/algosrch/web
	docker push registry.heroku.com/algosrch/web
	heroku container:release web --app algosrch

# ==============================================================================
# Running tests within the local computer

lint:
	staticcheck -checks=all ./backend/...

test:
	#go test ./... -count=1
	go test -v ./backend/... -count=1
	staticcheck -checks=all ./backend/...

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

# Stop and remove all containers (not only AlgoSearch)
docker-down-local:
	docker stop $(FILES)
	docker rm $(FILES)

# See logging of all containers
docker-logs-local:
	docker logs -f $(FILES)

# Remove all containers
docker-down:
	docker rm -f $(shell docker ps -aq)

# Clean and remove all docker images
docker-clean:
	docker system prune -f

# Remove all containers
docker-delete-all-containers:
	docker rm -f $(docker ps -a -q)
