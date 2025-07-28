package config

type LocalConfig struct {
	VaultURL        string `yaml:"vault_url" json:"vault_url,omitempty"`
	VaultRoleID     string `yaml:"vault_role_id" json:"vault_role_id,omitempty"`
	VaultRoleSecret string `yaml:"vault_role_secret" json:"vault_role_secret,omitempty"`
	VaultToken      string `yaml:"vault_token" json:"vault_token,omitempty"`
}
