package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go_challenge/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func addAuthorization(t *testing.T, r *http.Request, tokenMaker token.Maker, authorizationType string, username string, duration time.Duration) {

	accessToken, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, accessToken)
	r.Header.Add(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		Name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			Name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//Create Token And Add it to Header
				addAuthorization(t, request, tokenMaker, authorizationHeaderType, "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//Check the Token if its valid
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			Name: "No Authorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//Create Token And Add it to Header
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//Check the Token if its valid
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			Name: "Unsupported Authorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//Create Token And Add it to Header
				addAuthorization(t, request, tokenMaker, "unsupported type", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//Check the Token if its valid
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			Name: "Invalid Authorization Format",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//Create Token And Add it to Header
				addAuthorization(t, request, tokenMaker, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//Check the Token if its valid
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			Name: "Expired Token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//Create Token And Add it to Header + Send Negative Format !
				addAuthorization(t, request, tokenMaker, authorizationHeaderType, "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				//Check the Token if its valid
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.Name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
