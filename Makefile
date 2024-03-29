#!/usr/bin/env make

GO=go
PACKAGE_PREFIX=github.com/fcvarela/konig
BIN=konig
CLIENT_BIN=konig-client
PROTOS=rpc/graph.proto

all: $(BIN) $(CLIENT_BIN)

.PHONY: clean clean-protos protos lint vet

vet:
	$(GO) vet $(PACKAGE_PREFIX)/cmd/konig
	$(GO) vet $(PACKAGE_PREFIX)/cmd/client
	$(GO) vet $(PACKAGE_PREFIX)/rpc
	$(GO) vet $(PACKAGE_PREFIX)

lint:
	golint $(PACKAGE_PREFIX)/cmd/konig
	golint $(PACKAGE_PREFIX)/cmd/client
	golint $(PACKAGE_PREFIX)/rpc
	golint $(PACKAGE_PREFIX)

test:
	$(GO) test $(PACKAGE_PREFIX)/rpc $(PACKAGE_PREFIX)/graph

protos: $(PROTOS)
	$(GO) generate $(PACKAGE_PREFIX)/rpc

$(BIN): protos lint vet
	$(GO) build -o $(BIN) cmd/konig/main.go

$(CLIENT_BIN): protos lint vet
	$(GO) build -o $(CLIENT_BIN) cmd/client/main.go


clean-protos:
	find ./ -name \*.pb.go -exec rm -v {} +

clean: clean-protos
	rm -fr $(BIN) $(CLIENT_BIN)
