package lock

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestFileLock(t *testing.T) {
	if os.Getenv("ENVOY_LOCK_HELPER") == "1" {
		lockPath := os.Getenv("ENVOY_LOCK_PATH")

		lock, err := Acquire(lockPath)
		if err != nil {
			t.Fatal(err)
		}
		defer lock.Release()

		time.Sleep(2 * time.Second)
		return
	}

	lockPath := t.TempDir() + "/test.lock"

	cmdA := exec.Command(os.Args[0], "-test.run=TestFileLock")
	cmdA.Env = append(
		os.Environ(),
		"ENVOY_LOCK_HELPER=1",
		"ENVOY_LOCK_PATH="+lockPath,
	)

	if err := cmdA.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)

	start := time.Now()

	cmdB := exec.Command(os.Args[0], "-test.run=TestFileLock")
	cmdB.Env = append(
		os.Environ(),
		"ENVOY_LOCK_HELPER=1",
		"ENVOY_LOCK_PATH="+lockPath,
	)

	if err := cmdB.Run(); err != nil {
		t.Fatal(err)
	}

	elapsed := time.Since(start)

	if elapsed < 1*time.Second {
		t.Fatalf("Process B did not wait for the lock: %v", elapsed)
	}

	if err := cmdA.Wait(); err != nil {
		t.Fatal(err)
	}
}
