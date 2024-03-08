// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

// MoveResourceState calls the MoveResourceState method of the underlying
// provider serving the resource.
func (s *muxServer) MoveResourceState(ctx context.Context, req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	rpc := "MoveResourceState"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getResourceServer(ctx, req.TargetTypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov6.MoveResourceStateResponse{
			Diagnostics: diags,
		}, nil
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)

	// Remove and call server.MoveResourceState below directly.
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
	//nolint:staticcheck // Intentionally verifying interface implementation
	resourceServer, ok := server.(tfprotov6.ResourceServerWithMoveResourceState)

	if !ok {
		resp := &tfprotov6.MoveResourceStateResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "MoveResourceState Not Implemented",
					Detail: "A MoveResourceState call was received by the provider, however the provider does not implement MoveResourceState. " +
						"Either upgrade the provider to a version that implements MoveResourceState or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	logging.MuxTrace(ctx, "calling downstream server")

	// return server.MoveResourceState(ctx, req)
	return resourceServer.MoveResourceState(ctx, req)
}
