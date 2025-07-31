package massage

import (
	"context"
	"goCodetest/internal/models"
	"goCodetest/pkg/logger"
	"time"
)

func init() {
	// Initialize any necessary resources or configurations here
	logger.Info("Email notification service initialized")
}

func EmailNotificationForProductDeletion(ctx context.Context, productID int64) {
	traceID := ctx.Value(models.TraceIDKey)
	if traceID == nil {
		traceID = "unknown"
	}

	logger.Info("TraceId: %s - Sending email notification for product deletion with ID: %d", traceID, productID)

	// you can seed msg to kafka
	time.Sleep(2 * time.Second)

	logger.Info("TraceId: %s - Email notification sent successfully for product deletion with ID: %d", traceID, productID)
}
