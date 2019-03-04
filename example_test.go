package arukas_test

import (
	"context"
	"time"

	"fmt"

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

func Example_listApps() {

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

	// List all apps
	apps, err := client.ListApps()
	if err != nil {
		panic(err)
	}
	for _, app := range apps.Data {
		fmt.Printf("AppName: %s", app.Name())
	}

}

func Example_createContainer() {
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

	fmt.Printf("AppName: %s", app.Name())
}

func Example_readApp() {
	// Initialize API Client
	token := "<your-api-token>"
	secret := "<your-api-secret>"
	appID := "<your-app-id>" // UUID format is required for appID

	client, err := arukas.NewClient(&arukas.ClientParam{
		Token:  token,
		Secret: secret,
	})
	if err != nil {
		panic(err)
	}

	// Read app
	app, err := client.ReadApp(appID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("AppName: %s", app.Name())
}

func Example_deleteApp() {
	// Initialize API Client
	token := "<your-api-token>"
	secret := "<your-api-secret>"
	appID := "<your-app-id>" // UUID format is required for appID

	client, err := arukas.NewClient(&arukas.ClientParam{
		Token:  token,
		Secret: secret,
	})
	if err != nil {
		panic(err)
	}

	// Delete app
	err = client.DeleteApp(appID)
	if err != nil {
		panic(err)
	}
}

func Example_listServices() {
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

	// List all services
	services, err := client.ListServices()
	if err != nil {
		panic(err)
	}

	for _, service := range services.Data {
		fmt.Printf("ServiceID: %s\nEndPoint URL: %s", service.ID, service.EndPoint())
	}
}

func Example_updateService() {
	// Initialize API Client
	token := "<your-api-token>"
	secret := "<your-api-secret>"
	serviceID := "<your-service-id>" // UUID format is required for serviceID

	client, err := arukas.NewClient(&arukas.ClientParam{
		Token:  token,
		Secret: secret,
	})
	if err != nil {
		panic(err)
	}

	// Update service
	updatedService, err := client.UpdateService(serviceID, &arukas.RequestParam{
		Instances: 2,
		Image:     "nginx:latst",
		Ports: []*arukas.Port{
			{
				Protocol: "tcp",
				Number:   80,
			},
			{
				Protocol: "tcp",
				Number:   443,
			},
		},
		Environment: []*arukas.Env{
			{
				Key:   "MY_ENV",
				Value: "my_value",
			},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Services: %#v", updatedService)
}

func Example_powerOn() {
	// Initialize API Client
	token := "<your-api-token>"
	secret := "<your-api-secret>"
	serviceID := "<your-service-id>" // UUID format is required for serviceID

	client, err := arukas.NewClient(&arukas.ClientParam{
		Token:  token,
		Secret: secret,
	})
	if err != nil {
		panic(err)
	}

	// PowerOn
	if err := client.PowerOn(serviceID); err != nil {
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
}

func Example_powerOff() {
	// Initialize API Client
	token := "<your-api-token>"
	secret := "<your-api-secret>"
	serviceID := "<your-service-id>" // UUID format is required for serviceID

	client, err := arukas.NewClient(&arukas.ClientParam{
		Token:  token,
		Secret: secret,
	})
	if err != nil {
		panic(err)
	}

	// PowerOff
	if err := client.PowerOff(serviceID); err != nil {
		panic(err)
	}

	// Wait until container is stopped...
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	errChan := make(chan error)
	go func() {
		errChan <- client.WaitForState(ctx, serviceID, arukas.StatusStopped)
	}()

	select {
	case e := <-errChan:
		if e != nil {
			panic(e)
		}
	case <-ctx.Done():
		panic(ctx.Err())
	}
}
