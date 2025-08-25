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
	"github.com/rs/zerolog"
)

var (
	zerologLevels = map[Level]zerolog.Level{
		LevelTrace:  zerolog.TraceLevel,
		LevelDebug:  zerolog.DebugLevel,
		LevelInfo:   zerolog.InfoLevel,
		LevelWarn:   zerolog.WarnLevel,
		LevelNotice: zerolog.WarnLevel,
		LevelError:  zerolog.ErrorLevel,
		LevelFatal:  zerolog.FatalLevel,
	}

	hlogLevels = map[zerolog.Level]Level{
		zerolog.TraceLevel: LevelTrace,
		zerolog.DebugLevel: LevelDebug,
		zerolog.InfoLevel:  LevelInfo,
		zerolog.WarnLevel:  LevelWarn,
		zerolog.ErrorLevel: LevelError,
		zerolog.FatalLevel: LevelFatal,
	}
)

// matchHlogLevel map hlog.Level to zerolog.Level
func matchHlogLevel(level Level) zerolog.Level {
	zlvl, found := zerologLevels[level]

	if found {
		return zlvl
	}

	return zerolog.WarnLevel // Default level
}
func MatchLevel(lv interface{}) zerolog.Level {
	lve, ok := lv.(Level)
	if !ok {
		return zerolog.WarnLevel
	}
	zlvl, found := zerologLevels[lve]

	if found {
		return zlvl
	}

	return zerolog.WarnLevel // Default level
}

// matchZerologLevel map zerolog.Level to hlog.Level
func matchZerologLevel(level zerolog.Level) Level {
	hlvl, found := hlogLevels[level]

	if found {
		return hlvl
	}

	return LevelWarn // Default level
}
