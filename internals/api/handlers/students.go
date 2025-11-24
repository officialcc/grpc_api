package handlers

import (
	"context"
	"grpcapi/internals/repositories/mongodb"
	pb "grpcapi/proto/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddStudents(ctx context.Context, req *pb.Students) (*pb.Students, error) {

	for _, student := range req.GetStudents() {
		if student.Id != "" {
			return nil, status.Error(codes.InvalidArgument, "request is in incorrect format: non-empty ID field is not allowed")
		}
	}

	addedStudents, err := mongodb.AddStudentsToDb(ctx, req.GetStudents())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.Students{Students: addedStudents}, nil
}
