.PHONY: build run cov test clean

MAIN_TARGET = cmd/main.go
EXECUTABLE = main
COVERAGE = scripts/coverage.sh

all: build test cov run clean

build:
	go build ${MAIN_TARGET}

run:
	./${EXECUTABLE}

cov:
	./${COVERAGE}

test:
	go test ./...

clean:
	rm ${EXECUTABLE}
