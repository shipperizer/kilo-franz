package logging

import (
	"encoding/json"

	"go.uber.org/zap"
)

// NewLogger creates a new default logger
// it will need to be closed with
// ```
// defer logger.Desugar().Sync()
// ```
// to make sure all has been piped out before terminating
func NewLogger() *zap.SugaredLogger {
	rawJSON := []byte(`{
		"level": "error",
		"encoding": "json",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stdout","stderr"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "severity",
		  "levelEncoder": "lowercase",
		  "timeKey": "@timestamp",
		  "timeEncoder": "rfc3339nano"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
