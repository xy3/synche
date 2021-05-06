package cmd

import (
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/jobs"

	"github.com/spf13/cobra"
)

var serverCertIPs []string
var certsSaveDirPath string

var makecertsCmd = &cobra.Command{
	Use:   "makecerts",
	Short: "Generate HTTPS and TLS certificates and keys",
	Run: func(cmd *cobra.Command, args []string) {
		err := jobs.GenerateCertificates(serverCertIPs, certsSaveDirPath)
		if err != nil {
			log.WithError(err).Error("Failed to generate the certificates")
		}
	},
}

func init() {
	rootCmd.AddCommand(makecertsCmd)
	makecertsCmd.Flags().StringArrayVar(
		&serverCertIPs,
		"server-ips",
		[]string{"127.0.0.1", "0.0.0.0", "::1"},
		"the server IPs for which to generate the certificates for.",
	)
	makecertsCmd.Flags().StringVarP(
		&certsSaveDirPath,
		"output-dir",
		"o",
		".",
		"the directory to save the certificate files",
	)
}
