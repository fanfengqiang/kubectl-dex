package context

import (
	"log"
	"os/exec"
	"strings"
)

// Get is used to get the current context.
func Get() (user, cluster string) {

	cmd := exec.Command("kubectl", "config", "current-context")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	output := strings.Split(string(out), "/")

	return output[0], output[1]
}
