package log

import "go.uber.org/zap"

const (
	CtxField        = "span_context"
	ResourceField   = "resource"
	AttributesField = "attributes"
)

var emptyField zap.Field
