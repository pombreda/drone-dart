all: build

deps:
	go get github.com/franela/goblin
	go list github.com/drone/drone/... | xargs go get -t -v

build:
	go build

test:
	go test ./...

clean:
	rm -rf drone-dart

lessc:
	lessc ./website/styles/style.less > website/styles/style.css

run:
	go run main.go

deploy:
	scp -i $$DRONE_DART_KEY drone-dart $$DRONE_DART_USER@$$DRONE_DART_HOST:/tmp/drone-dart
	ssh -i $$DRONE_DART_KEY $$DRONE_DART_USER@$$DRONE_DART_HOST -- install -t /usr/local/bin /tmp/drone-dart
	ssh -i $$DRONE_DART_KEY $$DRONE_DART_USER@$$DRONE_DART_HOST -- rm -f /tmp/drone-dart
	ssh -i $$DRONE_DART_KEY $$DRONE_DART_USER@$$DRONE_DART_HOST -- restart drone-dart

