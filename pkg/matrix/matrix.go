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

	body := postRequest(s, bytes.NewBuffer(reqBody))

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
func PublishText(cfg config.Config, vars map[string]string, data []byte, token string) []byte {
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

	body := postRequest(s, bytes.NewBuffer(reqBody))

	return body
}

// JoinRoom will attempt to join a matrix rooom, assuming there is an invite pending.
func JoinRoom(cfg config.Config, vars map[string]string, token string) {
	s := fmt.Sprintf("%v:%v/_matrix/client/r0/rooms/%v/join?access_token=%v", cfg.Matrix.Homeserver, cfg.Matrix.Port, vars["matrixRoom"], token)
	_ = postRequest(s, bytes.NewBuffer(nil))
}

// postRequest performs the post requests to the Matrix server.
func postRequest(s string, data *bytes.Buffer) []byte {
	req, err := http.NewRequest(http.MethodPost, s, data)
	if err != nil {
		log.Error().Err(err).Msg("matrix.postRequest.req" + err.Error())
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("matrix.postRequest.resp" + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return body
}
