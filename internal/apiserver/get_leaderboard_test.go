package apiserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TanyEm/match-maker/v2/internal/lobby"
	"github.com/TanyEm/match-maker/v2/internal/match"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
)

func TestGetLeaderBoard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := NewAPIServer(lobby.NewMockLobbier(ctrl), match.NewMockKeeper(ctrl))

	tests := []struct {
		name                string
		reqURL              string
		expectedError       bool
		expectedCode        int
		expectedContentType string
		expectedBody        string
		expectedMockCalls   func()
	}{
		{
			name:                "valid request",
			reqURL:              "/leaderboard?match_id=72b33e85-e8cd-45e6-89f4-25bfdac584d8",
			expectedError:       false,
			expectedCode:        200,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody: `{
								"match_id":"110dc29f-dcd7-4bee-abea-2f7b24e47777", 
								"players":[
									{"player_id":"player2","level":2,"country":"USA","score":200},
									{"player_id":"player1","level":1,"country":"USA","score":100}
								]
							}`,
			expectedMockCalls: func() {
				srv.MatchKeeper.(*match.MockKeeper).EXPECT().
					GetLeaderBoard("72b33e85-e8cd-45e6-89f4-25bfdac584d8").
					Times(1).
					Return(&match.LeaderBoard{
						MatchID: "110dc29f-dcd7-4bee-abea-2f7b24e47777",
						Players: []match.PlayerInfo{
							{PlayerID: "player2", Level: 2, Country: "USA", Score: 200},
							{PlayerID: "player1", Level: 1, Country: "USA", Score: 100},
						},
					})
			},
		},
		{
			name:                "not valid request: empty match_id",
			reqURL:              "/leaderboard",
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"match_id is required"}`,
			expectedMockCalls:   func() {},
		},
		{
			name:                "not valid request: match_id is not valid UUID",
			reqURL:              "/leaderboard?match_id=not-valid-uuid",
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"match_id is not valid UUID"}`,
			expectedMockCalls:   func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectedMockCalls()
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tt.reqURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			srv.GinEngine.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedCode {
				t.Errorf("expected code %d, got %d", tt.expectedCode, recorder.Code)
			}

			if recorder.Header().Get("Content-Type") != tt.expectedContentType {
				t.Errorf("expected content type %s, got %s", tt.expectedContentType, recorder.Header().Get("Content-Type"))
			}

			if !tt.expectedError {
				var gotLeaderboardResponse GetLeaderBoardResponse
				if err := json.Unmarshal(recorder.Body.Bytes(), &gotLeaderboardResponse); err != nil {
					t.Fatal(err)
				}

				var expectedLeaderboardResponse GetLeaderBoardResponse
				if err := json.Unmarshal([]byte(tt.expectedBody), &expectedLeaderboardResponse); err != nil {
					t.Fatal(err)
				}

				if !cmp.Equal(gotLeaderboardResponse, expectedLeaderboardResponse) {
					t.Errorf("expected leaderboard response '%v', got '%v'", expectedLeaderboardResponse, gotLeaderboardResponse)
				}

			} else {
				if recorder.Body.String() != tt.expectedBody {
					t.Errorf("expected body '%s', got '%s'", tt.expectedBody, recorder.Body.String())
				}
			}
		})
	}
}
