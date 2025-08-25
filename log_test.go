package oceanlog

import (
	"context"
	"io"
	"testing"
)

// MockLumberjackLogger 模拟 lumberjack logger
type MockLumberjackLogger struct {
	io.Writer
}

func TestNew(t *testing.T) {
	a := InitOceanLog("test.log", "console", LevelDebug)
	a.CtxDebugf(context.Background(), "test")
}
