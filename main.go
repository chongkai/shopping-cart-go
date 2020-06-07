package main

import (
	"ckjiang/shopping-cart/example"
	"fmt"
	"github.com/cloudstateio/go-support/cloudstate"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"log"
)

func main1() {
	fmt.Println(proto.FileDescriptor("proto/shopping-cart.proto"))
	protoregistry.GlobalFiles.RangeFiles(func(d protoreflect.FileDescriptor) bool {
		fmt.Println(d.Path())
		return true
	})
}

func main() {
	server, err := cloudstate.New(cloudstate.Config{
		ServiceName:    "shopping-cart",
		ServiceVersion: "0.1.0",
	})
	if err != nil {
		log.Fatalf("CloudState.New failed: %v", err)
	}
	err = server.RegisterEventSourcedEntity(
		&cloudstate.EventSourcedEntity{
			ServiceName:   "example.shoppingcart.ShoppingCartService",
			PersistenceID: "ShoppingCart",
			EntityFunc:    example.NewShoppingCart,
		},
		cloudstate.DescriptorConfig{
			Service: "shopping-cart.proto",
		}.AddDomainDescriptor("domain.proto"),
	)
	if err != nil {
		log.Fatalf("CloudState failed to register entity: %v", err)
	}
	err = server.Run()
	if err != nil {
		log.Fatalf("CloudState failed to run: %v", err)
	}
}
