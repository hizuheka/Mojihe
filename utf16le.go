package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"unicode/utf16"

	"github.com/google/subcommands"
)

type Utf16leCmd struct {
	input  string
	output string
	henkan string
}

func (*Utf16leCmd) Name() string { return "utf16le" }
func (*Utf16leCmd) Synopsis() string {
	return "input ファイルを henkan ファイルを基に変換した結果を output ファイルに出力する"
}
func (*Utf16leCmd) Usage() string {
	return `utf16le -i 変換元ファイル -o 変換結果ファイル -g 変換定義ファイル:
	変換元ファイルを変換定義ファイルを基に変換した結果を変換結果ファイルに出力する。。
`
}

func (u *Utf16leCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&u.input, "i", "", "変換元ファイルのパス")
	f.StringVar(&u.output, "o", "", "変換結果ファイルのパス")
	f.StringVar(&u.henkan, "g", "", "変換定義ファイルのパス")
}

func (u *Utf16leCmd) validate() error {
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

func (u *Utf16leCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	var err error
	defer func() {
		if err != nil {
			slog.Error(err.Error())
		}
		slog.Info("END utf16le-Command")
	}()

	slog.Info("START find-Command")

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

	// 出力ファイルにBOMを書き込む
	outputFile.Write([]byte{0xFF, 0xFE})

	// 入力ファイルをバッファリングして読み込み
	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	// BOMをスキップ
	bom := make([]byte, 2)
	_, err = io.ReadFull(reader, bom)
	if err != nil {
		slog.Error("Error reading BOM:")
		return subcommands.ExitFailure
	}

	// UTF-16LEを1ルーンずつ読み込み、変換して出力
	for {
		// UTF-16LEで2バイトずつ読み込む
		var word uint16
		word, err = readUTF16WordLE(reader)
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Error reading input file:")
			return subcommands.ExitFailure
		}

		runes := utf16.Decode([]uint16{word})
		r := runes[0]

		// 変換を適用
		if newR, found := transMap[r]; found {
			r = newR
		}

		// 出力ファイルにUTF-16LEで書き込む
		writeUTF16RuneLE(writer, r)
	}

	fmt.Println("Transformation complete.")
	return subcommands.ExitSuccess
}

// 変換定義ファイルを読み込み、変換マップを作成する
func readTransformations(filePath string) (map[rune]rune, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	transMap := make(map[rune]rune)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) < 2 {
			continue // 不正な行をスキップ
		}

		// コードポイントを数値として解釈
		from := parseCodePoint(record[0])
		to := parseCodePoint(record[1])

		transMap[from] = to
	}

	return transMap, nil
}

// コードポイントの文字列をruneに変換
func parseCodePoint(cp string) rune {
	var r rune
	fmt.Sscanf(cp, "%X", &r)
	return r
}

// UTF-16LEで1ワード（2バイト）を読み込む
func readUTF16WordLE(reader *bufio.Reader) (uint16, error) {
	bytes := make([]byte, 2)
	_, err := io.ReadFull(reader, bytes)
	if err != nil {
		return 0, err
	}
	return uint16(bytes[1])<<8 | uint16(bytes[0]), nil
}

// 出力ファイルにUTF-16LEで1ルーン書き込む
func writeUTF16RuneLE(writer *bufio.Writer, r rune) {
	utf16Data := utf16.Encode([]rune{r})
	for _, word := range utf16Data {
		writer.WriteByte(byte(word))
		writer.WriteByte(byte(word >> 8))
	}
}
