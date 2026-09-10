package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAddDepCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:           "add-dep <target-ref> <source-ref>",
		Short:         "Declare a dependency on a target",
		Args:          exactArgs(2),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			targetRef, err := parseTargetRef(args[0])
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "bldoc add-dep: %v\n", err)
				return err
			}
			if _, err := parseSourceRef(args[1]); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "bldoc add-dep: %v\n", err)
				return err
			}
			if format != "" && targetRef.Field == "" {
				err := fmt.Errorf("--format requires a target-ref with a :field (got %q)", args[0])
				fmt.Fprintf(cmd.ErrOrStderr(), "bldoc add-dep: %v\n", err)
				return err
			}
			return notImplemented(cmd)
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "template used to render the field's value")
	return cmd
}
