package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getsentry/sentry-go"

	"github.com/edalford11/true-markets/config"
)

func InitSentry(environment config.Environment) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://34c2c760f8414c61b42c0a6227ca958e@o4504202057482240.ingest.sentry.io/4504215372300288",
		TracesSampleRate: 0.0, // Don't send traces
		Environment:      string(environment),
		AttachStacktrace: true,
	}); err != nil {
		panic(err)
	}
}

// FindNearestWdParentFolder traverses the working directory upwards until it finds the folderName and returns the absolute path.
func FindNearestWdParentFolder(folderName string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wdParts := strings.Split(wd, string(filepath.Separator))
	for i := len(wdParts) - 1; i > 0; i-- {
		if wdParts[i] == folderName {
			return string(filepath.Separator) + filepath.Join(wdParts[:i+1]...) + string(filepath.Separator)
		}
	}

	panic(fmt.Sprintf("Could not find %s in %s", folderName, wd))
}
