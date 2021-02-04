package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/config"
	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/matrix"
	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/router"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

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

		data := parseSlack(reqBody)

		// Attempt to join romo before publishing
		matrix.JoinRoom(cfg, vars, token)

		// Publish to Matrix
		resp := matrix.PublishText(cfg, vars, []byte(data), token)

		router.Respond(w, 200, resp)
	}

}

func parseSlack(body []byte) string {
	reqBody := struct {
		Text string
	}{}

	json.Unmarshal(body, &reqBody)

	return reqBody.Text
}
