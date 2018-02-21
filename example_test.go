package arukas_test

import (
	"context"
	"time"

	"github.com/yamamoto-febc/go-arukas"
)

func Example() {

	// Initialize API Client
	token := "<your-api-token>"
	secret := "<your-api-secret>"

	client, err := arukas.NewClient(&arukas.ClientParam{
		Token:  token,
		Secret: secret,
	})
	if err != nil {
		panic(err)
	}

	// Create container(app and service)
	app, err := client.CreateApp(&arukas.RequestParam{
		Name:      "<your-app-name>",
		Image:     "nginx:latest",  // supports only DockerHub public repositories
		Plan:      arukas.PlanFree, // allows [PlanFree / PlanHobby / PlanStandard1 / PlanStandard2]
		Instances: 1,               // number of instance count
		SubDomain: "",              // subdomain(EndPoint will be <subdomain>.arukascloud.io). if empty, using random domain.
		Ports: []*arukas.Port{
			{
				Protocol: "tcp",
				Number:   80,
			},
		},
		Environment: []*arukas.Env{
			{
				Key:   "FOO",
				Value: "BAR",
			},
		},
		Command: "nginx -g daemon off;",
	})

	if err != nil {
		panic(err)
	}

	// Get IDs
	appID := app.AppID()
	serviceID := app.ServiceID()

	// Power on
	err = client.PowerOn(serviceID)
	if err != nil {
		panic(err)
	}

	// Wait until container is running...
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	errChan := make(chan error)
	go func() {
		errChan <- client.WaitForState(ctx, serviceID, arukas.StatusRunning)
	}()

	select {
	case e := <-errChan:
		if e != nil {
			panic(e)
		}
	case <-ctx.Done():
		panic(ctx.Err())
	}

	// Update container(service)
	_, err = client.UpdateService(serviceID, &arukas.RequestParam{
		Image: "httpd:latest",
	})
	if err != nil {
		panic(err)
	}

	// Delete container(app)
	err = client.DeleteApp(appID)
	if err != nil {
		panic(err)
	}

}
