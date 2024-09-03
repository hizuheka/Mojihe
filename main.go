package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&Utf16leCmd{}, "")

	isDebug := flag.Bool("d", false, "debugログを出力")
	flag.Parse()

	// ログレベルの設定
	switch {
	case *isDebug:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
	ctx := context.Background()

	os.Exit(int(subcommands.Execute(ctx)))
}
