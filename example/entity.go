package example

import (
	"ckjiang/shopping-cart/example/domain"
	"ckjiang/shopping-cart/example/shoppingcart"
	"context"
	"fmt"
	"github.com/cloudstateio/go-support/cloudstate"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type ShoppingCart struct {
	// our domain object
	cart []*domain.LineItem
	// as an Emitter we can emit events
	cloudstate.EventEmitter
}

func NewShoppingCart() cloudstate.Entity {
	return &ShoppingCart{
		cart:         make([]*domain.LineItem, 0),
		EventEmitter: cloudstate.NewEmitter(), // TODO: the EventEmitter could be provided by the event sourced handler
	}
}

func (sc *ShoppingCart) HandleCommand(ctx context.Context, command interface{}) (handled bool, reply interface{}, err error) {
	switch cmd := command.(type) {
	case *shoppingcart.GetShoppingCart:
		reply, err := sc.GetCart(ctx, cmd)
		return true, reply, err
	case *shoppingcart.RemoveLineItem:
		reply, err := sc.RemoveItem(ctx, cmd)
		return true, reply, err
	case *shoppingcart.AddLineItem:
		reply, err := sc.AddItem(ctx, cmd)
		return true, reply, err
	default:
		return false, reply, err
	}
}

func (sc *ShoppingCart) GetCart(_ context.Context, _ *shoppingcart.GetShoppingCart) (*shoppingcart.Cart, error) {
	cart := &shoppingcart.Cart{}
	for _, item := range sc.cart {
		cart.Items = append(cart.Items, &shoppingcart.LineItem{
			ProductId: item.ProductId,
			Name:      item.Name,
			Quantity:  item.Quantity,
		})
	}
	return cart, nil
}

func (sc *ShoppingCart) AddItem(_ context.Context, li *shoppingcart.AddLineItem) (*empty.Empty, error) {
	if li.GetQuantity() <= 0 {
		return nil, fmt.Errorf("cannot add negative quantity of to item %s", li.GetProductId())
	}
	sc.Emit(&domain.ItemAdded{
		Item: &domain.LineItem{
			ProductId: li.ProductId,
			Name:      li.Name,
			Quantity:  li.Quantity,
		}})
	return &empty.Empty{}, nil
}

func (sc *ShoppingCart) HandleEvent(_ context.Context, event interface{}) (handled bool, err error) {
	switch e := event.(type) {
	case *domain.ItemAdded:
		return true, sc.ItemAdded(e)
	case *domain.ItemRemoved:
		return true, sc.ItemRemoved(e)
	default:
		return false, nil
	}
}

func (sc *ShoppingCart) ItemAdded(added *domain.ItemAdded) error { // TODO: enable handling for values
	if item, _ := sc.find(added.Item.ProductId); item != nil {
		item.Quantity += added.Item.Quantity
	} else {
		sc.cart = append(sc.cart, &domain.LineItem{
			ProductId: added.Item.ProductId,
			Name:      added.Item.Name,
			Quantity:  added.Item.Quantity,
		})
	}
	return nil
}

func (sc *ShoppingCart) find(productId string) (*domain.LineItem, int) {
	for i, item := range sc.cart {
		if item.ProductId == productId {
			return item, i
		}
	}
	return nil, -1
}

func (sc *ShoppingCart) Snapshot() (snapshot interface{}, err error) {
	return domain.Cart{
		Items: append(make([]*domain.LineItem, len(sc.cart)), sc.cart...),
	}, nil
}

func (sc *ShoppingCart) HandleSnapshot(snapshot interface{}) (handled bool, err error) {
	switch value := snapshot.(type) {
	case domain.Cart:
		sc.cart = append(sc.cart[:0], value.Items...)
		return true, nil
	default:
		return false, nil
	}
}

func (sc *ShoppingCart) ItemRemoved(e *domain.ItemRemoved) (err error) {
	if _, i := sc.find(e.ProductId); i >= 0 {
		sc.cart = append(sc.cart[:i], sc.cart[i+1:]...)
	}
	return nil
}

func (sc *ShoppingCart) RemoveItem(_ context.Context, cmd *shoppingcart.RemoveLineItem) (reply *empty.Empty, err error) {
	sc.Emit(&domain.ItemRemoved{ProductId: cmd.ProductId})
	return &empty.Empty{}, nil
}
