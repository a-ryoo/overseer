package services

import (
	"context"
	"fmt"
	"github.com/a-ryoo/overseer/config"
	json "github.com/json-iterator/go"
	"os"
	"strings"

	"github.com/gookit/goutil"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
	kauth "github.com/hashicorp/vault/api/auth/kubernetes"

	log "github.com/sirupsen/logrus"
)

const SATokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
const K8SVaultURL = "http://vault-active.vault.svc.cluster.local:8200"

var Role = os.Getenv("SA_ROLE")
var AppRoleID = os.Getenv("APP_ROLE_ID")
var AppRoleSecret = os.Getenv("APP_ROLE_SECRET")

type VaultSecretsManager[T any] struct {
	conf config.LocalConfig
	ctx  context.Context
}

func NewSecretsManager[T any](ctx context.Context, conf config.LocalConfig) *VaultSecretsManager[T] {
	return &VaultSecretsManager[T]{
		conf: conf,
		ctx:  ctx,
	}
}

func (s *VaultSecretsManager[T]) GetVaultEntity(store, path string) (T, error) {
	var result T
	var data, err = s.GetVaultEntityBytes(store, path)
	if err != nil {
		log.Errorf("[VAULT] Failed getting entity secret: %v", err)
		return result, err
	}

	unmarshalErr := json.Unmarshal(data, &result)
	if unmarshalErr != nil {
		log.Errorf("[VAULT] Failed unmarshalling entity secret: %v", unmarshalErr)
		return result, err
	}

	return result, err
}

func (s *VaultSecretsManager[T]) GetVaultEntityBytes(store, path string) ([]byte, error) {
	var dataMap, err = s.getClient(s.ctx).KVv2(store).Get(s.ctx, path)
	if err != nil {
		if strings.Contains(err.Error(), vault.ErrSecretNotFound.Error()) {
			log.Errorf("[VAULT] Failed getting entity secret: %v", err)
			return nil, err
		}
	}

	if dataMap == nil {
		return nil, fmt.Errorf("[VAULT] Failed getting entity secret: %v", err)
	}

	serialized, serializeErr := json.Marshal(dataMap.Data)
	if serializeErr != nil {
		log.Errorf("[VAULT] Failed serializing entity secret: %v", serializeErr)
		return nil, err
	}

	return serialized, nil
}

func (s *VaultSecretsManager[T]) getClient(ctx context.Context) *vault.Client {
	var conf = vault.DefaultConfig()
	conf.Address = s.conf.VaultURL

	switch {
	case !goutil.IsEmpty(Role):
		conf.Address = K8SVaultURL
		client, err := vault.NewClient(conf)
		if err != nil {
			log.Panicf("[VAULT] Unable to initialize Internal Vault client: %v", err)
		}
		k8sAuth, authErr := kauth.NewKubernetesAuth(
			Role,
			kauth.WithServiceAccountTokenPath(SATokenPath),
		)

		if authErr != nil {
			log.Panicf("[VAULT] Unable to initialize Kubernetes auth method: %v", err)
		}

		authInfo, authErr := client.Auth().Login(ctx, k8sAuth)
		if authErr != nil {
			log.Panicf("[VAULT] Unable to login with K8S auth method: %v", authErr)
		}

		if authInfo == nil {
			log.Panic("[VAULT] No auth info was returned after login with K8S")
		}

		client.SetMaxRetries(10)
		return client

	case !goutil.IsEmpty(s.conf.VaultToken):
		client, err := vault.NewClient(conf)
		if err != nil {
			log.Panicf("[VAULT] Unable to initialize Internal Vault client: %v", err)
		}
		client.SetToken(s.conf.VaultToken)
		client.SetMaxRetries(10)
		return client

	case !goutil.IsEmpty(AppRoleID) && !goutil.IsEmpty(AppRoleSecret):
		client, err := vault.NewClient(conf)
		if err != nil {
			log.Panicf("[VAULT] Unable to initialize Internal Vault client: %v", err)
		}
		appRoleAuth, appRoleErr := auth.NewAppRoleAuth(s.conf.VaultRoleID, &auth.SecretID{FromString: s.conf.VaultRoleSecret})
		if appRoleErr != nil {
			log.Panicf("[VAULT] Unable to initialize AppRole auth method: %v", appRoleErr)
		}

		authInfo, authErr := client.Auth().Login(ctx, appRoleAuth)
		if authErr != nil {
			log.Panicf("[VAULT] Unable to login to AppRole auth method: %v", authErr)
		}

		if authInfo == nil {
			log.Panic("[VAULT] No auth info was returned after login")
		}

		client.SetMaxRetries(10)
		return client

	default:
		log.Panic("[VAULT] No auth method was provided")
		return nil
	}
}
