package main

import (
	"encoding/json"
	"io/ioutil"
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
	"github.com/stretchr/testify/require"
)

func TestHandleGetConfigList(t *testing.T) {
	tests := map[string]struct {
		method         string
		statusCode     int
		patchFuncCalls func()
		resp           []model.AutocompleteListItem
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
			resp: []model.AutocompleteListItem{
				{
					Item: "https://test.com",
				},
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
	mockAPI.On("GetBundlePath").Return("/test/testBundlePath", nil)

	p := Plugin{}
	p.SetAPI(mockAPI)

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

			bodyBytes, err := ioutil.ReadAll(w.Body)
			require.Nil(t, err)
			out := []model.AutocompleteListItem{}
			if tc.statusCode == http.StatusOK {
				err = json.Unmarshal(bodyBytes, &out)
				require.Nil(t, err)
				assert.Equal(t, out, tc.resp)
			}

			assert.Equal(t, tc.statusCode, w.Result().StatusCode)
		})
	}
}
