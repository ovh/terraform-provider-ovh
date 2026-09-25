// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf6muxserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/logging"
)

func (s *muxServer) WriteStateBytes(ctx context.Context, req *tfprotov6.WriteStateBytesStream) (*tfprotov6.WriteStateBytesResponse, error) {
	rpc := "WriteStateBytes"
	ctx = logging.InitContext(ctx)
	ctx = logging.RpcContext(ctx, rpc)

	// The first chunk in the streamed request will contain the information we need to route (i.e. state store type name)
	firstChunk, firstDiags, wrapped := peekWriteStateBytesStream(req)
	if wrapped == nil {
		return &tfprotov6.WriteStateBytesResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "WriteStateBytes Invalid Request",
					Detail:   "A WriteStateBytes call was received without any state byte chunks.",
				},
			},
		}, nil
	}

	if diagnosticsHasError(firstDiags) {
		return &tfprotov6.WriteStateBytesResponse{
			Diagnostics: firstDiags,
		}, nil
	}

	if firstChunk == nil || firstChunk.Meta == nil || firstChunk.Meta.TypeName == "" {
		return &tfprotov6.WriteStateBytesResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "WriteStateBytes Invalid Request",
					Detail:   "A WriteStateBytes call was received without the required state store metadata in the first chunk.",
				},
			},
		}, nil
	}

	server, diags, err := s.getStateStoreServer(ctx, firstChunk.Meta.TypeName)

	if err != nil {
		return nil, err
	}

	if diagnosticsHasError(diags) {
		return &tfprotov6.WriteStateBytesResponse{
			Diagnostics: diags,
		}, nil
	}

	// TODO: Remove and call server.WriteStateBytes below directly once interface becomes required.
	stateStoreServer, ok := server.(tfprotov6.StateStoreServer)
	if !ok {
		resp := &tfprotov6.WriteStateBytesResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "WriteStateBytes Not Implemented",
					Detail: "A WriteStateBytes call was received by the provider, however the provider does not implement WriteStateBytes. " +
						"Either upgrade the provider to a version that implements WriteStateBytes or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return resp, nil
	}

	ctx = logging.Tfprotov6ProviderServerContext(ctx, server)
	logging.MuxTrace(ctx, "calling downstream server")

	return stateStoreServer.WriteStateBytes(ctx, wrapped)
}

type streamedChunk struct {
	chunk       *tfprotov6.WriteStateBytesChunk
	diagnostics []*tfprotov6.Diagnostic
}

// peekWriteStateBytesStream reads the first chunk from the WriteStateBytes streaming request to enable
// routing based on its metadata/type name, while preserving the complete stream for the downstream server.
//
// It returns the first chunk and its diagnostics for inspection, plus a wrapped stream that yields
// the first chunk followed by all remaining chunks.
func peekWriteStateBytesStream(req *tfprotov6.WriteStateBytesStream) (*tfprotov6.WriteStateBytesChunk, []*tfprotov6.Diagnostic, *tfprotov6.WriteStateBytesStream) {
	if req == nil || req.Chunks == nil {
		return nil, nil, nil
	}

	// Using a goroutine + channels here allows the mux server to route based on the first chunk's TypeName
	// without consuming the entire iterator.
	streamedChunks := make(chan streamedChunk, 1)
	done := make(chan struct{})

	go func() {
		defer close(streamedChunks)
		req.Chunks(func(chunk *tfprotov6.WriteStateBytesChunk, diags []*tfprotov6.Diagnostic) bool {
			select {
			case streamedChunks <- streamedChunk{chunk: chunk, diagnostics: diags}:
				return true
			case <-done:
				return false
			}
		})
	}()

	firstChunk, ok := <-streamedChunks
	if !ok {
		close(done)
		return nil, nil, nil
	}

	newStream := &tfprotov6.WriteStateBytesStream{
		Chunks: func(yield func(*tfprotov6.WriteStateBytesChunk, []*tfprotov6.Diagnostic) bool) {
			defer close(done)
			if !yield(firstChunk.chunk, firstChunk.diagnostics) {
				return
			}

			for streamedChunk := range streamedChunks {
				if !yield(streamedChunk.chunk, streamedChunk.diagnostics) {
					return
				}
			}
		},
	}

	return firstChunk.chunk, firstChunk.diagnostics, newStream
}
