// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov6tov5

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var ErrSchemaAttributeNestedTypeNotImplemented error = errors.New("SchemaAttribute NestedType is not implemented in protocol version 5")

func ApplyResourceChangeRequest(in *tfprotov6.ApplyResourceChangeRequest) *tfprotov5.ApplyResourceChangeRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ApplyResourceChangeRequest{
		Config:         DynamicValue(in.Config),
		PlannedPrivate: in.PlannedPrivate,
		PlannedState:   DynamicValue(in.PlannedState),
		PriorState:     DynamicValue(in.PriorState),
		ProviderMeta:   DynamicValue(in.ProviderMeta),
		TypeName:       in.TypeName,
	}
}

func ApplyResourceChangeResponse(in *tfprotov6.ApplyResourceChangeResponse) *tfprotov5.ApplyResourceChangeResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.ApplyResourceChangeResponse{
		Diagnostics:                 Diagnostics(in.Diagnostics),
		NewState:                    DynamicValue(in.NewState),
		Private:                     in.Private,
		UnsafeToUseLegacyTypeSystem: in.UnsafeToUseLegacyTypeSystem, //nolint:staticcheck
	}
}

func CallFunctionRequest(in *tfprotov6.CallFunctionRequest) *tfprotov5.CallFunctionRequest {
	if in == nil {
		return nil
	}

	out := &tfprotov5.CallFunctionRequest{
		Arguments: make([]*tfprotov5.DynamicValue, 0, len(in.Arguments)),
		Name:      in.Name,
	}

	for _, argument := range in.Arguments {
		out.Arguments = append(out.Arguments, DynamicValue(argument))
	}

	return out
}

func CallFunctionResponse(in *tfprotov6.CallFunctionResponse) *tfprotov5.CallFunctionResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.CallFunctionResponse{
		Error:  FunctionError(in.Error),
		Result: DynamicValue(in.Result),
	}
}

func ConfigureProviderRequest(in *tfprotov6.ConfigureProviderRequest) *tfprotov5.ConfigureProviderRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ConfigureProviderRequest{
		Config:           DynamicValue(in.Config),
		TerraformVersion: in.TerraformVersion,
	}
}

func ConfigureProviderResponse(in *tfprotov6.ConfigureProviderResponse) *tfprotov5.ConfigureProviderResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.ConfigureProviderResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func DataSourceMetadata(in tfprotov6.DataSourceMetadata) tfprotov5.DataSourceMetadata {
	return tfprotov5.DataSourceMetadata{
		TypeName: in.TypeName,
	}
}

func Diagnostics(in []*tfprotov6.Diagnostic) []*tfprotov5.Diagnostic {
	if in == nil {
		return nil
	}

	diags := make([]*tfprotov5.Diagnostic, 0, len(in))

	for _, diag := range in {
		if diag == nil {
			diags = append(diags, nil)
			continue
		}

		diags = append(diags, &tfprotov5.Diagnostic{
			Attribute: diag.Attribute,
			Detail:    diag.Detail,
			Severity:  tfprotov5.DiagnosticSeverity(diag.Severity),
			Summary:   diag.Summary,
		})
	}

	return diags
}

func DynamicValue(in *tfprotov6.DynamicValue) *tfprotov5.DynamicValue {
	if in == nil {
		return nil
	}

	return &tfprotov5.DynamicValue{
		JSON:    in.JSON,
		MsgPack: in.MsgPack,
	}
}

func Function(in *tfprotov6.Function) *tfprotov5.Function {
	if in == nil {
		return nil
	}

	out := &tfprotov5.Function{
		DeprecationMessage: in.DeprecationMessage,
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		Parameters:         make([]*tfprotov5.FunctionParameter, 0, len(in.Parameters)),
		Return:             FunctionReturn(in.Return),
		Summary:            in.Summary,
		VariadicParameter:  FunctionParameter(in.VariadicParameter),
	}

	for _, parameter := range in.Parameters {
		out.Parameters = append(out.Parameters, FunctionParameter(parameter))
	}

	return out
}

func FunctionError(in *tfprotov6.FunctionError) *tfprotov5.FunctionError {
	if in == nil {
		return nil
	}

	out := &tfprotov5.FunctionError{
		Text:             in.Text,
		FunctionArgument: in.FunctionArgument,
	}

	return out
}

func FunctionMetadata(in tfprotov6.FunctionMetadata) tfprotov5.FunctionMetadata {
	return tfprotov5.FunctionMetadata{
		Name: in.Name,
	}
}

func FunctionParameter(in *tfprotov6.FunctionParameter) *tfprotov5.FunctionParameter {
	if in == nil {
		return nil
	}

	return &tfprotov5.FunctionParameter{
		AllowNullValue:     in.AllowNullValue,
		AllowUnknownValues: in.AllowUnknownValues,
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		Name:               in.Name,
		Type:               in.Type,
	}
}

func FunctionReturn(in *tfprotov6.FunctionReturn) *tfprotov5.FunctionReturn {
	if in == nil {
		return nil
	}

	return &tfprotov5.FunctionReturn{
		Type: in.Type,
	}
}

func GetFunctionsRequest(in *tfprotov6.GetFunctionsRequest) *tfprotov5.GetFunctionsRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.GetFunctionsRequest{}
}

func GetFunctionsResponse(in *tfprotov6.GetFunctionsResponse) *tfprotov5.GetFunctionsResponse {
	if in == nil {
		return nil
	}

	functions := make(map[string]*tfprotov5.Function, len(in.Functions))

	for name, function := range in.Functions {
		functions[name] = Function(function)
	}

	return &tfprotov5.GetFunctionsResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
		Functions:   functions,
	}
}

func GetMetadataRequest(in *tfprotov6.GetMetadataRequest) *tfprotov5.GetMetadataRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.GetMetadataRequest{}
}

func GetMetadataResponse(in *tfprotov6.GetMetadataResponse) *tfprotov5.GetMetadataResponse {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.GetMetadataResponse{
		DataSources:        make([]tfprotov5.DataSourceMetadata, 0, len(in.DataSources)),
		Diagnostics:        Diagnostics(in.Diagnostics),
		Functions:          make([]tfprotov5.FunctionMetadata, 0, len(in.Functions)),
		Resources:          make([]tfprotov5.ResourceMetadata, 0, len(in.Resources)),
		ServerCapabilities: ServerCapabilities(in.ServerCapabilities),
	}

	for _, datasource := range in.DataSources {
		resp.DataSources = append(resp.DataSources, DataSourceMetadata(datasource))
	}

	for _, function := range in.Functions {
		resp.Functions = append(resp.Functions, FunctionMetadata(function))
	}

	for _, resource := range in.Resources {
		resp.Resources = append(resp.Resources, ResourceMetadata(resource))
	}

	return resp
}

func GetProviderSchemaRequest(in *tfprotov6.GetProviderSchemaRequest) *tfprotov5.GetProviderSchemaRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.GetProviderSchemaRequest{}
}

func GetProviderSchemaResponse(in *tfprotov6.GetProviderSchemaResponse) (*tfprotov5.GetProviderSchemaResponse, error) {
	if in == nil {
		return nil, nil
	}

	dataSourceSchemas := make(map[string]*tfprotov5.Schema, len(in.DataSourceSchemas))

	for k, v := range in.DataSourceSchemas {
		v5Schema, err := Schema(v)

		if err != nil {
			return nil, fmt.Errorf("unable to convert data source %q schema: %w", k, err)
		}

		dataSourceSchemas[k] = v5Schema
	}

	functions := make(map[string]*tfprotov5.Function, len(in.Functions))

	for name, function := range in.Functions {
		functions[name] = Function(function)
	}

	provider, err := Schema(in.Provider)

	if err != nil {
		return nil, fmt.Errorf("unable to convert provider schema: %w", err)
	}

	providerMeta, err := Schema(in.ProviderMeta)

	if err != nil {
		return nil, fmt.Errorf("unable to convert provider meta schema: %w", err)
	}

	resourceSchemas := make(map[string]*tfprotov5.Schema, len(in.ResourceSchemas))

	for k, v := range in.ResourceSchemas {
		v5Schema, err := Schema(v)

		if err != nil {
			return nil, fmt.Errorf("unable to convert resource %q schema: %w", k, err)
		}

		resourceSchemas[k] = v5Schema
	}

	return &tfprotov5.GetProviderSchemaResponse{
		DataSourceSchemas: dataSourceSchemas,
		Diagnostics:       Diagnostics(in.Diagnostics),
		Functions:         functions,
		Provider:          provider,
		ProviderMeta:      providerMeta,
		ResourceSchemas:   resourceSchemas,
	}, nil
}

func ImportResourceStateRequest(in *tfprotov6.ImportResourceStateRequest) *tfprotov5.ImportResourceStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ImportResourceStateRequest{
		ID:       in.ID,
		TypeName: in.TypeName,
	}
}

func ImportResourceStateResponse(in *tfprotov6.ImportResourceStateResponse) *tfprotov5.ImportResourceStateResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.ImportResourceStateResponse{
		Diagnostics:       Diagnostics(in.Diagnostics),
		ImportedResources: ImportedResources(in.ImportedResources),
	}
}

func ImportedResources(in []*tfprotov6.ImportedResource) []*tfprotov5.ImportedResource {
	if in == nil {
		return nil
	}

	res := make([]*tfprotov5.ImportedResource, 0, len(in))

	for _, imp := range in {
		if imp == nil {
			res = append(res, nil)
			continue
		}

		res = append(res, &tfprotov5.ImportedResource{
			Private:  imp.Private,
			State:    DynamicValue(imp.State),
			TypeName: imp.TypeName,
		})
	}

	return res
}

func MoveResourceStateRequest(in *tfprotov6.MoveResourceStateRequest) *tfprotov5.MoveResourceStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.MoveResourceStateRequest{
		SourcePrivate:         in.SourcePrivate,
		SourceProviderAddress: in.SourceProviderAddress,
		SourceSchemaVersion:   in.SourceSchemaVersion,
		SourceState:           RawState(in.SourceState),
		SourceTypeName:        in.SourceTypeName,
		TargetTypeName:        in.TargetTypeName,
	}
}

func MoveResourceStateResponse(in *tfprotov6.MoveResourceStateResponse) *tfprotov5.MoveResourceStateResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.MoveResourceStateResponse{
		Diagnostics:   Diagnostics(in.Diagnostics),
		TargetPrivate: in.TargetPrivate,
		TargetState:   DynamicValue(in.TargetState),
	}
}

func PlanResourceChangeRequest(in *tfprotov6.PlanResourceChangeRequest) *tfprotov5.PlanResourceChangeRequest {
	if in == nil {
		return nil
	}
	return &tfprotov5.PlanResourceChangeRequest{
		Config:           DynamicValue(in.Config),
		PriorPrivate:     in.PriorPrivate,
		PriorState:       DynamicValue(in.PriorState),
		ProposedNewState: DynamicValue(in.ProposedNewState),
		ProviderMeta:     DynamicValue(in.ProviderMeta),
		TypeName:         in.TypeName,
	}
}

func PlanResourceChangeResponse(in *tfprotov6.PlanResourceChangeResponse) *tfprotov5.PlanResourceChangeResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.PlanResourceChangeResponse{
		Diagnostics:                 Diagnostics(in.Diagnostics),
		PlannedPrivate:              in.PlannedPrivate,
		PlannedState:                DynamicValue(in.PlannedState),
		RequiresReplace:             in.RequiresReplace,
		UnsafeToUseLegacyTypeSystem: in.UnsafeToUseLegacyTypeSystem, //nolint:staticcheck
	}
}

func PrepareProviderConfigRequest(in *tfprotov6.ValidateProviderConfigRequest) *tfprotov5.PrepareProviderConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.PrepareProviderConfigRequest{
		Config: DynamicValue(in.Config),
	}
}

func PrepareProviderConfigResponse(in *tfprotov6.ValidateProviderConfigResponse) *tfprotov5.PrepareProviderConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.PrepareProviderConfigResponse{
		Diagnostics:    Diagnostics(in.Diagnostics),
		PreparedConfig: DynamicValue(in.PreparedConfig),
	}
}

func RawState(in *tfprotov6.RawState) *tfprotov5.RawState {
	if in == nil {
		return nil
	}

	return &tfprotov5.RawState{
		Flatmap: in.Flatmap,
		JSON:    in.JSON,
	}
}

func ReadDataSourceRequest(in *tfprotov6.ReadDataSourceRequest) *tfprotov5.ReadDataSourceRequest {
	if in == nil {
		return nil
	}
	return &tfprotov5.ReadDataSourceRequest{
		Config:       DynamicValue(in.Config),
		ProviderMeta: DynamicValue(in.ProviderMeta),
		TypeName:     in.TypeName,
	}
}

func ReadDataSourceResponse(in *tfprotov6.ReadDataSourceResponse) *tfprotov5.ReadDataSourceResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.ReadDataSourceResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
		State:       DynamicValue(in.State),
	}
}

func ReadResourceRequest(in *tfprotov6.ReadResourceRequest) *tfprotov5.ReadResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ReadResourceRequest{
		CurrentState: DynamicValue(in.CurrentState),
		Private:      in.Private,
		ProviderMeta: DynamicValue(in.ProviderMeta),
		TypeName:     in.TypeName,
	}
}

func ReadResourceResponse(in *tfprotov6.ReadResourceResponse) *tfprotov5.ReadResourceResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.ReadResourceResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
		NewState:    DynamicValue(in.NewState),
		Private:     in.Private,
	}
}

func ResourceMetadata(in tfprotov6.ResourceMetadata) tfprotov5.ResourceMetadata {
	return tfprotov5.ResourceMetadata{
		TypeName: in.TypeName,
	}
}

func Schema(in *tfprotov6.Schema) (*tfprotov5.Schema, error) {
	if in == nil {
		return nil, nil
	}

	block, err := SchemaBlock(in.Block)

	if err != nil {
		return nil, err
	}

	return &tfprotov5.Schema{
		Block:   block,
		Version: in.Version,
	}, nil
}

func SchemaAttribute(in *tfprotov6.SchemaAttribute) (*tfprotov5.SchemaAttribute, error) {
	if in == nil {
		return nil, nil
	}

	if in.NestedType != nil {
		return nil, fmt.Errorf("unable to convert attribute %q schema: %w", in.Name, ErrSchemaAttributeNestedTypeNotImplemented)
	}

	return &tfprotov5.SchemaAttribute{
		Computed:        in.Computed,
		Deprecated:      in.Deprecated,
		Description:     in.Description,
		DescriptionKind: StringKind(in.DescriptionKind),
		Name:            in.Name,
		Optional:        in.Optional,
		Required:        in.Required,
		Sensitive:       in.Sensitive,
		Type:            in.Type,
	}, nil
}

func SchemaBlock(in *tfprotov6.SchemaBlock) (*tfprotov5.SchemaBlock, error) {
	if in == nil {
		return nil, nil
	}

	var attrs []*tfprotov5.SchemaAttribute

	if in.Attributes != nil {
		attrs = make([]*tfprotov5.SchemaAttribute, 0, len(in.Attributes))

		for _, attr := range in.Attributes {
			v5Attr, err := SchemaAttribute(attr)

			if err != nil {
				return nil, err
			}

			attrs = append(attrs, v5Attr)
		}
	}

	var nestedBlocks []*tfprotov5.SchemaNestedBlock

	if in.BlockTypes != nil {
		nestedBlocks = make([]*tfprotov5.SchemaNestedBlock, 0, len(in.BlockTypes))

		for _, block := range in.BlockTypes {
			v5Block, err := SchemaNestedBlock(block)

			if err != nil {
				return nil, err
			}

			nestedBlocks = append(nestedBlocks, v5Block)
		}
	}

	return &tfprotov5.SchemaBlock{
		Attributes:      attrs,
		BlockTypes:      nestedBlocks,
		Deprecated:      in.Deprecated,
		Description:     in.Description,
		DescriptionKind: StringKind(in.DescriptionKind),
		Version:         in.Version,
	}, nil
}

func SchemaNestedBlock(in *tfprotov6.SchemaNestedBlock) (*tfprotov5.SchemaNestedBlock, error) {
	if in == nil {
		return nil, nil
	}

	block, err := SchemaBlock(in.Block)

	if err != nil {
		return nil, fmt.Errorf("unable to convert block %q schema: %w", in.TypeName, err)
	}

	return &tfprotov5.SchemaNestedBlock{
		Block:    block,
		MaxItems: in.MaxItems,
		MinItems: in.MinItems,
		Nesting:  tfprotov5.SchemaNestedBlockNestingMode(in.Nesting),
		TypeName: in.TypeName,
	}, nil
}

func ServerCapabilities(in *tfprotov6.ServerCapabilities) *tfprotov5.ServerCapabilities {
	if in == nil {
		return nil
	}

	return &tfprotov5.ServerCapabilities{
		GetProviderSchemaOptional: in.GetProviderSchemaOptional,
		PlanDestroy:               in.PlanDestroy,
	}
}

func StopProviderRequest(in *tfprotov6.StopProviderRequest) *tfprotov5.StopProviderRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.StopProviderRequest{}
}

func StopProviderResponse(in *tfprotov6.StopProviderResponse) *tfprotov5.StopProviderResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.StopProviderResponse{
		Error: in.Error,
	}
}

func StringKind(in tfprotov6.StringKind) tfprotov5.StringKind {
	return tfprotov5.StringKind(in)
}

func UpgradeResourceStateRequest(in *tfprotov6.UpgradeResourceStateRequest) *tfprotov5.UpgradeResourceStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.UpgradeResourceStateRequest{
		RawState: RawState(in.RawState),
		TypeName: in.TypeName,
		Version:  in.Version,
	}
}

func UpgradeResourceStateResponse(in *tfprotov6.UpgradeResourceStateResponse) *tfprotov5.UpgradeResourceStateResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.UpgradeResourceStateResponse{
		Diagnostics:   Diagnostics(in.Diagnostics),
		UpgradedState: DynamicValue(in.UpgradedState),
	}
}

func ValidateDataSourceConfigRequest(in *tfprotov6.ValidateDataResourceConfigRequest) *tfprotov5.ValidateDataSourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ValidateDataSourceConfigRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}
}

func ValidateDataSourceConfigResponse(in *tfprotov6.ValidateDataResourceConfigResponse) *tfprotov5.ValidateDataSourceConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.ValidateDataSourceConfigResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func ValidateResourceTypeConfigRequest(in *tfprotov6.ValidateResourceConfigRequest) *tfprotov5.ValidateResourceTypeConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov5.ValidateResourceTypeConfigRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}
}

func ValidateResourceTypeConfigResponse(in *tfprotov6.ValidateResourceConfigResponse) *tfprotov5.ValidateResourceTypeConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov5.ValidateResourceTypeConfigResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}
