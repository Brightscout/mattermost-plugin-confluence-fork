package main

import (
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

func TestHandleGetConfigList(t *testing.T) {
	tests := map[string]struct {
		method         string
		statusCode     int
		patchFuncCalls func()
	}{
		"success": {
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			patchFuncCalls: func() {
				monkey.PatchInstanceMethod(reflect.TypeOf(&Plugin{}), "GetConfigKeyList", func(_ *Plugin) ([]string, error) {
					return []string{
						"https://test.com",
					}, nil
				})
			},
		},
		"wrong api method": {
			method:     http.MethodPost,
			statusCode: http.StatusMethodNotAllowed,
		},
	}
	mockAPI := baseMock()
	mockAPI.On("LogError", mockAnythingOfTypeBatch("string", 13)...).Return(nil)

	mockAPI.On("LogDebug", mockAnythingOfTypeBatch("string", 11)...).Return(nil)

	p := Plugin{}
	p.SetAPI(mockAPI)

	mockAPI.On("GetBundlePath").Return("/test/testBundlePath", nil)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer monkey.UnpatchAll()

			mockAPI.On("GetUser", mock.AnythingOfType("string")).Return(&model.User{Id: "123", Roles: "system_admin"}, nil)

			if tc.patchFuncCalls != nil {
				tc.patchFuncCalls()
			}

			request := httptest.NewRequest(tc.method, "/api/v1/autocomplete/configs", nil)
			request.Header.Set(config.HeaderMattermostUserID, "test-user")
			w := httptest.NewRecorder()
			p.ServeHTTP(&plugin.Context{}, w, request)
			assert.Equal(t, tc.statusCode, w.Result().StatusCode)
		})
	}
}
