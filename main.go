package main

import (
	"log"
	"net/http"

	"github.com/ddtmachado/prom-backup-exporter/collector"
	"github.com/ddtmachado/prom-backup-exporter/config"

	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var globalConfig config.Config
var rootCmd = &cobra.Command{
	Use:   "backup-exporter",
	Short: "backup-exporter is a backup metric exporter for prometheus",
	Long: `Configurable backup metric exporter. Currently supported backup repositories are:
                - Restic
								- ElasticSearch
								- Tarball directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		startExporter()
	},
}

func main() {
	rootCmd.Execute()
}

func prometheusHandlerFunc(next http.Handler) http.Handler {
	return promhttp.Handler()
}

func rootHandler(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(`<html>
	<head><title>Backup Exporter</title></head>
	<body>
	<h1>Backup Exporter</h1>
	<p><a href='` + globalConfig.Path + `'>Metrics</a></p>
	</body>
	</html>`))
}

func startExporter() {
	backupCollector := collector.NewBackupCollector(globalConfig.Repos())
	prometheus.MustRegister(backupCollector)
	router := gin.Default()
	router.GET(globalConfig.Path, adapter.Wrap(prometheusHandlerFunc))
	router.GET("/", rootHandler)

	// By default it serves on :8080 unless a
	// Port value was defined in the config file.
	router.Run(":" + globalConfig.Port)
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.toml)")
	rootCmd.PersistentFlags().String("port", "--port", "http port to expose the backup exporter")
	rootCmd.PersistentFlags().String("path", "--path", "http path to expose the metrics")
	viper.BindPFlag("Port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("Path", rootCmd.PersistentFlags().Lookup("path"))
	viper.SetDefault("Port", "8080")
	viper.SetDefault("Path", "/metrics")
}

func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)
	} else {
		//Using default config file
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath("/etc/backup-exporter/")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println(err)
	}

	viper.Unmarshal(&globalConfig)
	log.Printf("Config: %s", globalConfig)
}
