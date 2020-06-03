package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/projectsyn/probatorem/clients/googleapis"
	"github.com/projectsyn/probatorem/validators/gcp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	InvalidProjectError = fmt.Errorf("project does not meet requirements")
)

type OpenShift4GCP struct {
	log       logr.Logger
	projectId string
	domain    string
}

func (cmd *OpenShift4GCP) Action(_ *kingpin.ParseContext) error {
	ctx := context.Background()
	c := googleapis.New(ctx, cmd.projectId)

	v := gcp.New(c, cmd.log)
	v.SetExpectedManagedZone(cmd.domain)

	valid, err := v.ValidateAll()
	if err != nil {
		return err
	}

	if !valid {
		return InvalidProjectError
	}

	return nil
}

func Logger() (logr.Logger, io.WriteCloser) {
	formatter := new(logrus.TextFormatter)
	formatter.DisableTimestamp = true
	logger := logrus.New()
	logger.SetFormatter(formatter)

	return logrusr.NewLogger(logger), logger.Writer()
}

func App(w io.Writer) *kingpin.Application {
	app := kingpin.New("purser", "Tool to check access permissions and preconditions on clouds.")
	app.ErrorWriter(w)
	app.UsageWriter(w)
	app.Version(fmt.Sprintf("%v, commit %v, built at %v", version, commit, date))

	return app
}

func main() {
	log, w := Logger()
	defer w.Close()

	app := App(w)

	cmd := &OpenShift4GCP{log: log}
	c := app.Command("gcp", "tests a Google Compute Platform project for OpenShift 4.")
	c.Arg("project", "Id or name of a project.").
		Required().
		StringVar(&cmd.projectId)
	c.Flag("domain", "Domain of a managed zone that must be public.").
		Short('d').
		StringVar(&cmd.domain)
	c.Action(cmd.Action)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		if err == InvalidProjectError {
			log.Error(nil, err.Error())
			os.Exit(3)
		}

		log.Error(err, "Failed to validate project.")
		os.Exit(1)
	}
}
