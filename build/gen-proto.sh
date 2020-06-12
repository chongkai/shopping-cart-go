#!/bin/bash

PROJ_DIR=$(dirname "${BASH_SOURCE[0]}")/..
GO_SUPPORT_PATH=$(go list -m -f {{.Dir}} github.com/cloudstateio/go-support)

cd $PROJ_DIR
protoc --go_out=. --proto_path=$GO_SUPPORT_PATH/protobuf/frontend:proto proto/shopping-cart.proto
protoc --go_out=. --proto_path=$GO_SUPPORT_PATH/protobuf/frontend:proto proto/domain/domain.proto