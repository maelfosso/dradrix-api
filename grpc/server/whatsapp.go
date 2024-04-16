package server

import (
	"context"
	"fmt"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	pb "stockinos.com/api/grpc/protos"
	"stockinos.com/api/models"
	"stockinos.com/api/services"
	"stockinos.com/api/storage"
)

type GrpcWhatsappServer struct {
	pb.UnimplementedWhatsappServer
	database *storage.Database
	log      *zap.Logger
}

type GrpcWhatsappOptions struct {
	Database *storage.Database
	Log      *zap.Logger
}

func NewGrpcWhatsappServer(opts GrpcWhatsappOptions) *GrpcWhatsappServer {
	if opts.Log == nil {
		opts.Log = zap.NewNop()
	}

	log.Println("log of opts is nil? ", opts.Log)

	return &GrpcWhatsappServer{
		database: opts.Database,
		log:      opts.Log,
	}
}

func (gs *GrpcWhatsappServer) InsertWozMessage(ctx context.Context, request *pb.InsertWozMessageRequest) (*pb.WhatsappMessageResponse, error) {
	waId, err := services.WASendTextMessage(request.To, request.Message)
	if err != nil {
		return nil, fmt.Errorf("error when sent whatsapp message: %w", err)
	}
	gs.log.Info("id after sending the wa message", zap.String("id", waId))

	textId := uuid.NewV4()
	waMsg := models.WhatsappMessage{
		ID:        waId,
		From:      request.From,
		To:        request.To,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),

		Type: "text",
		WhatsappMessageType: models.WhatsappMessageType{
			TextID: textId,
			Text: models.WhatsappMessageText{
				ID:   textId,
				Body: request.Message,
			},
		},
	}

	gs.log.Info("Whatsapp message to be inserted", zap.Any("wa message", waMsg))
	// err = gs.database.SaveWAMessages(ctx, []models.WhatsappMessage{waMsg})
	// if err != nil {
	// 	return nil, fmt.Errorf("error when saving the message in database: %w", err)
	// }

	textIdString := textId.String()
	return &pb.WhatsappMessageResponse{
		Id:        waMsg.ID,
		From:      waMsg.From,
		To:        waMsg.To,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),

		Type:   waMsg.Type,
		TextId: &textIdString,
		Text: &pb.WhatsappMessageText{
			Id:   waMsg.Text.ID.String(),
			Body: waMsg.Text.Body,
		},
	}, nil
}
