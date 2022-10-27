package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-plugin-confluence/server/serializer"
	"github.com/mattermost/mattermost-server/v6/model"
)

const (
	ParamUserID = "userID"
)

func (p *Plugin) handleConfluenceConfig(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	channelID := pathParams[ParamChannelID]
	userID := pathParams[ParamUserID]
	decoder := json.NewDecoder(r.Body)
	submitRequest := &model.SubmitDialogRequest{}
	if err := decoder.Decode(&submitRequest); err != nil {
		p.API.LogError("Error decoding SubmitDialogRequest.", "Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config := &serializer.ConfluenceConfig{
		ServerURL:    submitRequest.Submission[configServerURL].(string),
		ClientID:     submitRequest.Submission[configClientID].(string),
		ClientSecret: submitRequest.Submission[configClientSecret].(string),
	}

	if err := p.instanceStore.StoreInstanceConfig(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.API.LogError("Error storing instance config.", "Error", err.Error())
		return
	}

	post := &model.Post{
		UserId:    p.conf.botUserID,
		ChannelId: channelID,
		Message:   fmt.Sprintf("Your config is saved for confluence instance %s", config.ServerURL),
	}

	_ = p.API.SendEphemeralPost(userID, post)
	ReturnStatusOK(w)
}
