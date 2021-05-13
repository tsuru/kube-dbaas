package web

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ajg/form"
	echo "github.com/labstack/echo/v4"
	"github.com/tsuru/kube-dbaas/engine"
	"github.com/tsuru/kube-dbaas/engine/mongo"
	"github.com/tsuru/kube-dbaas/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (a *api) serviceCreate(c echo.Context) error {
	ctx := c.Request().Context()
	args := types.CreateArgs{
		// NOTE: using a different decoder for Parameters since the `r.PostForm()`
		// method does not understand the format used by github.com/ajf.form.
		Parameters: decodeFormParameters(c.Request()),
	}
	if err := c.Bind(&args); err != nil {
		return err
	}

	log.Printf("parameters: %#v", args)

	if args.Plan == "" {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Plan cannot be empty")
	}

	e, err := NewEngineFromPlan(a.Client, args.Plan)
	if err != nil {
		return err
	}

	err = e.CreateInstance(ctx, &args)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func decodeFormParameters(r *http.Request) map[string]interface{} {
	if r == nil {
		return nil
	}

	body := r.Body
	defer body.Close()

	var buffer bytes.Buffer
	reader := io.TeeReader(body, &buffer)

	var obj struct {
		Parameters map[string]interface{} `form:"parameters"`
	}
	newFormDecoder(reader).Decode(&obj)
	r.Body = ioutil.NopCloser(&buffer)
	return obj.Parameters
}

func newFormDecoder(r io.Reader) *form.Decoder {
	decoder := form.NewDecoder(r)
	decoder.IgnoreCase(true)
	decoder.IgnoreUnknownKeys(true)
	return decoder
}

func NewEngineFromPlan(cli client.Client, plan string) (engine.Engine, error) {
	e := mongo.New(cli)
	return e, nil
}
