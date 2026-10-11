package cmd

import (
	"fmt"
	"os"

	appkit "github.com/soundadam/appkit/go"
	"github.com/soundadam/appkit/go/cli"
	"github.com/spf13/cobra"

	"github.com/soundadam/tea/internal/app"
	"github.com/soundadam/tea/internal/teaerr"
	"github.com/soundadam/tea/internal/tui"
)

// NewRoot builds the command tree; info names this build for version and update.
func NewRoot(info appkit.App) *cobra.Command {
	root := &cobra.Command{
		Use:           "tea",
		Short:         "Keep this Mac awake, then restore its previous sleep setting when you're done.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if os.Geteuid() == 0 {
				return teaerr.RootRefused()
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := app.New(nil, cmd.OutOrStdout(), cmd.ErrOrStderr())
			if err != nil {
				return err
			}
			return tui.Run(application, info)
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true
	root.AddCommand(
		onCmd(),
		offCmd(),
		statusCmd(),
		shutdownCmd(),
		authCmd(),
		tuiCmd(info),
	)
	cli.AddCommands(root, info)
	return root
}

func Execute(info appkit.App) int {
	root := NewRoot(info)
	// As root the command is refused, and the update check must not create
	// root-owned files under the user's Application Support/tea, which also
	// holds state.json.
	var notice *cli.Notice
	if os.Geteuid() != 0 {
		notice = cli.StartNotice(info)
	}
	code := 0
	ran, err := root.ExecuteC()
	if err != nil {
		fmt.Fprintln(os.Stderr, app.FormatError(err))
		code = 1
		if teaerr.IsUsage(err) {
			code = 2
		}
	}
	// The client has its own update line and may already have updated.
	client := ran == nil || ran == root || ran.Name() == "tui"
	if notice != nil && !client {
		notice.Print(os.Stderr)
	}
	return code
}

func withApp(fn func(app.App, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		application, err := app.New(nil, cmd.OutOrStdout(), cmd.ErrOrStderr())
		if err != nil {
			return err
		}
		return fn(application, args)
	}
}

func onCmd() *cobra.Command {
	return &cobra.Command{Use: "on", Short: "Disable lid-close sleep and remember the previous setting", Args: cobra.NoArgs, RunE: withApp(func(a app.App, _ []string) error { return a.On() })}
}

func offCmd() *cobra.Command {
	return &cobra.Command{Use: "off", Short: "Restore the sleep setting tea owns", Args: cobra.NoArgs, RunE: withApp(func(a app.App, _ []string) error { return a.Off() })}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{Use: "status", Short: "Show awake mode, power source, and shutdown state", Args: cobra.NoArgs, RunE: withApp(func(a app.App, _ []string) error { return a.Status() })}
}

func tuiCmd(info appkit.App) *cobra.Command {
	run := withApp(func(a app.App, _ []string) error { return tui.Run(a, info) })
	return &cobra.Command{
		Use:        "tui",
		Aliases:    []string{"interactive"},
		Short:      "Open the tea client",
		Hidden:     true,
		Deprecated: "run tea with no arguments",
		Args:       cobra.NoArgs,
		RunE:       run,
	}
}

func shutdownCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "shutdown", Short: "Schedule or cancel one tea-owned shutdown"}
	cmd.AddCommand(
		&cobra.Command{Use: "after DURATION", Short: "Schedule a shutdown after a duration such as 30m or 2h", Args: cobra.ExactArgs(1), RunE: withApp(func(a app.App, args []string) error { return a.ShutdownAfter(args[0]) })},
		&cobra.Command{Use: "status", Short: "Show the tea-owned shutdown", Args: cobra.NoArgs, RunE: withApp(func(a app.App, _ []string) error { return a.ShutdownStatus() })},
		&cobra.Command{Use: "cancel [ACTION_ID]", Short: "Cancel the tea-owned shutdown", Args: cobra.MaximumNArgs(1), RunE: withApp(func(a app.App, args []string) error {
			id := ""
			if len(args) == 1 {
				id = args[0]
			}
			return a.ShutdownCancel(id)
		})},
	)
	return cmd
}

func authCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "auth", Short: "Inspect or register the narrow passwordless helper"}
	cmd.AddCommand(
		&cobra.Command{Use: "status", Short: "Show helper registration and sudo Touch ID", Args: cobra.NoArgs, RunE: withApp(func(a app.App, _ []string) error { return a.AuthStatus() })},
		&cobra.Command{Use: "register", Short: "Install the narrow passwordless helper", Args: cobra.NoArgs, RunE: withApp(func(a app.App, _ []string) error { return a.AuthRegister() })},
		&cobra.Command{Use: "unregister", Short: "Remove the helper and sudoers rule", Args: cobra.NoArgs, RunE: withApp(func(a app.App, _ []string) error { return a.AuthUnregister() })},
	)
	return cmd
}
