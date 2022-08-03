package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTest3(t *testing.T) {
	testCase := []struct {
		name string
		json string
		want string
	}{

		{
			name: "one",
			json: `[
				{
					"a": "12",
					"b": "12",
					"key": "x"
				},
				{
					"a": "12",
					"b": "12",
					"key": "y"
				}
			]`,
			want: `{"x":144,"y":144}`,
		},
		{
			name: "two",
			json: `[
				{
					"a": "13",
					"b": "13",
					"key": "x"
				},
				{
					"a": "13",
					"b": "13",
					"key": "y"
				}
			]`,
			want: `{"x":169,"y":169}`,
		},
		{
			name: "three",
			json: `[
				{
					"a": "55",
					"b": "52",
					"key": "ku"
				},
				{
					"a": "g",
					"b": "5",
					"key": "56"
				}
			]`,
			want: `{"x":2860,"y":0}`,
		},
	}

	handler := http.HandlerFunc(Test3)
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "http://localhost:8000/test3", bytes.NewBuffer([]byte(tc.json)))
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tc.want, rec.Body.String())

		})
	}

}
