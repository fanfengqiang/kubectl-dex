package credential

import (
	"fmt"
	"log"
	"os/exec"
)

// Create is used to create a new credential.
func Create(user, clientID, issuerURL, rawIDToken string) {
	cmd := exec.Command("kubectl", "config", "set-credentials", user, "--auth-provider", "oidc", "--auth-provider-arg", "idp-issuer-url="+issuerURL, "--auth-provider-arg", "client-id="+clientID, "--auth-provider-arg", "id-token="+rawIDToken)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("error occurred when create credential\n %s\n", err)
	}
	fmt.Printf("create credential output:\n%s\n", string(out))

}
