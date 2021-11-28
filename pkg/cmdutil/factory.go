package cmdutil

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
)

type Factory struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer

	Logger     *zap.Logger
	Executable string
}

func NewFactory(appVersion string, logger *zap.Logger) *Factory {

	executable := "kubectl-role-diff"
	if exe, err := os.Executable(); err == nil {
		executable = exe
	}

	return &Factory{
		In:     os.Stdin,
		Out:    colorable.NewColorable(os.Stdout),
		ErrOut: colorable.NewColorable(os.Stderr),

		Logger:     logger,
		Executable: executable,
	}
}
