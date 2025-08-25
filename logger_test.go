/*
 * Copyright 2022 CloudWeGo Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package oceanlog

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestFrom(t *testing.T) {
	b := &bytes.Buffer{}

	zl := zerolog.New(b).With().Str("key", "test").Logger()
	l := From(zl)

	l.Info("foo")

	assert.Equal(
		t,
		`{"level":"info","key":"test","message":"foo"}
`,
		b.String(),
	)
}

func TestGetLogger_notSet(t *testing.T) {
	_, err := GetLogger()

	assert.Error(t, err)
	assert.Equal(t, "GetDefaultLogger is not a zerolog logger", err.Error())
}

func TestGetLogger(t *testing.T) {
	//hlog.SetLogger(New())
	lo, err := GetLogger()

	assert.NoError(t, err)
	assert.IsType(t, DefaultLogger{}, lo)
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()
	l := New()
	c := l.WithContext(ctx)

	assert.NotNil(t, c)
	assert.IsType(t, zerolog.Ctx(c), &zerolog.Logger{})
}

func TestLoggerWithField(t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.SetOutput(b)
	l.WithField("service", "logging")

	l.Info("foobar")

	type Log struct {
		Level   string `json:"level"`
		Service string `json:"service"`
		Message string `json:"message"`
	}

	lo := &Log{}

	err := json.Unmarshal(b.Bytes(), lo)

	println(b.String())
	assert.NoError(t, err)
	assert.Equal(t, "logging", lo.Service)
}

func TestUnwrap(t *testing.T) {
	l := New()

	lo := l.Unwrap()

	assert.NotNil(t, lo)
	assert.IsType(t, zerolog.Logger{}, lo)
}

func TestLog(t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.SetOutput(b)

	l.Trace("foo")
	assert.Equal(
		t,
		`{"level":"debug","message":"foo"}
`,
		b.String(),
	)

	b.Reset()
	l.Debug("foo")
	assert.Equal(
		t,
		`{"level":"debug","message":"foo"}
`,
		b.String(),
	)

	b.Reset()
	l.Info("foo")
	assert.Equal(
		t,
		`{"level":"info","message":"foo"}
`,
		b.String(),
	)

	b.Reset()
	l.Notice("foo")
	assert.Equal(
		t,
		`{"level":"warn","message":"foo"}
`,
		b.String(),
	)

	b.Reset()
	l.Warn("foo")
	assert.Equal(
		t,
		`{"level":"warn","message":"foo"}
`,
		b.String(),
	)

	b.Reset()
	l.Error("foo")
	assert.Equal(
		t,
		`{"level":"error","message":"foo"}
`,
		b.String(),
	)
}

func TestLogf(t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.SetOutput(b)

	l.Tracef("foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"debug","message":"foobar"}
`,
		b.String(),
	)

	b.Reset()
	l.Debugf("foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"debug","message":"foobar"}
`,
		b.String(),
	)

	b.Reset()
	l.Infof("foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"info","message":"foobar"}
`,
		b.String(),
	)

	b.Reset()
	l.Noticef("foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"warn","message":"foobar"}
`,
		b.String(),
	)

	b.Reset()
	l.Warnf("foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"warn","message":"foobar"}
`,
		b.String(),
	)

	b.Reset()
	l.Errorf("foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"error","message":"foobar"}
`,
		b.String(),
	)
}

func TestCtxTracef(t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.SetOutput(b)
	ctx := l.log.WithContext(context.Background())

	l.CtxTracef(ctx, "foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"debug","message":"foobar"}
`,
		b.String(),
	)
	assert.NotNil(t, log.Ctx(ctx))
}

func TestCtxDebugf(t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.SetOutput(b)
	ctx := l.log.WithContext(context.Background())

	l.CtxDebugf(ctx, "foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"debug","message":"foobar"}
`,
		b.String(),
	)
	assert.NotNil(t, log.Ctx(ctx))
}

func TestCtxInfof(t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.SetOutput(b)
	ctx := l.log.WithContext(context.Background())

	l.CtxInfof(ctx, "foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"info","message":"foobar"}
`,
		b.String(),
	)
	assert.NotNil(t, log.Ctx(ctx))
}

func TestCtxNoticef(t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.SetOutput(b)
	ctx := l.log.WithContext(context.Background())

	l.CtxNoticef(ctx, "foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"warn","message":"foobar"}
`,
		b.String(),
	)
	assert.NotNil(t, log.Ctx(ctx))
}

func TestCtxWarnf(t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.SetOutput(b)
	ctx := l.log.WithContext(context.Background())

	l.CtxWarnf(ctx, "foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"warn","message":"foobar"}
`,
		b.String(),
	)
	assert.NotNil(t, log.Ctx(ctx))
}

func TestCtxErrorf(t *testing.T) {
	b := &bytes.Buffer{}
	l := New()
	l.SetOutput(b)
	ctx := l.log.WithContext(context.Background())

	l.CtxErrorf(ctx, "foo%s", "bar")
	assert.Equal(
		t,
		`{"level":"error","message":"foobar"}
`,
		b.String(),
	)
	assert.NotNil(t, log.Ctx(ctx))
}

func TestSetLevel(t *testing.T) {
	l := New()

	l.SetLevel(LevelDebug)
	assert.Equal(t, l.log.GetLevel(), zerolog.DebugLevel)

	l.SetLevel(LevelDebug)
	assert.Equal(t, l.log.GetLevel(), zerolog.DebugLevel)

	l.SetLevel(LevelError)
	assert.Equal(t, l.log.GetLevel(), zerolog.ErrorLevel)
}

// TestNewConsole_Integration 测试整个 ConsoleWriter 的集成行为
// TestNewConsole_FormatCaller_WrongType 测试 FormatCaller 函数处理错误类型输入
func TestNewConsole_FormatCaller_WrongType(t *testing.T) {
	// Arrange: 准备测试数据
	buffer := &bytes.Buffer{}
	console := NewConsole(buffer)

	testCases := []interface{}{
		123,        // 整数
		[]string{}, // 切片
		struct{}{}, // 结构体
		true,       // 布尔值
	}

	// Act & Assert: 对每个测试用例执行测试
	for _, tc := range testCases {
		result := console.FormatCaller(tc)
		if result != "" {
			t.Errorf("FormatCaller(%v): expected empty string for wrong type, but got %s", tc, result)
		}
	}

}

// TestNewConsole_FormatCaller_Normal 测试 FormatCaller 函数的正常情况
func TestNewConsole_FormatCaller_Normal(t *testing.T) {
	// Arrange: 准备测试数据
	buffer := &bytes.Buffer{}
	console := NewConsole(buffer)

	testCases := []struct {
		input    interface{}
		expected string
	}{
		{"/path/to/logger_test.go:76", "logger_test.go:"},
		{"main.go:123", "main.go:"},
		{"C:\\windows\\path\\test.go:456", "test.go:"}, // Windows 路径测试
	}

	// Act & Assert: 对每个测试用例执行测试
	for _, tc := range testCases {
		result := console.FormatCaller(tc.input)
		if result != tc.expected {
			t.Errorf("FormatCaller(%v): expected %s, but got %s", tc.input, tc.expected, result)
		}
	}
}
