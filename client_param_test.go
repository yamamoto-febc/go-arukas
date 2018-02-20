package arukas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientParam_Valiate(t *testing.T) {

	expects := []struct {
		scenario string
		expect   bool
		param    *ClientParam
	}{
		{
			scenario: "Token is empty",
			expect:   false,
			param:    &ClientParam{},
		},
		{
			scenario: "Secret is empty",
			expect:   false,
			param: &ClientParam{
				Token: "foo",
			},
		},
		{
			scenario: "Valid parameters",
			expect:   true,
			param: &ClientParam{
				Token:  "foo",
				Secret: "bar",
			},
		},
	}

	for _, expect := range expects {
		t.Run(expect.scenario, func(t *testing.T) {
			err := expect.param.validate()
			assert.Equal(t, expect.expect, err == nil)
		})
	}

}
