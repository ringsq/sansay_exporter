package main

import (
	"reflect"
	"sync"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/jarcoal/httpmock"
)

func TestScrapeTarget(t *testing.T) {
	var wg sync.WaitGroup
	testCollector := collector{target: "http://localhost:8888/", username: "user", password: "pass", logger: log.NewNopLogger()}

	type args struct {
		c      collector
		path   string
		result chan interface{}
		wg     *sync.WaitGroup
	}
	tests := []struct {
		name     string
		args     func(t *testing.T) args
		httpMock httpmock.Responder
		wantErr  bool
	}{
		{ //TODO: Add test cases
			name: "Test that errors are returned",
			args: func(t *testing.T) args {
				return args{
					c:      testCollector,
					path:   "",
					result: make(chan interface{}),
					wg:     &wg,
				}
			},
			httpMock: httpmock.NewXmlResponderOrPanic(200, "<xml>response</xml>"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", "http://localhost:8888/", tt.httpMock)
			wg.Add(1)
			go ScrapeTarget(tArgs.c, tArgs.path, tArgs.result, tArgs.wg)
			result := <-tArgs.result
			switch result.(type) {
			case error:
				if !tt.wantErr {
					t.Errorf("Expected an error, received %s", reflect.TypeOf(result))
				}
			}
		})
	}
}
