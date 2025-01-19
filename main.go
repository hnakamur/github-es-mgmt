package main

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/alecthomas/kong"
)

var cli struct {
	Certificate CertificateCmd `cmd:"" help:"Subdommand for the GHES certificate."`
	Settings    SettingsCmd    `cmd:"" help:"Subdommand for the GHES Settings."`
	Version     VersionCmd     `cmd:"" help:"Show version and exit."`
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
