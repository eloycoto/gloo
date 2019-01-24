package remove

import (
	"github.com/solo-io/go-utils/cliutils"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/flagutils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/surveyutils"
	"github.com/spf13/cobra"
)

func Route(opts *options.Options, optionsFunc ...cliutils.OptionsFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "route",
		Aliases: []string{"r", "routes"},
		Short:   "Remove a Route from a Virtual Service",
		Long: "Routes match patterns on requests and indicate the type of action to take when a proxy receives " +
			"a matching request. Requests can be broken down into their Match and Action components. " +
			"The order of routes within a Virtual Service matters. The first route in the virtual service " +
			"that matches a given request will be selected for routing. \n\n" +
			"If no virtual service is specified for this command, glooctl add route will attempt to add it to a " +
			"default virtualservice with domain '*'. if one does not exist, it will be created for you.\n\n" +
			"" +
			"Usage: `glooctl rm route [--name virtual-service-name] [--namespace namespace] [--index=X]`",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opts.Top.Interactive {
				if err := surveyutils.RemoveRouteFlagsInteractive(opts); err != nil {
					return err
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeRoute(opts)
		},
	}
	pflags := cmd.PersistentFlags()
	flagutils.AddOutputFlag(pflags, &opts.Top.Output)
	flagutils.RemoveRouteFlags(pflags, &opts.Remove.Route)
	cliutils.ApplyOptions(cmd, optionsFunc)
	return cmd
}

func removeRoute(opts *options.Options) error {
	index := int(opts.Remove.Route.RemoveIndex)

	virtualService, err := helpers.MustVirtualServiceClient().Read(opts.Metadata.Namespace, opts.Metadata.Name,
		clients.ReadOpts{Ctx: opts.Top.Ctx})
	if err != nil {
		return err
	}

	virtualService.VirtualHost.Routes = append(virtualService.VirtualHost.Routes[:index], virtualService.VirtualHost.Routes[:index]...)

	out, err := helpers.MustVirtualServiceClient().Write(virtualService, clients.WriteOpts{
		Ctx:               opts.Top.Ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return err
	}

	helpers.PrintVirtualServices(gatewayv1.VirtualServiceList{out}, opts.Top.Output)
	return nil
}
