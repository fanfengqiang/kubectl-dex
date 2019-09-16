package credential

import (
	"fmt"
	"log"
	"os/exec"
)

// Delete is used to delete a new credential.
func Delete(user string) {
	cmd := exec.Command("kubectl", "config", "set-credentials", user, "--auth-provider-arg", "id-token=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error occurred when delete credential\n %s\n", err)
	}
	fmt.Printf("delete credential output:\n%s\n", string(out))

}
