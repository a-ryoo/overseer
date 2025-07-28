package cmd

import (
	"github.com/a-ryoo/overseer/config"
	svc "github.com/a-ryoo/overseer/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	src  string
	dest string
)

func init() {
	renderCmd.Flags().StringVarP(&src, "src", "s", "", "Source file")
	renderCmd.Flags().StringVarP(&dest, "dest", "d", "", "Destination file")
}

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render a template",
	Run: func(cmd *cobra.Command, args []string) {
		var tempSvc = svc.NewTemplatingService()
		var conf = config.LocalConfig{
			VaultURL:        vaultURL,
			VaultToken:      vaultToken,
			VaultRoleID:     approleID,
			VaultRoleSecret: approleSecret,
		}
		var vaultSvc = svc.NewSecretsManager[map[string]string](cmd.Context(), conf)
		var rendered = tempSvc.RenderFile(src, dest, func(store, path, key string) string {
			var err error
			var result map[string]string
			result, err = vaultSvc.GetVaultEntity(store, path)
			if err != nil {
				log.Errorf("failed to get entity: %v", err)
				return "no_access"
			}
			return result[key]
		})

		err := os.WriteFile(dest, []byte(rendered), 0o777)
		if err != nil {
			log.Fatalf("failed to write file: %v", err)
		}
	},
}

func InitRender() {
	rootCmd.AddCommand(renderCmd)
}
