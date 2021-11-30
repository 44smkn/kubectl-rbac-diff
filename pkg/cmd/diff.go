package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/44smkn/kubectl-role-diff/pkg/cmdutil"
	"github.com/44smkn/kubectl-role-diff/pkg/kubernetes"
	"github.com/44smkn/kubectl-role-diff/pkg/policy"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	clioptions "k8s.io/cli-runtime/pkg/genericclioptions"
	k8ssdk "k8s.io/client-go/kubernetes"
)

const (
	diffUsage = "kubectl roll diff --from <Filename> --to <Filename>"

	fromFlag = "from"
	toFlag   = "to"
)

type DiffOptions struct {
	Out    io.Writer
	ErrOut io.Writer
	Logger *zap.Logger

	FromFilename string
	ToFilename   string
}

func NewCmdDiff(f *cmdutil.Factory, version, buildDate string) *cobra.Command {
	opts := &DiffOptions{
		Out:    f.Out,
		ErrOut: f.ErrOut,
		Logger: f.Logger,
	}
	cmd := &cobra.Command{
		Use:          diffUsage,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return diffRun(opts)
		},
	}

	cmd.SetOut(f.Out)
	cmd.SetErr(f.ErrOut)

	cmd.AddCommand(NewCmdVersion(version, buildDate))
	cmd.PersistentFlags().Bool("help", false, "Show help for command")

	cmd.Flags().StringVar(&opts.FromFilename, fromFlag, "", "from file name")
	cmd.Flags().StringVar(&opts.ToFilename, toFlag, "", "to file name")

	return cmd
}

func diffRun(opts *DiffOptions) error {
	restConfig, err := clioptions.NewConfigFlags(true).ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		return fmt.Errorf("failed to retrieve client config for kubernetes: %w", err)
	}
	client, err := k8ssdk.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("failed to build kubernetes client: %w", err)
	}

	srf := kubernetes.NewServerResourceFetcher(client)
	apiResources, err := srf.Fetch()
	if err != nil {
		return fmt.Errorf("failed to fetch apiResources: %w", err)
	}

	dir, err := os.MkdirTemp("", "kubectl-role-diff")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(dir) // clean up

	ptg := policy.NewPolicyTableGenerator(apiResources)
	fromContents, err := policy.ReadFile(opts.FromFilename)
	if err != nil {
		return fmt.Errorf("failed to --from read file: %w", err)
	}
	fromPolicies, err := ptg.Generate(fromContents)
	if err != nil {
		return fmt.Errorf("failed to generate \"from\" table: %w", err)
	}
	fromOut := filepath.Join(dir, "from.yaml")
	policy.WriteTableToFile(fromOut, fromPolicies.Render())

	toContents, err := policy.ReadFile(opts.ToFilename)
	if err != nil {
		return fmt.Errorf("failed to --to read file: %w", err)
	}
	toPolicies, err := ptg.Generate(toContents)
	if err != nil {
		return fmt.Errorf("failed to generate \"to\" table: %w", err)
	}
	toOut := filepath.Join(dir, "to.yaml")
	policy.WriteTableToFile(toOut, toPolicies.Render())

	// fmt.Fprintf(opts.Out, "%s - %s\n", fromOut, toOut)
	out, _ := exec.Command("colordiff", "-U", "2", fromOut, toOut).Output()
	if len(out) == 0 {
		return nil
	}
	fmt.Fprintln(opts.Out, string(out))

	return nil
}
