//go:generate protoc ./account.proto --go_out=plugins=grpc:./pb
package account

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Varun5711/fritzy/account/pb"
	kafkapkg "github.com/Varun5711/fritzy/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	kafkaProducer *kafkapkg.Producer
	pb.UnimplementedAccountServiceServer
}

func ListenGRPC(s Service, kafkaProducer *kafkapkg.Producer, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterAccountServiceServer(serv, &grpcServer{
		service:                           s,
		kafkaProducer:                     kafkaProducer,
		UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{},
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostAccount(ctx context.Context, r *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	a, err := s.service.PostAccount(ctx, r.Name)
	if err != nil {
		return nil, err
	}

	if s.kafkaProducer != nil {
		event := map[string]interface{}{
			"event_type": "account.created",
			"account_id": a.ID,
			"name":       a.Name,
		}
		if err := s.kafkaProducer.Publish(ctx, "account.events", a.ID, event); err != nil {
			log.Printf("Failed to publish account created event: %v", err)
		}
	}

	return &pb.PostAccountResponse{Account: &pb.Account{
		Id:   a.ID,
		Name: a.Name,
	}}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   a.ID,
			Name: a.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	accounts := []*pb.Account{}
	for _, p := range res {
		accounts = append(
			accounts,
			&pb.Account{
				Id:   p.ID,
				Name: p.Name,
			},
		)
	}
	return &pb.GetAccountsResponse{Accounts: accounts}, nil
}
