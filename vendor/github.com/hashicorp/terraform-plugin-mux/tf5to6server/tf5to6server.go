// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tf5to6server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov5tov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov6tov5"
)

// UpgradeServer wraps a protocol version 5 ProviderServer in a protocol
// version 6 server. Protocol version 6 is fully forwards compatible with
// protocol version 5, so no additional validation is performed.
//
// Protocol version 6 servers require Terraform CLI 1.0 or later.
//
// Terraform CLI 1.1.5 or later is required for terraform-provider-sdk based
// protocol version 5 servers to properly upgrade to protocol version 6.
func UpgradeServer(_ context.Context, v5server func() tfprotov5.ProviderServer) (tfprotov6.ProviderServer, error) {
	return v5tov6Server{
		v5Server: v5server(),
	}, nil
}

var _ tfprotov6.ProviderServer = v5tov6Server{}

// Temporarily verify that v5tov6Server implements new RPCs correctly.
// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
var (
	_ tfprotov6.FunctionServer = v5tov6Server{}
	//nolint:staticcheck // Intentional verification of interface implementation.
	_ tfprotov6.ResourceServerWithMoveResourceState = v5tov6Server{}
)

type v5tov6Server struct {
	v5Server tfprotov5.ProviderServer
}

func (s v5tov6Server) ApplyResourceChange(ctx context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	v5Req := tfprotov6tov5.ApplyResourceChangeRequest(req)
	v5Resp, err := s.v5Server.ApplyResourceChange(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ApplyResourceChangeResponse(v5Resp), nil
}

func (s v5tov6Server) CallFunction(ctx context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	// Remove and call s.v5Server.CallFunction below directly.
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	functionServer, ok := s.v5Server.(tfprotov5.FunctionServer)

	if !ok {
		v6Resp := &tfprotov6.CallFunctionResponse{
			Error: &tfprotov6.FunctionError{
				Text: "Provider Functions Not Implemented: A provider-defined function call was received by the provider, however the provider does not implement functions. " +
					"Either upgrade the provider to a version that implements provider-defined functions or this is a bug in Terraform that should be reported to the Terraform maintainers.",
			},
		}

		return v6Resp, nil
	}

	v5Req := tfprotov6tov5.CallFunctionRequest(req)
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	// v5Resp, err := s.v5Server.CallFunction(ctx, v5Req)
	v5Resp, err := functionServer.CallFunction(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.CallFunctionResponse(v5Resp), nil
}

func (s v5tov6Server) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	v5Req := tfprotov6tov5.ConfigureProviderRequest(req)
	v5Resp, err := s.v5Server.ConfigureProvider(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ConfigureProviderResponse(v5Resp), nil
}

func (s v5tov6Server) GetFunctions(ctx context.Context, req *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	// Remove and call s.v5Server.GetFunctions below directly.
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	functionServer, ok := s.v5Server.(tfprotov5.FunctionServer)

	if !ok {
		v6Resp := &tfprotov6.GetFunctionsResponse{
			Functions: map[string]*tfprotov6.Function{},
		}

		return v6Resp, nil
	}

	v5Req := tfprotov6tov5.GetFunctionsRequest(req)
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/210
	// v5Resp, err := s.v5Server.GetFunctions(ctx, v5Req)
	v5Resp, err := functionServer.GetFunctions(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.GetFunctionsResponse(v5Resp), nil
}

func (s v5tov6Server) GetMetadata(ctx context.Context, req *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	v5Req := tfprotov6tov5.GetMetadataRequest(req)
	v5Resp, err := s.v5Server.GetMetadata(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.GetMetadataResponse(v5Resp), nil
}

func (s v5tov6Server) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	v5Req := tfprotov6tov5.GetProviderSchemaRequest(req)
	v5Resp, err := s.v5Server.GetProviderSchema(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.GetProviderSchemaResponse(v5Resp), nil
}

func (s v5tov6Server) ImportResourceState(ctx context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	v5Req := tfprotov6tov5.ImportResourceStateRequest(req)
	v5Resp, err := s.v5Server.ImportResourceState(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ImportResourceStateResponse(v5Resp), nil
}

func (s v5tov6Server) MoveResourceState(ctx context.Context, req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	// Remove and call s.v5Server.MoveResourceState below directly.
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
	//nolint:staticcheck // Intentional verification of interface implementation.
	resourceServer, ok := s.v5Server.(tfprotov5.ResourceServerWithMoveResourceState)

	if !ok {
		v6Resp := &tfprotov6.MoveResourceStateResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "MoveResourceState Not Implemented",
					Detail: "A MoveResourceState call was received by the provider, however the provider does not implement the RPC. " +
						"Either upgrade the provider to a version that implements MoveResourceState or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return v6Resp, nil
	}

	v5Req := tfprotov6tov5.MoveResourceStateRequest(req)
	// Reference: https://github.com/hashicorp/terraform-plugin-mux/issues/219
	// v5Resp, err := s.v5Server.MoveResourceState(ctx, v5Req)
	v5Resp, err := resourceServer.MoveResourceState(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.MoveResourceStateResponse(v5Resp), nil
}

func (s v5tov6Server) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	v5Req := tfprotov6tov5.PlanResourceChangeRequest(req)
	v5Resp, err := s.v5Server.PlanResourceChange(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.PlanResourceChangeResponse(v5Resp), nil
}

// ProviderServer is a function compatible with tf6server.Serve.
func (s v5tov6Server) ProviderServer() tfprotov6.ProviderServer {
	return s
}

func (s v5tov6Server) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	v5Req := tfprotov6tov5.ReadDataSourceRequest(req)
	v5Resp, err := s.v5Server.ReadDataSource(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ReadDataSourceResponse(v5Resp), nil
}

func (s v5tov6Server) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	v5Req := tfprotov6tov5.ReadResourceRequest(req)
	v5Resp, err := s.v5Server.ReadResource(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ReadResourceResponse(v5Resp), nil
}

func (s v5tov6Server) StopProvider(ctx context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	v5Req := tfprotov6tov5.StopProviderRequest(req)
	v5Resp, err := s.v5Server.StopProvider(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.StopProviderResponse(v5Resp), nil
}

func (s v5tov6Server) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	v5Req := tfprotov6tov5.UpgradeResourceStateRequest(req)
	v5Resp, err := s.v5Server.UpgradeResourceState(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.UpgradeResourceStateResponse(v5Resp), nil
}

func (s v5tov6Server) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	v5Req := tfprotov6tov5.ValidateDataSourceConfigRequest(req)
	v5Resp, err := s.v5Server.ValidateDataSourceConfig(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ValidateDataResourceConfigResponse(v5Resp), nil
}

func (s v5tov6Server) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	v5Req := tfprotov6tov5.PrepareProviderConfigRequest(req)
	v5Resp, err := s.v5Server.PrepareProviderConfig(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ValidateProviderConfigResponse(v5Resp), nil
}

func (s v5tov6Server) ValidateResourceConfig(ctx context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	v5Req := tfprotov6tov5.ValidateResourceTypeConfigRequest(req)
	v5Resp, err := s.v5Server.ValidateResourceTypeConfig(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ValidateResourceConfigResponse(v5Resp), nil
}
