package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("redis-server", "./tools/rdb-test-server/rdb_testing.conf")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("Starting Redis server...")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start Redis server: %v", err)
	}

	log.Printf("Redis server started with PID %d", cmd.Process.Pid)

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Redis server exited with error: %v", err)
	}
}
