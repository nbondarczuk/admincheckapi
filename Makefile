include Makefile.defs

all: build

run:
	go run cmd/$(TARGET).go

build: #tidy fmt vet
	go build -ldflags=$(LDFLAGS) -o $(TARGET) cmd/$(TARGET)/$(TARGET).go

test:
	go test -count=1 ./... 

vtest:
	go test -v -count=1 ./... 

envtest:
	CONFIG=dev-config.yaml MSAD=1 POSTGRES=1 MYSQL=1 go test -count=1 ./... 

venvtest:
	CONFIG=dev-config.yaml MSAD=1 POSTGRES=1 MYSQL=1 go test -v -count=1 -v ./...

integrtest:
	(cd test/integr; ./run-all.sh)

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint ./...

admintokencheck:
	go build -o admintokencheck cmd/admintokencheck/admintokencheck.go

psql:
	psql --port=5432 --host=localhost --username=test --dbname=argonadmindb

images: imageapi imagepostgres imageswagger

imagesrun:imagepostgresrun imageapirun 

imageskill: imageapikill imagepostgreskill imageswaggerkill

imageapi: build
	docker build . -t $(TARGET):local

imageapirun:
	docker run -d -it --name $(TARGET)-local --network=host $(TARGET):local

imageapikill:
	docker stop $(docker ps -a -q --filter "name=admincheckapi-local" --format="{{.ID}}")

imagepostgres:
	cd db/postgres; docker build . -t postgres:local

imagepostgresrun:
	docker run -d -it --name postgres-local --network=host -e POSTGRES_PASSWORD=test postgres:local

imagepostgreskill:
	docker rm $(docker stop $(docker ps -a -q --filter "name=postgres-local" --format="{{.ID}}"))

imageswagger:
	cd doc/swager; docker build . -t swagger:local

imageswaggerrun:
	docker run --rm -d -p 8888:8080 -it --network host --name swagger-local swagger:local
	firefox http://localhost:8888
	docker rm $(docker stop $(docker ps -a -q --filter "name=swagger-local" --format="{{.ID}}"))

imageswaggerkill:
	docker rm $(docker stop $(docker ps -a -q --filter "name=swagger-local" --format="{{.ID}}"))

imageheartbeat:
	cd test/heartbeat; docker build . -t heartbeat:local

imageheartbeatrun:
	docker run --rm -d -it --network host --name heartbeat-local heartbeat:local

imageheartbeatkill:
	docker rm $(docker stop $(docker ps -a -q --filter "name=heartbeat-local" --format="{{.ID}}"))

imageclean:
	docker image prune -f -a

tar: clean
	tar -cvf /tmp/$(TARGET).tar *
	gzip -f /tmp/$(TARGET).tar

clean:
	go clean
	rm -f $(TARGET) admintokencheck
	find . -name "*~" -exec rm -f {} \;
	find . -name cache.json -exec rm -f {} \;

help:
	@echo 'Some popular management commands for admincheckapi'
	@echo
	@echo 'Usage:'
	@echo '    make all         Makes build target'
	@echo '    make build       Builds the executable'
	@echo '    make test        Start unit test harness'
	@echo '    make imageapi    Builds the docker image'
	@echo '    make imageapirun Executes the docker image'
	@echo '    make clean    Cleans the directory (and vendor directory).'
	@echo

.PHONY: help all build test vtest tidy fmt vet image imagerun imageclean clean tar setupdb cleandb sql doc swagger-ui swaer-editor
