SHELL := /bin/bash

all: install

install:
	GOPATH=${PWD} go install play-log

clean:
	rm -rf ${PWD}/bin/ ${PWD}/pkg/
