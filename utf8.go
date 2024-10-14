package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/google/subcommands"
)

type Utf8Cmd struct {
	input  string
	output string
	henkan string
}

func (*Utf8Cmd) Name() string { return "utf8" }
func (*Utf8Cmd) Synopsis() string {
	return "input ファイル(utf8エンコード)を henkan ファイルを基に変換した結果を output ファイル(utf8エンコード)に出力する"
}
func (*Utf8Cmd) Usage() string {
	return `utf8 -i 変換元ファイル -o 変換結果ファイル -g 変換定義ファイル:
	変換元ファイル(UTF8エンコード)を変換定義ファイル(UTF8エンコード)を基に変換した結果を変換結果ファイルに出力する。
`
}

func (u *Utf8Cmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&u.input, "i", "", "変換元ファイルのパス")
	f.StringVar(&u.output, "o", "", "変換結果ファイルのパス")
	f.StringVar(&u.henkan, "g", "", "変換定義ファイルのパス")
}

func (u *Utf8Cmd) validate() error {
	if u.input == "" {
		return fmt.Errorf("引数 -i が指定されていません。")
	}
	if u.output == "" {
		return fmt.Errorf("引数 -o が指定されていません。")
	}
	if u.henkan == "" {
		return fmt.Errorf("引数 -g が指定されていません。")
	}

	return nil
}

func (u *Utf8Cmd) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
		slog.Info("END utf8-Command")
	}()

	slog.Info("START utf8-Command")

	// 起動時引数のチェック
	if err = u.validate(); err != nil {
		return subcommands.ExitUsageError
	}

	// 変換定義ファイルを読み込む
	var transMap map[rune]rune
	transMap, err = readTransformations(u.henkan)
	if err != nil {
		slog.Error("Error reading transform file")
		return subcommands.ExitFailure
	}

	// 入力ファイルと出力ファイルをオープン
	var inputFile *os.File
	inputFile, err = os.Open(u.input)
	if err != nil {
		slog.Error("Error opening input file:")
		return subcommands.ExitFailure
	}
	defer inputFile.Close()

	var outputFile *os.File
	outputFile, err = os.Create(u.output)
	if err != nil {
		slog.Error("Error creating output file:")
		return subcommands.ExitFailure
	}
	defer outputFile.Close()

	// 入力ファイルをバッファリングして読み込み
	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	// 1ルーンずつ読み込み、変換して出力
	for {
		// 1文字ずつ読み込む
		var r rune
		r, _, err = reader.ReadRune()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			slog.Error("Error reading input file:")
			return subcommands.ExitFailure
		}

		// 変換を適用
		if newR, found := transMap[r]; found {
			r = newR
		}

		// 変換された文字を書き込む
		_, err = writer.WriteRune(r)
		if err != nil {
			slog.Error("Error write output file:")
			return subcommands.ExitFailure
		}
	}

	return subcommands.ExitSuccess
}
