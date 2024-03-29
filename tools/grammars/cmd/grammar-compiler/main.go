package main

import (
	"os"

	"github.com/github/linguist/tools/grammars/compiler"
	"github.com/urfave/cli/v2"
)

func cwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func wrap(err error) error {
	return cli.Exit(err, 255)
}

func main() {
	app := cli.NewApp()
	app.Name = "Linguist Grammars Compiler"
	app.Usage = "Compile user-submitted grammars and check them for errors"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "linguist-path",
			Value: cwd(),
			Usage: "path to Linguist root",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "add",
			Usage: "add a new grammar source",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "force",
					Aliases: []string{"f"},
					Usage: "ignore compilation errors",
				},
			},
			Action: func(c *cli.Context) error {
				conv, err := compiler.NewConverter(c.String("linguist-path"))
				if err != nil {
					return wrap(err)
				}
				if err := conv.AddGrammar(c.Args().First()); err != nil {
					if !c.Bool("force") {
						return wrap(err)
					}
				}
				if err := conv.WriteGrammarList(); err != nil {
					return wrap(err)
				}
				return nil
			},
		},
		{
			Name:  "update",
			Usage: "update grammars.yml with the contents of the grammars library",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "force",
					Aliases: []string{"f"},
					Usage: "write grammars.yml even if grammars fail to compile",
				},
				&cli.BoolFlag{
					Name:  "verbose",
					Aliases: []string{"v"},
					Usage: "show all warnings",
					Value: false,
				},
			},
			Action: func(c *cli.Context) error {
				conv, err := compiler.NewConverter(c.String("linguist-path"))
				if err != nil {
					return wrap(err)
				}
				if err := conv.ConvertGrammars(true); err != nil {
					return wrap(err)
				}
				if err := conv.Report(c.Bool("verbose")); err != nil {
					if !c.Bool("force") {
						return wrap(err)
					}
				}
				if err := conv.WriteGrammarList(); err != nil {
					return wrap(err)
				}
				return nil
			},
		},
		{
			Name:  "compile",
			Usage: "convert the grammars from the library",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "proto-out", Aliases: []string{"P"}},
				&cli.StringFlag{Name: "out", Aliases: []string{"o"}},
			},
			Action: func(c *cli.Context) error {
				conv, err := compiler.NewConverter(c.String("linguist-path"))
				if err != nil {
					return cli.Exit(err, 1)
				}
				if err := conv.ConvertGrammars(false); err != nil {
					return cli.Exit(err, 1)
				}
				if out := c.String("proto-out"); out != "" {
					if err := conv.WriteProto(out); err != nil {
						return cli.Exit(err, 1)
					}
				}
				if out := c.String("out"); out != "" {
					if err := conv.WriteJSON(out); err != nil {
						return cli.Exit(err, 1)
					}
				}
				if err := conv.Report(false); err != nil {
					return wrap(err)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
