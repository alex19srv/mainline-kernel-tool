build:
    go build -o ./target/usr/bin/mainline-kernel ./cmd/mainline_kernel/

test:
    go test ./...

install DESTDIR=".":
    cp -r ./target/. "{{DESTDIR}}/"
