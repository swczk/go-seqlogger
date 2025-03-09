package seqlogger

import (
	"context"
	"encoding/json"
	"log/slog"
	"runtime"
	"time"
)

func createSeqEvent(config Config, ctx context.Context, r slog.Record, baseAttrs []slog.Attr) map[string]interface{} {
	event := map[string]interface{}{
		"@t":  r.Time.UTC().Format(time.RFC3339Nano),
		"@m":  r.Message,
		"@mt": r.Message,
		"@l":  MapLogLevel(r.Level),
	}

	for _, attr := range baseAttrs {
		processAttribute(event, attr)
	}

	r.Attrs(func(a slog.Attr) bool {
		processAttribute(event, a)
		return true
	})

	if config.RequestIDKey != "" {
		if requestID, ok := ctx.Value(config.RequestIDKey).(string); ok {
			event["request_id"] = requestID
		}
	}

	if config.AddSource && r.PC != 0 {
		addSourceInfo(event, r.PC)
	}

	return event
}

func processAttribute(event map[string]interface{}, attr slog.Attr) {
	switch attr.Value.Kind() {
	case slog.KindGroup:
		groupMap := make(map[string]interface{})
		for _, a := range attr.Value.Group() {
			processAttribute(groupMap, a)
		}
		event[attr.Key] = groupMap
	case slog.KindTime:
		event[attr.Key] = attr.Value.Time().UTC().Format(time.RFC3339Nano)
	case slog.KindDuration:
		event[attr.Key] = attr.Value.Duration().Nanoseconds()
	case slog.KindAny:
		processAnyAttribute(event, attr)
	default:
		event[attr.Key] = attr.Value.String()
	}
}

func processAnyAttribute(event map[string]interface{}, attr slog.Attr) {
	switch v := attr.Value.Any().(type) {
	case error:
		event["@x"] = v.Error()
	case time.Time:
		event[attr.Key] = v.UTC().Format(time.RFC3339Nano)
	case time.Duration:
		event[attr.Key] = v.Nanoseconds()
	case []byte:
		parseJSONToEvent(event, attr.Key, v)
	case string:
		parseJSONToEvent(event, attr.Key, []byte(v))
	default:
		event[attr.Key] = v
	}
}

func parseJSONToEvent(event map[string]interface{}, baseKey string, jsonData []byte) {
	var parsedData map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsedData); err != nil {
		event[baseKey] = string(jsonData)
		return
	}

	flattenJSON(event, baseKey, parsedData)
}

func flattenJSON(event map[string]interface{}, baseKey string, data map[string]interface{}) {
	for key, value := range data {
		fullKey := key
		if baseKey != "" {
			fullKey = baseKey + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			flattenJSON(event, fullKey, v)
		case []interface{}:
			event[fullKey] = v
		default:
			event[fullKey] = v
		}
	}
}

func addSourceInfo(event map[string]interface{}, pc uintptr) {
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()
	event["source"] = map[string]interface{}{
		"file":     frame.File,
		"line":     frame.Line,
		"function": frame.Function,
	}
}
