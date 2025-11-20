package handlers

import (
	"context"
	"grpcapi/internals/models"
	"grpcapi/internals/repositories/mongodb"
	pb "grpcapi/proto/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {

	for _, teacher := range req.GetTeachers() {
		if teacher.Id != "" {
			return nil, status.Error(codes.InvalidArgument, "request is in incorrect format: non-empty ID field is not allowed")
		}
	}

	addedTeachers, err := mongodb.AddTeachersToDb(ctx, req.GetTeachers())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Teachers{Teachers: addedTeachers}, nil
}

func (s *Server) GetTeachers(ctx context.Context, req *pb.GetTeachersRequest) (*pb.Teachers, error) {
	// Filtering - Getting the filters from the request -> Another function
	filter, err := buildFilter(req.Teacher, &models.Teacher{})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Sorting - Getting the sort options from the request -> Another function
	sortOptions := buildSortOptions(req.GetSortBy())
	// Access the database to fetch data - Another function

	teachers, err := mongodb.GetTeachersFromDb(ctx, sortOptions, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Teachers{Teachers: teachers}, nil
}
