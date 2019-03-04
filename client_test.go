package arukas

import (
	"encoding/json"
	"errors"
	"testing"

	"context"
	"time"

	"github.com/stretchr/testify/assert"
)

type testHTTPAPI struct {
	getResult   []byte
	patchResult []byte
	postResult  []byte
	putResult   []byte

	getError    error
	patchError  error
	postError   error
	putError    error
	deleteError error
}

func (c *testHTTPAPI) get(path string) ([]byte, error) {
	return c.getResult, c.getError
}
func (c *testHTTPAPI) patch(path string, body interface{}) ([]byte, error) {
	return c.patchResult, c.patchError
}

func (c *testHTTPAPI) post(path string, body interface{}) ([]byte, error) {
	return c.postResult, c.postError
}

func (c *testHTTPAPI) put(path string, body interface{}) ([]byte, error) {
	return c.putResult, c.putError
}

func (c *testHTTPAPI) delete(path string) error {
	return c.deleteError
}

var validCreateAppParam = &RequestParam{
	Name:      "foobar",
	Image:     "nginx:latest",
	Instances: 1,
	Ports: Ports{
		{
			Protocol: "tcp",
			Number:   80,
		},
	},
	Plan: PlanFree,
}

var validUpdateServiceParam = &RequestParam{
	Image:     "httpd:latest",
	Instances: 1,
}

var testServiceID = "01BEF829-72E4-48F9-81DA-E3B41A1EDAC9" // valid UUID

func TestCreateApp(t *testing.T) {
	t.Run("Invalid parameter", func(t *testing.T) {
		c := &client{
			httpAPI: &testHTTPAPI{},
		}

		p := &RequestParam{}

		res, err := c.CreateApp(p)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("POST /apps returns error", func(t *testing.T) {
		expectError := errors.New("dummy")
		c := &client{
			httpAPI: &testHTTPAPI{
				postError: expectError,
			},
		}

		res, err := c.CreateApp(validCreateAppParam)
		assert.Error(t, err)
		assert.Equal(t, expectError, err)
		assert.Nil(t, res)
	})

	t.Run("POST /apps returns invalid JSON", func(t *testing.T) {
		c := &client{
			httpAPI: &testHTTPAPI{
				postResult: []byte(`{ "invalid": `), // JSON missed right }
			},
		}

		res, err := c.CreateApp(validCreateAppParam)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("POST /apps succeed", func(t *testing.T) {
		createdApp := &AppData{
			Data: &App{
				ID:   testServiceID,
				Type: TypeApps,
			},
		}
		data, err := json.Marshal(createdApp)
		if err != nil {
			t.Fatal(err)
		}

		c := &client{
			httpAPI: &testHTTPAPI{
				postResult: data,
			},
		}

		res, err := c.CreateApp(validCreateAppParam)
		assert.NoError(t, err)
		assert.Equal(t, createdApp, res)
	})
}

func TestWaitForStatus(t *testing.T) {
	getServiceData := func(status string) []byte {
		service := &ServiceData{
			Data: &Service{
				ID:   testServiceID,
				Type: TypeServices,
				Attributes: &ServiceAttr{
					Status: status,
				},
			},
		}
		data, err := json.Marshal(service)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}

	t.Run("Should block until returns nil", func(t *testing.T) {
		data := getServiceData(StatusStopped)
		c := &client{
			httpAPI: &testHTTPAPI{
				getResult: data,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.WaitForState(ctx, testServiceID, StatusStopped)
		assert.NoError(t, err)
	})
	t.Run("Should block until returns error", func(t *testing.T) {
		dummyErr := errors.New("dummy")
		c := &client{
			httpAPI: &testHTTPAPI{
				getError: dummyErr,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.WaitForState(ctx, testServiceID, StatusStopped)
		assert.Equal(t, dummyErr, err)
	})

	t.Run("Should block until timeout", func(t *testing.T) {
		data := getServiceData(StatusStopped)
		c := &client{
			httpAPI: &testHTTPAPI{
				getResult: data,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.WaitForState(ctx, testServiceID, StatusRunning)
		assert.Equal(t, ctx.Err(), err)
	})

}

func TestAccAppCRUD(t *testing.T) {
	if !isAccTest() {
		t.SkipNow()
	}

	defer initialize()()

	c := &client{
		httpAPI: realHTTPClient,
	}

	var id string

	t.Run("Create", func(t *testing.T) {
		res, err := c.CreateApp(validCreateAppParam)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("List", func(ti *testing.T) {
		list, err := c.ListApps()
		assert.NoError(t, err)
		assert.NotNil(t, list)
		if list == nil || len(list.Data) == 0 {
			t.Fatal(err)
		}
		assert.True(t, len(list.Data) > 0)
		id = list.Data[0].ID
	})

	t.Run("Read", func(t *testing.T) {
		app, err := c.ReadApp(id)
		assert.NoError(t, err)
		assert.NotNil(t, app)
		assert.Equal(t, id, app.Data.ID)
	})

	t.Run("Delete", func(t *testing.T) {
		// DELETE
		err := c.DeleteApp(id)
		assert.NoError(t, err)

		list, err := c.ListApps()
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.True(t, len(list.Data) == 0)
	})
}

func TestAccResourceNotFound(t *testing.T) {
	if !isAccTest() {
		t.SkipNow()
	}

	defer initialize()()
	c := &client{
		httpAPI: realHTTPClient,
	}

	service, err := c.ReadService(testServiceID)
	assert.Error(t, err)
	_, ok := err.(ErrorNotFound)
	assert.True(t, ok)
	assert.Nil(t, service)
}

func TestAccServiceCRUD(t *testing.T) {
	if !isAccTest() {
		t.SkipNow()
	}

	defer initialize()()

	c := &client{
		httpAPI: realHTTPClient,
	}

	var id string

	t.Run("Create(App)", func(t *testing.T) {
		res, err := c.CreateApp(validCreateAppParam)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("List", func(t *testing.T) {
		list, err := c.ListServices()
		assert.NoError(t, err)
		assert.NotNil(t, list)
		assert.True(t, len(list.Data) > 0)

		id = list.Data[0].ID
	})

	t.Run("Read", func(t *testing.T) {
		s, err := c.ReadService(id)
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.Equal(t, id, s.Data.ID)
	})

	t.Run("Update", func(t *testing.T) {

		res, err := c.UpdateService(id, validUpdateServiceParam)

		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, int32(1), res.Data.Attributes.Instances)
		assert.Equal(t, "httpd:latest", res.Data.Attributes.Image)
	})

	t.Run("PowerOn", func(t *testing.T) {
		err := c.PowerOn(id)
		assert.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		errChan := make(chan error)

		go func() {
			errChan <- c.WaitForState(ctx, id, StatusRunning)
		}()

		select {
		case e := <-errChan:
			assert.NoError(t, e)
		case <-ctx.Done():
			assert.Fail(t, ctx.Err().Error())
		}
	})

	t.Run("PowerOff", func(t *testing.T) {
		err := c.PowerOff(id)
		assert.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		errChan := make(chan error)

		go func() {
			errChan <- c.WaitForState(ctx, id, StatusStopped)
		}()

		select {
		case e := <-errChan:
			assert.NoError(t, e)
		case <-ctx.Done():
			assert.Fail(t, ctx.Err().Error())
		}

	})

}

func initialize() func() {
	cleanup()
	return cleanup
}

func cleanup() {
	c := &client{
		httpAPI: realHTTPClient,
	}
	list, err := c.ListApps()
	if err != nil {
		return
	}
	for _, app := range list.Data {
		id := app.ID
		c.DeleteApp(id)
	}
}
