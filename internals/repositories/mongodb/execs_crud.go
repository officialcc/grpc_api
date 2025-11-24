package mongodb

import (
	"context"
	"grpcapi/internals/models"
	"grpcapi/pkg/utils"
	pb "grpcapi/proto/gen"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddExecsToDb(ctx context.Context, execsFromReq []*pb.Exec) ([]*pb.Exec, error) {
	client, err := CreateMongoClient()
	if err != nil {
		return nil, utils.ErrorHandler(err, "internal error")
	}
	defer client.Disconnect(ctx)

	newExecs := make([]*models.Exec, len(execsFromReq))
	for i, pbExec := range execsFromReq {
		newExecs[i] = mapPbExecToModelExec(pbExec)

	}

	var addedExecs []*pb.Exec
	for _, exec := range newExecs {
		result, err := client.Database("school").Collection("execs").InsertOne(ctx, exec)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error adding value to database")
		}
		objectId, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			exec.Id = objectId.Hex()
		}

		pbExec := mapModelExecToPb(*exec)
		addedExecs = append(addedExecs, pbExec)
	}
	return addedExecs, nil
}
