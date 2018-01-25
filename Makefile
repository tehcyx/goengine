default: builddocker

buildgo:
	go build -i

builddocker:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o goenginedocker
	docker build -t goengine .

run:
	docker run --rm -p 8080:8080 goengine

test:
	go test ./...

cover:
	go test ./.. -cover