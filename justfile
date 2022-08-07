run:
    go run . -s 2020 -e 2022

test:
    go test ./...

test-accept:
    GOLDEN_UPDATE=true go test ./...
