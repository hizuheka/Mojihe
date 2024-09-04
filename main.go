package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/subcommands"
)

var version string

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&Utf16leCmd{}, "")

	isDebug := flag.Bool("d", false, "debugログを出力")
	isVersion := flag.Bool("v", false, "バージョンを出力")
	flag.Parse()

	// ログレベルの設定
	switch {
	case *isDebug:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
	ctx := context.Background()

	if *isVersion {
		fmt.Printf("mojihe version %s\n", version)
		os.Exit(int(subcommands.ExitSuccess))
	}
	os.Exit(int(subcommands.Execute(ctx)))
}
