package http

import (
	"strconv"

	"github.com/aruncs31s/skvms/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetLimitAndOffset(
	c *gin.Context,
) (int, int) {
	limit := c.Query("limit")

	limitInt, err := strconv.Atoi(limit)

	if err != nil || limitInt == 0 {

		limitInt = 10
		logger.GetLogger().Warn(
			"Applaying Default Limit",
			zap.Int(
				"limit", limitInt,
			),
		)
	}

	offset := c.Query("offset")

	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt == 0 {
		offsetInt = 0
		logger.GetLogger().Warn(
			"Applaying Default Offset",
			zap.Int(
				"offset", offsetInt,
			),
		)
	}
	return limitInt, offsetInt
}
