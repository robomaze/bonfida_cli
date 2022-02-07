package common_test

import (
	"testing"

	"github.com/robomaze/bonfida_cli/api/common"
	"github.com/stretchr/testify/assert"
)

func TestErrReqParam_Error(t *testing.T) {
	err := &common.ErrReqParam{
		Name: "testParam",
		Msg:  "sneaky error",
	}

	assert.EqualError(t, err, "param name 'testParam' with error 'sneaky error'")
}
