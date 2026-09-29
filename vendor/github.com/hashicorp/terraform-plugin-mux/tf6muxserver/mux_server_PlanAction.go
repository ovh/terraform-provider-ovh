// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) PlanAction(ctx context.Context, req *tfprotov6.PlanActionRequest) (*tfprotov6.PlanActionResponse, error) {
	rpc := "PlanAction"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getActionServer(ctx, req.ActionType)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov6.PlanActionResponse{
			Diagnostics: diags,
		}, nil
	}

	// TODO: Remove and call server.PlanAction below directly once interface becomes required.
	actionServer, ok := server.(tfprotov6.ActionServer)
	if !ok {
		resp := &tfprotov6.PlanActionResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "PlanAction Not Implemented",
					Detail: "A PlanAction call was received by the provider, however the provider does not implement PlanAction. " +
						"Either upgrade the provider to a version that implements PlanAction or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return actionServer.PlanAction(ctx, req)
}
