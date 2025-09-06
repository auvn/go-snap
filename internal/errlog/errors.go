package errlog

import "log/slog"

func SwallowError(err error, msg string) {
	if err == nil {
		return
	}

	slog.Warn(msg, "error", err)
}
