package cmd

import (
	"fmt"

	"github.com/fanfengqiang/kubectl-dex/pkg/idtoken"
	kubecontext "github.com/fanfengqiang/kubectl-dex/pkg/kubeconfig/context"
	"github.com/fanfengqiang/kubectl-dex/pkg/kubeconfig/credential"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	u, p string
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
		fmt.Println("user:", u)
		fmt.Println("password:", p)
		fmt.Println("client_id:", viper.GetString("oidc.client.id"))
		fmt.Println("client_secret:", viper.GetString("oidc.client.secret"))
		fmt.Println("client_redirect_url:", viper.GetString("oidc.client.redirect_url"))
		fmt.Println("issuer:", viper.GetString("oidc.issuer.url"))
		fmt.Println("extra_scopes", viper.GetStringSlice("oidc.extra_scopes"))

		clientID := viper.GetString("oidc.client.id")
		clientSecret := viper.GetString("oidc.client.secret")
		redirectURI := viper.GetString("oidc.client.redirect_url")
		issuerURL := viper.GetString("oidc.issuer.url")
		scopes := []string{"openid", "profile", "email", "offline_access", "groups"}
		scopes = append(scopes, viper.GetStringSlice("oidc.extra_scopes")...)

		app := idtoken.Config(clientID, clientSecret, redirectURI, issuerURL, scopes)
		token, _ := app.GetIDToken(u, p)
		rawIDToken := app.ParseToken(token)
		fmt.Println(rawIDToken)
		credential.Create(u, clientID, issuerURL, rawIDToken)
		kubecontext.Create(u, "kubernetes")
		kubecontext.Use(u, "kubernetes")

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&u, "user", "u", "fanfengqiang@qq.com", "dex user")
	loginCmd.Flags().StringVarP(&p, "password", "p", "fanfengqiang", "dex password")

}
