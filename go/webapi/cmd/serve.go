package cmd

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nanoteck137/{{ .ProjectName }}/database"
	"github.com/nanoteck137/{{ .ProjectName }}/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getDatabaseUrl() string {
	if !viper.IsSet("database_url") {
		log.Fatal("'database_url' not set")
	}

	return viper.GetString("database_url")
}

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl := getDatabaseUrl()
		conn, err := pgxpool.New(context.Background(), dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		db := database.New(conn)
		e := server.New(db)

		listenAddr := viper.GetString("listen_addr")
		err = e.Start(listenAddr)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	// serveCmd.Flags().IntP("port", "p", 3000, "Server Port")
	// viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))

	rootCmd.AddCommand(serveCmd)
}
