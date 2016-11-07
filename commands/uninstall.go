package commands

import (
	"github.com/SvenDowideit/jig/util"

	//	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var Uninstall = cli.Command{
	Name:  "uninstall",
	Usage: "uninstall jig",
	Flags: []cli.Flag{},
	Action: func(context *cli.Context) error {
		util.Run("docker-machine", "-D", "stop", "rancher")
		util.Run("docker-machine", "-D", "rm", "-y", "rancher")

		return nil
	},
}
