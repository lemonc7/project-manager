package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Oudwins/zog"
)

func flattenIssues(flattened map[string][]string) string {
	parts := make([]string, 0, len(flattened))
	for field, issues := range flattened {
		if len(issues) == 0 {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s: %s", field, strings.Join(issues, ", ")))
	}
	return strings.Join(parts, "; ")
}

func SchemaValidate(schema *zog.StructSchema, dstPtr any) error {
	if errs := schema.Validate(dstPtr); errs != nil {
		flattened := zog.Issues.Flatten(errs)
		return errors.New(flattenIssues(flattened))
	}
	return nil
}
