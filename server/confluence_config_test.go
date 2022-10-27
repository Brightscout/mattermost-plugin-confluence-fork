package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/mattermost/mattermost-plugin-confluence/server/config"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/mattermost/mattermost-server/v6/plugin/plugintest/mock"
	"github.com/stretchr/testify/assert"
)

const (
	validUserID    = "iu73atknztnctef8b8ey9gm6zc"
	validChannelID = "tgniw3kmrjd93qns11cboditme"
)

func TestHandleConfluenceConfig(t *testing.T) {
	tests := map[string]struct {
		method         string
		statusCode     int
		body           string
		userID         string
		channelID      string
		patchFuncCalls func()
	}{
		"success": {
			method:     http.MethodPost,
			statusCode: http.StatusOK,
			body: `{
				"type": "dialog_submission",
				"callback_id": "callbackID",
				"state": "", 
				"submission": {
					"Client ID": "mock-ClientID",
					"Client Secret": "mock-ClientSecret",
					"Server URL": "https://test.com"
				},
				"cancelled": false
			}`,
			userID:    "iu73atknztnctef8b8ey9gm6zc",
			channelID: "tgniw3kmrjd93qns11cboditme",
			patchFuncCalls: func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(&Plugin{}), "GetConfigKeyList", func(_ *Plugin) ([]string, error) {
					return []string{
						"https://test.com",
					}, nil
				})
			},
		},
		"wrong api method": {
			method:     http.MethodGet,
			statusCode: http.StatusMethodNotAllowed,
		},
		"invalid body": {
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			body:       `{`,
			userID:     "iu73atknztnctef8b8ey9gm6zc",
			channelID:  "tgniw3kmrjd93qns11cboditme",
			patchFuncCalls: func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(&Plugin{}), "GetConfigKeyList", func(_ *Plugin) ([]string, error) {
					return []string{
						"https://test.com",
					}, nil
				})
			},
		},
		"invalid userID or channelID": {
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			body:       `{`,
			userID:     "mock-userID",
			channelID:  "mock-channelID",
			patchFuncCalls: func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(&Plugin{}), "GetConfigKeyList", func(_ *Plugin) ([]string, error) {
					return []string{
						"https://test.com",
					}, nil
				})
			},
		},
	}
	mockAPI := baseMock()
	mockAPI.On("LogError", mockAnythingOfTypeBatch("string", 13)...).Return(nil)

	mockAPI.On("LogDebug", mockAnythingOfTypeBatch("string", 11)...).Return(nil)

	p := Plugin{}
	p.SetAPI(mockAPI)

	mockAPI.On("GetBundlePath").Return("/test/testBundlePath", nil)

	p.userStore = getMockUserStoreKV()
	p.instanceStore = p.getMockInstanceStoreKV(1)

	p.otsStore = p.getMockOTSStoreKV()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()

			mockAPI.On("SendEphemeralPost", mock.AnythingOfType("string"), mock.AnythingOfType("*model.Post")).Once().Return(&model.Post{})

			mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Id: "123", Roles: "system_admin"}, nil)

			if tc.patchFuncCalls != nil {
				tc.patchFuncCalls()
			}

			request := httptest.NewRequest(tc.method, fmt.Sprintf("/api/v1/config/%s/%s", validUserID, validUserID), bytes.NewBufferString(tc.body))

			request.Header.Set(config.HeaderMattermostUserID, "test-user")
			w := httptest.NewRecorder()
			p.ServeHTTP(&plugin.Context{}, w, request)
			assert.Equal(t, tc.statusCode, w.Result().StatusCode)
		})
	}
}
