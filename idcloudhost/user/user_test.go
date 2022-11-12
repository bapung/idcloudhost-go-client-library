package user

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

var (
	mockHttpClient = &HTTPClientMock{}
	u              = UserAPI{}
	loc            = "jkt01"
)

func TestGetUser(t *testing.T) {
	u.Init(mockHttpClient, userAuthToken, loc)
	testCases := []struct {
		RequestData map[string]interface{}
		Body        string
		StatusCode  int
		Error       error
	}{
		{
			RequestData: map[string]interface{}{
				"uuid": "this-is-a-supposed-to-be-valid",
			},
			Body:       `{"cookie_id":"61b0378574974ae88dbfec0feb9917bc","id":8,"last_activity":"2018-02-22 14:18:47","name":"user@example.com","profile":null,"profile_data":{"avatar":"https://s.gravatar.com/avatar/bbb?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fsv.png","created_at":"2018-10-25 11:02:59","email":"user@example.com","first_name":"Cloudia","id":22,"last_name":"Iaas","personal_id_number":"123456","phone_number":"+111111111","updated_at":"2021-05-18 11:07:00","user_id":8},"state":{}}`,
			StatusCode: http.StatusOK,
			Error:      nil,
		},
	}
	for _, test := range testCases {
		mockHttpClient.DoFunc = func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				Body:       io.NopCloser(strings.NewReader(test.Body)),
				StatusCode: test.StatusCode,
			}, nil
		}

		err := u.Get(test.RequestData["uuid"].(string))
		if err != nil && err.Error() != test.Error.Error() {
			t.Fatalf("want %v, got %v", err, test.Error.Error())
		}
	}
}
