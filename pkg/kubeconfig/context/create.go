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
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}
