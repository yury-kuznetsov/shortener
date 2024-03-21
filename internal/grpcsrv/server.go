package grpcsrv

import (
	"context"
	"errors"

	"github.com/yury-kuznetsov/shortener/api/pb"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/models"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contextKey string

const KeyUserID contextKey = "USER_ID"

// NewCoderServer creates a new instance of CoderServer with the provided Coder instance.
func NewCoderServer(coder *uricoder.Coder) *CoderServer {
	return &CoderServer{coder: coder}
}

// CoderServer is a struct that represents a server implementing the Coder gRPC service.
// It embeds pb.UnimplementedServiceServer and contains a reference to the uricoder.Coder instance.
type CoderServer struct {
	pb.UnimplementedServiceServer
	coder *uricoder.Coder
}

// Decode is a method of CoderServer that decodes the provided code into a URI.
// It requires a context object and a DecodeRequest as input parameters.
// It returns a DecodeResponse and an error.
// The context object is used to get the user ID from the context value.
// It calls the ToURI method of the coder instance to get the URI for the provided code.
// If an error occurs while decoding, it checks if the error is a "ErrRowDeleted" error and returns a status error with the relevant code and message.
// If the error is not a "ErrRowDeleted" error, it returns a status error with the internal server error code and the error message.
// It returns the URI in a DecodeResponse if decoding is successful.
// Example usage:
//
//	ctx := context.Background()
//	request := &pb.DecodeRequest{Code: "abc123"}
//	response, err := coderServer.Decode(ctx, request)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Decoded URI:", response.Uri)
func (s *CoderServer) Decode(ctx context.Context, in *pb.DecodeRequest) (*pb.DecodeResponse, error) {
	userID := ctx.Value(KeyUserID).(int)
	uri, err := s.coder.ToURI(ctx, in.GetCode(), userID)
	if err != nil {
		if errors.Is(err, models.ErrRowDeleted) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DecodeResponse{Uri: uri}, nil
}

// Encode is a method of CoderServer that encodes the provided URI into a code.
// It requires a context object and an EncodeRequest as input parameters.
// It returns an EncodeResponse and an error.
// The context object is used to get the user ID from the context value.
// It calls the ToCode method of the coder instance to get the code for the provided URI.
// If the code is empty and an error occurs, it returns a status error with the invalid argument code and the error message.
// It constructs an EncodeResponse with the base address and the code.
// If an error occurs while encoding, it returns the response along with a status error with the already exists code and the error message.
// Otherwise, it returns the response and nil error.
// Example usage:
//
//	ctx := context.Background()
//	request := &pb.EncodeRequest{Uri: "https://example.com"}
//	response, err := coderServer.Encode(ctx, request)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println("Encoded Code:", response.Code)
func (s *CoderServer) Encode(ctx context.Context, in *pb.EncodeRequest) (*pb.EncodeResponse, error) {
	userID := ctx.Value(KeyUserID).(int)
	code, err := s.coder.ToCode(ctx, in.GetUri(), userID)
	if code == "" && err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	response := &pb.EncodeResponse{Code: config.Options.BaseAddr + "/" + code}
	if err != nil {
		return response, status.Error(codes.AlreadyExists, err.Error())
	}
	return response, nil
}

// EncodeByID is a method of CoderServer that encodes the provided URI into a code based on the user ID.
// It requires a context object and an EncodeByIDRequest as input parameters.
// It returns an EncodeByIDResponse and an error.
// The context object is used to get the user ID from the context value.
// It calls the ToCode method of the coder instance to get the code for the provided URI.
// If the code is empty and err is not nil, it returns a status error with the invalid argument code and the error message.
// It constructs an EncodeByIDResponse with the provided ID and the code prefixed with the base address from the config options.
// If err is not nil, it returns the response with a status error with the already exists code and the error message.
// Otherwise, it returns the response and nil as the error.
func (s *CoderServer) EncodeByID(ctx context.Context, in *pb.EncodeByIDRequest) (*pb.EncodeByIDResponse, error) {
	userID := ctx.Value(KeyUserID).(int)
	code, err := s.coder.ToCode(ctx, in.GetUri(), userID)
	if code == "" && err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	response := &pb.EncodeByIDResponse{
		Id:   in.GetId(),
		Code: config.Options.BaseAddr + "/" + code,
	}
	if err != nil {
		return response, status.Error(codes.AlreadyExists, err.Error())
	}
	return response, nil
}

// GetHistory is a method of CoderServer that retrieves the history of encoded URLs for a user.
// It requires a context object as an input parameter.
// It returns a GetHistoryResponse and an error.
// The context object is used to get the user ID from the context value.
// It calls the GetHistory method of the coder instance to retrieve the history data for the user.
// If an error occurs while retrieving the history data, it returns a status error with the internal server error code and the error message.
// It creates a list of History objects based on the retrieved data and appends them to the Histories slice in the GetHistoryResponse.
// It returns the GetHistoryResponse with the list of histories if retrieval is successful.
//
// Example usage:
//
//	ctx := context.Background()
//	response, err := coderServer.GetHistory(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, history := range response.Histories {
//	    fmt.Println("Code:", history.Code)
//	    fmt.Println("URI :", history.Uri)
//	}
func (s *CoderServer) GetHistory(ctx context.Context) (*pb.GetHistoryResponse, error) {
	userID := ctx.Value(KeyUserID).(int)
	data, err := s.coder.GetHistory(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	var histories []*pb.History
	for _, v := range data {
		histories = append(histories, &pb.History{
			Code: v.ShortURL,
			Uri:  v.OriginalURL,
		})
	}
	return &pb.GetHistoryResponse{Histories: histories}, nil
}

// Delete is a method of CoderServer that deletes the URLs associated with the provided codes.
// It requires a context object and a DeleteRequest as input parameters.
// It returns a DeleteResponse and an error.
// The context object is used to get the user ID from the context value.
// It calls the DeleteUrls method of the coder instance to delete the URLs associated with the provided codes.
// It returns an empty DeleteResponse if deletion is successful.
// Example usage:
// ctx := context.Background()
// request := &pb.DeleteRequest{Codes: []string{"abc123", "def456"}}
// response, err := coderServer.Delete(ctx, request)
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// fmt.Println("Deletion completed successfully")
func (s *CoderServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	userID := ctx.Value(KeyUserID).(int)
	_ = s.coder.DeleteUrls(in.Codes, userID)
	return &pb.DeleteResponse{}, nil
}
