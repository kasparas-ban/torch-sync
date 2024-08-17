package pkg

import (
	"log/slog"
	"os"
)

func GetLogTextHandler() *slog.TextHandler {
	th := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}

			return a
		},
	})

	return th
}
