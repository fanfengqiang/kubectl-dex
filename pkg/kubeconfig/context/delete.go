package context

import (
	"fmt"
	"log"
	"os/exec"
)

// Delete is used to delete a context.
func Delete(user, cluster string) {

	cmd := exec.Command("kubectl", "config", "delete-context", user+"/"+cluster)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error occurred when delete context:\n %s\n", err)
	}
	fmt.Printf("Dlete context output:\n%s\n", string(out))
}
