package cmd

import (
	"errors"
	"github.com/go-openapi/loads"
	"github.com/goftp/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/ftp"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/restapi/operations"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

func startHttpServer(flags *flag.FlagSet, sig chan struct{}) error {
	var (
		apiServer   *restapi.Server // make sure init is called
		swaggerSpec *loads.Document
		err         error
		port        int
	)

	swaggerSpec, err = loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatal(err)
	}

	port, err = flags.GetInt("port")
	if err != nil {
		log.Fatal(err)
	}
	if port == 0 {
		restapi.OverrideFlags()
	}

	// get server with flag values filled out
	apiServer = restapi.NewServer(operations.NewSyncheAPI(swaggerSpec))
	defer apiServer.Shutdown()

	apiServer.ConfigureAPI()
	apiServer.Host = c.Config.Server.Host
	if err = apiServer.Serve(); err != nil {
		return err
	}

	<-sig
	log.Debug("Shutting down HTTP Server ...")
	return nil
}

func startFtpServer(flags *flag.FlagSet, sig chan struct{}) error {
	cfg := c.Config.Ftp
	if len(cfg.PassivePorts) > 0 {
		portRange := strings.Split(cfg.PassivePorts, "-")
		if len(portRange) != 2 {
			return errors.New("wrong ftp-passive-port-range format, eg: 52013-52114")
		}
		if _, err := strconv.Atoi(portRange[0]); err != nil {
			return err
		}
		if _, err := strconv.Atoi(portRange[1]); err != nil {
			return err
		}
		if portRange[0] >= portRange[1] {
			return errors.New("ftp port start should be less than port end")
		}
	}

	ftpLogger := log.New()

	options := &server.ServerOpts{
		Factory:        &ftp.Factory{Logger: ftpLogger},
		Auth:           &ftp.Auth{},
		Hostname:       cfg.Hostname,
		PublicIp:       cfg.PublicIp,
		PassivePorts:   cfg.PassivePorts,
		Port:           cfg.Port,
		TLS:            false,
		CertFile:       cfg.CertFile,
		KeyFile:        cfg.KeyFile,
		ExplicitFTPS:   true,
		WelcomeMessage: cfg.WelcomeMessage,
		Logger:         &ftp.Logger{},
	}

	go func() {
		ftpServer := server.NewServer(options)
		if err := ftpServer.ListenAndServe(); err != nil {
			log.Errorf("ftp server start with error: %s", err)
			panic(err)
		}
	}()
	<-sig
	log.Debug("FTP Server shutting down...")
	return nil
}

func NewServeCmd() *cobra.Command {
	// serveCmd represents the serve command
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the Synche API server",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
			_, err := database.InitSyncheData()
			if err != nil {
				log.WithError(err).Fatal("Failed to start Synche Data requirements")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			var (
				wg  sync.WaitGroup
				sig = make(chan struct{})
			)

			wg.Add(3)
			go func() {
				defer wg.Done()
				_ = startHttpServer(cmd.Flags(), sig)
			}()
			go func() {
				defer wg.Done()
				_ = startFtpServer(cmd.Flags(), sig)
			}()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			close(sig)
			wg.Done()
			wg.Wait()
		},
	}
	return serveCmd
}

func init() {
	serveCmd := NewServeCmd()
	rootCmd.AddCommand(serveCmd)
	if ServerFlags.HasFlags() {
		serveCmd.Flags().AddFlagSet(ServerFlags)
	}

	serveCmd.Flags().IntVar(&c.Config.Ftp.Port, "ftp-port", 2121, "port for the ftp server to run on")
	serveCmd.Flags().StringVar(&c.Config.Ftp.Hostname, "ftp-hostname", "127.0.0.1", "ftp service listen IP")
	serveCmd.Flags().StringVar(&c.Config.Ftp.PublicIp, "ftp-public-ip", "", "ftp client connect to ftp server to transfer data in passive mode")
	serveCmd.Flags().StringVar(&c.Config.Ftp.PassivePorts, "ftp-passive-port-range", "52013-52114", "ftp server will pick a port random in the range, open the data connection tunnel in passive mode")
	serveCmd.Flags().StringVar(&c.Config.Ftp.WelcomeMessage, "ftp-welcome-message", "Welcome to the Synche FTP server", "ftp server welcome message")

	err := viper.BindPFlag("config.server.port", serveCmd.Flags().Lookup("port"))
	if err != nil {
		log.Fatal(err)
	}
}
