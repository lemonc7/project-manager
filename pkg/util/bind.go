package util

import (
	"io"
	"net/http"

	"github.com/Oudwins/zog"
	"github.com/Oudwins/zog/parsers/zjson"
	"github.com/labstack/echo/v4"
)

func BindRequest(r io.Reader, schema *zog.StructSchema, reqPtr any) error {
	errs := schema.Parse(zjson.Decode(r), reqPtr)
	if errs != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			zog.Issues.Flatten(errs),
		)
	}
	return nil
}
