// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tf5to6server

import (
	"context"
	"fmt"
	"slices"

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
	v5Req := tfprotov6tov5.CallFunctionRequest(req)

	v5Resp, err := s.v5Server.CallFunction(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.CallFunctionResponse(v5Resp), nil
}

func (s v5tov6Server) CloseEphemeralResource(ctx context.Context, req *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	v5Req := tfprotov6tov5.CloseEphemeralResourceRequest(req)

	v5Resp, err := s.v5Server.CloseEphemeralResource(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.CloseEphemeralResourceResponse(v5Resp), nil
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
	v5Req := tfprotov6tov5.GetFunctionsRequest(req)

	v5Resp, err := s.v5Server.GetFunctions(ctx, v5Req)
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

func (s v5tov6Server) GetResourceIdentitySchemas(ctx context.Context, req *tfprotov6.GetResourceIdentitySchemasRequest) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {

	v5Req := tfprotov6tov5.GetResourceIdentitySchemasRequest(req)
	v5Resp, err := s.v5Server.GetResourceIdentitySchemas(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.GetResourceIdentitySchemasResponse(v5Resp), nil
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
	v5Req := tfprotov6tov5.MoveResourceStateRequest(req)

	v5Resp, err := s.v5Server.MoveResourceState(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.MoveResourceStateResponse(v5Resp), nil
}

func (s v5tov6Server) OpenEphemeralResource(ctx context.Context, req *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	v5Req := tfprotov6tov5.OpenEphemeralResourceRequest(req)

	v5Resp, err := s.v5Server.OpenEphemeralResource(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.OpenEphemeralResourceResponse(v5Resp), nil
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

func (s v5tov6Server) RenewEphemeralResource(ctx context.Context, req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	v5Req := tfprotov6tov5.RenewEphemeralResourceRequest(req)

	v5Resp, err := s.v5Server.RenewEphemeralResource(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.RenewEphemeralResourceResponse(v5Resp), nil
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

func (s v5tov6Server) UpgradeResourceIdentity(ctx context.Context, req *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	v5Req := tfprotov6tov5.UpgradeResourceIdentityRequest(req)
	v5Resp, err := s.v5Server.UpgradeResourceIdentity(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.UpgradeResourceIdentityResponse(v5Resp), nil
}

func (s v5tov6Server) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	v5Req := tfprotov6tov5.ValidateDataSourceConfigRequest(req)
	v5Resp, err := s.v5Server.ValidateDataSourceConfig(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ValidateDataResourceConfigResponse(v5Resp), nil
}

func (s v5tov6Server) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	v5Req := tfprotov6tov5.ValidateEphemeralResourceConfigRequest(req)

	v5Resp, err := s.v5Server.ValidateEphemeralResourceConfig(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ValidateEphemeralResourceConfigResponse(v5Resp), nil
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

func (s v5tov6Server) ValidateListResourceConfig(ctx context.Context, req *tfprotov6.ValidateListResourceConfigRequest) (*tfprotov6.ValidateListResourceConfigResponse, error) {
	// TODO: Remove and call s.v5Server.ValidateListResourceConfig below directly once interface becomes required
	listResourceServer, ok := s.v5Server.(tfprotov5.ListResourceServer)
	if !ok {
		v6Resp := &tfprotov6.ValidateListResourceConfigResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "ValidateListResourceConfig Not Implemented",
					Detail: "A ValidateListResourceConfig call was received by the provider, however the provider does not implement the RPC. " +
						"Either upgrade the provider to a version that implements ValidateListResourceConfig or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return v6Resp, nil
	}

	v5Req := tfprotov6tov5.ValidateListResourceConfigRequest(req)

	v5Resp, err := listResourceServer.ValidateListResourceConfig(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ValidateListResourceConfigResponse(v5Resp), nil
}

func (s v5tov6Server) ListResource(ctx context.Context, req *tfprotov6.ListResourceRequest) (*tfprotov6.ListResourceServerStream, error) {
	// TODO: Remove and call s.v5Server.ListResource below directly once interface becomes required
	listResourceServer, ok := s.v5Server.(tfprotov5.ListResourceServer)
	if !ok {
		v6Resp := &tfprotov6.ListResourceServerStream{
			Results: slices.Values([]tfprotov6.ListResourceResult{
				{
					Diagnostics: []*tfprotov6.Diagnostic{
						{
							Severity: tfprotov6.DiagnosticSeverityError,
							Summary:  "ListResource Not Implemented",
							Detail: "A ListResource call was received by the provider, however the provider does not implement the RPC. " +
								"Either upgrade the provider to a version that implements ListResource or this is a bug in Terraform that should be reported to the Terraform maintainers.",
						},
					},
				},
			}),
		}

		return v6Resp, nil
	}

	v5Req := tfprotov6tov5.ListResourceRequest(req)

	v5Resp, err := listResourceServer.ListResource(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ListResourceServerStream(v5Resp), nil
}

func (s v5tov6Server) ValidateActionConfig(ctx context.Context, req *tfprotov6.ValidateActionConfigRequest) (*tfprotov6.ValidateActionConfigResponse, error) {
	// TODO: Remove and call s.v5Server.ValidateActionConfig below directly once interface becomes required
	actionServer, ok := s.v5Server.(tfprotov5.ActionServer)
	if !ok {
		v6Resp := &tfprotov6.ValidateActionConfigResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "ValidateActionConfig Not Implemented",
					Detail: "A ValidateActionConfig call was received by the provider, however the provider does not implement the RPC. " +
						"Either upgrade the provider to a version that implements ValidateActionConfig or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return v6Resp, nil
	}

	v5Req := tfprotov6tov5.ValidateActionConfigRequest(req)

	// v5Resp, err := s.v5Server.ValidateActionConfig(ctx, v5Req)
	v5Resp, err := actionServer.ValidateActionConfig(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.ValidateActionConfigResponse(v5Resp), nil
}

func (s v5tov6Server) PlanAction(ctx context.Context, req *tfprotov6.PlanActionRequest) (*tfprotov6.PlanActionResponse, error) {
	// TODO: Remove and call s.v5Server.PlanAction below directly once interface becomes required
	actionServer, ok := s.v5Server.(tfprotov5.ActionServer)
	if !ok {
		v6Resp := &tfprotov6.PlanActionResponse{
			Diagnostics: []*tfprotov6.Diagnostic{
				{
					Severity: tfprotov6.DiagnosticSeverityError,
					Summary:  "PlanAction Not Implemented",
					Detail: "A PlanAction call was received by the provider, however the provider does not implement the RPC. " +
						"Either upgrade the provider to a version that implements PlanAction or this is a bug in Terraform that should be reported to the Terraform maintainers.",
				},
			},
		}

		return v6Resp, nil
	}

	v5Req := tfprotov6tov5.PlanActionRequest(req)

	// v5Resp, err := s.v5Server.PlanAction(ctx, v5Req)
	v5Resp, err := actionServer.PlanAction(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.PlanActionResponse(v5Resp), nil
}

func (s v5tov6Server) InvokeAction(ctx context.Context, req *tfprotov6.InvokeActionRequest) (*tfprotov6.InvokeActionServerStream, error) {
	// TODO: Remove and call s.v5Server.InvokeAction below directly once interface becomes required
	actionServer, ok := s.v5Server.(tfprotov5.ActionServer)
	if !ok {
		v6Resp := &tfprotov6.InvokeActionServerStream{
			Events: slices.Values([]tfprotov6.InvokeActionEvent{
				{
					Type: tfprotov6.CompletedInvokeActionEventType{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
								Summary:  "InvokeAction Not Implemented",
								Detail: "An InvokeAction call was received by the provider, however the provider does not implement the RPC. " +
									"Either upgrade the provider to a version that implements InvokeAction or this is a bug in Terraform that should be reported to the Terraform maintainers.",
							},
						},
					},
				},
			}),
		}

		return v6Resp, nil
	}

	v5Req := tfprotov6tov5.InvokeActionRequest(req)

	// v5Resp, err := s.v5Server.InvokeAction(ctx, v5Req)
	v5Resp, err := actionServer.InvokeAction(ctx, v5Req)
	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.InvokeActionServerStream(v5Resp), nil
}

func (s v5tov6Server) ValidateStateStoreConfig(ctx context.Context, req *tfprotov6.ValidateStateStoreConfigRequest) (*tfprotov6.ValidateStateStoreConfigResponse, error) {
	return &tfprotov6.ValidateStateStoreConfigResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "ValidateStateStoreConfig Not Implemented",
				Detail: fmt.Sprintf(
					"A ValidateStateStoreConfig call was received by the provider for state store %q, however the underlying provider server was built with protocol version 5 and does not support state stores. "+
						"State stores are a protocol version 6 feature, and this call should not have been made. This is a bug in Terraform or terraform-plugin-mux and should be reported to the provider developers.",
					req.TypeName),
			},
		},
	}, nil
}

func (s v5tov6Server) ConfigureStateStore(ctx context.Context, req *tfprotov6.ConfigureStateStoreRequest) (*tfprotov6.ConfigureStateStoreResponse, error) {
	return &tfprotov6.ConfigureStateStoreResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "ConfigureStateStore Not Implemented",
				Detail: fmt.Sprintf(
					"A ConfigureStateStore call was received by the provider for state store %q, however the underlying provider server was built with protocol version 5 and does not support state stores. "+
						"State stores are a protocol version 6 feature, and this call should not have been made. This is a bug in Terraform or terraform-plugin-mux and should be reported to the provider developers.",
					req.TypeName),
			},
		},
	}, nil
}

func (s v5tov6Server) ReadStateBytes(ctx context.Context, req *tfprotov6.ReadStateBytesRequest) (*tfprotov6.ReadStateBytesStream, error) {
	return &tfprotov6.ReadStateBytesStream{
		Chunks: slices.Values([]tfprotov6.ReadStateByteChunk{
			{
				Diagnostics: []*tfprotov6.Diagnostic{
					{
						Severity: tfprotov6.DiagnosticSeverityError,
						Summary:  "ReadStateBytes Not Implemented",
						Detail: fmt.Sprintf(
							"A ReadStateBytes call was received by the provider for state store %q, however the underlying provider server was built with protocol version 5 and does not support state stores. "+
								"State stores are a protocol version 6 feature, and this call should not have been made. This is a bug in Terraform or terraform-plugin-mux and should be reported to the provider developers.",
							req.TypeName),
					},
				},
			},
		}),
	}, nil
}

func (s v5tov6Server) WriteStateBytes(ctx context.Context, req *tfprotov6.WriteStateBytesStream) (*tfprotov6.WriteStateBytesResponse, error) {
	return &tfprotov6.WriteStateBytesResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "WriteStateBytes Not Implemented",
				Detail: "A WriteStateBytes call was received by the provider, however the underlying provider server was built with protocol version 5 and does not support state stores. " +
					"State stores are a protocol version 6 feature, and this call should not have been made. This is a bug in Terraform or terraform-plugin-mux and should be reported to the provider developers.",
			},
		},
	}, nil
}

func (s v5tov6Server) GetStates(ctx context.Context, req *tfprotov6.GetStatesRequest) (*tfprotov6.GetStatesResponse, error) {
	return &tfprotov6.GetStatesResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "GetStates Not Implemented",
				Detail: fmt.Sprintf(
					"A GetStates call was received by the provider for state store %q, however the underlying provider server was built with protocol version 5 and does not support state stores. "+
						"State stores are a protocol version 6 feature, and this call should not have been made. This is a bug in Terraform or terraform-plugin-mux and should be reported to the provider developers.",
					req.TypeName),
			},
		},
	}, nil
}

func (s v5tov6Server) DeleteState(ctx context.Context, req *tfprotov6.DeleteStateRequest) (*tfprotov6.DeleteStateResponse, error) {
	return &tfprotov6.DeleteStateResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "DeleteState Not Implemented",
				Detail: fmt.Sprintf(
					"A DeleteState call was received by the provider for state store %q, however the underlying provider server was built with protocol version 5 and does not support state stores. "+
						"State stores are a protocol version 6 feature, and this call should not have been made. This is a bug in Terraform or terraform-plugin-mux and should be reported to the provider developers.",
					req.TypeName),
			},
		},
	}, nil
}

func (s v5tov6Server) LockState(ctx context.Context, req *tfprotov6.LockStateRequest) (*tfprotov6.LockStateResponse, error) {
	return &tfprotov6.LockStateResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "LockState Not Implemented",
				Detail: fmt.Sprintf(
					"A LockState call was received by the provider for state store %q, however the underlying provider server was built with protocol version 5 and does not support state stores. "+
						"State stores are a protocol version 6 feature, and this call should not have been made. This is a bug in Terraform or terraform-plugin-mux and should be reported to the provider developers.",
					req.TypeName),
			},
		},
	}, nil
}

func (s v5tov6Server) UnlockState(ctx context.Context, req *tfprotov6.UnlockStateRequest) (*tfprotov6.UnlockStateResponse, error) {
	return &tfprotov6.UnlockStateResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "UnlockState Not Implemented",
				Detail: fmt.Sprintf(
					"An UnlockState call was received by the provider for state store %q, however the underlying provider server was built with protocol version 5 and does not support state stores. "+
						"State stores are a protocol version 6 feature, and this call should not have been made. This is a bug in Terraform or terraform-plugin-mux and should be reported to the provider developers.",
					req.TypeName),
			},
		},
	}, nil
}

func (s v5tov6Server) GenerateResourceConfig(ctx context.Context, req *tfprotov6.GenerateResourceConfigRequest) (*tfprotov6.GenerateResourceConfigResponse, error) {
	v5Req := tfprotov6tov5.GenerateResourceConfigRequest(req)
	v5Resp, err := s.v5Server.GenerateResourceConfig(ctx, v5Req)

	if err != nil {
		return nil, err
	}

	return tfprotov5tov6.GenerateResourceConfigResponse(v5Resp), nil
}
