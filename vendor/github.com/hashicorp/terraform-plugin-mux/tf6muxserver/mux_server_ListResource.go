// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) ListResource(ctx context.Context, req *tfprotov6.ListResourceRequest) (*tfprotov6.ListResourceServerStream, error) {
	rpc := "ListResource"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getListResourceServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	// If there is an error diagnostic, the stream will return a single ListResourceResult with the error diagnostic
	// this should help to make the error more readable and keep the stream from starting if there is a problem.

	if diagnosticsHasError(diags) {
		return &tfprotov6.ListResourceServerStream{
			Results: slices.Values([]tfprotov6.ListResourceResult{
				{
					Diagnostics: diags,
				},
			}),
		}, nil
	}

	// TODO: Remove and call server.ListResource below directly once interface becomes required.
	listResourceServer, ok := server.(tfprotov6.ListResourceServer)
	if !ok {
		resp := &tfprotov6.ListResourceServerStream{
			Results: slices.Values([]tfprotov6.ListResourceResult{
				{
					Diagnostics: []*tfprotov6.Diagnostic{
						{
							Severity: tfprotov6.DiagnosticSeverityError,
							Summary:  "ListResource Not Implemented",
							Detail: "A ListResource call was received by the provider, however the provider does not implement ListResource. " +
								"Either upgrade the provider to a version that implements ListResource or this is a bug in Terraform that should be reported to the Terraform maintainers.",
						},
					},
				},
			}),
		}

		return resp, nil
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return listResourceServer.ListResource(ctx, req)
}
