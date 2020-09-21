package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseReq(t *testing.T) {
	resp, err := ParseReq("")
	assert.EqualError(t, err, ErrEmptyInput.Error())
	assert.Nil(t, resp)

	resp, err = ParseReq("://")
	assert.EqualError(t, err, ErrMissingProto.Error())
	assert.Nil(t, resp)

	resp, err = ParseReq("!!!:://!!~@localhost:90")
	assert.EqualError(t, err, ErrInvalidFormat.Error())
	assert.Nil(t, resp)

	resp, err = ParseReq("scp://")
	assert.EqualError(t, err, "protocol SCP is not supported")
	assert.Nil(t, resp)

	resp, err = ParseReq("http://")
	assert.EqualError(t, err, ErrMissingHost.Error())
	assert.Nil(t, resp)

	resp, err = ParseReq("http://localhost")
	assert.EqualError(t, err, ErrMissingPort.Error())
	assert.Nil(t, resp)

	resp, err = ParseReq("http://localhost:8080/path/ignored")
	assert.NoError(t, err)
	assert.Equal(t, &Req{
		Host:  "LOCALHOST:8080",
		Proto: ProtoHTTP,
	}, resp)
}
