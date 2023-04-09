package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/your-username/dirfuzz/fuzz"
	"github.com/your-username/dirfuzz/output"
)

func main() {
	var options fuzz.Options

	var rootCmd = &cobra.Command{
		Use:   "dirfuzz",
		Short: "A directory fuzzer written in Go",
		Long:  `dirfuzz is a tool for recursively scanning web directories and files for hidden content and misconfigurations.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := options.Validate(); err != nil {
				log.Fatal(err)
			}
			scanner := fuzz.NewScanner(options)
			results, err := scanner.Scan()
			if err != nil {
				log.Fatal(err)
			}
			if err := output.WriteResults(options.OutputFormat, results, options.OutputFile); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.Flags().StringVarP(&options.TargetURL, "url", "u", "", "The target URL to scan")
	rootCmd.Flags().StringVarP(&options.WordlistFile, "wordlist", "w", "", "The path to the wordlist file")
	rootCmd.Flags().IntVarP(&options.Threads, "threads", "t", 10, "The number of threads to use")
	rootCmd.Flags().IntVarP(&options.Timeout, "timeout", "T", 10, "The request timeout in seconds")
	rootCmd.Flags().StringVarP(&options.Extensions, "extensions", "x", "", "A comma-separated list of file extensions to scan")
	rootCmd.Flags().StringVarP(&options.IgnoreRegex, "ignore", "i", "", "A regular expression to ignore certain responses")
	rootCmd.Flags().StringVarP(&options.OutputFile, "output", "o", "", "The path to the output file")
	rootCmd.Flags().StringVarP(&options.OutputFormat, "format", "f", "csv", "The output format (csv, json)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
