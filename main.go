package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

var (
	vaultToken string
	vaultAddr  string
	outputFile string
)

var rootCmd = &cobra.Command{
	Use:   "vault-scanner",
	Short: "A CLI tool to scan and save Vault secrets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Vault Scanner starting...")

		client, err := newVaultClient()
		if err != nil {
			fmt.Printf("Error creating Vault client: %v\n", err)
			os.Exit(1)
		}

		secrets, err := scanSecrets(client)
		if err != nil {
			fmt.Printf("Error scanning secrets: %v\n", err)
			os.Exit(1)
		}

		if err := writeSecretsToFile(secrets); err != nil {
			fmt.Printf("Error writing secrets to file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully saved secrets to %s\n", outputFile)
	},
}

func newVaultClient() (*api.Client, error) {
	config := &api.Config{
		Address: vaultAddr,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(vaultToken)
	return client, nil
}

func scanSecrets(client *api.Client) (map[string]interface{}, error) {
	secrets := make(map[string]interface{})
	mounts, err := client.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	for path, mount := range mounts {
		if mount.Type == "kv" {
			fmt.Printf("Scanning %s...\n", path)
			listSecretsRecursive(client, path, mount, secrets)
		}
	}

	return secrets, nil
}

func listSecretsRecursive(client *api.Client, path string, mount *api.MountOutput, secrets map[string]interface{}) error {
	listPath := path
	if mount.Options["version"] == "2" {
		listPath = path + "metadata/"
	}

	secretValues, err := client.Logical().List(listPath)
	if err != nil {
		return err
	}

	if secretValues == nil {
		return nil
	}

	for _, key := range secretValues.Data["keys"].([]interface{}) {
		secretPath := path + key.(string)
		// Check if it's a directory
		if secretPath[len(secretPath)-1] == '/' {
			if err := listSecretsRecursive(client, secretPath, mount, secrets); err != nil {
				fmt.Printf("Error scanning %s: %v\n", secretPath, err)
			}
		} else {
			readPath := ""
			if mount.Options["version"] == "2" {
				readPath = path + "data/" + key.(string)
			} else {
				readPath = secretPath
			}

			secret, err := client.Logical().Read(readPath)
			if err != nil {
				fmt.Printf("Error reading secret %s: %v\n", secretPath, err)
				continue
			}
			if secret != nil {
				if mount.Options["version"] == "2" {
					secrets[secretPath] = secret.Data["data"]
				} else {
					secrets[secretPath] = secret.Data
				}
			}
		}
	}

	return nil
}

func writeSecretsToFile(secrets map[string]interface{}) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(secrets)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultToken, "token", "", "Vault token to use for authentication")
	rootCmd.PersistentFlags().StringVar(&vaultAddr, "addr", "http://127.0.0.1:8200", "Vault address")
	rootCmd.PersistentFlags().StringVar(&outputFile, "output", "secrets.json", "Output file to save the secrets")
	rootCmd.MarkPersistentFlagRequired("token")
}

func main() {
	Execute()
}
