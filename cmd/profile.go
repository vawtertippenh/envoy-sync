package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-sync/internal/envfile"
	"envoy-sync/internal/profile"

	"github.com/spf13/cobra"
)

func init() {
	var storePath string

	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage named environment profiles",
	}

	setCmd := &cobra.Command{
		Use:   "set <name> <envfile>",
		Short: "Save an env file as a named profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, envPath := args[0], args[1]
			env, err := envfile.Parse(envPath)
			if err != nil {
				return fmt.Errorf("parse env file: %w", err)
			}
			s, err := profile.LoadStore(storePath)
			if err != nil {
				return err
			}
			s.Set(name, env)
			if err := profile.SaveStore(storePath, s); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Profile %q saved to %s\n", name, storePath)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := profile.LoadStore(storePath)
			if err != nil {
				return err
			}
			names := s.List()
			if len(names) == 0 {
				fmt.Println("No profiles found.")
				return nil
			}
			fmt.Println(strings.Join(names, "\n"))
			return nil
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a named profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := profile.LoadStore(storePath)
			if err != nil {
				return err
			}
			if !s.Delete(args[0]) {
				return fmt.Errorf("profile %q not found", args[0])
			}
			if err := profile.SaveStore(storePath, s); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Profile %q deleted\n", args[0])
			return nil
		},
	}

	profileCmd.PersistentFlags().StringVar(&storePath, "store", ".envoy-profiles.json", "Path to profile store file")
	profileCmd.AddCommand(setCmd, listCmd, deleteCmd)
	rootCmd.AddCommand(profileCmd)
}
