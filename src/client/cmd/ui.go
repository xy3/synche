package cmd

import (
	"encoding/json"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/config"
	"net/http"
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
	PreRun: authenticateUserPreRun,
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
		log.Info("Client Server running on: 127.0.0.1:9448")
		startClientServer()
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().UintVarP(&guiPort, "gui-port", "p", 3000, "port the Synche GUI is running on")
}

func startClientServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		type request struct {
			FilePath    string
			DirectoryID uint
		}
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = FileUpload(req.FilePath, req.DirectoryID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
		return
	})

	mux.HandleFunc("/progress", func(w http.ResponseWriter, r *http.Request) {

	})

	log.Fatal(http.ListenAndServe(":9448", cors.AllowAll().Handler(mux)))
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
