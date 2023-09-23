// Package plugin assists at building the plugin's main
package plugin

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CmdName returns the arg[0] of this executable
func CmdName() string {
	return filepath.Base(os.Args[0])
}

type runCmd func(stdin io.ReadCloser, stdout io.WriteCloser) error

// Config specifies how the plugin operates
type Config struct {
	Name    string // Name is the intended executable's name
	Short   string // Short is the short description shown on help
	Version string // Version indicates the version of the generator

	// Run represents the main loop of the generator, returning the
	// exit code
	Run func(io.ReadCloser, io.WriteCloser) int

	// RunE is an alternative to Run that returns an error directly
	RunE func(io.ReadCloser, io.WriteCloser) error
}

// SetDefaults attempts to fill possible gaps in the config
func (cfg *Config) SetDefaults() {
	if cfg.Name == "" {
		cfg.Name = CmdName()
	}

	switch {
	case cfg.RunE != nil:
		// ready
	case cfg.Run != nil:
		// convert Run() to RunE()
		cfg.RunE = func(in io.ReadCloser, out io.WriteCloser) error {
			code := cfg.Run(in, out)

			return &ExitError{Code: code}
		}
	default:
		// generate RunE() placeholder
		cfg.RunE = func(in io.ReadCloser, out io.WriteCloser) error {
			return fmt.Errorf("%s protoc plugin not implemented", cfg.Name)
		}
	}
}

func openFileFlag(flags *pflag.FlagSet, name string, flag int, perm fs.FileMode) (*os.File, error) {
	if flags.Changed(name) {
		s, err := flags.GetString(name)
		switch {
		case err != nil:
			return nil, err
		case s == "-":
			return nil, nil
		default:
			return os.OpenFile(s, flag, perm)
		}
	}
	return nil, nil
}

// NewRoot generates a root [cobra.Command] for the plugin's main
func NewRoot(cfg *Config) (*cobra.Command, error) {
	if cfg == nil {
		cfg = new(Config)
	}

	cfg.SetDefaults()

	cmd := &cobra.Command{
		Use:     cfg.Name,
		Short:   cfg.Short,
		Version: cfg.Version,
		RunE:    newRootRun(cfg.RunE),
	}

	// stdin/stdout
	flags := cmd.LocalFlags()
	flags.StringP("input", "f", "", "file to use instead of stdin")
	flags.StringP("output", "o", "", "file to use instead of stdout")

	return cmd, nil
}

func newRootRun(runE runCmd) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		flags := cmd.LocalFlags()

		// stdin
		in, err := openFileFlag(flags, "input", os.O_RDONLY, 0)
		switch {
		case err != nil:
			return err
		case in != nil:
			defer in.Close()
		default:
			in = os.Stdin
		}

		// stdout
		out, err := openFileFlag(flags, "output", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		switch {
		case err != nil:
			return err
		case out != nil:
			defer out.Close()
		default:
			out = os.Stdout
		}

		// run plugin
		return runE(in, out)
	}
}
