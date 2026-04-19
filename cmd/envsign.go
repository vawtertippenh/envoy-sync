package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/envsign"
)

func init() {
	var secret string
	var keys string
	var sigFile string

	cmd := &cobra.Command{
		Use:   "envsign",
		Short: "Sign or verify an .env file using HMAC-SHA256",
	}

	signCmd := &cobra.Command{
		Use:   "sign [file]",
		Short: "Sign an .env file and print the signature JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return err
			}
			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					keyList = append(keyList, strings.TrimSpace(k))
				}
			}
			sig, err := envsign.Sign(env, secret, keyList)
			if err != nil {
				return err
			}
			return json.NewEncoder(os.Stdout).Encode(sig)
		},
	}
	signCmd.Flags().StringVar(&secret, "secret", "", "HMAC secret (required)")
	signCmd.Flags().StringVar(&keys, "keys", "", "Comma-separated keys to sign (default: all)")
	signCmd.MarkFlagRequired("secret")

	verifyCmd := &cobra.Command{
		Use:   "verify [file]",
		Short: "Verify an .env file against a signature JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := envfile.Parse(args[0])
			if err != nil {
				return err
			}
			f, err := os.Open(sigFile)
			if err != nil {
				return err
			}
			defer f.Close()
			var sig envsign.Signature
			if err := json.NewDecoder(f).Decode(&sig); err != nil {
				return err
			}
			ok, err := envsign.Verify(env, secret, sig)
			if err != nil {
				return err
			}
			if ok {
				fmt.Println("signature valid ✓")
			} else {
				fmt.Fprintln(os.Stderr, "signature INVALID ✗")
				os.Exit(1)
			}
			return nil
		},
	}
	verifyCmd.Flags().StringVar(&secret, "secret", "", "HMAC secret (required)")
	verifyCmd.Flags().StringVar(&sigFile, "sig", "", "Path to signature JSON file (required)")
	verifyCmd.MarkFlagRequired("secret")
	verifyCmd.MarkFlagRequired("sig")

	cmd.AddCommand(signCmd, verifyCmd)
	rootCmd.AddCommand(cmd)
}
