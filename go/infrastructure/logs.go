package infrastructure

import (
	"fmt"
	"time"
	"os"
	"flag"

	log "github.com/sirupsen/logrus"
	"github.com/afiskon/promtail-client/promtail"
)

type promtailHook struct {
	conf   promtail.ClientConfig
	client promtail.Client
}

func newPromtailHook(lokiPushURL, jobName string) (*promtailHook, error) {
	conf := promtail.ClientConfig{
		PushURL:            lokiPushURL,
		Labels:             fmt.Sprintf(`{job="%s"}`,jobName),
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel: promtail.INFO,
		PrintLevel: promtail.ERROR,
	}
	loki, err := promtail.NewClientProto(conf)
	if err != nil {
		return nil, err
	}
	return &promtailHook{conf:conf,client:loki}, err
}

func (hook *promtailHook) Levels() []log.Level {
	return []log.Level{
	log.PanicLevel,
	log.FatalLevel,
	log.ErrorLevel,
	log.WarnLevel,
	log.InfoLevel,
	}
}

func (hook *promtailHook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case log.PanicLevel:
		hook.client.Errorf("PANIC %s", line)
	case log.FatalLevel:
		hook.client.Errorf("FATAL %s",line)
	case log.ErrorLevel:
		hook.client.Errorf("%s",line)
	case log.WarnLevel:
		hook.client.Warnf("%s",line)
	case log.InfoLevel:
		hook.client.Infof("%s",line)
	case log.DebugLevel, log.TraceLevel:
		hook.client.Debugf("%s",line)
	}
	return nil
}

func MustRegisterPromtail(lokiPushURL,jobName string) {
	debug := flag.Bool("debug", false, "debug")
        log.Info("Starting promtail with debug=",*debug)
	if *debug {
                log.SetLevel(log.DebugLevel)
        }

	hook, err := newPromtailHook(lokiPushURL,jobName)
	if err != nil {
		panic(err)
	}
	log.AddHook(hook)
}
