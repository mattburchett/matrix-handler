package prometheus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/config"
	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/matrix"
	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/router"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type prometheusMessage struct {
	Alerts []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname string `json:"alertname"`
			Instance  string `json:"instance"`
			Severity  string `json:"severity"`
		} `json:"labels,omitempty"`
		Annotations struct {
			Description string `json:"description"`
			Summary     string `json:"summary"`
		} `json:"annotations"`
		StartsAt time.Time `json:"startsAt"`
	} `json:"alerts"`
}

// Handle is the incoming handler for Slack-type requests.
func Handle(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// Get Matrix Token for User/Pass in path
		token := matrix.GetToken(cfg, vars)

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("An error has occurred")
		}

		data := parsePrometheus(reqBody)

		// Attempt to join room before publishing
		matrix.JoinRoom(cfg, vars, token)

		// Publish to Matrix
		for _, message := range data {
			_ = matrix.PublishText(cfg, vars, []byte(message), token)
		}

		router.Respond(w, 200, nil)
	}

}

func parsePrometheus(body []byte) []string {

	reqBody := prometheusMessage{}

	json.Unmarshal(body, &reqBody)

	// NodeClockNotSynchronising : firing
	// Instance : 10.234.62.22:9100
	// Severity : warning
	// Started : 2021-09-08T15:10:49.704865181Z
	// Description : Clock on 10.234.62.22:9100 is not synchronising. Ensure NTP is configured on this host.
	var message []string
	for _, i := range reqBody.Alerts {
		message = append(message, fmt.Sprintf("%s : %s\nInstance : %s\nSeverity : %s\nStarted : %s\nDescription : %s\n", i.Labels.Alertname, i.Status,
			i.Labels.Instance, i.Labels.Severity, i.StartsAt, i.Annotations.Description))
	}

	return message
}
