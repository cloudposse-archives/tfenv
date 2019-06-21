SHELL = /bin/bash

PATH:=$(PATH):$(GOPATH)/bin

-include $(shell curl -sSL -o .build-harness "https://git.io/build-harness"; echo .build-harness)

.PHONY : test
test:
	$(MAKE) -C $(@)

release/tfenv: main.go
	$(MAKE) go/build
