package actions

import (
	"context"
	"errors"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/platform/provider"
	"github.com/AlexsJones/gravitywell/platform/provider/digitalocean"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	"github.com/digitalocean/godo"
	"github.com/fatih/color"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"os"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
func NewDigitalOceanConfig() (*digitalocean.DigitalOceanProvider, error) {

	token := os.Getenv("DIGITAL_OCEAN_TOKEN")
	if token == "" {
		return nil, errors.New("no DIGITAL_OCEAN_TOKEN")
	}
	dop := digitalocean.DigitalOceanProvider{}
	tokenSource := &TokenSource{
		AccessToken: token,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)
	dop.ClusterManagerClient = client
	dop.Context = context.Background()
	return &dop, nil
}
func DigitalOceanClusterProcessor(digitalocean *digitalocean.DigitalOceanProvider,
	commandFlag configuration.CommandFlag,
	cluster kinds.ProviderCluster) error {

	create := func() error {

		err := provider.Create(digitalocean, cluster)

		if err != nil {
			return err
		}
		// Run post install -----------------------------------------------------
		for _, executeCommand := range cluster.PostInstallHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}

			}
		}
		return nil
	}
	delete := func() error {
		err := provider.Delete(digitalocean, cluster)
		if err != nil {
			color.Red(err.Error())
		}
		// Run post delete -----------------------------------------------------
		for _, executeCommand := range cluster.PostDeleteHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	// Run Command ------------------------------------------------------------------
	switch commandFlag {
	case configuration.Create:
		return create()
	case configuration.Apply:
		logger.Info("Ignoring apply on cluster - no such option")
		return nil
	case configuration.Replace:
		if err := delete(); err != nil {
			return err
		}
		return create()
	case configuration.Delete:
		return delete()
	}
	return nil
}
