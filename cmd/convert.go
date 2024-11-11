package cmd

import (
	"terraseq/internal"
	"github.com/spf13/cobra"
	"os"
	"fmt"
)

var inFile, inFormat, outFile string

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts a DNA file to another format.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "[INFO] Converting to %s...\n", outFormat)
		if err := convert(inFile, inFormat, outFile, outFormat); err != nil {
			fmt.Fprintln(os.Stderr, "[WARNING] Error during conversion: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stderr, "[INFO] Conversion completed successfully.\n")
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if inFile == "" || inFormat == "" || outFile == "" || outFormat == "" {
			cmd.Help()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&inFile, "inFile", "i", "", "Path to input file")
	convertCmd.Flags().StringVarP(&inFormat, "inFormat", "f", "", "Format of input file")
	convertCmd.Flags().StringVarP(&outFile, "outFile", "o", "", "Path to output file")
	convertCmd.Flags().StringVarP(&outFormat, "outFormat", "t", "23andme", "Format of output file")
	convertCmd.MarkFlagRequired("inFile")
	convertCmd.MarkFlagRequired("inFormat")
	convertCmd.MarkFlagRequired("outFile")
	convertCmd.MarkFlagRequired("outFormat")

	convertCmd.SetHelpFunc(ConvertHelp)
	convertCmd.SilenceUsage = true
}

func convert(inFile, inFormat, outFile, outFormat string) error {
	var result internal.ParseResult

	switch inFormat {
		case "ancestry":
			result = internal.ParseAncestryDNA(inFile)
		case "23andme":
			result = internal.Parse23andMe(inFile)
		case "ftdnav2":
			result = internal.ParseFTDNA(inFile)
		case "ftdnav1", "myheritage":
			result = internal.ParseMyHeritage(inFile)
		default:
			return fmt.Errorf("unsupported input format: %s", inFormat)
	}

	if result.Err != nil {
		return result.Err
	}

	return internal.WriteDNAData(result.Data, outFile, outFormat)
}

func ConvertHelp(cmd *cobra.Command, args []string) {
	fmt.Fprintln(cmd.OutOrStdout(), "Converts a DNA file to another format.")
	fmt.Fprintln(cmd.OutOrStdout(), "https://github.com/enelsr/terraseq")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "usage: terraseq convert [-i|--inFile FILE] [-f|--inFormat FORMAT]")
	fmt.Fprintln(cmd.OutOrStdout(), "                      [-o|--outFile FILE] [-t|--outFormat FORMAT]")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "Parse optional command line arguments.")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "options:")
	fmt.Fprintln(cmd.OutOrStdout(), "  -h, --help                  Display this help message and exit")
	fmt.Fprintln(cmd.OutOrStdout(), "  -i, --inFile FILE           Specify the path to the input file")
	fmt.Fprintln(cmd.OutOrStdout(), "                              (e.g., input.txt)")
	fmt.Fprintln(cmd.OutOrStdout(), "  -f, --inFormat FORMAT       Define the input file format")
	fmt.Fprintln(cmd.OutOrStdout(), "                              (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage)")
	fmt.Fprintln(cmd.OutOrStdout(), "  -o, --outFile FILE          Specify the path for the output file")
	fmt.Fprintln(cmd.OutOrStdout(), "                              (e.g., output.txt)")
	fmt.Fprintln(cmd.OutOrStdout(), "  -t, --outFormat FORMAT      Define the output file format")
	fmt.Fprintln(cmd.OutOrStdout(), "                              (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage)")
}
