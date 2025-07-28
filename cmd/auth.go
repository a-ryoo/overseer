package cmd

import (
	"fmt"
	"github.com/a-ryoo/overseer/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func init() {
	loginCmd.Flags().StringVar(&vaultURL, "vault_url", "", "Vault URL")
	loginCmd.Flags().StringVar(&vaultToken, "vault_token", "", "Vault token")
	loginCmd.Flags().StringVar(&approleID, "approle_id", "", "Vault AppRole ID")
	loginCmd.Flags().StringVar(&approleSecret, "approle_secret", "", "Vault AppRole Secret")
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Vault",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LocalConfig{
			VaultURL:        vaultURL,
			VaultToken:      vaultToken,
			VaultRoleID:     approleID,
			VaultRoleSecret: approleSecret,
		}

		savePath := filepath.Join(os.Getenv("HOME"), ".vault", "creds.yaml")

		if err := os.MkdirAll(filepath.Dir(savePath), 0o755); err != nil {
			log.Fatalf("failed to create config directory: %v\n", err)
		}

		file, err := os.Create(savePath)
		if err != nil {
			log.Fatalf("failed to create config file: %v", err)
		}
		defer file.Close()

		encoder := yaml.NewEncoder(file)
		encoder.SetIndent(2)
		if err := encoder.Encode(&cfg); err != nil {
			log.Fatalf("failed to write config: %v", err)
		}

		fmt.Printf("Vault credentials saved to %s\n", savePath)
	},
}

func InitAuth() {
	authCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(authCmd)
}
