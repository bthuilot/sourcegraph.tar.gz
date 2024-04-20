package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/bthuilot/sourcegraph.tar.gz/pkg/sourcegraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"time"
)

func init() {
	rootCmd.PersistentFlags().StringP("query", "q", "", "The SourceGraph query to run. Will be restricted to files only.")
	rootCmd.PersistentFlags().StringP("out", "o", "", "THe name of the file to output the tarball to.")
	rootCmd.PersistentFlags().BoolP("compress", "c", false, "Compress the tarball with gzip")
	_ = rootCmd.MarkPersistentFlagRequired("query")
	_ = viper.BindPFlag("query", rootCmd.PersistentFlags().Lookup("query"))
	_ = viper.BindPFlag("out", rootCmd.PersistentFlags().Lookup("out"))
	viper.MustBindEnv("sourcegraph-token", "SOURCEGRAPH_TOKEN")
}

var rootCmd = &cobra.Command{
	Use:   "sourcegraph-tar",
	Short: "Export SourceGraph file matches to a gunzipped tarball",
	Long:  `Run a SourceGraph query and export each file result to a gunzipped tarball.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		log.Printf("running sourcegraph-tar")

		log.Printf("intilializing SourceGraph client")
		sgCLient := sourcegraph.NewClient(viper.GetString("sourcegraph-token"))

		var output io.Writer
		if viper.IsSet("out") {
			log.Printf("outputting to: %s", viper.GetString("out"))
			output, err = os.Create(viper.GetString("out"))
			if err != nil {
				log.Fatalf("error creating output file: %s", err)
			}
		} else {
			output = os.Stdout
			log.Printf("outputting to stdout")
		}

		if viper.GetBool("compress") {
			log.Printf("compressing tarball with gzip")
			output = gzip.NewWriter(output)
		}

		log.Printf("initializing tarball writer")
		tarWriter := tar.NewWriter(output)

		log.Printf("running query: %s", viper.GetString("query"))

		searchResults, err := sgCLient.SearchFiles(viper.GetString("query"))
		if err != nil {
			log.Fatalf("error initializing SourceGraph client: %s", err)
		}

		for _, file := range searchResults {
			log.Printf("downloading file: %s", file.Path)
			file, err := sgCLient.GetFile(file.Repository, file.Path, file.TotalLines)
			if err != nil {
				log.Fatalf("error downloading file: %s", err)
			}
			path := fmt.Sprintf("%s/%s", file.Repository, file.Path)
			header := &tar.Header{
				Name:    path,
				Size:    int64(file.Size),
				Mode:    0644,
				ModTime: time.Now(),
			}
			err = tarWriter.WriteHeader(header)
			if err != nil {
				log.Printf("error writing tarball header '%s': %s, skipping", path, err)
				continue
			}
			_, err = tarWriter.Write([]byte(file.Contents))
			if err != nil {
				log.Printf("error writing tarball contents '%s': %s, skipping", path, err)
				continue
			}
		}
		// export
		log.Printf("closing tarball writer")
		err = tarWriter.Close()
		if err != nil {
			log.Fatalf("error closing tarball writer: %s", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
