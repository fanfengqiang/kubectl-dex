package context

import (
	"fmt"
	"log"
	"os/exec"
)

// Use is used to set a context.
func Use(user, cluster string) {

	cmd := exec.Command("kubectl", "config", "use-context", user+"/"+cluster)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}
