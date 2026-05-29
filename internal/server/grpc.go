package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/mandacode-labs/dodream/pkg/proto/dodream/engine/v1"
)

// EngineGRPCServer implements the gRPC EngineService.
type EngineGRPCServer struct {
	pb.UnimplementedEngineServiceServer
	processor *EventProcessor
}

// NewEngineGRPCServer creates a new gRPC server for the engine.
func NewEngineGRPCServer(processor *EventProcessor) *EngineGRPCServer {
	return &EngineGRPCServer{processor: processor}
}

// GetCardState returns the current state for a user-card pair.
func (s *EngineGRPCServer) GetCardState(ctx context.Context, req *pb.GetCardStateRequest) (*pb.CardStateResponse, error) {
	state, err := s.processor.stateService.GetByUserAndCard(ctx, req.UserId, req.CardId)
	if err != nil {
		return nil, fmt.Errorf("get card state: %w", err)
	}

	return &pb.CardStateResponse{
		StateId:      state.ID().String(),
		NextReviewAt: state.NextReviewAt().Format("2006-01-02T15:04:05Z"),
	}, nil
}

// GetNextCards returns the next recommended cards for a user in a collection.
func (s *EngineGRPCServer) GetNextCards(ctx context.Context, req *pb.GetNextCardsRequest) (*pb.NextCardsResponse, error) {
	if s.processor.redisClient == nil {
		return nil, fmt.Errorf("redis not available")
	}

	cards, err := s.processor.redisClient.GetNextCards(ctx, req.UserId, req.CollectionId, int64(req.Limit))
	if err != nil {
		return nil, fmt.Errorf("get next cards: %w", err)
	}

	recommendations := make([]*pb.CardRecommendation, len(cards))
	for i, c := range cards {
		cardID, _ := c.Member.(string)
		recommendations[i] = &pb.CardRecommendation{
			CardId:     cardID,
			FinalScore: c.Score,
		}
	}

	return &pb.NextCardsResponse{Cards: recommendations}, nil
}

// StartGRPCServer starts the gRPC server and returns the server instance.
func StartGRPCServer(addr string, srv *EngineGRPCServer) (*grpc.Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("create grpc listener: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterEngineServiceServer(grpcServer, srv)

	log.Printf("Starting gRPC server on %s", addr)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	return grpcServer, nil
}
