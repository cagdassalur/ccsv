package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"

	"github.com/zcag/ccsv/util"

	"github.com/spf13/cobra"
)

var (
	left_column string
	right_column string
	column_flag string
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Get diff of csv files based on specified columns, outputs uniq values of left side",
	Long: `Get diff of csv files based on specified columns, outputs uniq values of left side 

	ccsv diff -l 1 -r 4 left.csv right.csv
	ccsv diff -l id -r userid left.csv right.csv
	ccsv diff -c id left.csv right.csv`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("diff command needs two csv files")
		}

		if column_flag == "" && (left_column == "" || right_column == "") {
			return fmt.Errorf("You need to provide either -c for both columns or -l and -r for each column")
		}

		if column_flag != "" {
			left_column = column_flag
			right_column = column_flag
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		right_hashes, err := util.HashCSV(right_column, args[1])
		if err != nil { return err }

		file, err := os.Open(args[0])
		if err != nil { return fmt.Errorf("Failed to open file: %s\n", err) }
		defer file.Close()
		reader := csv.NewReader(file)

		headers, err := reader.Read()
		if err != nil { return err }

		col_index, err := util.ParseColumnFlag(left_column, headers)
		if err != nil { return err }

		writer := csv.NewWriter(os.Stdout)
		record := headers
		for {
			if !slices.Contains(right_hashes, util.Hash(record[col_index])) {
				if err := writer.Write(record); err != nil { return err }
				writer.Flush()
			}

			record, err = reader.Read()
			if err != nil && err.Error() == "EOF" { break }
			if err != nil { return err }
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringVarP(
		&left_column,
		"left_column",
		"l",
		"",
		"name or index of column to diff for left file",
	)
	diffCmd.Flags().StringVarP(
		&right_column,
		"right_column",
		"r",
		"",
		"name or index of column to diff for right file",
	)

	diffCmd.Flags().StringVarP(
		&column_flag,
		"column",
		"c",
		"",
		"name or index of column to diff for both files",
	)
}
