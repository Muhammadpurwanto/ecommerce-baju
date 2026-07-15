package util

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcCodeToHTTPStatus maps gRPC status codes to appropriate HTTP status codes.
func GrpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return fiber.StatusOK
	case codes.InvalidArgument:
		return fiber.StatusBadRequest
	case codes.NotFound:
		return fiber.StatusNotFound
	case codes.AlreadyExists:
		return fiber.StatusConflict
	case codes.PermissionDenied:
		return fiber.StatusForbidden
	case codes.Unauthenticated:
		return fiber.StatusUnauthorized
	case codes.ResourceExhausted:
		return fiber.StatusTooManyRequests
	case codes.FailedPrecondition:
		return fiber.StatusBadRequest
	case codes.Unimplemented:
		return fiber.StatusNotImplemented
	case codes.Unavailable:
		return fiber.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return fiber.StatusGatewayTimeout
	default:
		return fiber.StatusInternalServerError
	}
}

// HandleGrpcError converts a gRPC error into a proper Fiber JSON response.
// It extracts the gRPC status message for client-facing response and maps
// the gRPC status code to the appropriate HTTP status code.
// Internal/unknown errors return a generic message to avoid leaking details.
func HandleGrpcError(c *fiber.Ctx, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC status error — return generic 500
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "An internal error occurred",
		})
	}

	httpStatus := GrpcCodeToHTTPStatus(st.Code())
	message := st.Message()

	// Untuk error internal/unknown, jangan expose detail ke client
	if st.Code() == codes.Internal || st.Code() == codes.Unknown {
		logger.Error("gRPC internal error",
        zap.Error(err),
        zap.String("grpc_code", st.Code().String()),
    )
		message = "An internal error occurred"
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}
