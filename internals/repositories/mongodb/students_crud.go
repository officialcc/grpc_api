package mongodb

import (
	"context"
	"grpcapi/internals/models"
	"grpcapi/pkg/utils"
	pb "grpcapi/proto/gen"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddStudentsToDb(ctx context.Context, studentsFromReq []*pb.Student) ([]*pb.Student, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	newStudents := make([]*models.Student, len(studentsFromReq))
	for i, pbStudent := range studentsFromReq {
		newStudents[i] = mapPbStudentToModelStudent(pbStudent)

	}

	var addedStudents []*pb.Student
	for _, student := range newStudents {
		result, err := client.Database("school").Collection("students").InsertOne(ctx, student)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error adding value to database")
		}
		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			student.Id = objectId.Hex()
		}

		pbStudent := mapModelStudentToPb(*student)
		addedStudents = append(addedStudents, pbStudent)
	}
	return addedStudents, nil
}
