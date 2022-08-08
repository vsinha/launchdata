run:
    go run . cache -s 2020 -e 2022

test:
    go test ./...

test-accept:
    GOLDEN_UPDATE=true go test ./...
