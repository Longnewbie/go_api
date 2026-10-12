package middleware

import (
	"context"
	"user-management-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.GetHeader("X-Trace-Id")
		if traceId == "" {
			traceId = uuid.New().String()
		}

		contextValue := context.WithValue(ctx.Request.Context(), logger.TraceIdKey, traceId)
		ctx.Request = ctx.Request.WithContext(contextValue)

		ctx.Writer.Header().Set("X-Trace-Id", traceId) // attach response header

		ctx.Set(string(logger.TraceIdKey), traceId)

		ctx.Next()
	}
}
