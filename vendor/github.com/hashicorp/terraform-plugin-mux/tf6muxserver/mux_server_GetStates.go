// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) GetStates(ctx context.Context, req *tfprotov6.GetStatesRequest) (*tfprotov6.GetStatesResponse, error) {
	rpc := "GetStates"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getStateStoreServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov6.GetStatesResponse{
			Diagnostics: diags,
		}, nil
	}

	// TODO: Remove and call server.GetStates below directly once interface becomes required.
	stateStoreServer, ok := server.(tfprotov6.StateStoreServer)
	if !ok {
		resp := &tfprotov6.GetStatesResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "GetStates Not Implemented",
					Detail: "A GetStates call was received by the provider, however the provider does not implement GetStates. " +
						"Either upgrade the provider to a version that implements GetStates or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return stateStoreServer.GetStates(ctx, req)
}
