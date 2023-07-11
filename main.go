package main

import (
	"context"
	"log"
	"net"

	"github.com/psonhoang/go-grpc/invoicer"
	"google.golang.org/grpc"
)

type myInvoicerServer struct {
	invoicer.UnimplementedInvoicerServer
}

func (s myInvoicerServer) Create(ctx context.Context, req *invoicer.CreateRequest) (*invoicer.CreateResponse, error) {
	taxAmount := invoicer.Amount{
		Amount:   req.Amount.GetAmount() * 11 / 100,
		Currency: req.Amount.GetCurrency(),
	}

	tipAmount := invoicer.Amount{
		Amount:   taxAmount.GetAmount() * 18 / 100,
		Currency: taxAmount.GetCurrency(),
	}

	return &invoicer.CreateResponse{
		From: req.GetFrom(),
		To:   req.GetTo(),
		Tax:  &taxAmount,
		Tip:  &tipAmount,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}

	serverRegistrar := grpc.NewServer()
	service := &myInvoicerServer{}

	invoicer.RegisterInvoicerServer(serverRegistrar, service)

	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to server: %s", err)
	}
}
