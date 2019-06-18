package handler

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dbmonstar/dbmond/model"

	"github.com/dbmonstar/dbmond/common"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/alertmanager/template"
)

// startHook hook API
func startHookAPI(r *gin.RouterGroup) {

	chanHook := make(chan model.Hook, 100)
	sendMessage(chanHook)
	crc32q := crc32.MakeTable(0xD5828281)

	// new
	r.POST("/hook/send", func(c *gin.Context) {
		var err error
		var params template.Data

		// bind template json data
		err = c.BindJSON(&params)
		if ErrorIf(c, err) {
			return
		}

		rows := 0
		for _, alert := range params.Alerts {

			// Generate hash code for this alert
			var result uint64
			for key, value := range alert.Labels {
				result += uint64(crc32.Checksum([]byte(key+":"+value), crc32q))
			}

			hook := model.Hook{
				Name:        alert.Labels["alertname"],
				Hash:        fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d", result)))),
				Instance:    alert.Labels["instance"],
				Level:       alert.Labels["level"],
				Status:      strings.ToLower(alert.Status),
				Subject:     alert.Annotations["summary"],
				Description: alert.Annotations["description"],
			}

			switch strings.ToLower(alert.Status) {
			case "firing":
				hook.StartsAt = alert.StartsAt

			case "resolved":
				hook.StartsAt = alert.StartsAt
				hook.EndsAt = alert.EndsAt

			default:
				log.Printf("no action on status: %s", alert.Status)
				continue
			}

			hook.Insert()
			chanHook <- hook
			rows++
		}

		Success(c, rows)
	})
}

func sendMessage(chanHook chan model.Hook) {

	go func() {

		httpClient := &http.Client{Timeout: 3 * time.Second}
		timezone := common.Timezone

		loc, err := time.LoadLocation(timezone)
		if err != nil {
			log.Printf("set timezone failed, set to UTC")
			loc, _ = time.LoadLocation("UTC")
		}

		for {
			hook := <-chanHook
			log.Printf("Message routine..")

			var message string
			switch strings.ToLower(hook.Status) {
			case "firing":
				// ========================
				// Message
				// ========================
				message = "[Firing] " + hook.Subject + "\n"
				message += "* Instance: " + hook.Instance + "\n"
				message += "* Level: " + hook.Level + "\n"
				message += "* Started: " + hook.StartsAt.In(loc).Format("01/02 15:04:05") + "\n"
				message += "* Description: " + hook.Description + "\n"

			case "resolved":
				// duration minutes
				minutes := fmt.Sprintf(" (%.1f min)", hook.EndsAt.Sub(hook.StartsAt).Minutes())

				// ========================
				// Message
				// ========================
				message = "[Resolved] " + hook.Subject + "\n"
				message += "* Instance: " + hook.Instance + "\n"
				message += "* Ended: " + hook.EndsAt.In(loc).Format("01/02 15:04:05") + minutes + "\n"
				message += "* Level: " + hook.Level + "\n"
			}

			api := common.AlarmAPI[hook.Level]
			if api != "" {
				resp, err := httpClient.Get(api + url.QueryEscape(message))
				if err != nil {
					log.Printf("Send alert to [%s] failed: %s", api, err.Error())
					continue
				}
				defer resp.Body.Close()
			}
		}
	}()
}
