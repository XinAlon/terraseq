package cmd

import (
	"terraseq/internal"
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var alignFile, outFormat string

var flip bool

var alignCmd = &cobra.Command{
	Use:   "align",
	Short: "Aligns DNA sequences with a reference.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "[INFO] Aligning...")
		if err := align(inFile, inFormat, outFile, outFormat, alignFile); err != nil {
			fmt.Fprintln(os.Stderr, "[WARNING] Error during alignment: %v\n", err)
			return
		}
		fmt.Fprintln(os.Stderr, "[INFO] Alignment completed successfully.")
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if inFile == "" || inFormat == "" || outFile == "" || alignFile == "" {
			cmd.Help()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(alignCmd)

	alignCmd.Flags().StringVarP(&inFile, "inFile", "i", "", "")
	alignCmd.Flags().StringVarP(&inFormat, "inFormat", "f", "", "")
	alignCmd.Flags().StringVarP(&outFile, "outFile", "o", "", "")
	alignCmd.Flags().StringVarP(&outFormat, "outFormat", "t", "23andme", "")
	alignCmd.Flags().StringVarP(&alignFile, "alignFile", "a", "", "")
	alignCmd.Flags().BoolVar(&flip, "flip", false, "")
	alignCmd.MarkFlagRequired("inFile")
	alignCmd.MarkFlagRequired("inFormat")
	alignCmd.MarkFlagRequired("outFile")
	alignCmd.MarkFlagRequired("alignFile")

	alignCmd.SetHelpFunc(AlignHelp)
	alignCmd.SilenceUsage = true
}

func align(inFile, inFormat, outFile, outFormat, alignFile string) error {

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

	templateRecords, err := internal.ParseTemplate(alignFile)
	if err != nil {
		return fmt.Errorf("error parsing template file: %v", err)
	}

	return internal.AlignDNA(result.Data, templateRecords, outFile, outFormat, flip)
}

func AlignHelp(cmd *cobra.Command, args []string) {
	fmt.Fprintln(cmd.OutOrStdout(), "Aligns DNA sequences with a reference.")
	fmt.Fprintln(cmd.OutOrStdout(), "https://github.com/enelsr/terraseq")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "usage: terraseq align [-a|--alignFile FILE] [-i|--inFile FILE] [-f|--inFormat FORMAT]")
	fmt.Fprintln(cmd.OutOrStdout(), "                      [-o|--outFile FILE] (-t|--outFormat FORMAT) (--flip)")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "Parse optional command line arguments.")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "options:")
	fmt.Fprintln(cmd.OutOrStdout(), "  -h, --help                  Display this help message and exit")
	fmt.Fprintln(cmd.OutOrStdout(), "  -a, --alignFile FILE        Specify the path to the alignment file")
	fmt.Fprintln(cmd.OutOrStdout(), "                              (e.g., alignment.bim)")
	fmt.Fprintln(cmd.OutOrStdout(), "  -i, --inFile FILE           Specify the path to the input file")
	fmt.Fprintln(cmd.OutOrStdout(), "                              (e.g., input.txt)")
	fmt.Fprintln(cmd.OutOrStdout(), "  -f, --inFormat FORMAT       Define the format of the input file")
	fmt.Fprintln(cmd.OutOrStdout(), "                              (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage)")
	fmt.Fprintln(cmd.OutOrStdout(), "  -t, --outFormat FORMAT      Define the format of the output file")
	fmt.Fprintln(cmd.OutOrStdout(), "                              (options: 23andme, ancestry, ftdnav1, ftdnav2, myheritage)")
	fmt.Fprintln(cmd.OutOrStdout(), "  -o, --outFile FILE          Specify the path for the output file")
	fmt.Fprintln(cmd.OutOrStdout(), "                              (e.g., output.txt)")
	fmt.Fprintln(cmd.OutOrStdout(), "  --flip                      Flips the alleles in accordance with the reference"")
}
