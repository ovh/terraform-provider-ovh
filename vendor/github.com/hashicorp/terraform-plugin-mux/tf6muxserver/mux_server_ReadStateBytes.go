// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) ReadStateBytes(ctx context.Context, req *tfprotov6.ReadStateBytesRequest) (*tfprotov6.ReadStateBytesStream, error) {
	rpc := "ReadStateBytes"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	server, diags, err := s.getStateStoreServer(ctx, req.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov6.ReadStateBytesStream{
			Chunks: slices.Values([]tfprotov6.ReadStateByteChunk{{
				Diagnostics: diags,
			}}),
		}, nil
	}

	// TODO: Remove and call server.ReadStateBytes below directly once interface becomes required.
	stateStoreServer, ok := server.(tfprotov6.StateStoreServer)
	if !ok {
		resp := &tfprotov6.ReadStateBytesStream{
			Chunks: slices.Values([]tfprotov6.ReadStateByteChunk{{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "ReadStateBytes Not Implemented",
						Detail: "A ReadStateBytes call was received by the provider, however the provider does not implement ReadStateBytes. " +
							"Either upgrade the provider to a version that implements ReadStateBytes or this is a bug in Terraform that should be reported to the Terraform maintainers.",
					},
				},
			}}),
		}

		return resp, nil
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return stateStoreServer.ReadStateBytes(ctx, req)
}
