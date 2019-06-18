package handler

import (
	"encoding/json"

	"github.com/dbmonstar/dbmond/common"
	"github.com/dbmonstar/dbmond/model"

	"github.com/gin-gonic/gin"
)

// StartErrorLogAPI error log api
func startErrorLogAPI(r *gin.RouterGroup) {

	r.POST("/errorlog", func(c *gin.Context) {
		body, err := c.GetRawData()
		if ErrorIf(c, err) {
			return
		}

		var errorLogs []model.ErrorLog
		json.Unmarshal(body, &errorLogs)

		successCount := 0
		for _, errorLog := range errorLogs {
			if err := errorLog.Insert(); err != nil {
				common.Log.Error(err)
				continue
			}
			successCount++
		}

		Success(c, successCount)
	})
}
