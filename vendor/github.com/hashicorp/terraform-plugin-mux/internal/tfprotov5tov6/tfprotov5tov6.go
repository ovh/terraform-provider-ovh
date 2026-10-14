// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package tfprotov5tov6

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func ApplyResourceChangeRequest(in *tfprotov5.ApplyResourceChangeRequest) *tfprotov6.ApplyResourceChangeRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ApplyResourceChangeRequest{
		Config:          DynamicValue(in.Config),
		PlannedPrivate:  in.PlannedPrivate,
		PlannedState:    DynamicValue(in.PlannedState),
		PriorState:      DynamicValue(in.PriorState),
		ProviderMeta:    DynamicValue(in.ProviderMeta),
		TypeName:        in.TypeName,
		PlannedIdentity: ResourceIdentityData(in.PlannedIdentity),
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
		NewIdentity:                 ResourceIdentityData(in.NewIdentity),
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

func CloseEphemeralResourceRequest(in *tfprotov5.CloseEphemeralResourceRequest) *tfprotov6.CloseEphemeralResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.CloseEphemeralResourceRequest{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}

func CloseEphemeralResourceResponse(in *tfprotov5.CloseEphemeralResourceResponse) *tfprotov6.CloseEphemeralResourceResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.CloseEphemeralResourceResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func ConfigureProviderRequest(in *tfprotov5.ConfigureProviderRequest) *tfprotov6.ConfigureProviderRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ConfigureProviderRequest{
		ClientCapabilities: ConfigureProviderClientCapabilities(in.ClientCapabilities),
		Config:             DynamicValue(in.Config),
		TerraformVersion:   in.TerraformVersion,
	}
}

func ConfigureProviderClientCapabilities(in *tfprotov5.ConfigureProviderClientCapabilities) *tfprotov6.ConfigureProviderClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ConfigureProviderClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
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

func Deferred(in *tfprotov5.Deferred) *tfprotov6.Deferred {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.Deferred{
		Reason: tfprotov6.DeferredReason(in.Reason),
	}

	return resp
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

func ResourceIdentityData(in *tfprotov5.ResourceIdentityData) *tfprotov6.ResourceIdentityData {
	if in == nil {
		return nil
	}

	return &tfprotov6.ResourceIdentityData{
		IdentityData: DynamicValue(in.IdentityData),
	}
}

func EphemeralResourceMetadata(in tfprotov5.EphemeralResourceMetadata) tfprotov6.EphemeralResourceMetadata {
	return tfprotov6.EphemeralResourceMetadata{
		TypeName: in.TypeName,
	}
}

func ListResourceMetadata(in tfprotov5.ListResourceMetadata) tfprotov6.ListResourceMetadata {
	return tfprotov6.ListResourceMetadata{
		TypeName: in.TypeName,
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
		Actions:            make([]tfprotov6.ActionMetadata, 0, len(in.Actions)),
		DataSources:        make([]tfprotov6.DataSourceMetadata, 0, len(in.DataSources)),
		Diagnostics:        Diagnostics(in.Diagnostics),
		EphemeralResources: make([]tfprotov6.EphemeralResourceMetadata, 0, len(in.Resources)),
		ListResources:      make([]tfprotov6.ListResourceMetadata, 0, len(in.ListResources)),
		Functions:          make([]tfprotov6.FunctionMetadata, 0, len(in.Functions)),
		Resources:          make([]tfprotov6.ResourceMetadata, 0, len(in.Resources)),
		ServerCapabilities: ServerCapabilities(in.ServerCapabilities),
	}

	for _, datasource := range in.DataSources {
		resp.DataSources = append(resp.DataSources, DataSourceMetadata(datasource))
	}

	for _, ephemeralResource := range in.EphemeralResources {
		resp.EphemeralResources = append(resp.EphemeralResources, EphemeralResourceMetadata(ephemeralResource))
	}

	for _, listResource := range in.ListResources {
		resp.ListResources = append(resp.ListResources, ListResourceMetadata(listResource))
	}

	for _, function := range in.Functions {
		resp.Functions = append(resp.Functions, FunctionMetadata(function))
	}

	for _, resource := range in.Resources {
		resp.Resources = append(resp.Resources, ResourceMetadata(resource))
	}

	for _, action := range in.Actions {
		resp.Actions = append(resp.Actions, ActionMetadata(action))
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

	ephemeralResourceSchemas := make(map[string]*tfprotov6.Schema, len(in.EphemeralResourceSchemas))

	for k, v := range in.EphemeralResourceSchemas {
		ephemeralResourceSchemas[k] = Schema(v)
	}

	listResourceSchemas := make(map[string]*tfprotov6.Schema, len(in.ListResourceSchemas))

	for k, v := range in.ListResourceSchemas {
		listResourceSchemas[k] = Schema(v)
	}

	functions := make(map[string]*tfprotov6.Function, len(in.Functions))

	for name, function := range in.Functions {
		functions[name] = Function(function)
	}

	resourceSchemas := make(map[string]*tfprotov6.Schema, len(in.ResourceSchemas))

	for k, v := range in.ResourceSchemas {
		resourceSchemas[k] = Schema(v)
	}

	actionSchemas := make(map[string]*tfprotov6.ActionSchema, len(in.ActionSchemas))

	for k, v := range in.ActionSchemas {
		actionSchemas[k] = ActionSchema(v)
	}

	return &tfprotov6.GetProviderSchemaResponse{
		ActionSchemas:            actionSchemas,
		DataSourceSchemas:        dataSourceSchemas,
		Diagnostics:              Diagnostics(in.Diagnostics),
		EphemeralResourceSchemas: ephemeralResourceSchemas,
		ListResourceSchemas:      listResourceSchemas,
		Functions:                functions,
		Provider:                 Schema(in.Provider),
		ProviderMeta:             Schema(in.ProviderMeta),
		ResourceSchemas:          resourceSchemas,
		ServerCapabilities:       ServerCapabilities(in.ServerCapabilities),
	}
}

func GetResourceIdentitySchemasRequest(in *tfprotov5.GetResourceIdentitySchemasRequest) *tfprotov6.GetResourceIdentitySchemasRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.GetResourceIdentitySchemasRequest{}
}

func GetResourceIdentitySchemasResponse(in *tfprotov5.GetResourceIdentitySchemasResponse) *tfprotov6.GetResourceIdentitySchemasResponse {
	if in == nil {
		return nil
	}

	identitySchemas := make(map[string]*tfprotov6.ResourceIdentitySchema, len(in.IdentitySchemas))

	for k, v := range in.IdentitySchemas {
		identitySchemas[k] = ResourceIdentitySchema(v)
	}

	return &tfprotov6.GetResourceIdentitySchemasResponse{
		Diagnostics:     Diagnostics(in.Diagnostics),
		IdentitySchemas: identitySchemas,
	}
}

func ImportResourceStateRequest(in *tfprotov5.ImportResourceStateRequest) *tfprotov6.ImportResourceStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ImportResourceStateRequest{
		ClientCapabilities: ImportResourceStateClientCapabilities(in.ClientCapabilities),
		ID:                 in.ID,
		TypeName:           in.TypeName,
		Identity:           ResourceIdentityData(in.Identity),
	}
}

func ImportResourceStateClientCapabilities(in *tfprotov5.ImportResourceStateClientCapabilities) *tfprotov6.ImportResourceStateClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ImportResourceStateClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ImportResourceStateResponse(in *tfprotov5.ImportResourceStateResponse) *tfprotov6.ImportResourceStateResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ImportResourceStateResponse{
		Deferred:          Deferred(in.Deferred),
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
			Identity: ResourceIdentityData(imp.Identity),
		})
	}

	return res
}

func MoveResourceStateRequest(in *tfprotov5.MoveResourceStateRequest) *tfprotov6.MoveResourceStateRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.MoveResourceStateRequest{
		SourcePrivate:               in.SourcePrivate,
		SourceProviderAddress:       in.SourceProviderAddress,
		SourceSchemaVersion:         in.SourceSchemaVersion,
		SourceState:                 RawState(in.SourceState),
		SourceTypeName:              in.SourceTypeName,
		TargetTypeName:              in.TargetTypeName,
		SourceIdentity:              RawState(in.SourceIdentity),
		SourceIdentitySchemaVersion: in.SourceIdentitySchemaVersion,
	}
}

func MoveResourceStateResponse(in *tfprotov5.MoveResourceStateResponse) *tfprotov6.MoveResourceStateResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.MoveResourceStateResponse{
		Diagnostics:    Diagnostics(in.Diagnostics),
		TargetPrivate:  in.TargetPrivate,
		TargetState:    DynamicValue(in.TargetState),
		TargetIdentity: ResourceIdentityData(in.TargetIdentity),
	}
}

func OpenEphemeralResourceRequest(in *tfprotov5.OpenEphemeralResourceRequest) *tfprotov6.OpenEphemeralResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.OpenEphemeralResourceRequest{
		TypeName:           in.TypeName,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: OpenEphemeralResourceClientCapabilities(in.ClientCapabilities),
	}
}

func OpenEphemeralResourceClientCapabilities(in *tfprotov5.OpenEphemeralResourceClientCapabilities) *tfprotov6.OpenEphemeralResourceClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.OpenEphemeralResourceClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func OpenEphemeralResourceResponse(in *tfprotov5.OpenEphemeralResourceResponse) *tfprotov6.OpenEphemeralResourceResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.OpenEphemeralResourceResponse{
		Result:      DynamicValue(in.Result),
		Diagnostics: Diagnostics(in.Diagnostics),
		Private:     in.Private,
		RenewAt:     in.RenewAt,
		Deferred:    Deferred(in.Deferred),
	}
}

func PlanResourceChangeRequest(in *tfprotov5.PlanResourceChangeRequest) *tfprotov6.PlanResourceChangeRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.PlanResourceChangeRequest{
		ClientCapabilities: PlanResourceChangeClientCapabilities(in.ClientCapabilities),
		Config:             DynamicValue(in.Config),
		PriorPrivate:       in.PriorPrivate,
		PriorState:         DynamicValue(in.PriorState),
		ProposedNewState:   DynamicValue(in.ProposedNewState),
		ProviderMeta:       DynamicValue(in.ProviderMeta),
		TypeName:           in.TypeName,
		PriorIdentity:      ResourceIdentityData(in.PriorIdentity),
	}
}

func PlanResourceChangeClientCapabilities(in *tfprotov5.PlanResourceChangeClientCapabilities) *tfprotov6.PlanResourceChangeClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.PlanResourceChangeClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func PlanResourceChangeResponse(in *tfprotov5.PlanResourceChangeResponse) *tfprotov6.PlanResourceChangeResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.PlanResourceChangeResponse{
		Deferred:                    Deferred(in.Deferred),
		Diagnostics:                 Diagnostics(in.Diagnostics),
		PlannedPrivate:              in.PlannedPrivate,
		PlannedState:                DynamicValue(in.PlannedState),
		RequiresReplace:             in.RequiresReplace,
		UnsafeToUseLegacyTypeSystem: in.UnsafeToUseLegacyTypeSystem, //nolint:staticcheck
		PlannedIdentity:             ResourceIdentityData(in.PlannedIdentity),
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
		ClientCapabilities: ReadDataSourceClientCapabilities(in.ClientCapabilities),
		Config:             DynamicValue(in.Config),
		ProviderMeta:       DynamicValue(in.ProviderMeta),
		TypeName:           in.TypeName,
	}
}

func ReadDataSourceClientCapabilities(in *tfprotov5.ReadDataSourceClientCapabilities) *tfprotov6.ReadDataSourceClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ReadDataSourceClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadDataSourceResponse(in *tfprotov5.ReadDataSourceResponse) *tfprotov6.ReadDataSourceResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ReadDataSourceResponse{
		Deferred:    Deferred(in.Deferred),
		Diagnostics: Diagnostics(in.Diagnostics),
		State:       DynamicValue(in.State),
	}
}

func ReadResourceRequest(in *tfprotov5.ReadResourceRequest) *tfprotov6.ReadResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ReadResourceRequest{
		ClientCapabilities: ReadResourceClientCapabilities(in.ClientCapabilities),
		CurrentState:       DynamicValue(in.CurrentState),
		Private:            in.Private,
		ProviderMeta:       DynamicValue(in.ProviderMeta),
		TypeName:           in.TypeName,
		CurrentIdentity:    ResourceIdentityData(in.CurrentIdentity),
	}
}

func ReadResourceClientCapabilities(in *tfprotov5.ReadResourceClientCapabilities) *tfprotov6.ReadResourceClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ReadResourceClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func ReadResourceResponse(in *tfprotov5.ReadResourceResponse) *tfprotov6.ReadResourceResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ReadResourceResponse{
		Deferred:    Deferred(in.Deferred),
		Diagnostics: Diagnostics(in.Diagnostics),
		NewState:    DynamicValue(in.NewState),
		Private:     in.Private,
		NewIdentity: ResourceIdentityData(in.NewIdentity),
	}
}

func RenewEphemeralResourceRequest(in *tfprotov5.RenewEphemeralResourceRequest) *tfprotov6.RenewEphemeralResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.RenewEphemeralResourceRequest{
		TypeName: in.TypeName,
		Private:  in.Private,
	}
}

func RenewEphemeralResourceResponse(in *tfprotov5.RenewEphemeralResourceResponse) *tfprotov6.RenewEphemeralResourceResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.RenewEphemeralResourceResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
		Private:     in.Private,
		RenewAt:     in.RenewAt,
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
		Computed:           in.Computed,
		Deprecated:         in.Deprecated,
		DeprecationMessage: in.DeprecationMessage,
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		Name:               in.Name,
		Optional:           in.Optional,
		Required:           in.Required,
		Sensitive:          in.Sensitive,
		Type:               in.Type,
		WriteOnly:          in.WriteOnly,
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
		Attributes:         attrs,
		BlockTypes:         nestedBlocks,
		Deprecated:         in.Deprecated,
		DeprecationMessage: in.DeprecationMessage,
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		Version:            in.Version,
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

func ResourceIdentitySchema(in *tfprotov5.ResourceIdentitySchema) *tfprotov6.ResourceIdentitySchema {
	if in == nil {
		return nil
	}

	var attrs []*tfprotov6.ResourceIdentitySchemaAttribute

	if in.IdentityAttributes != nil {
		attrs = make([]*tfprotov6.ResourceIdentitySchemaAttribute, 0, len(in.IdentityAttributes))

		for _, attr := range in.IdentityAttributes {
			attrs = append(attrs, ResourceIdentitySchemaAttribute(attr))
		}
	}

	return &tfprotov6.ResourceIdentitySchema{
		Version:            in.Version,
		IdentityAttributes: attrs,
	}
}

func ResourceIdentitySchemaAttribute(in *tfprotov5.ResourceIdentitySchemaAttribute) *tfprotov6.ResourceIdentitySchemaAttribute {
	if in == nil {
		return nil
	}

	return &tfprotov6.ResourceIdentitySchemaAttribute{
		Name:              in.Name,
		Type:              in.Type,
		RequiredForImport: in.RequiredForImport,
		OptionalForImport: in.OptionalForImport,
		Description:       in.Description,
	}
}

func ServerCapabilities(in *tfprotov5.ServerCapabilities) *tfprotov6.ServerCapabilities {
	if in == nil {
		return nil
	}

	return &tfprotov6.ServerCapabilities{
		GetProviderSchemaOptional: in.GetProviderSchemaOptional,
		MoveResourceState:         in.MoveResourceState,
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

func UpgradeResourceIdentityRequest(in *tfprotov5.UpgradeResourceIdentityRequest) *tfprotov6.UpgradeResourceIdentityRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.UpgradeResourceIdentityRequest{
		TypeName:    in.TypeName,
		Version:     in.Version,
		RawIdentity: RawState(in.RawIdentity),
	}
}

func UpgradeResourceIdentityResponse(in *tfprotov5.UpgradeResourceIdentityResponse) *tfprotov6.UpgradeResourceIdentityResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.UpgradeResourceIdentityResponse{
		Diagnostics:      Diagnostics(in.Diagnostics),
		UpgradedIdentity: ResourceIdentityData(in.UpgradedIdentity),
	}
}

func ValidateEphemeralResourceConfigRequest(in *tfprotov5.ValidateEphemeralResourceConfigRequest) *tfprotov6.ValidateEphemeralResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateEphemeralResourceConfigRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}
}

func ValidateEphemeralResourceConfigResponse(in *tfprotov5.ValidateEphemeralResourceConfigResponse) *tfprotov6.ValidateEphemeralResourceConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateEphemeralResourceConfigResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
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
		ClientCapabilities: ValidateResourceConfigClientCapabilities(in.ClientCapabilities),
		Config:             DynamicValue(in.Config),
		TypeName:           in.TypeName,
	}
}

func ValidateResourceConfigClientCapabilities(in *tfprotov5.ValidateResourceTypeConfigClientCapabilities) *tfprotov6.ValidateResourceConfigClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ValidateResourceConfigClientCapabilities{
		WriteOnlyAttributesAllowed: in.WriteOnlyAttributesAllowed,
	}

	return resp
}

func ValidateResourceConfigResponse(in *tfprotov5.ValidateResourceTypeConfigResponse) *tfprotov6.ValidateResourceConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateResourceConfigResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func ValidateListResourceConfigRequest(in *tfprotov5.ValidateListResourceConfigRequest) *tfprotov6.ValidateListResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateListResourceConfigRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}
}

func ValidateListResourceConfigResponse(in *tfprotov5.ValidateListResourceConfigResponse) *tfprotov6.ValidateListResourceConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateListResourceConfigResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func ListResourceRequest(in *tfprotov5.ListResourceRequest) *tfprotov6.ListResourceRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ListResourceRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}
}

func ListResourceServerStream(in *tfprotov5.ListResourceServerStream) *tfprotov6.ListResourceServerStream {
	if in == nil {
		return nil
	}

	return &tfprotov6.ListResourceServerStream{
		Results: func(yield func(tfprotov6.ListResourceResult) bool) {
			for res := range in.Results {
				if !yield(ListResourceResult(res)) {
					break
				}
			}
		},
	}
}

func ListResourceResult(in tfprotov5.ListResourceResult) tfprotov6.ListResourceResult {
	return tfprotov6.ListResourceResult{
		DisplayName: in.DisplayName,
		Resource:    DynamicValue(in.Resource),
		Identity:    ResourceIdentityData(in.Identity),
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func ActionMetadata(in tfprotov5.ActionMetadata) tfprotov6.ActionMetadata {
	return tfprotov6.ActionMetadata{
		TypeName: in.TypeName,
	}
}

func ActionSchema(in *tfprotov5.ActionSchema) *tfprotov6.ActionSchema {
	if in == nil {
		return nil
	}

	actionSchema := &tfprotov6.ActionSchema{
		Schema: Schema(in.Schema),
	}
	return actionSchema
}

func ValidateActionConfigRequest(in *tfprotov5.ValidateActionConfigRequest) *tfprotov6.ValidateActionConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateActionConfigRequest{
		Config:     DynamicValue(in.Config),
		ActionType: in.ActionType,
	}
}

func ValidateActionConfigResponse(in *tfprotov5.ValidateActionConfigResponse) *tfprotov6.ValidateActionConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.ValidateActionConfigResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}

func PlanActionRequest(in *tfprotov5.PlanActionRequest) *tfprotov6.PlanActionRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.PlanActionRequest{
		ActionType:         in.ActionType,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: PlanActionClientCapabilities(in.ClientCapabilities),
	}
}

func PlanActionClientCapabilities(in *tfprotov5.PlanActionClientCapabilities) *tfprotov6.PlanActionClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.PlanActionClientCapabilities{
		DeferralAllowed: in.DeferralAllowed,
	}

	return resp
}

func PlanActionResponse(in *tfprotov5.PlanActionResponse) *tfprotov6.PlanActionResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.PlanActionResponse{
		Diagnostics: Diagnostics(in.Diagnostics),
		Deferred:    Deferred(in.Deferred),
	}
}

func InvokeActionClientCapabilities(in *tfprotov5.InvokeActionClientCapabilities) *tfprotov6.InvokeActionClientCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.InvokeActionClientCapabilities{}

	return resp
}

func InvokeActionRequest(in *tfprotov5.InvokeActionRequest) *tfprotov6.InvokeActionRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.InvokeActionRequest{
		ActionType:         in.ActionType,
		ClientCapabilities: InvokeActionClientCapabilities(in.ClientCapabilities),
		Config:             DynamicValue(in.Config),
	}
}

func InvokeActionServerStream(in *tfprotov5.InvokeActionServerStream) *tfprotov6.InvokeActionServerStream {
	if in == nil {
		return nil
	}

	return &tfprotov6.InvokeActionServerStream{
		Events: func(yield func(tfprotov6.InvokeActionEvent) bool) {
			for res := range in.Events {
				if !yield(InvokeActionEvent(res)) {
					break
				}
			}
		},
	}
}

func InvokeActionEvent(in tfprotov5.InvokeActionEvent) tfprotov6.InvokeActionEvent {
	switch event := (in.Type).(type) {
	case tfprotov5.ProgressInvokeActionEventType:
		return tfprotov6.InvokeActionEvent{
			Type: tfprotov6.ProgressInvokeActionEventType{
				Message: event.Message,
			},
		}
	case tfprotov5.CompletedInvokeActionEventType:
		return tfprotov6.InvokeActionEvent{
			Type: tfprotov6.CompletedInvokeActionEventType{
				Diagnostics: Diagnostics(event.Diagnostics),
			},
		}
	}

	// It is not currently possible to create tfprotov5.InvokeActionEventType
	// implementations outside the terraform-plugin-go module. If this panic was reached,
	// it implies that a new event type was introduced and needs to be implemented
	// as a new case above.
	panic(fmt.Sprintf("unimplemented tfprotov5.InvokeActionEventType type: %T", in.Type))
}

func GenerateResourceConfigRequest(in *tfprotov5.GenerateResourceConfigRequest) *tfprotov6.GenerateResourceConfigRequest {
	if in == nil {
		return nil
	}

	return &tfprotov6.GenerateResourceConfigRequest{
		TypeName: in.TypeName,
		State:    DynamicValue(in.State),
	}
}

func GenerateResourceConfigResponse(in *tfprotov5.GenerateResourceConfigResponse) *tfprotov6.GenerateResourceConfigResponse {
	if in == nil {
		return nil
	}

	return &tfprotov6.GenerateResourceConfigResponse{
		Config:      DynamicValue(in.Config),
		Diagnostics: Diagnostics(in.Diagnostics),
	}
}
