package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/techschool/simplebank/db/mock"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
	"github.com/techschool/simplebank/worker"
)

func TestGetAccountApi(t *testing.T) {
	user, _ := RandomUser(t)
	account := RandomAccount(user.Username)

	taskDisibutor := worker.NewRedisTaskDistributor(&asynq.RedisClientOpt{})
	testCase := []struct {
		name          string
		ID            int64
		setupAuth     func(t *testing.T, request *http.Request, maker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "TEST OK",
			ID:   account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, maker, request, authorizationBearerType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkMatchBodyRequest(t, recorder.Body, account)
			},
		},
		{
			name: "NOT FOUND",
			ID:   account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, maker, request, authorizationBearerType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "INTERNAL SV ERROR",
			ID:   account.ID,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, maker, request, authorizationBearerType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BADD REQUEST",
			ID:   0,
			setupAuth: func(t *testing.T, request *http.Request, maker token.Maker) {
				addAuthorization(t, maker, request, authorizationBearerType, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			//create stubs
			tc.buildStubs(store)
			server := newTestServer(t, store, taskDisibutor)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", tc.ID)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})

	}

}

func checkMatchBodyRequest(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func RandomAccount(username string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Currency: util.RandomCurrency(),
		Owner:    username,
		Balance:  util.RandomMoney(),
	}
}
