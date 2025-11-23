package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/ya-metrics/pkg/database"
	"github.com/yogenyslav/ya-metrics/tests/mocks"
)

func TestHandler_Ping(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		db       func() database.DB
		w        http.ResponseWriter
		r        *http.Request
		wantCode int
	}{
		{
			name: "ping ok",
			db: func() database.DB {
				m := mocks.NewMockDB(ctrl)
				m.EXPECT().Ping(gomock.Any()).Return(nil)
				return m
			},
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "/ping", nil),
			wantCode: http.StatusOK,
		},
		{
			name: "ping db error",
			db: func() database.DB {
				m := mocks.NewMockDB(ctrl)
				m.EXPECT().Ping(gomock.Any()).Return(assert.AnError)
				return m
			},
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "/ping", nil),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(nil, tt.db())
			h.Ping(tt.w, tt.r)

			resp := tt.w.(*httptest.ResponseRecorder)
			assert.Equal(t, tt.wantCode, resp.Code)
		})
	}
}
