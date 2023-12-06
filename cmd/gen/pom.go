package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	cmd := exec.Command("mvn", "dependency:copy-dependencies", "-DoutputDirectory=target")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	matches, err := filepath.Glob("target/*")
	if err != nil {
		panic(err)
	}

	jars := make([]string, 0, len(matches))
	for _, match := range matches {
		if strings.Contains(match, ".aar.jar") {
			continue
		}

		if filepath.Ext(match) == ".jar" {
			if strings.Contains(match, "extracted") {
				if err := os.Remove(match); err != nil {
					panic(err)
				}
				continue
			} else {
				jars = append(jars, match)
			}
		}

		jar, err := os.Open(match)
		if err != nil {
			panic(err)
		}

		stats, err := jar.Stat()
		if err != nil {
			panic(err)
		}

		if stats.IsDir() {
			continue
		}

		if filepath.Ext(match) == ".aar" {
			//unzip
			zip, err := zip.OpenReader(match)
			if err != nil {
				panic(err)
			}

			for _, f := range zip.File {
				if f.FileInfo().IsDir() {
					continue
				}

				if filepath.Ext(f.Name) != ".jar" {
					continue
				}

				jar, err := f.Open()
				if err != nil {
					panic(err)
				}

				d := filepath.Join("target", filepath.Base(match)+".jar")
				out, err := os.Create(d)
				if err != nil {
					panic(err)
				}

				if _, err := io.Copy(out, jar); err != nil {
					panic(err)
				}

				jars = append(jars, d)
			}
		}

		out, err := os.Create(filepath.Join("jar", filepath.Base(match)))
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(out, jar); err != nil {
			panic(err)
		}

		jar.Close()
		out.Close()
	}

	for i := range jars {
		jars[i] = "./" + filepath.Join("internal", "android", "target", filepath.Base(jars[i]))
	}

	fmt.Println(strings.Join(jars, ":"))
}
