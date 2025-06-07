package log

import (
	"context"
	"testing"

	"github.com/welltop-cn/common/cloud/config"
)

func Test_logger(t *testing.T) {
	config.InitConfig()
	InitLogger()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "ttttttttt")

	CtxErrorf(ctx, "debug format %s", "this is a debug log")
}
