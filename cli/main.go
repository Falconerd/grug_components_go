package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	grug_components "github.com/Falconerd/grug_components_go"
)

func compileFile(inputPath, outputPath string) {
	htmlContent := grug_components.CompileHtmlFromFile(inputPath)
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		fmt.Println("Error creating directories:", err)
		return
	}
	if err := os.WriteFile(outputPath, []byte(htmlContent), 0644); err != nil {
		fmt.Println("Error writing to file:", outputPath, err)
	}
}

func compileDirectory(srcDir, destDir string, recursive bool) {
	filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if d.IsDir() && !recursive {
			return filepath.SkipDir
		}
		if filepath.Ext(path) == ".html" {
			relPath, _ := filepath.Rel(srcDir, path)
			outputPath := filepath.Join(destDir, relPath)
			compileFile(path, outputPath)
		}
		return nil
	})
}

func main() {
	dFlag := flag.Bool("d", false, "Compile a whole directory recursively")
	nFlag := flag.Bool("n", false, "Disable recursion")
	flag.Parse()

	if *dFlag {
		if len(flag.Args()) != 2 {
			fmt.Println("Usage: grugc -d src/ build/")
			return
		}
		compileDirectory(flag.Arg(0), flag.Arg(1), !*nFlag)
	} else {
		if len(flag.Args()) != 2 {
			fmt.Println("Usage: grugc index.html build/index.html")
			return
		}
		compileFile(flag.Arg(0), flag.Arg(1))
	}
}
