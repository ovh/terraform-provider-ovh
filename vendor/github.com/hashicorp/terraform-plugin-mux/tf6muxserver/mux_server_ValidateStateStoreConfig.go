// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) ValidateStateStoreConfig(ctx context.Context, req *tfprotov6.ValidateStateStoreConfigRequest) (*tfprotov6.ValidateStateStoreConfigResponse, error) {
	rpc := "ValidateStateStoreConfig"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getStateStoreServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov6.ValidateStateStoreConfigResponse{
			Diagnostics: diags,
		}, nil
	}

	// TODO: Remove and call server.ValidateStateStoreConfig below directly once interface becomes required.
	stateStoreServer, ok := server.(tfprotov6.StateStoreServer)
	if !ok {
		resp := &tfprotov6.ValidateStateStoreConfigResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "ValidateStateStoreConfig Not Implemented",
					Detail: "A ValidateStateStoreConfig call was received by the provider, however the provider does not implement ValidateStateStoreConfig. " +
						"Either upgrade the provider to a version that implements ValidateStateStoreConfig or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return stateStoreServer.ValidateStateStoreConfig(ctx, req)
}
