build:
     go vet ./...
     staticcheck ./...
     go build .

alias run := smoke
smoke:
    go run . cache -s 2020 -e 2022

test:
    go test ./...

cache-all:
    go run . cache all --output-dir ./data

test-accept:
    GOLDEN_UPDATE=true go test ./...
