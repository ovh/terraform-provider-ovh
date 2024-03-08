// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfprotov5tov6

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func ApplyResourceChangeRequest(in *tfprotov5.ApplyResourceChangeRequest) *tfprotov6.ApplyResourceChangeRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ApplyResourceChangeRequest{
		Config:         DynamicValue(in.Config),
		PlannedPrivate: in.PlannedPrivate,
		PlannedState:   DynamicValue(in.PlannedState),
		PriorState:     DynamicValue(in.PriorState),
		ProviderMeta:   DynamicValue(in.ProviderMeta),
		TypeName:       in.TypeName,
	}
}

func ApplyResourceChangeResponse(in *tfprotov5.ApplyResourceChangeResponse) *tfprotov6.ApplyResourceChangeResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ApplyResourceChangeResponse{
		Diagnostics:                 Diagnostics(in.Diagnostics),
		NewState:                    DynamicValue(in.NewState),
		Private:                     in.Private,
		UnsafeToUseLegacyTypeSystem: in.UnsafeToUseLegacyTypeSystem, //nolint:staticcheck
	}
}

func CallFunctionRequest(in *tfprotov5.CallFunctionRequest) *tfprotov6.CallFunctionRequest {
	if in == nil {
		return nil
	}

	out := &tfprotov6.CallFunctionRequest{
		Arguments: make([]*tfprotov6.DynamicValue, 0, len(in.Arguments)),
		Name:      in.Name,
	}

	for _, argument := range in.Arguments {
		out.Arguments = append(out.Arguments, DynamicValue(argument))
	}

	return out
}

func CallFunctionResponse(in *tfprotov5.CallFunctionResponse) *tfprotov6.CallFunctionResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.CallFunctionResponse{
		Error:  FunctionError(in.Error),
		Result: DynamicValue(in.Result),
	}
}

func ConfigureProviderRequest(in *tfprotov5.ConfigureProviderRequest) *tfprotov6.ConfigureProviderRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ConfigureProviderRequest{
		Config:           DynamicValue(in.Config),
		TerraformVersion: in.TerraformVersion,
	}
}

func ConfigureProviderResponse(in *tfprotov5.ConfigureProviderResponse) *tfprotov6.ConfigureProviderResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ConfigureProviderResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func DataSourceMetadata(in tfprotov5.DataSourceMetadata) tfprotov6.DataSourceMetadata {
	return tfprotov6.DataSourceMetadata{
		TypeName: in.TypeName,
	}
}

func Diagnostics(in []*tfprotov5.Diagnostic) []*tfprotov6.Diagnostic {
	if in == nil {
		return nil
	}

	diags := make([]*tfprotov6.Diagnostic, 0, len(in))

	for _, diag := range in {
		if diag == nil {
			diags = append(diags, nil)
			continue
		}

		diags = append(diags, &tfprotov6.Diagnostic{
			Severity:  tfprotov6.DiagnosticSeverity(diag.Severity),
			Summary:   diag.Summary,
			Detail:    diag.Detail,
			Attribute: diag.Attribute,
		})
	}

	return diags
}

func DynamicValue(in *tfprotov5.DynamicValue) *tfprotov6.DynamicValue {
	if in == nil {
		return nil
	}

	return &tfprotov6.DynamicValue{
		MsgPack: in.MsgPack,
		JSON:    in.JSON,
	}
}

func Function(in *tfprotov5.Function) *tfprotov6.Function {
	if in == nil {
		return nil
	}

	out := &tfprotov6.Function{
		DeprecationMessage: in.DeprecationMessage,
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		Parameters:         make([]*tfprotov6.FunctionParameter, 0, len(in.Parameters)),
		Return:             FunctionReturn(in.Return),
		Summary:            in.Summary,
		VariadicParameter:  FunctionParameter(in.VariadicParameter),
	}

	for _, parameter := range in.Parameters {
		out.Parameters = append(out.Parameters, FunctionParameter(parameter))
	}

	return out
}

func FunctionError(in *tfprotov5.FunctionError) *tfprotov6.FunctionError {
	if in == nil {
		return nil
	}

	out := &tfprotov6.FunctionError{
		Text:             in.Text,
		FunctionArgument: in.FunctionArgument,
	}

	return out
}

func FunctionMetadata(in tfprotov5.FunctionMetadata) tfprotov6.FunctionMetadata {
	return tfprotov6.FunctionMetadata{
		Name: in.Name,
	}
}

func FunctionParameter(in *tfprotov5.FunctionParameter) *tfprotov6.FunctionParameter {
	if in == nil {
		return nil
	}

	return &tfprotov6.FunctionParameter{
		AllowNullValue:     in.AllowNullValue,
		AllowUnknownValues: in.AllowUnknownValues,
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		Name:               in.Name,
		Type:               in.Type,
	}
}

func FunctionReturn(in *tfprotov5.FunctionReturn) *tfprotov6.FunctionReturn {
	if in == nil {
		return nil
	}

	return &tfprotov6.FunctionReturn{
		Type: in.Type,
	}
}

func GetFunctionsRequest(in *tfprotov5.GetFunctionsRequest) *tfprotov6.GetFunctionsRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.GetFunctionsRequest{}
}

func GetFunctionsResponse(in *tfprotov5.GetFunctionsResponse) *tfprotov6.GetFunctionsResponse {
	if in == nil {
		return nil
	}

	functions := make(map[string]*tfprotov6.Function, len(in.Functions))

	for name, function := range in.Functions {
		functions[name] = Function(function)
	}

	return &tfprotov6.GetFunctionsResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
		Functions:   functions,
	}
}

func GetMetadataRequest(in *tfprotov5.GetMetadataRequest) *tfprotov6.GetMetadataRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.GetMetadataRequest{}
}

func GetMetadataResponse(in *tfprotov5.GetMetadataResponse) *tfprotov6.GetMetadataResponse {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.GetMetadataResponse{
		DataSources:        make([]tfprotov6.DataSourceMetadata, 0, len(in.DataSources)),
		Diagnostics:        Diagnostics(in.Diagnostics),
		Functions:          make([]tfprotov6.FunctionMetadata, 0, len(in.Functions)),
		Resources:          make([]tfprotov6.ResourceMetadata, 0, len(in.Resources)),
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

func GetProviderSchemaRequest(in *tfprotov5.GetProviderSchemaRequest) *tfprotov6.GetProviderSchemaRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.GetProviderSchemaRequest{}
}

func GetProviderSchemaResponse(in *tfprotov5.GetProviderSchemaResponse) *tfprotov6.GetProviderSchemaResponse {
	if in == nil {
		return nil
	}

	dataSourceSchemas := make(map[string]*tfprotov6.Schema, len(in.DataSourceSchemas))

	for k, v := range in.DataSourceSchemas {
		dataSourceSchemas[k] = Schema(v)
	}

	functions := make(map[string]*tfprotov6.Function, len(in.Functions))

	for name, function := range in.Functions {
		functions[name] = Function(function)
	}

	resourceSchemas := make(map[string]*tfprotov6.Schema, len(in.ResourceSchemas))

	for k, v := range in.ResourceSchemas {
		resourceSchemas[k] = Schema(v)
	}

	return &tfprotov6.GetProviderSchemaResponse{
		DataSourceSchemas:  dataSourceSchemas,
		Diagnostics:        Diagnostics(in.Diagnostics),
		Functions:          functions,
		Provider:           Schema(in.Provider),
		ProviderMeta:       Schema(in.ProviderMeta),
		ResourceSchemas:    resourceSchemas,
		ServerCapabilities: ServerCapabilities(in.ServerCapabilities),
	}
}

func ImportResourceStateRequest(in *tfprotov5.ImportResourceStateRequest) *tfprotov6.ImportResourceStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ImportResourceStateRequest{
		ID:       in.ID,
		TypeName: in.TypeName,
	}
}

func ImportResourceStateResponse(in *tfprotov5.ImportResourceStateResponse) *tfprotov6.ImportResourceStateResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ImportResourceStateResponse{
		Diagnostics:       Diagnostics(in.Diagnostics),
		ImportedResources: ImportedResources(in.ImportedResources),
	}
}

func ImportedResources(in []*tfprotov5.ImportedResource) []*tfprotov6.ImportedResource {
	if in == nil {
		return nil
	}

	res := make([]*tfprotov6.ImportedResource, 0, len(in))

	for _, imp := range in {
		if imp == nil {
			res = append(res, nil)
			continue
		}

		res = append(res, &tfprotov6.ImportedResource{
			Private:  imp.Private,
			State:    DynamicValue(imp.State),
			TypeName: imp.TypeName,
		})
	}

	return res
}

func MoveResourceStateRequest(in *tfprotov5.MoveResourceStateRequest) *tfprotov6.MoveResourceStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.MoveResourceStateRequest{
		SourcePrivate:         in.SourcePrivate,
		SourceProviderAddress: in.SourceProviderAddress,
		SourceSchemaVersion:   in.SourceSchemaVersion,
		SourceState:           RawState(in.SourceState),
		SourceTypeName:        in.SourceTypeName,
		TargetTypeName:        in.TargetTypeName,
	}
}

func MoveResourceStateResponse(in *tfprotov5.MoveResourceStateResponse) *tfprotov6.MoveResourceStateResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.MoveResourceStateResponse{
		Diagnostics:   Diagnostics(in.Diagnostics),
		TargetPrivate: in.TargetPrivate,
		TargetState:   DynamicValue(in.TargetState),
	}
}

func PlanResourceChangeRequest(in *tfprotov5.PlanResourceChangeRequest) *tfprotov6.PlanResourceChangeRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.PlanResourceChangeRequest{
		Config:           DynamicValue(in.Config),
		PriorPrivate:     in.PriorPrivate,
		PriorState:       DynamicValue(in.PriorState),
		ProposedNewState: DynamicValue(in.ProposedNewState),
		ProviderMeta:     DynamicValue(in.ProviderMeta),
		TypeName:         in.TypeName,
	}
}

func PlanResourceChangeResponse(in *tfprotov5.PlanResourceChangeResponse) *tfprotov6.PlanResourceChangeResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.PlanResourceChangeResponse{
		Diagnostics:                 Diagnostics(in.Diagnostics),
		PlannedPrivate:              in.PlannedPrivate,
		PlannedState:                DynamicValue(in.PlannedState),
		RequiresReplace:             in.RequiresReplace,
		UnsafeToUseLegacyTypeSystem: in.UnsafeToUseLegacyTypeSystem, //nolint:staticcheck
	}
}

func RawState(in *tfprotov5.RawState) *tfprotov6.RawState {
	if in == nil {
		return nil
	}

	return &tfprotov6.RawState{
		Flatmap: in.Flatmap,
		JSON:    in.JSON,
	}
}

func ReadDataSourceRequest(in *tfprotov5.ReadDataSourceRequest) *tfprotov6.ReadDataSourceRequest {
	if in == nil {
		return nil
	}
	return &tfprotov6.ReadDataSourceRequest{
		Config:       DynamicValue(in.Config),
		ProviderMeta: DynamicValue(in.ProviderMeta),
		TypeName:     in.TypeName,
	}
}

func ReadDataSourceResponse(in *tfprotov5.ReadDataSourceResponse) *tfprotov6.ReadDataSourceResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ReadDataSourceResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
		State:       DynamicValue(in.State),
	}
}

func ReadResourceRequest(in *tfprotov5.ReadResourceRequest) *tfprotov6.ReadResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ReadResourceRequest{
		CurrentState: DynamicValue(in.CurrentState),
		Private:      in.Private,
		ProviderMeta: DynamicValue(in.ProviderMeta),
		TypeName:     in.TypeName,
	}
}

func ReadResourceResponse(in *tfprotov5.ReadResourceResponse) *tfprotov6.ReadResourceResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ReadResourceResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
		NewState:    DynamicValue(in.NewState),
		Private:     in.Private,
	}
}

func ResourceMetadata(in tfprotov5.ResourceMetadata) tfprotov6.ResourceMetadata {
	return tfprotov6.ResourceMetadata{
		TypeName: in.TypeName,
	}
}

func Schema(in *tfprotov5.Schema) *tfprotov6.Schema {
	if in == nil {
		return nil
	}

	return &tfprotov6.Schema{
		Block:   SchemaBlock(in.Block),
		Version: in.Version,
	}
}

func SchemaAttribute(in *tfprotov5.SchemaAttribute) *tfprotov6.SchemaAttribute {
	if in == nil {
		return nil
	}

	return &tfprotov6.SchemaAttribute{
		Computed:        in.Computed,
		Deprecated:      in.Deprecated,
		Description:     in.Description,
		DescriptionKind: StringKind(in.DescriptionKind),
		Name:            in.Name,
		Optional:        in.Optional,
		Required:        in.Required,
		Sensitive:       in.Sensitive,
		Type:            in.Type,
	}
}

func SchemaBlock(in *tfprotov5.SchemaBlock) *tfprotov6.SchemaBlock {
	if in == nil {
		return nil
	}

	var attrs []*tfprotov6.SchemaAttribute

	if in.Attributes != nil {
		attrs = make([]*tfprotov6.SchemaAttribute, 0, len(in.Attributes))

		for _, attr := range in.Attributes {
			attrs = append(attrs, SchemaAttribute(attr))
		}
	}

	var nestedBlocks []*tfprotov6.SchemaNestedBlock

	if in.BlockTypes != nil {
		nestedBlocks = make([]*tfprotov6.SchemaNestedBlock, 0, len(in.BlockTypes))

		for _, block := range in.BlockTypes {
			nestedBlocks = append(nestedBlocks, SchemaNestedBlock(block))
		}
	}

	return &tfprotov6.SchemaBlock{
		Attributes:      attrs,
		BlockTypes:      nestedBlocks,
		Deprecated:      in.Deprecated,
		Description:     in.Description,
		DescriptionKind: StringKind(in.DescriptionKind),
		Version:         in.Version,
	}
}

func SchemaNestedBlock(in *tfprotov5.SchemaNestedBlock) *tfprotov6.SchemaNestedBlock {
	if in == nil {
		return nil
	}

	return &tfprotov6.SchemaNestedBlock{
		Block:    SchemaBlock(in.Block),
		MaxItems: in.MaxItems,
		MinItems: in.MinItems,
		Nesting:  tfprotov6.SchemaNestedBlockNestingMode(in.Nesting),
		TypeName: in.TypeName,
	}
}

func ServerCapabilities(in *tfprotov5.ServerCapabilities) *tfprotov6.ServerCapabilities {
	if in == nil {
		return nil
	}

	return &tfprotov6.ServerCapabilities{
		GetProviderSchemaOptional: in.GetProviderSchemaOptional,
		PlanDestroy:               in.PlanDestroy,
	}
}

func StopProviderRequest(in *tfprotov5.StopProviderRequest) *tfprotov6.StopProviderRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.StopProviderRequest{}
}

func StopProviderResponse(in *tfprotov5.StopProviderResponse) *tfprotov6.StopProviderResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.StopProviderResponse{
		Error: in.Error,
	}
}

func StringKind(in tfprotov5.StringKind) tfprotov6.StringKind {
	return tfprotov6.StringKind(in)
}

func UpgradeResourceStateRequest(in *tfprotov5.UpgradeResourceStateRequest) *tfprotov6.UpgradeResourceStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.UpgradeResourceStateRequest{
		RawState: RawState(in.RawState),
		TypeName: in.TypeName,
		Version:  in.Version,
	}
}

func UpgradeResourceStateResponse(in *tfprotov5.UpgradeResourceStateResponse) *tfprotov6.UpgradeResourceStateResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.UpgradeResourceStateResponse{
		Diagnostics:   Diagnostics(in.Diagnostics),
		UpgradedState: DynamicValue(in.UpgradedState),
	}
}

func ValidateDataResourceConfigRequest(in *tfprotov5.ValidateDataSourceConfigRequest) *tfprotov6.ValidateDataResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateDataResourceConfigRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}
}

func ValidateDataResourceConfigResponse(in *tfprotov5.ValidateDataSourceConfigResponse) *tfprotov6.ValidateDataResourceConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateDataResourceConfigResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func ValidateProviderConfigRequest(in *tfprotov5.PrepareProviderConfigRequest) *tfprotov6.ValidateProviderConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateProviderConfigRequest{
		Config: DynamicValue(in.Config),
	}
}

func ValidateProviderConfigResponse(in *tfprotov5.PrepareProviderConfigResponse) *tfprotov6.ValidateProviderConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateProviderConfigResponse{
		Diagnostics:    Diagnostics(in.Diagnostics),
		PreparedConfig: DynamicValue(in.PreparedConfig),
	}
}

func ValidateResourceConfigRequest(in *tfprotov5.ValidateResourceTypeConfigRequest) *tfprotov6.ValidateResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateResourceConfigRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}
}

func ValidateResourceConfigResponse(in *tfprotov5.ValidateResourceTypeConfigResponse) *tfprotov6.ValidateResourceConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateResourceConfigResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}
