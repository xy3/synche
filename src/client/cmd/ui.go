package cmd

import (
	log "github.com/sirupsen/logrus"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"
)

var guiPort uint

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Open the Synche UI in your browser",
	Run: func(cmd *cobra.Command, args []string) {
		uiUrl := "http://127.0.0.1:" + strconv.Itoa(int(guiPort))

		log.Infof("Opening GUI at %s", uiUrl)

		baseUrl := path.Join(c.Config.Server.Host, c.Config.Server.BasePath)
		var prefix = "http://"
		if c.Config.Server.Https {
			prefix = "https://"
		}
		err := os.Setenv("NEXT_PUBLIC_BASE_URL", prefix+baseUrl)
		if err != nil {
			log.WithError(err).Fatal("Failed to set the base URL environment variable")
			return
		}

		openBrowser(uiUrl)
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().UintVarP(&guiPort, "gui-port", "p", 3000, "port the Synche GUI is running on")
}

func openBrowser(url string) bool {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}
