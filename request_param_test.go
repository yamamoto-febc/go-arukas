package arukas

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestParam_ValidateForCreate(t *testing.T) {

	expects := []struct {
		expect   bool
		scenario string
		param    *RequestParam
	}{
		{
			expect:   false,
			scenario: "Required values are empty",
			param:    &RequestParam{},
		},
		{
			expect:   true,
			scenario: "Valid minimum values",
			param: &RequestParam{
				Name:      "foobar",
				Image:     "foobar:latest",
				Instances: 1,
				Ports: Ports{
					{
						Protocol: "tcp",
						Number:   80,
					},
				},
				Plan: PlanFree,
			},
		},
		{
			expect:   false,
			scenario: "Invalid protocol",
			param: &RequestParam{
				Name:      "foobar",
				Image:     "foobar:latest",
				Instances: 1,
				Ports: Ports{
					{
						Protocol: "ssh",
						Number:   80,
					},
				},
				Plan: PlanFree,
			},
		},
		{
			expect:   false,
			scenario: "Invalid port number",
			param: &RequestParam{
				Name:      "foobar",
				Image:     "foobar:latest",
				Instances: 1,
				Ports: Ports{
					{
						Protocol: "ssh",
						Number:   65536,
					},
				},
				Plan: PlanFree,
			},
		},
		{
			expect:   false,
			scenario: "Invalid Region",
			param: &RequestParam{
				Name:      "foobar",
				Image:     "foobar:latest",
				Instances: 1,
				Ports: Ports{
					{
						Protocol: "tcp",
						Number:   80,
					},
				},
				Plan:   PlanFree,
				Region: "foobar",
			},
		},
		{
			expect:   false,
			scenario: "Invalid Plan",
			param: &RequestParam{
				Name:      "foobar",
				Image:     "foobar:latest",
				Instances: 1,
				Ports: Ports{
					{
						Protocol: "tcp",
						Number:   80,
					},
				},
				Plan: "foobar",
			},
		},
	}

	for _, expect := range expects {
		t.Run(expect.scenario, func(t *testing.T) {
			err := expect.param.ValidateForCreate()
			assert.Equal(t, expect.expect, err == nil)
		})
	}

}

func TestRequestParam_ValidateForUpdate(t *testing.T) {

	expects := []struct {
		expect   bool
		scenario string
		param    *RequestParam
	}{
		{
			expect:   false,
			scenario: "Required values are empty",
			param:    &RequestParam{},
		},
		{
			expect:   true,
			scenario: "Valid minimum values",
			param: &RequestParam{
				Image:     "foobar:latest",
				Instances: 1,
			},
		},
	}

	for _, expect := range expects {
		t.Run(expect.scenario, func(t *testing.T) {
			err := expect.param.ValidateForUpdate()
			assert.Equal(t, expect.expect, err == nil)
		})
	}

}
