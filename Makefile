GO = go
LDFLAGS =
CFLAGS = -gcflags "-N -l"

PKG_SRCS = $(shell find {pkg,internal} -name '*.go')

PROTOC := protoc
PROTO_SRCS := $(shell find {pkg,internal} -name '*.proto')
PROTO_INCS = -I $$GOPATH/src -I vendor -I .
#-I $$GOPATH/src -I vendor

GO_PROTO_TARGETS := $(PROTO_SRCS:.proto=.pb.go)
GO_PROTO_OPTIONS = --gogo_out=plugins=grpc,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,paths=source_relative:.

DATAGATHER_SRCS := $(shell find cmd/datagather -name '*.go')
DATAGATHER_MAIN := cmd/datagather/main.go
DATAGATHER_TARGET := datagather

all: proto $(DATAGATHER_TARGET)

$(DATAGATHER_TARGET): $(DATAGATHER_SRCS) $(PKG_SRCS)
	$(GO) build $(CFLAGS) $(LDFLAGS) -o $(DATAGATHER_TARGET) $(DATAGATHER_MAIN)

proto: $(GO_PROTO_TARGETS)

$(GO_PROTO_TARGETS): $(PROTO_SRCS)
	for proto_src in $(PROTO_SRCS) ; do \
		$(PROTOC) $(PROTO_INCS) $(GO_PROTO_OPTIONS) $$proto_src ; \
	done

clean:
	$(RM) $(GO_PROTO_TARGETS)
	$(GO) clean
