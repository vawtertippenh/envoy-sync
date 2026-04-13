package cmd

import (
	"fmt"
	"os"

	"envoy-sync/internal/encrypt"
	"envoy-sync/internal/envfile"

	"github.com/spf13/cobra"
)

func init() {
	var passphrase string
	var decrypt bool
	var outputFile string

	encryptCmd := &cobra.Command{
		Use:   "encrypt [file]",
		Short: "Encrypt or decrypt values in a .env file using AES-GCM",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if passphrase == "" {
				return fmt.Errorf("--passphrase is required")
			}

			env, err := envfile.Parse(args[0])
			if err != nil {
				return fmt.Errorf("parse env file: %w", err)
			}

			var result map[string]string
			if decrypt {
				result, err = encrypt.DecryptMap(env, passphrase)
				if err != nil {
					return fmt.Errorf("decrypt: %w", err)
				}
			} else {
				result, err = encrypt.EncryptMap(env, passphrase)
				if err != nil {
					return fmt.Errorf("encrypt: %w", err)
				}
			}

			out := os.Stdout
			if outputFile != "" {
				f, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("create output file: %w", err)
				}
				defer f.Close()
				out = f
			}

			for k, v := range result {
				fmt.Fprintf(out, "%s=%s\n", k, v)
			}
			return nil
		},
	}

	encryptCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "Passphrase used for AES-GCM encryption/decryption (required)")
	encryptCmd.Flags().BoolVarP(&decrypt, "decrypt", "d", false, "Decrypt values instead of encrypting")
	encryptCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")

	rootCmd.AddCommand(encryptCmd)
}
