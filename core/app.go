package core

import "github.com/urfave/cli/v2"

const (
	CLI_NAME = "c2-chopper"
)

func CreateApp(commands []*cli.Command) *cli.App {
	app := &cli.App{
		Name:     CLI_NAME,
		Commands: commands,
		Authors: []*cli.Author{
			{
				Name:  "Yassine Belkhadem",
				Email: "yassine.belkhadem@insat.rnu.tn",
			},
			{
				Name:  "Mongi Saidane",
				Email: "mongi.saidane@insat.ucar.tn",
			},
		},
		Usage:                "A futuristic C2 framework where you can achieve all your dreams",
		Copyright:            "(c) 2023 Soter14",
		EnableBashCompletion: true,
	}

	return app
}
