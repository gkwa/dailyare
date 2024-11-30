package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gkwa/dailyare/core"
	"github.com/gkwa/dailyare/internal/logger"
)

var (
	cfgFile   string
	verbose   int
	logFormat string
	cliLogger logr.Logger
	since     string
	noCache   bool
)

var rootCmd = &cobra.Command{
	Use:   "dailyare",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cliLogger.IsZero() {
			cliLogger = logger.NewConsoleLogger(verbose, logFormat == "json")
		}

		ctx := logr.NewContext(context.Background(), cliLogger)
		cmd.SetContext(ctx)
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.Info("Running command")

		client, err := api.DefaultRESTClient()
		if err != nil {
			logger.Error(err, "Failed to create REST client")
			return
		}

		notificationRepo := core.NewGithubRepository(client)
		prService := core.NewGithubPRService(client)
		cacheService := core.NewFileCacheService(viper.GetString("home"))
		service := core.NewNotificationService(notificationRepo, prService, cacheService)

		err = service.FetchNotifications(logger, since, noCache)
		if err != nil {
			logger.Error(err, "Failed to fetch notifications")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dailyare.yaml)")
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "increase verbosity")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "", "json or text (default is text)")
	rootCmd.Flags().StringVar(&since, "since", "7d", "Filter notifications by time (default: 7d)")
	rootCmd.Flags().BoolVar(&noCache, "no-cache", false, "Bypass the cache and fetch fresh data")

	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		fmt.Printf("Error binding verbose flag: %v\n", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format")); err != nil {
		fmt.Printf("Error binding log-format flag: %v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".dailyare")

		viper.Set("home", home)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	logFormat = viper.GetString("log-format")
	verbose = viper.GetInt("verbose")
}

func LoggerFrom(ctx context.Context, keysAndValues ...interface{}) logr.Logger {
	if cliLogger.IsZero() {
		cliLogger = logger.NewConsoleLogger(verbose, logFormat == "json")
	}
	newLogger := cliLogger
	if ctx != nil {
		if l, err := logr.FromContext(ctx); err == nil {
			newLogger = l
		}
	}
	return newLogger.WithValues(keysAndValues...)
}
