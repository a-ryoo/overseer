package cmd

import (
	"github.com/a-ryoo/overseer/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var (
	vaultURL      string
	vaultToken    string
	approleID     string
	approleSecret string
)

var rootCmd = &cobra.Command{
	Use:   "overseer",
	Short: "Overseer is a tool to manage secrets in Vault",
}

func Start() error {
	var cfg config.LocalConfig
	configPath := filepath.Join(os.Getenv("HOME"), ".vault", "creds.yaml")
	_ = os.MkdirAll(filepath.Dir(configPath), 0o755)
	if _, err := os.ReadFile(configPath); err != nil {
		if out, err := yaml.Marshal(&cfg); err == nil {
			_ = os.WriteFile(configPath, out, 0o777)
		} else {
			log.Errorf("warning: failed to save config: %v\n", err)
		}
	}

	if data, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			log.Fatalf("failed to parse existing config: %v", err)
		}
	}
	vaultURL = cfg.VaultURL
	vaultToken = cfg.VaultToken
	approleID = cfg.VaultRoleID
	approleSecret = cfg.VaultRoleSecret

	InitAuth()
	InitRender()
	return rootCmd.Execute()
}
