package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Ishan27g/shpMod/pkg"
	"github.com/urfave/cli/v2"
)

var config *pkg.Config
var modules pkg.Modules

var enable = &cli.Command{
	Name:    "enable",
	Aliases: []string{"e"},
	Usage:   "enable modules",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "all",
			Value:    false,
			Usage:    "enable all modules",
			Required: false,
		},
	},
	BashComplete: func(cCtx *cli.Context) {
		for _, t := range filterSlice(config.FilterShellCompletions(modules.Names()...), cCtx.Args().Slice()...) {
			fmt.Println(t)
		}
	},
	Action: func(cCtx *cli.Context) error {
		var enable = append(cCtx.Args().Slice(), config.Enable...)
		if cCtx.Bool("all") {
			enable = modules.Names()
		}
		pkg.Rewrite(defaultHcl(), enable...)
		return nil
	},
}

var start = &cli.Command{
	Name:    "start",
	Aliases: []string{"r"},
	Usage:   "shipyard run",
	Action:  func(cCtx *cli.Context) error { return osExecShipyardCmd("run") },
}
var stop = &cli.Command{
	Name:    "stop",
	Aliases: []string{"d"},
	Usage:   "shipyard destroy",
	Action:  func(cCtx *cli.Context) error { return osExecShipyardCmd("destroy") },
}
var configShow = &cli.Command{
	Name:    "showConfig",
	Aliases: []string{"sh"},
	Usage:   "print $HOME/.shpMod/cfg.hcl",
	Action: func(cCtx *cli.Context) error {
		fmt.Println(fmt.Sprintf("%s=%s", shipyardModulesEnvKey, defaultHcl()))
		fmt.Println(fmt.Sprintf("config : %s", defaultConfig()))
		if config != nil {
			fmt.Printf(fmt.Sprintf("%+v\n", *config))
		} else {
			fmt.Printf("empty\n")
		}
		return nil
	},
}

var configSet = &cli.Command{
	Name:    "setConfig",
	Aliases: []string{"sc"},
	Usage:   "set $HOME/.shpMod/cfg.hcl",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "clear",
			Value:    false,
			Usage:    "clear the config",
			Required: false,
		},
	},
	Action: func(cCtx *cli.Context) error {
		if cCtx.Bool("clear") {
			fmt.Println("cleaning", defaultConfig())
			return os.Remove(defaultConfig())
		}
		if cCtx.Args().Len() != 1 {
			return errors.New("missing filename " + "`shpMod setConfig config.hcl`")
		}

		_, err := pkg.ReadConfig(cCtx.Args().Get(0))
		if err != nil {
			return err
		}
		input, err := ioutil.ReadFile(cCtx.Args().Get(0))
		if err != nil {
			fmt.Println(err)
			return err
		}

		err = ioutil.WriteFile(defaultConfig(), input, 0644)
		return err
	},
}

func main() {

	var err error
	modules, err = pkg.ReadModulesFile(defaultHcl())
	if err != nil {
		return
	}

	config, _ = pkg.ReadConfig(defaultConfig())

	if err := (&cli.App{
		HideHelpCommand:      true,
		EnableBashCompletion: true,
		Commands:             []*cli.Command{enable, start, stop, configShow, configSet},
	}).Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
