package main

import (
	"github.com/Shubhaankar-sharma/todoapp/api"
	"github.com/Shubhaankar-sharma/todoapp/server"
	"github.com/go-chi/chi/v5"
	"log"
	// _ "net/http/pprof"  // only in use when profiling
	"os"

	"github.com/Shubhaankar-sharma/todoapp/utils"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

func main() {
	config, err := utils.LoadConfig("./", "app")
	if err != nil {
		log.Fatalln(err.Error())
	}
	commands(config)

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig(testConf bool) (conf utils.Config, err error) {
	if testConf {
		conf, err = utils.LoadConfig("./", "test")
	} else {
		conf, err = utils.LoadConfig("./", "app")
	}
	return
}

func commands(config utils.Config) {
	app.Commands = []*cli.Command{
		{
			Name:  "migrate_up",
			Usage: "Migrate DB to latest version",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "test",
					Aliases: []string{"t"},
					Usage:   "loads test.env instead of app.env",
				},
			},
			Action: func(c *cli.Context) error {
				conf, err := loadConfig(c.Bool("test"))
				if err != nil {
					return err
				}
				err = utils.MigrateUp(conf, "./db/migration/")
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "dropdb",
			Usage: "Drop the DB",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "test",
					Aliases: []string{"t"},
					Usage:   "loads test.env instead of app.env",
				},
			},
			Action: func(c *cli.Context) error {
				conf, err := loadConfig(c.Bool("test"))
				if err != nil {
					return err
				}
				err = utils.MigrateDown(conf, "./db/migration/")
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "migrate_steps",
			Usage: "Migrate with Steps",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "steps",
					Usage: "Number of steps of migrations to run",
				},
				&cli.BoolFlag{
					Name:    "test",
					Aliases: []string{"t"},
					Usage:   "loads test.env instead of app.env",
				},
			},
			Action: func(c *cli.Context) error {
				conf, err := loadConfig(c.Bool("test"))
				if err != nil {
					return err
				}
				err = utils.MigrateSteps(c.Int("steps"), conf, "./db/migration/")
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "runserver",
			Usage: "Run Api Server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "host",
					Usage:   "Host on which server has to be run",
					Value:   "localhost",
					Aliases: []string{"H"},
				},
				&cli.IntFlag{
					Name:    "port",
					Usage:   "Port on which server has to be run",
					Value:   5000,
					Aliases: []string{"P"},
				},
			},
			Action: func(c *cli.Context) error {
				s := server.NewServer(config)
				//Create Routers Here
				apiRouter := chi.NewRouter()

				//Add Routes to Routers Here
				api.MainRouter(apiRouter, s.Queries, config)
				//Mount Routers here
				s.Router.Mount("/", apiRouter)
				// r.Mount("/debug/", middleware.Profiler()) // Only in use when profiling
				//Store Router in Struct
				err := s.RunServer(c.String("host"), c.Int("port"))
				return err
			},
		},
	}
}
