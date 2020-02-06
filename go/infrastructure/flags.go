package infrastructure

import (
	"flag"
)

var (
   flagDebug *bool
   flagLokiPushURL  *string
   flagLokiJobName *string
)

func MustRegisterFlags() {
	flagDebug = flag.Bool("debug", false, "debug")
        flagLokiPushURL = flag.String("loki.pushurl", ""," loki url to push logs")
        flagLokiJobName = flag.String("loki.jobname", "", "name of the loki job")
}

