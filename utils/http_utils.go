package utils

import (
	"github.com/aruncs31s/skvms/utils/http"
	"github.com/gin-gonic/gin"
)

func GetLimitAndOffset(
	c *gin.Context,
) (int, int) {
	return http.GetLimitAndOffset(c)
}
