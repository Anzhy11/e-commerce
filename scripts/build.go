package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	_ = os.MkdirAll("bin", 0o755)

	entries, err := os.ReadDir("cmd")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		name := e.Name()
		log.Println("Building", name)

		cmd := exec.Command("go", "build", "-o", filepath.Join("bin", name), "./cmd/"+name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
