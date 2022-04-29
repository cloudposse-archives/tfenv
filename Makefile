SHELL = /bin/bash

PATH:=$(PATH):$(GOPATH)/bin

-include $(shell curl -sSL -o .build-harness "https://cloudposse.tools/build-harness"; echo .build-harness)

.PHONY : test
test:
	$(MAKE) -C $(@)

release/tfenv: main.go
	$(MAKE) go/build
