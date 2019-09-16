package context

import (
	"fmt"
	"log"
	"os/exec"
)

// Create is used to create a new context.
func Create(user, cluster string) {

	cmd := exec.Command("kubectl", "config", "set-context", user+"/"+cluster, "--user="+user, "--cluster="+cluster)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error occurred when create context\n %s\n", err)
	}
	fmt.Printf("Create context output:\n%s\n", string(out))
}
