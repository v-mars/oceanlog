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

package main

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMatchHlogLevel(t *testing.T) {
	assert.Equal(t, zerolog.TraceLevel, matchHlogLevel(LevelTrace))
	assert.Equal(t, zerolog.DebugLevel, matchHlogLevel(LevelDebug))
	assert.Equal(t, zerolog.InfoLevel, matchHlogLevel(LevelInfo))
	assert.Equal(t, zerolog.WarnLevel, matchHlogLevel(LevelWarn))
	assert.Equal(t, zerolog.ErrorLevel, matchHlogLevel(LevelError))
	assert.Equal(t, zerolog.FatalLevel, matchHlogLevel(LevelFatal))
}

func TestMatchZerologLevel(t *testing.T) {
	assert.Equal(t, LevelTrace, matchZerologLevel(zerolog.TraceLevel))
	assert.Equal(t, LevelDebug, matchZerologLevel(zerolog.DebugLevel))
	assert.Equal(t, LevelInfo, matchZerologLevel(zerolog.InfoLevel))
	assert.Equal(t, LevelWarn, matchZerologLevel(zerolog.WarnLevel))
	assert.Equal(t, LevelError, matchZerologLevel(zerolog.ErrorLevel))
	assert.Equal(t, LevelFatal, matchZerologLevel(zerolog.FatalLevel))
}
