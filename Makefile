SHELL=/usr/bin/env bash

.PHONY: clean
clean:
	rm ic-auth

.PHONY: all
all:
	go build -o ic-auth *.go