package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/johnmackenzie91/commonlogger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type stubEndpoint string

func (s stubEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s))
}

func Test_NewServer(t *testing.T) {
	l := commonlogger.New(logrus.New(), commonlogger.Config{})

	endpoint := stubEndpoint("hello world")

	sut := NewHTTP(WithAddress("localhost:0"), WithHandler(endpoint), WithLogger(l))

	var stubCtx, cancel = context.WithCancel(context.Background())

	go func() {
		err := sut.Run(stubCtx)
		assert.Nil(t, err)
	}()

	time.Sleep(2 * time.Second)

	res, err := http.Get(sut.URL())

	assert.Nil(t, err)

	b, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Nil(t, res.Body.Close())

	assert.Equal(t, "hello world", string(b))

	cancel()

	assert.Nil(t, err)

	res, err = http.Get(sut.URL())
	assert.NotNil(t, err)
}
