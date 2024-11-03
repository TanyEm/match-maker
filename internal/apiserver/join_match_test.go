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

func TestJoinMatch(t *testing.T) {
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
			reqURL:              "/match?join_id=72b33e85-e8cd-45e6-89f4-25bfdac584d8",
			expectedError:       false,
			expectedCode:        200,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"match_id":"d1b18698-f7eb-4cb1-b7f2-97e4b44c2c1d"}`,
			expectedMockCalls: func() {
				srv.Lobby.(*lobby.MockLobbier).EXPECT().
					GetMatchMakingTime().
					Times(1)
				srv.Lobby.(*lobby.MockLobbier).EXPECT().
					GetMatchByJoinID("72b33e85-e8cd-45e6-89f4-25bfdac584d8").
					AnyTimes().
					Return("d1b18698-f7eb-4cb1-b7f2-97e4b44c2c1d")
			},
		},
		{
			name:                "not valid request: empty join_id",
			reqURL:              "/match",
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"join_id is required"}`,
			expectedMockCalls:   func() {},
		},
		{
			name:                "not valid request: join_id is not valid UUID",
			reqURL:              "/match?join_id=not-valid-uuid",
			expectedError:       true,
			expectedCode:        400,
			expectedContentType: "application/json; charset=utf-8",
			expectedBody:        `{"error":"join_id is not valid UUID"}`,
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
				var gotMatchResponse MatchResponse
				if err := json.Unmarshal(recorder.Body.Bytes(), &gotMatchResponse); err != nil {
					t.Fatal(err)
				}

				var expectedMatchResponse MatchResponse
				if err := json.Unmarshal([]byte(tt.expectedBody), &expectedMatchResponse); err != nil {
					t.Fatal(err)
				}

				if !cmp.Equal(gotMatchResponse, expectedMatchResponse) {
					t.Errorf("expected match response '%v', got '%v'", expectedMatchResponse, gotMatchResponse)
				}

			} else {
				if recorder.Body.String() != tt.expectedBody {
					t.Errorf("expected body '%s', got '%s'", tt.expectedBody, recorder.Body.String())
				}
			}
		})
	}
}
