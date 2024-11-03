package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TanyEm/match-maker/v2/internal/lobby"
	"github.com/TanyEm/match-maker/v2/internal/match"
	"github.com/TanyEm/match-maker/v2/internal/player"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestJoinLobby(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := NewAPIServer(lobby.NewMockLobbier(ctrl), match.NewMockKeeper(ctrl))

	tests := []struct {
		name                string
		req                 []byte
		expectedError       bool
		expectedCode        int
		expectedContentType string
		expectedBody        string
		expectedMockCalls   func()
	}{
		{
			name:                "valid request",
			req:                 []byte(`{"player_id": "player1", "level": 1, "country": "USA"}`),
			expectedError:       false,
			expectedCode:        200,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"join_id":"00000000-0000-0000-0000-000000000000"}`,
			expectedMockCalls: func() {
				p := player.Player{
					PlayerID: "player1",
					Level:    1,
					Country:  "USA",
					JoinID:   gomock.Any().String(),
				}
				srv.Lobby.(*lobby.MockLobbier).EXPECT().AddPlayer(player.EqPlayer(p)).Times(1)
			},
		},
		{
			name:                "not valid request: empty player_id",
			req:                 []byte(`{"player_id": "", "level": 1, "country": "USA"}`),
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"Key: 'LobbyRequest.PlayerID' Error:Field validation for 'PlayerID' failed on the 'required' tag"}`,
			expectedMockCalls:   func() {},
		},
		{
			name:                "not valid request: incorrect level negative number",
			req:                 []byte(`{"player_id": "player1", "level": -1, "country": "USA"}`),
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"Key: 'LobbyRequest.Level' Error:Field validation for 'Level' failed on the 'min' tag"}`,
			expectedMockCalls:   func() {},
		},
		{
			name:                "not valid request: incorrect level zero",
			req:                 []byte(`{"player_id": "player1", "level": 0, "country": "USA"}`),
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"Key: 'LobbyRequest.Level' Error:Field validation for 'Level' failed on the 'min' tag"}`,
			expectedMockCalls:   func() {},
		},
		{
			name:                "not valid request: incorrect level wrong type",
			req:                 []byte(`{"player_id": "player1", "level": "1", "country": "USA"}`),
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"json: cannot unmarshal string into Go struct field LobbyRequest.level of type int"}`,
			expectedMockCalls:   func() {},
		},
		{
			name:                "not valid request: incorrect country code",
			req:                 []byte(`{"player_id": "player1", "level": 1, "country": "US"}`),
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"Key: 'LobbyRequest.Country' Error:Field validation for 'Country' failed on the 'isocountry' tag"}`,
			expectedMockCalls:   func() {},
		},
		{
			name:                "not valid request: invalid json",
			req:                 []byte(`{"player_id": "player1",`),
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"unexpected EOF"}`,
			expectedMockCalls:   func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectedMockCalls()
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/lobby", bytes.NewReader(tt.req))
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Content-Type", "application/json")

			srv.GinEngine.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, recorder.Code)
			}

			if recorder.Header().Get("Content-Type") != tt.expectedContentType {
				t.Errorf("expected content type %s, got %s", tt.expectedContentType, recorder.Header().Get("Content-Type"))
			}

			if !tt.expectedError {
				var lobbyResponse LobbyResponse
				if err := json.Unmarshal(recorder.Body.Bytes(), &lobbyResponse); err != nil {
					t.Fatal(err)
				}

				if _, err := uuid.Parse(lobbyResponse.JoinID); err != nil {
					t.Errorf("join_id is not valid UUID: '%s'", lobbyResponse.JoinID)
				}
			} else {
				if recorder.Body.String() != tt.expectedBody {
					t.Errorf("expected body '%s', got '%s'", tt.expectedBody, recorder.Body.String())
				}
			}
		})
	}
}
