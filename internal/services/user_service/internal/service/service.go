package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	
	"we_ride/internal/services/user_service/internal/jwt"
	"we_ride/internal/services/user_service/internal/repository"
	pb "we_ride/internal/services/user_service/protoc/gen/go"
)

type ServerAPI struct {
	pb.UnimplementedAuthServer
	repo repository.Repository
}

func New(repo repository.Repository) *ServerAPI {
	return &ServerAPI{repo: repo}
}

func (s *ServerAPI) Login(
	ctx context.Context,
	req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}
	token, err := s.repo.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	resp := &pb.LoginResponse{Token: token}

	return resp, nil
}

func (s *ServerAPI) Register(
	ctx context.Context,
	req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}
	if req.GetFirstName() == "" {
		return nil, status.Error(codes.InvalidArgument, "First name is required")
	}
	if req.GetLastName() == "" {
		return nil, status.Error(codes.InvalidArgument, "Last name is required")
	}
	if req.GetGender() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Gender is required")
	}
	UserId, err := s.repo.SaveUser(
		ctx,
		req.GetEmail(),
		req.GetPassword(), req.GetFirstName(),
		req.GetLastName(),
		req.GetGender(),
	)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	resp := &pb.RegisterResponse{UserId: UserId}
	return resp, nil
}

func (s *ServerAPI) HistoryOrRoutes(
	ctx context.Context,
	req *pb.HistoryOfRoutesRequest) (*pb.HistoryOfRoutesResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Metadata is required")
	}
	token := strings.TrimPrefix(md["authorization"][0], "Bearer ")
	claims, err := jwt.ValidateToken(token, s.repo.Secret)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	userID := claims.UserID
	routes, err := s.repo.GetUserRoutes(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get routes: %v", err)
	}

	resp := pb.HistoryOfRoutesResponse{}
	for _, route := range routes {
		resp.Routes = append(resp.Routes, &pb.Route{
			RouteId:    route.RouteId,
			DriverId:   route.DriverId,
			TotalPrice: route.TotalPrice,
			StartPoint: route.StartPoint,
			EndPoint:   route.EndPoint,
			Distance:   route.Distance,
		})
	}
	return &resp, nil
}
