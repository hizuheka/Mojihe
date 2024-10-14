package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
)

// 変換定義ファイルを読み込み、変換マップを作成する
func readTransformations(filePath string) (map[rune]rune, error) {
	slog.Debug("START readTransformations")
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
