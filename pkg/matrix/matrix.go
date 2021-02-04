package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.linuxrocker.com/mattburchett/matrix-handler/pkg/config"
	"github.com/rs/zerolog/log"
)

// GetToken will get the access token from Matrix to perform communications.
func GetToken(cfg config.Config, vars map[string]string) string {
	matrixConfig := struct {
		Type     string `json:"type"`
		Username string `json:"user"`
		Password string `json:"password"`
	}{
		Type:     "m.login.password",
		Username: vars["matrixUser"],
		Password: vars["matrixPassword"],
	}

	reqBody, err := json.Marshal(matrixConfig)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	s := fmt.Sprintf("%v:%v/_matrix/client/r0/login", cfg.Matrix.Homeserver, cfg.Matrix.Port)

	req, err := http.NewRequest(http.MethodPost, s, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Error().Err(err).Msg("matrix.GetToken.req" + err.Error())
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("matrix.GetToken.resp" + err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("matrix.GetToken.body" + err.Error())
	}

	respBody := struct {
		AccessToken string `json:"access_token"`
	}{}

	err = json.Unmarshal(body, &respBody)
	if err != nil {
		log.Error().Err(err).Msg("matrix.GetToken.respBody" + err.Error())
	}

	return respBody.AccessToken
}

// PublishText will publish the data to Matrix using the specified vars.
func PublishText(cfg config.Config, vars map[string]string, data []byte, token string) {
	matrixPublish := struct {
		MsgType string `json:"msgtype"`
		Body    string `json:"body"`
	}{
		MsgType: "m.text",
		Body:    string(data),
	}

	reqBody, err := json.Marshal(matrixPublish)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	s := fmt.Sprintf("%v:%v/_matrix/client/r0/rooms/%v/send/m.room.message?access_token=%v", cfg.Matrix.Homeserver, cfg.Matrix.Port, vars["matrixRoom"], token)

	req, err := http.NewRequest(http.MethodPost, s, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Error().Err(err).Msg("matrix.PublishText.req" + err.Error())
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("matrix.PublishText.resp" + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	defer resp.Body.Close()

}
