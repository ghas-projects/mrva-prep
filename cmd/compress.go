package cmd

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Gzip-compress the database for serving to the browser",
	Long: `Creates a .gz copy of the SQLite database using maximum compression (gzip
level 9). The Blazor WASM UI fetches the compressed file and decompresses
it client-side via the browser's native DecompressionStream API.

The original uncompressed database is left in place.`,
	RunE: runCompress,
}

func init() {
	rootCmd.AddCommand(compressCmd)
}

func runCompress(cmd *cobra.Command, args []string) error {
	outPath := dbPath + ".gz"

	fmt.Printf("  compressing %s -> %s ... ", dbPath, outPath)

	src, err := os.Open(dbPath)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("open source: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(outPath)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("create destination: %w", err)
	}
	defer dst.Close()

	gw, err := gzip.NewWriterLevel(dst, gzip.BestCompression)
	if err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("create gzip writer: %w", err)
	}

	if _, err := io.Copy(gw, src); err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("compress: %w", err)
	}

	if err := gw.Close(); err != nil {
		fmt.Println("FAILED")
		return fmt.Errorf("finalize gzip: %w", err)
	}

	srcInfo, _ := src.Stat()
	dstInfo, _ := dst.Stat()
	ratio := float64(dstInfo.Size()) / float64(srcInfo.Size()) * 100

	fmt.Printf("ok (%.1f%% of original)\n", ratio)
	return nil
}
