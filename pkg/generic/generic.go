package generic

import (
	"io/ioutil"
	"net/http"

	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/config"
	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/matrix"
	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/router"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// Handle is the incoming handler for Generic-type requests.
func Handle(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// Get Matrix Token for User/Pass in path
		token := matrix.GetToken(cfg, vars)

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("An error has occurred")
		}

		resp := matrix.PublishText(cfg, vars, reqBody, token)

		router.Respond(w, 200, resp)
	}

}
