package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/alecthomas/kong"
)

var cli struct {
	Certificate CertificateCmd `cmd:"" help:"Subdommand for the GHES certificate."`
	Settings    SettingsCmd    `cmd:"" help:"Subdommand for the GHES Settings."`
	Version     VersionCmd     `cmd:"" help:"Show version and exit."`
}

var programSlogLevel = new(slog.LevelVar)

func setSlogLevelDebug() {
	programSlogLevel.Set(slog.LevelDebug)
}

type Context struct {
	context.Context
}

type VersionCmd struct{}

func (v *VersionCmd) Run(ctx *Context) error {
	fmt.Println(Version())
	return nil
}

func main() {
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programSlogLevel})
	slog.SetDefault(slog.New(h))

	ctx := kong.Parse(&cli)
	err := ctx.Run(&Context{
		Context: context.WithValue(context.Background(), "key1", "value1"),
	})
	ctx.FatalIfErrorf(err)
}

func Version() string {
	// https://blog.lufia.org/entry/2020/12/18/002238
	info, ok := debug.ReadBuildInfo()
	if !ok {
		// Goモジュールが無効など
		return "(devel)"
	}
	return info.Main.Version
}
