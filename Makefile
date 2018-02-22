default: build

build: test cover
	go build -i -o bin/app

docker:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o bin/appdocker
	docker build -t goengine .

buildwin: test cover
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows CGO_LDFLAGS="-L/usr/local/Cellar/mingw-w64/5.0.3/toolchain-x86_64/x86_64-w64-mingw32/lib -lSDL2" CGO_CFLAGS="-I/usr/local/Cellar/mingw-w64/5.0.3/toolchain-x86_64/x86_64-w64-mingw32/include -D_REENTRANT" go build -i -o bin/app.exe

run:
	docker run --rm -p 8080:8080 goengine

test:
	go test ./...

cover:
	go test ./... -cover

clean:
	rm -rf bin