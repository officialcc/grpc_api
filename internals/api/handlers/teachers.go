package handlers

import (
	"context"
	"grpcapi/internals/models"
	"grpcapi/internals/repositories/mongodb"
	"grpcapi/pkg/utils"
	pb "grpcapi/proto/gen"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (s *Server) UpdateTeachers(ctx context.Context, req *pb.Teachers) (*pb.Teachers, error) {
	updatedTeachers, err := mongodb.ModifyTeachersInDb(ctx, req.Teachers)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Teachers{Teachers: updatedTeachers}, nil
}

func (s *Server) DeleteTeachers(ctx context.Context, req *pb.TeacherIds) (*pb.DeleteTeachersConfirmation, error) {
	ids := req.GetIds()
	var teacherIdsToDelete []string
	for _, teacher := range ids {
		teacherIdsToDelete = append(teacherIdsToDelete, teacher.Id)
	}

	deletedIds, err := mongodb.DeleteTeachersFromDb(ctx, teacherIdsToDelete)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteTeachersConfirmation{
		Status:     "Teachers successfully deleted",
		DeletedIds: deletedIds,
	}, nil
}

func (*Server) GetStudentsbyClassTeacher(ctx context.Context, req *pb.TeacherId) (*pb.Students, error) {
	teacherId := req.GetId()

	client, err := mongodb.CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	objId, err := primitive.ObjectIDFromHex(teacherId)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Invalid Teacher ID")
	}

	var teacher models.Teacher
	err = client.Database("school").Collection("teachers").FindOne(ctx, bson.M{"_id": objId}).Decode(&teacher)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrorHandler(err, "teacher not found")
		}
		return nil, utils.ErrorHandler(err, "internal error")
	}

	cursor, err := client.Database("school").Collection("students").Find(ctx, bson.M{"class": teacher.Class})
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer cursor.Close(ctx)
	
}
