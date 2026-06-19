package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ basetypes.ObjectTypable = CloudNetworkPrivateSubnetCurrentStateType{}

type CloudNetworkPrivateSubnetCurrentStateType struct {
	basetypes.ObjectType
}

func (t CloudNetworkPrivateSubnetCurrentStateType) Equal(o attr.Type) bool {
	other, ok := o.(CloudNetworkPrivateSubnetCurrentStateType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CloudNetworkPrivateSubnetCurrentStateType) String() string {
	return "CloudNetworkPrivateSubnetCurrentStateType"
}

func (t CloudNetworkPrivateSubnetCurrentStateType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	allocationPoolsAttribute, ok := attributes["allocation_pools"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allocation_pools is missing from object`)

		return nil, diags
	}

	allocationPoolsVal, ok := allocationPoolsAttribute.(ovhtypes.TfListNestedValue[CurrentStateAllocationPoolsValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allocation_pools expected to be ovhtypes.TfListNestedValue[CurrentStateAllocationPoolsValue], was: %T`, allocationPoolsAttribute))
	}

	cidrAttribute, ok := attributes["cidr"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`cidr is missing from object`)

		return nil, diags
	}

	cidrVal, ok := cidrAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`cidr expected to be ovhtypes.TfStringValue, was: %T`, cidrAttribute))
	}

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return nil, diags
	}

	descriptionVal, ok := descriptionAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be ovhtypes.TfStringValue, was: %T`, descriptionAttribute))
	}

	dhcpEnabledAttribute, ok := attributes["dhcp_enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dhcp_enabled is missing from object`)

		return nil, diags
	}

	dhcpEnabledVal, ok := dhcpEnabledAttribute.(ovhtypes.TfBoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dhcp_enabled expected to be ovhtypes.TfBoolValue, was: %T`, dhcpEnabledAttribute))
	}

	dnsNameserversAttribute, ok := attributes["dns_nameservers"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dns_nameservers is missing from object`)

		return nil, diags
	}

	dnsNameserversVal, ok := dnsNameserversAttribute.(ovhtypes.TfListNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dns_nameservers expected to be ovhtypes.TfListNestedValue[ovhtypes.TfStringValue], was: %T`, dnsNameserversAttribute))
	}

	gatewayIpAttribute, ok := attributes["gateway_ip"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`gateway_ip is missing from object`)

		return nil, diags
	}

	gatewayIpVal, ok := gatewayIpAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`gateway_ip expected to be ovhtypes.TfStringValue, was: %T`, gatewayIpAttribute))
	}

	hostRoutesAttribute, ok := attributes["host_routes"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`host_routes is missing from object`)

		return nil, diags
	}

	hostRoutesVal, ok := hostRoutesAttribute.(ovhtypes.TfListNestedValue[CurrentStateHostRoutesValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`host_routes expected to be ovhtypes.TfListNestedValue[CurrentStateHostRoutesValue], was: %T`, hostRoutesAttribute))
	}

	locationAttribute, ok := attributes["location"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`location is missing from object`)

		return nil, diags
	}

	locationVal, ok := locationAttribute.(CurrentStateLocationValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`location expected to be CurrentStateLocationValue, was: %T`, locationAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return nil, diags
	}

	nameVal, ok := nameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, nameAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CloudNetworkPrivateSubnetCurrentStateValue{
		AllocationPools: allocationPoolsVal,
		Cidr:            cidrVal,
		Description:     descriptionVal,
		DhcpEnabled:     dhcpEnabledVal,
		DnsNameservers:  dnsNameserversVal,
		GatewayIp:       gatewayIpVal,
		HostRoutes:      hostRoutesVal,
		Location:        locationVal,
		Name:            nameVal,
		state:           attr.ValueStateKnown,
	}, diags
}

func NewCloudNetworkPrivateSubnetCurrentStateValueNull() CloudNetworkPrivateSubnetCurrentStateValue {
	return CloudNetworkPrivateSubnetCurrentStateValue{
		state: attr.ValueStateNull,
	}
}

func NewCloudNetworkPrivateSubnetCurrentStateValueUnknown() CloudNetworkPrivateSubnetCurrentStateValue {
	return CloudNetworkPrivateSubnetCurrentStateValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCloudNetworkPrivateSubnetCurrentStateValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CloudNetworkPrivateSubnetCurrentStateValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CloudNetworkPrivateSubnetCurrentStateValue Attribute Value",
				"While creating a CloudNetworkPrivateSubnetCurrentStateValue value, a missing attribute value was detected. "+
					"A CloudNetworkPrivateSubnetCurrentStateValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudNetworkPrivateSubnetCurrentStateValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CloudNetworkPrivateSubnetCurrentStateValue Attribute Type",
				"While creating a CloudNetworkPrivateSubnetCurrentStateValue value, an invalid attribute value was detected. "+
					"A CloudNetworkPrivateSubnetCurrentStateValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudNetworkPrivateSubnetCurrentStateValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CloudNetworkPrivateSubnetCurrentStateValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CloudNetworkPrivateSubnetCurrentStateValue Attribute Value",
				"While creating a CloudNetworkPrivateSubnetCurrentStateValue value, an extra attribute value was detected. "+
					"A CloudNetworkPrivateSubnetCurrentStateValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CloudNetworkPrivateSubnetCurrentStateValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	allocationPoolsAttribute, ok := attributes["allocation_pools"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allocation_pools is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	allocationPoolsVal, ok := allocationPoolsAttribute.(ovhtypes.TfListNestedValue[CurrentStateAllocationPoolsValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allocation_pools expected to be ovhtypes.TfListNestedValue[CurrentStateAllocationPoolsValue], was: %T`, allocationPoolsAttribute))
	}

	cidrAttribute, ok := attributes["cidr"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`cidr is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	cidrVal, ok := cidrAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`cidr expected to be ovhtypes.TfStringValue, was: %T`, cidrAttribute))
	}

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	descriptionVal, ok := descriptionAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be ovhtypes.TfStringValue, was: %T`, descriptionAttribute))
	}

	dhcpEnabledAttribute, ok := attributes["dhcp_enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dhcp_enabled is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	dhcpEnabledVal, ok := dhcpEnabledAttribute.(ovhtypes.TfBoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dhcp_enabled expected to be ovhtypes.TfBoolValue, was: %T`, dhcpEnabledAttribute))
	}

	dnsNameserversAttribute, ok := attributes["dns_nameservers"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dns_nameservers is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	dnsNameserversVal, ok := dnsNameserversAttribute.(ovhtypes.TfListNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dns_nameservers expected to be ovhtypes.TfListNestedValue[ovhtypes.TfStringValue], was: %T`, dnsNameserversAttribute))
	}

	gatewayIpAttribute, ok := attributes["gateway_ip"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`gateway_ip is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	gatewayIpVal, ok := gatewayIpAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`gateway_ip expected to be ovhtypes.TfStringValue, was: %T`, gatewayIpAttribute))
	}

	hostRoutesAttribute, ok := attributes["host_routes"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`host_routes is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	hostRoutesVal, ok := hostRoutesAttribute.(ovhtypes.TfListNestedValue[CurrentStateHostRoutesValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`host_routes expected to be ovhtypes.TfListNestedValue[CurrentStateHostRoutesValue], was: %T`, hostRoutesAttribute))
	}

	locationAttribute, ok := attributes["location"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`location is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	locationVal, ok := locationAttribute.(CurrentStateLocationValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`location expected to be CurrentStateLocationValue, was: %T`, locationAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, nameAttribute))
	}

	if diags.HasError() {
		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), diags
	}

	return CloudNetworkPrivateSubnetCurrentStateValue{
		AllocationPools: allocationPoolsVal,
		Cidr:            cidrVal,
		Description:     descriptionVal,
		DhcpEnabled:     dhcpEnabledVal,
		DnsNameservers:  dnsNameserversVal,
		GatewayIp:       gatewayIpVal,
		HostRoutes:      hostRoutesVal,
		Location:        locationVal,
		Name:            nameVal,
		state:           attr.ValueStateKnown,
	}, diags
}

func NewCloudNetworkPrivateSubnetCurrentStateValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CloudNetworkPrivateSubnetCurrentStateValue {
	object, diags := NewCloudNetworkPrivateSubnetCurrentStateValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCloudNetworkPrivateSubnetCurrentStateValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CloudNetworkPrivateSubnetCurrentStateType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCloudNetworkPrivateSubnetCurrentStateValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCloudNetworkPrivateSubnetCurrentStateValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCloudNetworkPrivateSubnetCurrentStateValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCloudNetworkPrivateSubnetCurrentStateValueMust(CloudNetworkPrivateSubnetCurrentStateValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CloudNetworkPrivateSubnetCurrentStateType) ValueType(ctx context.Context) attr.Value {
	return CloudNetworkPrivateSubnetCurrentStateValue{}
}

var _ basetypes.ObjectValuable = CloudNetworkPrivateSubnetCurrentStateValue{}

type CloudNetworkPrivateSubnetCurrentStateValue struct {
	AllocationPools ovhtypes.TfListNestedValue[CurrentStateAllocationPoolsValue] `tfsdk:"allocation_pools" json:"allocationPools"`
	Cidr            ovhtypes.TfStringValue                                       `tfsdk:"cidr" json:"cidr"`
	Description     ovhtypes.TfStringValue                                       `tfsdk:"description" json:"description"`
	DhcpEnabled     ovhtypes.TfBoolValue                                         `tfsdk:"dhcp_enabled" json:"dhcpEnabled"`
	DnsNameservers  ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]           `tfsdk:"dns_nameservers" json:"dnsNameservers"`
	GatewayIp       ovhtypes.TfStringValue                                       `tfsdk:"gateway_ip" json:"gatewayIp"`
	HostRoutes      ovhtypes.TfListNestedValue[CurrentStateHostRoutesValue]      `tfsdk:"host_routes" json:"hostRoutes"`
	Location        CurrentStateLocationValue                                    `tfsdk:"location" json:"location"`
	Name            ovhtypes.TfStringValue                                       `tfsdk:"name" json:"name"`
	state           attr.ValueState
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) ToCreate() *CloudNetworkPrivateSubnetCurrentStateValue {
	res := &CloudNetworkPrivateSubnetCurrentStateValue{}

	if !v.DhcpEnabled.IsNull() {
		res.DhcpEnabled = v.DhcpEnabled
	}

	res.state = attr.ValueStateKnown

	return res
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) MarshalJSON() ([]byte, error) {
	toMarshal := map[string]any{}
	if !v.DhcpEnabled.IsNull() && !v.DhcpEnabled.IsUnknown() {
		toMarshal["dhcpEnabled"] = v.DhcpEnabled
	}

	return json.Marshal(toMarshal)
}

func (v *CloudNetworkPrivateSubnetCurrentStateValue) UnmarshalJSON(data []byte) error {
	type JsonCurrentStateValue CloudNetworkPrivateSubnetCurrentStateValue

	var tmp JsonCurrentStateValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.AllocationPools = tmp.AllocationPools
	v.Cidr = tmp.Cidr
	v.Description = tmp.Description
	v.DhcpEnabled = tmp.DhcpEnabled
	v.DnsNameservers = tmp.DnsNameservers
	v.GatewayIp = tmp.GatewayIp
	v.HostRoutes = tmp.HostRoutes
	v.Location = tmp.Location
	v.Name = tmp.Name

	v.state = attr.ValueStateKnown

	return nil
}

func (v *CloudNetworkPrivateSubnetCurrentStateValue) MergeWith(other *CloudNetworkPrivateSubnetCurrentStateValue) {

	if (v.AllocationPools.IsUnknown() || v.AllocationPools.IsNull()) && !other.AllocationPools.IsUnknown() {
		v.AllocationPools = other.AllocationPools
	} else if !other.AllocationPools.IsUnknown() {
		newSlice := make([]attr.Value, 0)
		elems := v.AllocationPools.Elements()
		newElems := other.AllocationPools.Elements()

		if len(elems) != len(newElems) {
			v.AllocationPools = other.AllocationPools
		} else {
			for idx, e := range elems {
				tmp := e.(CurrentStateAllocationPoolsValue)
				tmp2 := newElems[idx].(CurrentStateAllocationPoolsValue)
				tmp.MergeWith(&tmp2)
				newSlice = append(newSlice, tmp)
			}

			v.AllocationPools = ovhtypes.TfListNestedValue[CurrentStateAllocationPoolsValue]{
				ListValue: basetypes.NewListValueMust(CurrentStateAllocationPoolsValue{}.Type(context.Background()), newSlice),
			}
		}
	}

	if (v.Cidr.IsUnknown() || v.Cidr.IsNull()) && !other.Cidr.IsUnknown() {
		v.Cidr = other.Cidr
	}

	if (v.Description.IsUnknown() || v.Description.IsNull()) && !other.Description.IsUnknown() {
		v.Description = other.Description
	}

	if (v.DhcpEnabled.IsUnknown() || v.DhcpEnabled.IsNull()) && !other.DhcpEnabled.IsUnknown() {
		v.DhcpEnabled = other.DhcpEnabled
	}

	if (v.DnsNameservers.IsUnknown() || v.DnsNameservers.IsNull()) && !other.DnsNameservers.IsUnknown() {
		v.DnsNameservers = other.DnsNameservers
	}

	if (v.GatewayIp.IsUnknown() || v.GatewayIp.IsNull()) && !other.GatewayIp.IsUnknown() {
		v.GatewayIp = other.GatewayIp
	}

	if (v.HostRoutes.IsUnknown() || v.HostRoutes.IsNull()) && !other.HostRoutes.IsUnknown() {
		v.HostRoutes = other.HostRoutes
	} else if !other.HostRoutes.IsUnknown() {
		newSlice := make([]attr.Value, 0)
		elems := v.HostRoutes.Elements()
		newElems := other.HostRoutes.Elements()

		if len(elems) != len(newElems) {
			v.HostRoutes = other.HostRoutes
		} else {
			for idx, e := range elems {
				tmp := e.(CurrentStateHostRoutesValue)
				tmp2 := newElems[idx].(CurrentStateHostRoutesValue)
				tmp.MergeWith(&tmp2)
				newSlice = append(newSlice, tmp)
			}

			v.HostRoutes = ovhtypes.TfListNestedValue[CurrentStateHostRoutesValue]{
				ListValue: basetypes.NewListValueMust(CurrentStateHostRoutesValue{}.Type(context.Background()), newSlice),
			}
		}
	}

	if v.Location.IsUnknown() && !other.Location.IsUnknown() {
		v.Location = other.Location
	} else if !other.Location.IsUnknown() {
		v.Location.MergeWith(&other.Location)
	}

	if (v.Name.IsUnknown() || v.Name.IsNull()) && !other.Name.IsUnknown() {
		v.Name = other.Name
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"allocationPools": v.AllocationPools,
		"cidr":            v.Cidr,
		"description":     v.Description,
		"dhcpEnabled":     v.DhcpEnabled,
		"dnsNameservers":  v.DnsNameservers,
		"gatewayIp":       v.GatewayIp,
		"hostRoutes":      v.HostRoutes,
		"location":        v.Location,
		"name":            v.Name,
	}
}
func (v CloudNetworkPrivateSubnetCurrentStateValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 9)

	var val tftypes.Value
	var err error

	attrTypes["allocation_pools"] = basetypes.ListType{
		ElemType: CurrentStateAllocationPoolsValue{}.Type(ctx),
	}.TerraformType(ctx)
	attrTypes["cidr"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["description"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["dhcp_enabled"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["dns_nameservers"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["gateway_ip"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["host_routes"] = basetypes.ListType{
		ElemType: CurrentStateHostRoutesValue{}.Type(ctx),
	}.TerraformType(ctx)
	attrTypes["location"] = basetypes.ObjectType{
		AttrTypes: CurrentStateLocationValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 9)

		val, err = v.AllocationPools.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["allocation_pools"] = val

		val, err = v.Cidr.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["cidr"] = val

		val, err = v.Description.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["description"] = val

		val, err = v.DhcpEnabled.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["dhcp_enabled"] = val

		val, err = v.DnsNameservers.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["dns_nameservers"] = val

		val, err = v.GatewayIp.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["gateway_ip"] = val

		val, err = v.HostRoutes.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["host_routes"] = val

		val, err = v.Location.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["location"] = val

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) String() string {
	return "CloudNetworkPrivateSubnetCurrentStateValue"
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"allocation_pools": ovhtypes.NewTfListNestedType[CurrentStateAllocationPoolsValue](ctx),
			"cidr":             ovhtypes.TfStringType{},
			"description":      ovhtypes.TfStringType{},
			"dhcp_enabled":     ovhtypes.TfBoolType{},
			"dns_nameservers":  ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			"gateway_ip":       ovhtypes.TfStringType{},
			"host_routes":      ovhtypes.NewTfListNestedType[CurrentStateHostRoutesValue](ctx),
			"location": CurrentStateLocationType{
				basetypes.ObjectType{
					AttrTypes: CurrentStateLocationValue{}.AttributeTypes(ctx),
				},
			},
			"name": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"allocation_pools": v.AllocationPools,
			"cidr":             v.Cidr,
			"description":      v.Description,
			"dhcp_enabled":     v.DhcpEnabled,
			"dns_nameservers":  v.DnsNameservers,
			"gateway_ip":       v.GatewayIp,
			"host_routes":      v.HostRoutes,
			"location":         v.Location,
			"name":             v.Name,
		})

	return objVal, diags
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) Equal(o attr.Value) bool {
	other, ok := o.(CloudNetworkPrivateSubnetCurrentStateValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.AllocationPools.Equal(other.AllocationPools) {
		return false
	}

	if !v.Cidr.Equal(other.Cidr) {
		return false
	}

	if !v.Description.Equal(other.Description) {
		return false
	}

	if !v.DhcpEnabled.Equal(other.DhcpEnabled) {
		return false
	}

	if !v.DnsNameservers.Equal(other.DnsNameservers) {
		return false
	}

	if !v.GatewayIp.Equal(other.GatewayIp) {
		return false
	}

	if !v.HostRoutes.Equal(other.HostRoutes) {
		return false
	}

	if !v.Location.Equal(other.Location) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	return true
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) Type(ctx context.Context) attr.Type {
	return CloudNetworkPrivateSubnetCurrentStateType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CloudNetworkPrivateSubnetCurrentStateValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"allocation_pools": ovhtypes.NewTfListNestedType[CurrentStateAllocationPoolsValue](ctx),
		"cidr":             ovhtypes.TfStringType{},
		"description":      ovhtypes.TfStringType{},
		"dhcp_enabled":     ovhtypes.TfBoolType{},
		"dns_nameservers":  ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
		"gateway_ip":       ovhtypes.TfStringType{},
		"host_routes":      ovhtypes.NewTfListNestedType[CurrentStateHostRoutesValue](ctx),
		"location":         CurrentStateLocationValue{}.Type(ctx),
		"name":             ovhtypes.TfStringType{},
	}
}

var _ basetypes.ObjectTypable = CurrentStateAllocationPoolsType{}

type CurrentStateAllocationPoolsType struct {
	basetypes.ObjectType
}

func (t CurrentStateAllocationPoolsType) Equal(o attr.Type) bool {
	other, ok := o.(CurrentStateAllocationPoolsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CurrentStateAllocationPoolsType) String() string {
	return "CurrentStateAllocationPoolsType"
}

func (t CurrentStateAllocationPoolsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	endAttribute, ok := attributes["end"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`end is missing from object`)

		return nil, diags
	}

	endVal, ok := endAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`end expected to be ovhtypes.TfStringValue, was: %T`, endAttribute))
	}

	startAttribute, ok := attributes["start"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`start is missing from object`)

		return nil, diags
	}

	startVal, ok := startAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`start expected to be ovhtypes.TfStringValue, was: %T`, startAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CurrentStateAllocationPoolsValue{
		End:   endVal,
		Start: startVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewCurrentStateAllocationPoolsValueNull() CurrentStateAllocationPoolsValue {
	return CurrentStateAllocationPoolsValue{
		state: attr.ValueStateNull,
	}
}

func NewCurrentStateAllocationPoolsValueUnknown() CurrentStateAllocationPoolsValue {
	return CurrentStateAllocationPoolsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCurrentStateAllocationPoolsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CurrentStateAllocationPoolsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CurrentStateAllocationPoolsValue Attribute Value",
				"While creating a CurrentStateAllocationPoolsValue value, a missing attribute value was detected. "+
					"A CurrentStateAllocationPoolsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CurrentStateAllocationPoolsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CurrentStateAllocationPoolsValue Attribute Type",
				"While creating a CurrentStateAllocationPoolsValue value, an invalid attribute value was detected. "+
					"A CurrentStateAllocationPoolsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CurrentStateAllocationPoolsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CurrentStateAllocationPoolsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CurrentStateAllocationPoolsValue Attribute Value",
				"While creating a CurrentStateAllocationPoolsValue value, an extra attribute value was detected. "+
					"A CurrentStateAllocationPoolsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CurrentStateAllocationPoolsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCurrentStateAllocationPoolsValueUnknown(), diags
	}

	endAttribute, ok := attributes["end"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`end is missing from object`)

		return NewCurrentStateAllocationPoolsValueUnknown(), diags
	}

	endVal, ok := endAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`end expected to be ovhtypes.TfStringValue, was: %T`, endAttribute))
	}

	startAttribute, ok := attributes["start"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`start is missing from object`)

		return NewCurrentStateAllocationPoolsValueUnknown(), diags
	}

	startVal, ok := startAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`start expected to be ovhtypes.TfStringValue, was: %T`, startAttribute))
	}

	if diags.HasError() {
		return NewCurrentStateAllocationPoolsValueUnknown(), diags
	}

	return CurrentStateAllocationPoolsValue{
		End:   endVal,
		Start: startVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewCurrentStateAllocationPoolsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CurrentStateAllocationPoolsValue {
	object, diags := NewCurrentStateAllocationPoolsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCurrentStateAllocationPoolsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CurrentStateAllocationPoolsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCurrentStateAllocationPoolsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCurrentStateAllocationPoolsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCurrentStateAllocationPoolsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCurrentStateAllocationPoolsValueMust(CurrentStateAllocationPoolsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CurrentStateAllocationPoolsType) ValueType(ctx context.Context) attr.Value {
	return CurrentStateAllocationPoolsValue{}
}

var _ basetypes.ObjectValuable = CurrentStateAllocationPoolsValue{}

type CurrentStateAllocationPoolsValue struct {
	End   ovhtypes.TfStringValue `tfsdk:"end" json:"end"`
	Start ovhtypes.TfStringValue `tfsdk:"start" json:"start"`
	state attr.ValueState
}

func (v *CurrentStateAllocationPoolsValue) UnmarshalJSON(data []byte) error {
	type JsonCurrentStateAllocationPoolsValue CurrentStateAllocationPoolsValue

	var tmp JsonCurrentStateAllocationPoolsValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.End = tmp.End
	v.Start = tmp.Start

	v.state = attr.ValueStateKnown

	return nil
}

func (v *CurrentStateAllocationPoolsValue) MergeWith(other *CurrentStateAllocationPoolsValue) {

	if (v.End.IsUnknown() || v.End.IsNull()) && !other.End.IsUnknown() {
		v.End = other.End
	}

	if (v.Start.IsUnknown() || v.Start.IsNull()) && !other.Start.IsUnknown() {
		v.Start = other.Start
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v CurrentStateAllocationPoolsValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"end":   v.End,
		"start": v.Start,
	}
}
func (v CurrentStateAllocationPoolsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["end"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["start"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.End.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["end"] = val

		val, err = v.Start.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["start"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CurrentStateAllocationPoolsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CurrentStateAllocationPoolsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CurrentStateAllocationPoolsValue) String() string {
	return "CurrentStateAllocationPoolsValue"
}

func (v CurrentStateAllocationPoolsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"end":   ovhtypes.TfStringType{},
			"start": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"end":   v.End,
			"start": v.Start,
		})

	return objVal, diags
}

func (v CurrentStateAllocationPoolsValue) Equal(o attr.Value) bool {
	other, ok := o.(CurrentStateAllocationPoolsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.End.Equal(other.End) {
		return false
	}

	if !v.Start.Equal(other.Start) {
		return false
	}

	return true
}

func (v CurrentStateAllocationPoolsValue) Type(ctx context.Context) attr.Type {
	return CurrentStateAllocationPoolsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CurrentStateAllocationPoolsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"end":   ovhtypes.TfStringType{},
		"start": ovhtypes.TfStringType{},
	}
}

var _ basetypes.ObjectTypable = CurrentStateHostRoutesType{}

type CurrentStateHostRoutesType struct {
	basetypes.ObjectType
}

func (t CurrentStateHostRoutesType) Equal(o attr.Type) bool {
	other, ok := o.(CurrentStateHostRoutesType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CurrentStateHostRoutesType) String() string {
	return "CurrentStateHostRoutesType"
}

func (t CurrentStateHostRoutesType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	destinationAttribute, ok := attributes["destination"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`destination is missing from object`)

		return nil, diags
	}

	destinationVal, ok := destinationAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`destination expected to be ovhtypes.TfStringValue, was: %T`, destinationAttribute))
	}

	nextHopAttribute, ok := attributes["next_hop"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`next_hop is missing from object`)

		return nil, diags
	}

	nextHopVal, ok := nextHopAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`next_hop expected to be ovhtypes.TfStringValue, was: %T`, nextHopAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CurrentStateHostRoutesValue{
		Destination: destinationVal,
		NextHop:     nextHopVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewCurrentStateHostRoutesValueNull() CurrentStateHostRoutesValue {
	return CurrentStateHostRoutesValue{
		state: attr.ValueStateNull,
	}
}

func NewCurrentStateHostRoutesValueUnknown() CurrentStateHostRoutesValue {
	return CurrentStateHostRoutesValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCurrentStateHostRoutesValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CurrentStateHostRoutesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CurrentStateHostRoutesValue Attribute Value",
				"While creating a CurrentStateHostRoutesValue value, a missing attribute value was detected. "+
					"A CurrentStateHostRoutesValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CurrentStateHostRoutesValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CurrentStateHostRoutesValue Attribute Type",
				"While creating a CurrentStateHostRoutesValue value, an invalid attribute value was detected. "+
					"A CurrentStateHostRoutesValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CurrentStateHostRoutesValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CurrentStateHostRoutesValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CurrentStateHostRoutesValue Attribute Value",
				"While creating a CurrentStateHostRoutesValue value, an extra attribute value was detected. "+
					"A CurrentStateHostRoutesValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CurrentStateHostRoutesValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCurrentStateHostRoutesValueUnknown(), diags
	}

	destinationAttribute, ok := attributes["destination"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`destination is missing from object`)

		return NewCurrentStateHostRoutesValueUnknown(), diags
	}

	destinationVal, ok := destinationAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`destination expected to be ovhtypes.TfStringValue, was: %T`, destinationAttribute))
	}

	nextHopAttribute, ok := attributes["next_hop"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`next_hop is missing from object`)

		return NewCurrentStateHostRoutesValueUnknown(), diags
	}

	nextHopVal, ok := nextHopAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`next_hop expected to be ovhtypes.TfStringValue, was: %T`, nextHopAttribute))
	}

	if diags.HasError() {
		return NewCurrentStateHostRoutesValueUnknown(), diags
	}

	return CurrentStateHostRoutesValue{
		Destination: destinationVal,
		NextHop:     nextHopVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewCurrentStateHostRoutesValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CurrentStateHostRoutesValue {
	object, diags := NewCurrentStateHostRoutesValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCurrentStateHostRoutesValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CurrentStateHostRoutesType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCurrentStateHostRoutesValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCurrentStateHostRoutesValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCurrentStateHostRoutesValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCurrentStateHostRoutesValueMust(CurrentStateHostRoutesValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CurrentStateHostRoutesType) ValueType(ctx context.Context) attr.Value {
	return CurrentStateHostRoutesValue{}
}

var _ basetypes.ObjectValuable = CurrentStateHostRoutesValue{}

type CurrentStateHostRoutesValue struct {
	Destination ovhtypes.TfStringValue `tfsdk:"destination" json:"destination"`
	NextHop     ovhtypes.TfStringValue `tfsdk:"next_hop" json:"nextHop"`
	state       attr.ValueState
}

func (v *CurrentStateHostRoutesValue) UnmarshalJSON(data []byte) error {
	type JsonCurrentStateHostRoutesValue CurrentStateHostRoutesValue

	var tmp JsonCurrentStateHostRoutesValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Destination = tmp.Destination
	v.NextHop = tmp.NextHop

	v.state = attr.ValueStateKnown

	return nil
}

func (v *CurrentStateHostRoutesValue) MergeWith(other *CurrentStateHostRoutesValue) {

	if (v.Destination.IsUnknown() || v.Destination.IsNull()) && !other.Destination.IsUnknown() {
		v.Destination = other.Destination
	}

	if (v.NextHop.IsUnknown() || v.NextHop.IsNull()) && !other.NextHop.IsUnknown() {
		v.NextHop = other.NextHop
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v CurrentStateHostRoutesValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"destination": v.Destination,
		"nextHop":     v.NextHop,
	}
}
func (v CurrentStateHostRoutesValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["destination"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["next_hop"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Destination.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["destination"] = val

		val, err = v.NextHop.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["next_hop"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CurrentStateHostRoutesValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CurrentStateHostRoutesValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CurrentStateHostRoutesValue) String() string {
	return "CurrentStateHostRoutesValue"
}

func (v CurrentStateHostRoutesValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"destination": ovhtypes.TfStringType{},
			"next_hop":    ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"destination": v.Destination,
			"next_hop":    v.NextHop,
		})

	return objVal, diags
}

func (v CurrentStateHostRoutesValue) Equal(o attr.Value) bool {
	other, ok := o.(CurrentStateHostRoutesValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Destination.Equal(other.Destination) {
		return false
	}

	if !v.NextHop.Equal(other.NextHop) {
		return false
	}

	return true
}

func (v CurrentStateHostRoutesValue) Type(ctx context.Context) attr.Type {
	return CurrentStateHostRoutesType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CurrentStateHostRoutesValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"destination": ovhtypes.TfStringType{},
		"next_hop":    ovhtypes.TfStringType{},
	}
}

var _ basetypes.ObjectTypable = CurrentStateLocationType{}

type CurrentStateLocationType struct {
	basetypes.ObjectType
}

func (t CurrentStateLocationType) Equal(o attr.Type) bool {
	other, ok := o.(CurrentStateLocationType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CurrentStateLocationType) String() string {
	return "CurrentStateLocationType"
}

func (t CurrentStateLocationType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	availabilityZoneAttribute, ok := attributes["availability_zone"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`availability_zone is missing from object`)

		return nil, diags
	}

	availabilityZoneVal, ok := availabilityZoneAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`availability_zone expected to be ovhtypes.TfStringValue, was: %T`, availabilityZoneAttribute))
	}

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return nil, diags
	}

	regionVal, ok := regionAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be ovhtypes.TfStringValue, was: %T`, regionAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CurrentStateLocationValue{
		AvailabilityZone: availabilityZoneVal,
		Region:           regionVal,
		state:            attr.ValueStateKnown,
	}, diags
}

func NewCurrentStateLocationValueNull() CurrentStateLocationValue {
	return CurrentStateLocationValue{
		state: attr.ValueStateNull,
	}
}

func NewCurrentStateLocationValueUnknown() CurrentStateLocationValue {
	return CurrentStateLocationValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCurrentStateLocationValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CurrentStateLocationValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CurrentStateLocationValue Attribute Value",
				"While creating a CurrentStateLocationValue value, a missing attribute value was detected. "+
					"A CurrentStateLocationValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CurrentStateLocationValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CurrentStateLocationValue Attribute Type",
				"While creating a CurrentStateLocationValue value, an invalid attribute value was detected. "+
					"A CurrentStateLocationValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CurrentStateLocationValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CurrentStateLocationValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CurrentStateLocationValue Attribute Value",
				"While creating a CurrentStateLocationValue value, an extra attribute value was detected. "+
					"A CurrentStateLocationValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CurrentStateLocationValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCurrentStateLocationValueUnknown(), diags
	}

	availabilityZoneAttribute, ok := attributes["availability_zone"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`availability_zone is missing from object`)

		return NewCurrentStateLocationValueUnknown(), diags
	}

	availabilityZoneVal, ok := availabilityZoneAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`availability_zone expected to be ovhtypes.TfStringValue, was: %T`, availabilityZoneAttribute))
	}

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return NewCurrentStateLocationValueUnknown(), diags
	}

	regionVal, ok := regionAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be ovhtypes.TfStringValue, was: %T`, regionAttribute))
	}

	if diags.HasError() {
		return NewCurrentStateLocationValueUnknown(), diags
	}

	return CurrentStateLocationValue{
		AvailabilityZone: availabilityZoneVal,
		Region:           regionVal,
		state:            attr.ValueStateKnown,
	}, diags
}

func NewCurrentStateLocationValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CurrentStateLocationValue {
	object, diags := NewCurrentStateLocationValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCurrentStateLocationValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CurrentStateLocationType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCurrentStateLocationValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCurrentStateLocationValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCurrentStateLocationValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCurrentStateLocationValueMust(CurrentStateLocationValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CurrentStateLocationType) ValueType(ctx context.Context) attr.Value {
	return CurrentStateLocationValue{}
}

var _ basetypes.ObjectValuable = CurrentStateLocationValue{}

type CurrentStateLocationValue struct {
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone" json:"availabilityZone"`
	Region           ovhtypes.TfStringValue `tfsdk:"region" json:"region"`
	state            attr.ValueState
}

func (v *CurrentStateLocationValue) UnmarshalJSON(data []byte) error {
	type JsonCurrentStateLocationValue CurrentStateLocationValue

	var tmp JsonCurrentStateLocationValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.AvailabilityZone = tmp.AvailabilityZone
	v.Region = tmp.Region

	v.state = attr.ValueStateKnown

	return nil
}

func (v *CurrentStateLocationValue) MergeWith(other *CurrentStateLocationValue) {

	if (v.AvailabilityZone.IsUnknown() || v.AvailabilityZone.IsNull()) && !other.AvailabilityZone.IsUnknown() {
		v.AvailabilityZone = other.AvailabilityZone
	}

	if (v.Region.IsUnknown() || v.Region.IsNull()) && !other.Region.IsUnknown() {
		v.Region = other.Region
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v CurrentStateLocationValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"availabilityZone": v.AvailabilityZone,
		"region":           v.Region,
	}
}
func (v CurrentStateLocationValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["availability_zone"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["region"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.AvailabilityZone.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["availability_zone"] = val

		val, err = v.Region.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["region"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CurrentStateLocationValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CurrentStateLocationValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CurrentStateLocationValue) String() string {
	return "CurrentStateLocationValue"
}

func (v CurrentStateLocationValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"availability_zone": ovhtypes.TfStringType{},
			"region":            ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"availability_zone": v.AvailabilityZone,
			"region":            v.Region,
		})

	return objVal, diags
}

func (v CurrentStateLocationValue) Equal(o attr.Value) bool {
	other, ok := o.(CurrentStateLocationValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.AvailabilityZone.Equal(other.AvailabilityZone) {
		return false
	}

	if !v.Region.Equal(other.Region) {
		return false
	}

	return true
}

func (v CurrentStateLocationValue) Type(ctx context.Context) attr.Type {
	return CurrentStateLocationType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CurrentStateLocationValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"availability_zone": ovhtypes.TfStringType{},
		"region":            ovhtypes.TfStringType{},
	}
}

var _ basetypes.ObjectTypable = CloudNetworkPrivateSubnetCurrentTasksType{}

type CloudNetworkPrivateSubnetCurrentTasksType struct {
	basetypes.ObjectType
}

func (t CloudNetworkPrivateSubnetCurrentTasksType) Equal(o attr.Type) bool {
	other, ok := o.(CloudNetworkPrivateSubnetCurrentTasksType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CloudNetworkPrivateSubnetCurrentTasksType) String() string {
	return "CloudNetworkPrivateSubnetCurrentTasksType"
}

func (t CloudNetworkPrivateSubnetCurrentTasksType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	errorsAttribute, ok := attributes["errors"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`errors is missing from object`)

		return nil, diags
	}

	errorsVal, ok := errorsAttribute.(ovhtypes.TfListNestedValue[CurrentTasksErrorsValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`errors expected to be ovhtypes.TfListNestedValue[CurrentTasksErrorsValue], was: %T`, errorsAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return nil, diags
	}

	idVal, ok := idAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be ovhtypes.TfStringValue, was: %T`, idAttribute))
	}

	linkAttribute, ok := attributes["link"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`link is missing from object`)

		return nil, diags
	}

	linkVal, ok := linkAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`link expected to be ovhtypes.TfStringValue, was: %T`, linkAttribute))
	}

	statusAttribute, ok := attributes["status"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`status is missing from object`)

		return nil, diags
	}

	statusVal, ok := statusAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`status expected to be ovhtypes.TfStringValue, was: %T`, statusAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return nil, diags
	}

	typeVal, ok := typeAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be ovhtypes.TfStringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CloudNetworkPrivateSubnetCurrentTasksValue{
		Errors: errorsVal,
		Id:     idVal,
		Link:   linkVal,
		Status: statusVal,
		CloudNetworkPrivateSubnetCurrentTasksType: typeVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewCloudNetworkPrivateSubnetCurrentTasksValueNull() CloudNetworkPrivateSubnetCurrentTasksValue {
	return CloudNetworkPrivateSubnetCurrentTasksValue{
		state: attr.ValueStateNull,
	}
}

func NewCloudNetworkPrivateSubnetCurrentTasksValueUnknown() CloudNetworkPrivateSubnetCurrentTasksValue {
	return CloudNetworkPrivateSubnetCurrentTasksValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCloudNetworkPrivateSubnetCurrentTasksValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CloudNetworkPrivateSubnetCurrentTasksValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CloudNetworkPrivateSubnetCurrentTasksValue Attribute Value",
				"While creating a CloudNetworkPrivateSubnetCurrentTasksValue value, a missing attribute value was detected. "+
					"A CloudNetworkPrivateSubnetCurrentTasksValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudNetworkPrivateSubnetCurrentTasksValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CloudNetworkPrivateSubnetCurrentTasksValue Attribute Type",
				"While creating a CloudNetworkPrivateSubnetCurrentTasksValue value, an invalid attribute value was detected. "+
					"A CloudNetworkPrivateSubnetCurrentTasksValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudNetworkPrivateSubnetCurrentTasksValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CloudNetworkPrivateSubnetCurrentTasksValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CloudNetworkPrivateSubnetCurrentTasksValue Attribute Value",
				"While creating a CloudNetworkPrivateSubnetCurrentTasksValue value, an extra attribute value was detected. "+
					"A CloudNetworkPrivateSubnetCurrentTasksValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CloudNetworkPrivateSubnetCurrentTasksValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCloudNetworkPrivateSubnetCurrentTasksValueUnknown(), diags
	}

	errorsAttribute, ok := attributes["errors"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`errors is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentTasksValueUnknown(), diags
	}

	errorsVal, ok := errorsAttribute.(ovhtypes.TfListNestedValue[CurrentTasksErrorsValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`errors expected to be ovhtypes.TfListNestedValue[CurrentTasksErrorsValue], was: %T`, errorsAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentTasksValueUnknown(), diags
	}

	idVal, ok := idAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be ovhtypes.TfStringValue, was: %T`, idAttribute))
	}

	linkAttribute, ok := attributes["link"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`link is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentTasksValueUnknown(), diags
	}

	linkVal, ok := linkAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`link expected to be ovhtypes.TfStringValue, was: %T`, linkAttribute))
	}

	statusAttribute, ok := attributes["status"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`status is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentTasksValueUnknown(), diags
	}

	statusVal, ok := statusAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`status expected to be ovhtypes.TfStringValue, was: %T`, statusAttribute))
	}

	typeAttribute, ok := attributes["type"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`type is missing from object`)

		return NewCloudNetworkPrivateSubnetCurrentTasksValueUnknown(), diags
	}

	typeVal, ok := typeAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`type expected to be ovhtypes.TfStringValue, was: %T`, typeAttribute))
	}

	if diags.HasError() {
		return NewCloudNetworkPrivateSubnetCurrentTasksValueUnknown(), diags
	}

	return CloudNetworkPrivateSubnetCurrentTasksValue{
		Errors: errorsVal,
		Id:     idVal,
		Link:   linkVal,
		Status: statusVal,
		CloudNetworkPrivateSubnetCurrentTasksType: typeVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewCloudNetworkPrivateSubnetCurrentTasksValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CloudNetworkPrivateSubnetCurrentTasksValue {
	object, diags := NewCloudNetworkPrivateSubnetCurrentTasksValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCloudNetworkPrivateSubnetCurrentTasksValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CloudNetworkPrivateSubnetCurrentTasksType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCloudNetworkPrivateSubnetCurrentTasksValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCloudNetworkPrivateSubnetCurrentTasksValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCloudNetworkPrivateSubnetCurrentTasksValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCloudNetworkPrivateSubnetCurrentTasksValueMust(CloudNetworkPrivateSubnetCurrentTasksValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CloudNetworkPrivateSubnetCurrentTasksType) ValueType(ctx context.Context) attr.Value {
	return CloudNetworkPrivateSubnetCurrentTasksValue{}
}

var _ basetypes.ObjectValuable = CloudNetworkPrivateSubnetCurrentTasksValue{}

type CloudNetworkPrivateSubnetCurrentTasksValue struct {
	Errors                                    ovhtypes.TfListNestedValue[CurrentTasksErrorsValue] `tfsdk:"errors" json:"errors"`
	Id                                        ovhtypes.TfStringValue                              `tfsdk:"id" json:"id"`
	Link                                      ovhtypes.TfStringValue                              `tfsdk:"link" json:"link"`
	Status                                    ovhtypes.TfStringValue                              `tfsdk:"status" json:"status"`
	CloudNetworkPrivateSubnetCurrentTasksType ovhtypes.TfStringValue                              `tfsdk:"type" json:"type"`
	state                                     attr.ValueState
}

func (v *CloudNetworkPrivateSubnetCurrentTasksValue) UnmarshalJSON(data []byte) error {
	type JsonCurrentTasksValue CloudNetworkPrivateSubnetCurrentTasksValue

	var tmp JsonCurrentTasksValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Errors = tmp.Errors
	v.Id = tmp.Id
	v.Link = tmp.Link
	v.Status = tmp.Status
	v.CloudNetworkPrivateSubnetCurrentTasksType = tmp.CloudNetworkPrivateSubnetCurrentTasksType

	v.state = attr.ValueStateKnown

	return nil
}

func (v *CloudNetworkPrivateSubnetCurrentTasksValue) MergeWith(other *CloudNetworkPrivateSubnetCurrentTasksValue) {

	if (v.Errors.IsUnknown() || v.Errors.IsNull()) && !other.Errors.IsUnknown() {
		v.Errors = other.Errors
	} else if !other.Errors.IsUnknown() {
		newSlice := make([]attr.Value, 0)
		elems := v.Errors.Elements()
		newElems := other.Errors.Elements()

		if len(elems) != len(newElems) {
			v.Errors = other.Errors
		} else {
			for idx, e := range elems {
				tmp := e.(CurrentTasksErrorsValue)
				tmp2 := newElems[idx].(CurrentTasksErrorsValue)
				tmp.MergeWith(&tmp2)
				newSlice = append(newSlice, tmp)
			}

			v.Errors = ovhtypes.TfListNestedValue[CurrentTasksErrorsValue]{
				ListValue: basetypes.NewListValueMust(CurrentTasksErrorsValue{}.Type(context.Background()), newSlice),
			}
		}
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.Link.IsUnknown() || v.Link.IsNull()) && !other.Link.IsUnknown() {
		v.Link = other.Link
	}

	if (v.Status.IsUnknown() || v.Status.IsNull()) && !other.Status.IsUnknown() {
		v.Status = other.Status
	}

	if (v.CloudNetworkPrivateSubnetCurrentTasksType.IsUnknown() || v.CloudNetworkPrivateSubnetCurrentTasksType.IsNull()) && !other.CloudNetworkPrivateSubnetCurrentTasksType.IsUnknown() {
		v.CloudNetworkPrivateSubnetCurrentTasksType = other.CloudNetworkPrivateSubnetCurrentTasksType
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v CloudNetworkPrivateSubnetCurrentTasksValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"errors": v.Errors,
		"id":     v.Id,
		"link":   v.Link,
		"status": v.Status,
		"type":   v.CloudNetworkPrivateSubnetCurrentTasksType,
	}
}
func (v CloudNetworkPrivateSubnetCurrentTasksValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 5)

	var val tftypes.Value
	var err error

	attrTypes["errors"] = basetypes.ListType{
		ElemType: CurrentTasksErrorsValue{}.Type(ctx),
	}.TerraformType(ctx)
	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["link"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["status"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["type"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 5)

		val, err = v.Errors.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["errors"] = val

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		val, err = v.Link.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["link"] = val

		val, err = v.Status.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["status"] = val

		val, err = v.CloudNetworkPrivateSubnetCurrentTasksType.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["type"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CloudNetworkPrivateSubnetCurrentTasksValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CloudNetworkPrivateSubnetCurrentTasksValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CloudNetworkPrivateSubnetCurrentTasksValue) String() string {
	return "CloudNetworkPrivateSubnetCurrentTasksValue"
}

func (v CloudNetworkPrivateSubnetCurrentTasksValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"errors": ovhtypes.NewTfListNestedType[CurrentTasksErrorsValue](ctx),
			"id":     ovhtypes.TfStringType{},
			"link":   ovhtypes.TfStringType{},
			"status": ovhtypes.TfStringType{},
			"type":   ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"errors": v.Errors,
			"id":     v.Id,
			"link":   v.Link,
			"status": v.Status,
			"type":   v.CloudNetworkPrivateSubnetCurrentTasksType,
		})

	return objVal, diags
}

func (v CloudNetworkPrivateSubnetCurrentTasksValue) Equal(o attr.Value) bool {
	other, ok := o.(CloudNetworkPrivateSubnetCurrentTasksValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Errors.Equal(other.Errors) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Link.Equal(other.Link) {
		return false
	}

	if !v.Status.Equal(other.Status) {
		return false
	}

	if !v.CloudNetworkPrivateSubnetCurrentTasksType.Equal(other.CloudNetworkPrivateSubnetCurrentTasksType) {
		return false
	}

	return true
}

func (v CloudNetworkPrivateSubnetCurrentTasksValue) Type(ctx context.Context) attr.Type {
	return CloudNetworkPrivateSubnetCurrentTasksType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CloudNetworkPrivateSubnetCurrentTasksValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"errors": ovhtypes.NewTfListNestedType[CurrentTasksErrorsValue](ctx),
		"id":     ovhtypes.TfStringType{},
		"link":   ovhtypes.TfStringType{},
		"status": ovhtypes.TfStringType{},
		"type":   ovhtypes.TfStringType{},
	}
}

var _ basetypes.ObjectTypable = CurrentTasksErrorsType{}

type CurrentTasksErrorsType struct {
	basetypes.ObjectType
}

func (t CurrentTasksErrorsType) Equal(o attr.Type) bool {
	other, ok := o.(CurrentTasksErrorsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CurrentTasksErrorsType) String() string {
	return "CurrentTasksErrorsType"
}

func (t CurrentTasksErrorsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	messageAttribute, ok := attributes["message"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`message is missing from object`)

		return nil, diags
	}

	messageVal, ok := messageAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`message expected to be ovhtypes.TfStringValue, was: %T`, messageAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CurrentTasksErrorsValue{
		Message: messageVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewCurrentTasksErrorsValueNull() CurrentTasksErrorsValue {
	return CurrentTasksErrorsValue{
		state: attr.ValueStateNull,
	}
}

func NewCurrentTasksErrorsValueUnknown() CurrentTasksErrorsValue {
	return CurrentTasksErrorsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCurrentTasksErrorsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CurrentTasksErrorsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CurrentTasksErrorsValue Attribute Value",
				"While creating a CurrentTasksErrorsValue value, a missing attribute value was detected. "+
					"A CurrentTasksErrorsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CurrentTasksErrorsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CurrentTasksErrorsValue Attribute Type",
				"While creating a CurrentTasksErrorsValue value, an invalid attribute value was detected. "+
					"A CurrentTasksErrorsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CurrentTasksErrorsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CurrentTasksErrorsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CurrentTasksErrorsValue Attribute Value",
				"While creating a CurrentTasksErrorsValue value, an extra attribute value was detected. "+
					"A CurrentTasksErrorsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CurrentTasksErrorsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCurrentTasksErrorsValueUnknown(), diags
	}

	messageAttribute, ok := attributes["message"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`message is missing from object`)

		return NewCurrentTasksErrorsValueUnknown(), diags
	}

	messageVal, ok := messageAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`message expected to be ovhtypes.TfStringValue, was: %T`, messageAttribute))
	}

	if diags.HasError() {
		return NewCurrentTasksErrorsValueUnknown(), diags
	}

	return CurrentTasksErrorsValue{
		Message: messageVal,
		state:   attr.ValueStateKnown,
	}, diags
}

func NewCurrentTasksErrorsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CurrentTasksErrorsValue {
	object, diags := NewCurrentTasksErrorsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCurrentTasksErrorsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CurrentTasksErrorsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCurrentTasksErrorsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCurrentTasksErrorsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCurrentTasksErrorsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCurrentTasksErrorsValueMust(CurrentTasksErrorsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CurrentTasksErrorsType) ValueType(ctx context.Context) attr.Value {
	return CurrentTasksErrorsValue{}
}

var _ basetypes.ObjectValuable = CurrentTasksErrorsValue{}

type CurrentTasksErrorsValue struct {
	Message ovhtypes.TfStringValue `tfsdk:"message" json:"message"`
	state   attr.ValueState
}

func (v *CurrentTasksErrorsValue) UnmarshalJSON(data []byte) error {
	type JsonCurrentTasksErrorsValue CurrentTasksErrorsValue

	var tmp JsonCurrentTasksErrorsValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Message = tmp.Message

	v.state = attr.ValueStateKnown

	return nil
}

func (v *CurrentTasksErrorsValue) MergeWith(other *CurrentTasksErrorsValue) {

	if (v.Message.IsUnknown() || v.Message.IsNull()) && !other.Message.IsUnknown() {
		v.Message = other.Message
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v CurrentTasksErrorsValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"message": v.Message,
	}
}
func (v CurrentTasksErrorsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 1)

	var val tftypes.Value
	var err error

	attrTypes["message"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 1)

		val, err = v.Message.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["message"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CurrentTasksErrorsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CurrentTasksErrorsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CurrentTasksErrorsValue) String() string {
	return "CurrentTasksErrorsValue"
}

func (v CurrentTasksErrorsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"message": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"message": v.Message,
		})

	return objVal, diags
}

func (v CurrentTasksErrorsValue) Equal(o attr.Value) bool {
	other, ok := o.(CurrentTasksErrorsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Message.Equal(other.Message) {
		return false
	}

	return true
}

func (v CurrentTasksErrorsValue) Type(ctx context.Context) attr.Type {
	return CurrentTasksErrorsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CurrentTasksErrorsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"message": ovhtypes.TfStringType{},
	}
}

var _ basetypes.ObjectTypable = CloudNetworkPrivateSubnetTargetSpecType{}

type CloudNetworkPrivateSubnetTargetSpecType struct {
	basetypes.ObjectType
}

func (t CloudNetworkPrivateSubnetTargetSpecType) Equal(o attr.Type) bool {
	other, ok := o.(CloudNetworkPrivateSubnetTargetSpecType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CloudNetworkPrivateSubnetTargetSpecType) String() string {
	return "CloudNetworkPrivateSubnetTargetSpecType"
}

func (t CloudNetworkPrivateSubnetTargetSpecType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	allocationPoolsAttribute, ok := attributes["allocation_pools"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allocation_pools is missing from object`)

		return nil, diags
	}

	allocationPoolsVal, ok := allocationPoolsAttribute.(ovhtypes.TfListNestedValue[TargetSpecAllocationPoolsValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allocation_pools expected to be ovhtypes.TfListNestedValue[TargetSpecAllocationPoolsValue], was: %T`, allocationPoolsAttribute))
	}

	cidrAttribute, ok := attributes["cidr"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`cidr is missing from object`)

		return nil, diags
	}

	cidrVal, ok := cidrAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`cidr expected to be ovhtypes.TfStringValue, was: %T`, cidrAttribute))
	}

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return nil, diags
	}

	descriptionVal, ok := descriptionAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be ovhtypes.TfStringValue, was: %T`, descriptionAttribute))
	}

	dhcpEnabledAttribute, ok := attributes["dhcp_enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dhcp_enabled is missing from object`)

		return nil, diags
	}

	dhcpEnabledVal, ok := dhcpEnabledAttribute.(ovhtypes.TfBoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dhcp_enabled expected to be ovhtypes.TfBoolValue, was: %T`, dhcpEnabledAttribute))
	}

	dnsNameserversAttribute, ok := attributes["dns_nameservers"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dns_nameservers is missing from object`)

		return nil, diags
	}

	dnsNameserversVal, ok := dnsNameserversAttribute.(ovhtypes.TfListNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dns_nameservers expected to be ovhtypes.TfListNestedValue[ovhtypes.TfStringValue], was: %T`, dnsNameserversAttribute))
	}

	gatewayIpAttribute, ok := attributes["gateway_ip"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`gateway_ip is missing from object`)

		return nil, diags
	}

	gatewayIpVal, ok := gatewayIpAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`gateway_ip expected to be ovhtypes.TfStringValue, was: %T`, gatewayIpAttribute))
	}

	locationAttribute, ok := attributes["location"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`location is missing from object`)

		return nil, diags
	}

	locationVal, ok := locationAttribute.(TargetSpecLocationValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`location expected to be TargetSpecLocationValue, was: %T`, locationAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return nil, diags
	}

	nameVal, ok := nameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, nameAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CloudNetworkPrivateSubnetTargetSpecValue{
		AllocationPools: allocationPoolsVal,
		Cidr:            cidrVal,
		Description:     descriptionVal,
		DhcpEnabled:     dhcpEnabledVal,
		DnsNameservers:  dnsNameserversVal,
		GatewayIp:       gatewayIpVal,
		Location:        locationVal,
		Name:            nameVal,
		state:           attr.ValueStateKnown,
	}, diags
}

func NewCloudNetworkPrivateSubnetTargetSpecValueNull() CloudNetworkPrivateSubnetTargetSpecValue {
	return CloudNetworkPrivateSubnetTargetSpecValue{
		state: attr.ValueStateNull,
	}
}

func NewCloudNetworkPrivateSubnetTargetSpecValueUnknown() CloudNetworkPrivateSubnetTargetSpecValue {
	return CloudNetworkPrivateSubnetTargetSpecValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCloudNetworkPrivateSubnetTargetSpecValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CloudNetworkPrivateSubnetTargetSpecValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CloudNetworkPrivateSubnetTargetSpecValue Attribute Value",
				"While creating a CloudNetworkPrivateSubnetTargetSpecValue value, a missing attribute value was detected. "+
					"A CloudNetworkPrivateSubnetTargetSpecValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudNetworkPrivateSubnetTargetSpecValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CloudNetworkPrivateSubnetTargetSpecValue Attribute Type",
				"While creating a CloudNetworkPrivateSubnetTargetSpecValue value, an invalid attribute value was detected. "+
					"A CloudNetworkPrivateSubnetTargetSpecValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudNetworkPrivateSubnetTargetSpecValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CloudNetworkPrivateSubnetTargetSpecValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CloudNetworkPrivateSubnetTargetSpecValue Attribute Value",
				"While creating a CloudNetworkPrivateSubnetTargetSpecValue value, an extra attribute value was detected. "+
					"A CloudNetworkPrivateSubnetTargetSpecValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CloudNetworkPrivateSubnetTargetSpecValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	allocationPoolsAttribute, ok := attributes["allocation_pools"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`allocation_pools is missing from object`)

		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	allocationPoolsVal, ok := allocationPoolsAttribute.(ovhtypes.TfListNestedValue[TargetSpecAllocationPoolsValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`allocation_pools expected to be ovhtypes.TfListNestedValue[TargetSpecAllocationPoolsValue], was: %T`, allocationPoolsAttribute))
	}

	cidrAttribute, ok := attributes["cidr"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`cidr is missing from object`)

		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	cidrVal, ok := cidrAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`cidr expected to be ovhtypes.TfStringValue, was: %T`, cidrAttribute))
	}

	descriptionAttribute, ok := attributes["description"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`description is missing from object`)

		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	descriptionVal, ok := descriptionAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`description expected to be ovhtypes.TfStringValue, was: %T`, descriptionAttribute))
	}

	dhcpEnabledAttribute, ok := attributes["dhcp_enabled"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dhcp_enabled is missing from object`)

		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	dhcpEnabledVal, ok := dhcpEnabledAttribute.(ovhtypes.TfBoolValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dhcp_enabled expected to be ovhtypes.TfBoolValue, was: %T`, dhcpEnabledAttribute))
	}

	dnsNameserversAttribute, ok := attributes["dns_nameservers"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dns_nameservers is missing from object`)

		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	dnsNameserversVal, ok := dnsNameserversAttribute.(ovhtypes.TfListNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dns_nameservers expected to be ovhtypes.TfListNestedValue[ovhtypes.TfStringValue], was: %T`, dnsNameserversAttribute))
	}

	gatewayIpAttribute, ok := attributes["gateway_ip"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`gateway_ip is missing from object`)

		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	gatewayIpVal, ok := gatewayIpAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`gateway_ip expected to be ovhtypes.TfStringValue, was: %T`, gatewayIpAttribute))
	}

	locationAttribute, ok := attributes["location"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`location is missing from object`)

		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	locationVal, ok := locationAttribute.(TargetSpecLocationValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`location expected to be TargetSpecLocationValue, was: %T`, locationAttribute))
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, nameAttribute))
	}

	if diags.HasError() {
		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), diags
	}

	return CloudNetworkPrivateSubnetTargetSpecValue{
		AllocationPools: allocationPoolsVal,
		Cidr:            cidrVal,
		Description:     descriptionVal,
		DhcpEnabled:     dhcpEnabledVal,
		DnsNameservers:  dnsNameserversVal,
		GatewayIp:       gatewayIpVal,
		Location:        locationVal,
		Name:            nameVal,
		state:           attr.ValueStateKnown,
	}, diags
}

func NewCloudNetworkPrivateSubnetTargetSpecValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CloudNetworkPrivateSubnetTargetSpecValue {
	object, diags := NewCloudNetworkPrivateSubnetTargetSpecValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewCloudNetworkPrivateSubnetTargetSpecValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CloudNetworkPrivateSubnetTargetSpecType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCloudNetworkPrivateSubnetTargetSpecValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCloudNetworkPrivateSubnetTargetSpecValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCloudNetworkPrivateSubnetTargetSpecValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCloudNetworkPrivateSubnetTargetSpecValueMust(CloudNetworkPrivateSubnetTargetSpecValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CloudNetworkPrivateSubnetTargetSpecType) ValueType(ctx context.Context) attr.Value {
	return CloudNetworkPrivateSubnetTargetSpecValue{}
}

var _ basetypes.ObjectValuable = CloudNetworkPrivateSubnetTargetSpecValue{}

type CloudNetworkPrivateSubnetTargetSpecValue struct {
	AllocationPools ovhtypes.TfListNestedValue[TargetSpecAllocationPoolsValue] `tfsdk:"allocation_pools" json:"allocationPools"`
	Cidr            ovhtypes.TfStringValue                                     `tfsdk:"cidr" json:"cidr"`
	Description     ovhtypes.TfStringValue                                     `tfsdk:"description" json:"description"`
	DhcpEnabled     ovhtypes.TfBoolValue                                       `tfsdk:"dhcp_enabled" json:"dhcpEnabled"`
	DnsNameservers  ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]         `tfsdk:"dns_nameservers" json:"dnsNameservers"`
	GatewayIp       ovhtypes.TfStringValue                                     `tfsdk:"gateway_ip" json:"gatewayIp"`
	Location        TargetSpecLocationValue                                    `tfsdk:"location" json:"location"`
	Name            ovhtypes.TfStringValue                                     `tfsdk:"name" json:"name"`
	state           attr.ValueState
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) ToCreate() *CloudNetworkPrivateSubnetTargetSpecValue {
	res := &CloudNetworkPrivateSubnetTargetSpecValue{}

	if !v.Description.IsNull() {
		res.Description = v.Description
	}

	if !v.DhcpEnabled.IsNull() {
		res.DhcpEnabled = v.DhcpEnabled
	}

	if !v.DnsNameservers.IsNull() {
		res.DnsNameservers = v.DnsNameservers
	}

	if !v.GatewayIp.IsNull() {
		res.GatewayIp = v.GatewayIp
	}

	if !v.Location.IsNull() {
		res.Location = v.Location
	}

	if !v.Name.IsNull() {
		res.Name = v.Name
	}

	if !v.AllocationPools.IsNull() {
		res.AllocationPools = v.AllocationPools
	}

	if !v.Cidr.IsNull() {
		res.Cidr = v.Cidr
	}

	res.state = attr.ValueStateKnown

	return res
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) ToUpdate() *CloudNetworkPrivateSubnetTargetSpecValue {
	res := &CloudNetworkPrivateSubnetTargetSpecValue{}

	if !v.Description.IsNull() {
		res.Description = v.Description
	}

	if !v.DhcpEnabled.IsNull() {
		res.DhcpEnabled = v.DhcpEnabled
	}

	if !v.DnsNameservers.IsNull() {
		res.DnsNameservers = v.DnsNameservers
	}

	if !v.GatewayIp.IsNull() {
		res.GatewayIp = v.GatewayIp
	}

	if !v.Name.IsNull() {
		res.Name = v.Name
	}

	if !v.AllocationPools.IsNull() {
		res.AllocationPools = v.AllocationPools
	}

	res.state = attr.ValueStateKnown

	return res
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) MarshalJSON() ([]byte, error) {
	toMarshal := map[string]any{}
	if !v.AllocationPools.IsNull() && !v.AllocationPools.IsUnknown() {
		toMarshal["allocationPools"] = v.AllocationPools
	}
	if !v.Cidr.IsNull() && !v.Cidr.IsUnknown() {
		toMarshal["cidr"] = v.Cidr
	}
	if !v.Description.IsNull() && !v.Description.IsUnknown() {
		toMarshal["description"] = v.Description
	}
	if !v.DhcpEnabled.IsNull() && !v.DhcpEnabled.IsUnknown() {
		toMarshal["dhcpEnabled"] = v.DhcpEnabled
	}
	if !v.DnsNameservers.IsNull() && !v.DnsNameservers.IsUnknown() {
		toMarshal["dnsNameservers"] = v.DnsNameservers
	}
	if !v.GatewayIp.IsNull() && !v.GatewayIp.IsUnknown() {
		toMarshal["gatewayIp"] = v.GatewayIp
	}
	if !v.Location.IsNull() && !v.Location.IsUnknown() {
		toMarshal["location"] = v.Location
	}
	if !v.Name.IsNull() && !v.Name.IsUnknown() {
		toMarshal["name"] = v.Name
	}

	return json.Marshal(toMarshal)
}

func (v *CloudNetworkPrivateSubnetTargetSpecValue) UnmarshalJSON(data []byte) error {
	type JsonTargetSpecValue CloudNetworkPrivateSubnetTargetSpecValue

	var tmp JsonTargetSpecValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.AllocationPools = tmp.AllocationPools
	v.Cidr = tmp.Cidr
	v.Description = tmp.Description
	v.DhcpEnabled = tmp.DhcpEnabled
	v.DnsNameservers = tmp.DnsNameservers
	v.GatewayIp = tmp.GatewayIp
	v.Location = tmp.Location
	v.Name = tmp.Name

	v.state = attr.ValueStateKnown

	return nil
}

func (v *CloudNetworkPrivateSubnetTargetSpecValue) MergeWith(other *CloudNetworkPrivateSubnetTargetSpecValue) {

	if (v.AllocationPools.IsUnknown() || v.AllocationPools.IsNull()) && !other.AllocationPools.IsUnknown() {
		v.AllocationPools = other.AllocationPools
	} else if !other.AllocationPools.IsUnknown() {
		newSlice := make([]attr.Value, 0)
		elems := v.AllocationPools.Elements()
		newElems := other.AllocationPools.Elements()

		if len(elems) != len(newElems) {
			v.AllocationPools = other.AllocationPools
		} else {
			for idx, e := range elems {
				tmp := e.(TargetSpecAllocationPoolsValue)
				tmp2 := newElems[idx].(TargetSpecAllocationPoolsValue)
				tmp.MergeWith(&tmp2)
				newSlice = append(newSlice, tmp)
			}

			v.AllocationPools = ovhtypes.TfListNestedValue[TargetSpecAllocationPoolsValue]{
				ListValue: basetypes.NewListValueMust(TargetSpecAllocationPoolsValue{}.Type(context.Background()), newSlice),
			}
		}
	}

	if (v.Cidr.IsUnknown() || v.Cidr.IsNull()) && !other.Cidr.IsUnknown() {
		v.Cidr = other.Cidr
	}

	if (v.Description.IsUnknown() || v.Description.IsNull()) && !other.Description.IsUnknown() {
		v.Description = other.Description
	}

	if (v.DhcpEnabled.IsUnknown() || v.DhcpEnabled.IsNull()) && !other.DhcpEnabled.IsUnknown() {
		v.DhcpEnabled = other.DhcpEnabled
	}

	if (v.DnsNameservers.IsUnknown() || v.DnsNameservers.IsNull()) && !other.DnsNameservers.IsUnknown() {
		v.DnsNameservers = other.DnsNameservers
	}

	if (v.GatewayIp.IsUnknown() || v.GatewayIp.IsNull()) && !other.GatewayIp.IsUnknown() {
		v.GatewayIp = other.GatewayIp
	}

	if v.Location.IsUnknown() && !other.Location.IsUnknown() {
		v.Location = other.Location
	} else if !other.Location.IsUnknown() {
		v.Location.MergeWith(&other.Location)
	}

	if (v.Name.IsUnknown() || v.Name.IsNull()) && !other.Name.IsUnknown() {
		v.Name = other.Name
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"allocationPools": v.AllocationPools,
		"cidr":            v.Cidr,
		"description":     v.Description,
		"dhcpEnabled":     v.DhcpEnabled,
		"dnsNameservers":  v.DnsNameservers,
		"gatewayIp":       v.GatewayIp,
		"location":        v.Location,
		"name":            v.Name,
	}
}
func (v CloudNetworkPrivateSubnetTargetSpecValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 8)

	var val tftypes.Value
	var err error

	attrTypes["allocation_pools"] = basetypes.ListType{
		ElemType: TargetSpecAllocationPoolsValue{}.Type(ctx),
	}.TerraformType(ctx)
	attrTypes["cidr"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["description"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["dhcp_enabled"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["dns_nameservers"] = basetypes.ListType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["gateway_ip"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["location"] = basetypes.ObjectType{
		AttrTypes: TargetSpecLocationValue{}.AttributeTypes(ctx),
	}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 8)

		val, err = v.AllocationPools.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["allocation_pools"] = val

		val, err = v.Cidr.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["cidr"] = val

		val, err = v.Description.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["description"] = val

		val, err = v.DhcpEnabled.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["dhcp_enabled"] = val

		val, err = v.DnsNameservers.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["dns_nameservers"] = val

		val, err = v.GatewayIp.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["gateway_ip"] = val

		val, err = v.Location.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["location"] = val

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) String() string {
	return "CloudNetworkPrivateSubnetTargetSpecValue"
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"allocation_pools": ovhtypes.NewTfListNestedType[TargetSpecAllocationPoolsValue](ctx),
			"cidr":             ovhtypes.TfStringType{},
			"description":      ovhtypes.TfStringType{},
			"dhcp_enabled":     ovhtypes.TfBoolType{},
			"dns_nameservers":  ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			"gateway_ip":       ovhtypes.TfStringType{},
			"location": TargetSpecLocationType{
				basetypes.ObjectType{
					AttrTypes: TargetSpecLocationValue{}.AttributeTypes(ctx),
				},
			},
			"name": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"allocation_pools": v.AllocationPools,
			"cidr":             v.Cidr,
			"description":      v.Description,
			"dhcp_enabled":     v.DhcpEnabled,
			"dns_nameservers":  v.DnsNameservers,
			"gateway_ip":       v.GatewayIp,
			"location":         v.Location,
			"name":             v.Name,
		})

	return objVal, diags
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) Equal(o attr.Value) bool {
	other, ok := o.(CloudNetworkPrivateSubnetTargetSpecValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.AllocationPools.Equal(other.AllocationPools) {
		return false
	}

	if !v.Cidr.Equal(other.Cidr) {
		return false
	}

	if !v.Description.Equal(other.Description) {
		return false
	}

	if !v.DhcpEnabled.Equal(other.DhcpEnabled) {
		return false
	}

	if !v.DnsNameservers.Equal(other.DnsNameservers) {
		return false
	}

	if !v.GatewayIp.Equal(other.GatewayIp) {
		return false
	}

	if !v.Location.Equal(other.Location) {
		return false
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	return true
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) Type(ctx context.Context) attr.Type {
	return CloudNetworkPrivateSubnetTargetSpecType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CloudNetworkPrivateSubnetTargetSpecValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"allocation_pools": ovhtypes.NewTfListNestedType[TargetSpecAllocationPoolsValue](ctx),
		"cidr":             ovhtypes.TfStringType{},
		"description":      ovhtypes.TfStringType{},
		"dhcp_enabled":     ovhtypes.TfBoolType{},
		"dns_nameservers":  ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
		"gateway_ip":       ovhtypes.TfStringType{},
		"location":         TargetSpecLocationValue{}.Type(ctx),
		"name":             ovhtypes.TfStringType{},
	}
}

var _ basetypes.ObjectTypable = TargetSpecAllocationPoolsType{}

type TargetSpecAllocationPoolsType struct {
	basetypes.ObjectType
}

func (t TargetSpecAllocationPoolsType) Equal(o attr.Type) bool {
	other, ok := o.(TargetSpecAllocationPoolsType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t TargetSpecAllocationPoolsType) String() string {
	return "TargetSpecAllocationPoolsType"
}

func (t TargetSpecAllocationPoolsType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	endAttribute, ok := attributes["end"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`end is missing from object`)

		return nil, diags
	}

	endVal, ok := endAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`end expected to be ovhtypes.TfStringValue, was: %T`, endAttribute))
	}

	startAttribute, ok := attributes["start"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`start is missing from object`)

		return nil, diags
	}

	startVal, ok := startAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`start expected to be ovhtypes.TfStringValue, was: %T`, startAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return TargetSpecAllocationPoolsValue{
		End:   endVal,
		Start: startVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewTargetSpecAllocationPoolsValueNull() TargetSpecAllocationPoolsValue {
	return TargetSpecAllocationPoolsValue{
		state: attr.ValueStateNull,
	}
}

func NewTargetSpecAllocationPoolsValueUnknown() TargetSpecAllocationPoolsValue {
	return TargetSpecAllocationPoolsValue{
		state: attr.ValueStateUnknown,
	}
}

func NewTargetSpecAllocationPoolsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (TargetSpecAllocationPoolsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing TargetSpecAllocationPoolsValue Attribute Value",
				"While creating a TargetSpecAllocationPoolsValue value, a missing attribute value was detected. "+
					"A TargetSpecAllocationPoolsValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TargetSpecAllocationPoolsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid TargetSpecAllocationPoolsValue Attribute Type",
				"While creating a TargetSpecAllocationPoolsValue value, an invalid attribute value was detected. "+
					"A TargetSpecAllocationPoolsValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TargetSpecAllocationPoolsValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("TargetSpecAllocationPoolsValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra TargetSpecAllocationPoolsValue Attribute Value",
				"While creating a TargetSpecAllocationPoolsValue value, an extra attribute value was detected. "+
					"A TargetSpecAllocationPoolsValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra TargetSpecAllocationPoolsValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewTargetSpecAllocationPoolsValueUnknown(), diags
	}

	endAttribute, ok := attributes["end"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`end is missing from object`)

		return NewTargetSpecAllocationPoolsValueUnknown(), diags
	}

	endVal, ok := endAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`end expected to be ovhtypes.TfStringValue, was: %T`, endAttribute))
	}

	startAttribute, ok := attributes["start"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`start is missing from object`)

		return NewTargetSpecAllocationPoolsValueUnknown(), diags
	}

	startVal, ok := startAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`start expected to be ovhtypes.TfStringValue, was: %T`, startAttribute))
	}

	if diags.HasError() {
		return NewTargetSpecAllocationPoolsValueUnknown(), diags
	}

	return TargetSpecAllocationPoolsValue{
		End:   endVal,
		Start: startVal,
		state: attr.ValueStateKnown,
	}, diags
}

func NewTargetSpecAllocationPoolsValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) TargetSpecAllocationPoolsValue {
	object, diags := NewTargetSpecAllocationPoolsValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewTargetSpecAllocationPoolsValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t TargetSpecAllocationPoolsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewTargetSpecAllocationPoolsValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewTargetSpecAllocationPoolsValueUnknown(), nil
	}

	if in.IsNull() {
		return NewTargetSpecAllocationPoolsValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewTargetSpecAllocationPoolsValueMust(TargetSpecAllocationPoolsValue{}.AttributeTypes(ctx), attributes), nil
}

func (t TargetSpecAllocationPoolsType) ValueType(ctx context.Context) attr.Value {
	return TargetSpecAllocationPoolsValue{}
}

var _ basetypes.ObjectValuable = TargetSpecAllocationPoolsValue{}

type TargetSpecAllocationPoolsValue struct {
	End   ovhtypes.TfStringValue `tfsdk:"end" json:"end"`
	Start ovhtypes.TfStringValue `tfsdk:"start" json:"start"`
	state attr.ValueState
}

func (v TargetSpecAllocationPoolsValue) ToCreate() *TargetSpecAllocationPoolsValue {
	res := &TargetSpecAllocationPoolsValue{}

	if !v.End.IsNull() {
		res.End = v.End
	}

	if !v.Start.IsNull() {
		res.Start = v.Start
	}

	res.state = attr.ValueStateKnown

	return res
}

func (v TargetSpecAllocationPoolsValue) ToUpdate() *TargetSpecAllocationPoolsValue {
	res := &TargetSpecAllocationPoolsValue{}

	if !v.End.IsNull() {
		res.End = v.End
	}

	if !v.Start.IsNull() {
		res.Start = v.Start
	}

	res.state = attr.ValueStateKnown

	return res
}

func (v TargetSpecAllocationPoolsValue) MarshalJSON() ([]byte, error) {
	toMarshal := map[string]any{}
	if !v.End.IsNull() && !v.End.IsUnknown() {
		toMarshal["end"] = v.End
	}
	if !v.Start.IsNull() && !v.Start.IsUnknown() {
		toMarshal["start"] = v.Start
	}

	return json.Marshal(toMarshal)
}

func (v *TargetSpecAllocationPoolsValue) UnmarshalJSON(data []byte) error {
	type JsonTargetSpecAllocationPoolsValue TargetSpecAllocationPoolsValue

	var tmp JsonTargetSpecAllocationPoolsValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.End = tmp.End
	v.Start = tmp.Start

	v.state = attr.ValueStateKnown

	return nil
}

func (v *TargetSpecAllocationPoolsValue) MergeWith(other *TargetSpecAllocationPoolsValue) {

	if (v.End.IsUnknown() || v.End.IsNull()) && !other.End.IsUnknown() {
		v.End = other.End
	}

	if (v.Start.IsUnknown() || v.Start.IsNull()) && !other.Start.IsUnknown() {
		v.Start = other.Start
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v TargetSpecAllocationPoolsValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"end":   v.End,
		"start": v.Start,
	}
}
func (v TargetSpecAllocationPoolsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["end"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["start"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.End.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["end"] = val

		val, err = v.Start.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["start"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v TargetSpecAllocationPoolsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v TargetSpecAllocationPoolsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v TargetSpecAllocationPoolsValue) String() string {
	return "TargetSpecAllocationPoolsValue"
}

func (v TargetSpecAllocationPoolsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"end":   ovhtypes.TfStringType{},
			"start": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"end":   v.End,
			"start": v.Start,
		})

	return objVal, diags
}

func (v TargetSpecAllocationPoolsValue) Equal(o attr.Value) bool {
	other, ok := o.(TargetSpecAllocationPoolsValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.End.Equal(other.End) {
		return false
	}

	if !v.Start.Equal(other.Start) {
		return false
	}

	return true
}

func (v TargetSpecAllocationPoolsValue) Type(ctx context.Context) attr.Type {
	return TargetSpecAllocationPoolsType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v TargetSpecAllocationPoolsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"end":   ovhtypes.TfStringType{},
		"start": ovhtypes.TfStringType{},
	}
}

var _ basetypes.ObjectTypable = TargetSpecLocationType{}

type TargetSpecLocationType struct {
	basetypes.ObjectType
}

func (t TargetSpecLocationType) Equal(o attr.Type) bool {
	other, ok := o.(TargetSpecLocationType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t TargetSpecLocationType) String() string {
	return "TargetSpecLocationType"
}

func (t TargetSpecLocationType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	availabilityZoneAttribute, ok := attributes["availability_zone"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`availability_zone is missing from object`)

		return nil, diags
	}

	availabilityZoneVal, ok := availabilityZoneAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`availability_zone expected to be ovhtypes.TfStringValue, was: %T`, availabilityZoneAttribute))
	}

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return nil, diags
	}

	regionVal, ok := regionAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be ovhtypes.TfStringValue, was: %T`, regionAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return TargetSpecLocationValue{
		AvailabilityZone: availabilityZoneVal,
		Region:           regionVal,
		state:            attr.ValueStateKnown,
	}, diags
}

func NewTargetSpecLocationValueNull() TargetSpecLocationValue {
	return TargetSpecLocationValue{
		state: attr.ValueStateNull,
	}
}

func NewTargetSpecLocationValueUnknown() TargetSpecLocationValue {
	return TargetSpecLocationValue{
		state: attr.ValueStateUnknown,
	}
}

func NewTargetSpecLocationValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (TargetSpecLocationValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing TargetSpecLocationValue Attribute Value",
				"While creating a TargetSpecLocationValue value, a missing attribute value was detected. "+
					"A TargetSpecLocationValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TargetSpecLocationValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid TargetSpecLocationValue Attribute Type",
				"While creating a TargetSpecLocationValue value, an invalid attribute value was detected. "+
					"A TargetSpecLocationValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("TargetSpecLocationValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("TargetSpecLocationValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra TargetSpecLocationValue Attribute Value",
				"While creating a TargetSpecLocationValue value, an extra attribute value was detected. "+
					"A TargetSpecLocationValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra TargetSpecLocationValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewTargetSpecLocationValueUnknown(), diags
	}

	availabilityZoneAttribute, ok := attributes["availability_zone"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`availability_zone is missing from object`)

		return NewTargetSpecLocationValueUnknown(), diags
	}

	availabilityZoneVal, ok := availabilityZoneAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`availability_zone expected to be ovhtypes.TfStringValue, was: %T`, availabilityZoneAttribute))
	}

	regionAttribute, ok := attributes["region"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`region is missing from object`)

		return NewTargetSpecLocationValueUnknown(), diags
	}

	regionVal, ok := regionAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`region expected to be ovhtypes.TfStringValue, was: %T`, regionAttribute))
	}

	if diags.HasError() {
		return NewTargetSpecLocationValueUnknown(), diags
	}

	return TargetSpecLocationValue{
		AvailabilityZone: availabilityZoneVal,
		Region:           regionVal,
		state:            attr.ValueStateKnown,
	}, diags
}

func NewTargetSpecLocationValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) TargetSpecLocationValue {
	object, diags := NewTargetSpecLocationValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewTargetSpecLocationValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t TargetSpecLocationType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewTargetSpecLocationValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewTargetSpecLocationValueUnknown(), nil
	}

	if in.IsNull() {
		return NewTargetSpecLocationValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewTargetSpecLocationValueMust(TargetSpecLocationValue{}.AttributeTypes(ctx), attributes), nil
}

func (t TargetSpecLocationType) ValueType(ctx context.Context) attr.Value {
	return TargetSpecLocationValue{}
}

var _ basetypes.ObjectValuable = TargetSpecLocationValue{}

type TargetSpecLocationValue struct {
	AvailabilityZone ovhtypes.TfStringValue `tfsdk:"availability_zone" json:"availabilityZone"`
	Region           ovhtypes.TfStringValue `tfsdk:"region" json:"region"`
	state            attr.ValueState
}

func (v TargetSpecLocationValue) ToCreate() *TargetSpecLocationValue {
	res := &TargetSpecLocationValue{}

	if !v.AvailabilityZone.IsNull() {
		res.AvailabilityZone = v.AvailabilityZone
	}

	if !v.Region.IsNull() {
		res.Region = v.Region
	}

	res.state = attr.ValueStateKnown

	return res
}

func (v TargetSpecLocationValue) MarshalJSON() ([]byte, error) {
	toMarshal := map[string]any{}
	if !v.AvailabilityZone.IsNull() && !v.AvailabilityZone.IsUnknown() {
		toMarshal["availabilityZone"] = v.AvailabilityZone
	}
	if !v.Region.IsNull() && !v.Region.IsUnknown() {
		toMarshal["region"] = v.Region
	}

	return json.Marshal(toMarshal)
}

func (v *TargetSpecLocationValue) UnmarshalJSON(data []byte) error {
	type JsonTargetSpecLocationValue TargetSpecLocationValue

	var tmp JsonTargetSpecLocationValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.AvailabilityZone = tmp.AvailabilityZone
	v.Region = tmp.Region

	v.state = attr.ValueStateKnown

	return nil
}

func (v *TargetSpecLocationValue) MergeWith(other *TargetSpecLocationValue) {

	if (v.AvailabilityZone.IsUnknown() || v.AvailabilityZone.IsNull()) && !other.AvailabilityZone.IsUnknown() {
		v.AvailabilityZone = other.AvailabilityZone
	}

	if (v.Region.IsUnknown() || v.Region.IsNull()) && !other.Region.IsUnknown() {
		v.Region = other.Region
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v TargetSpecLocationValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"availabilityZone": v.AvailabilityZone,
		"region":           v.Region,
	}
}
func (v TargetSpecLocationValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["availability_zone"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["region"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.AvailabilityZone.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["availability_zone"] = val

		val, err = v.Region.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["region"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v TargetSpecLocationValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v TargetSpecLocationValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v TargetSpecLocationValue) String() string {
	return "TargetSpecLocationValue"
}

func (v TargetSpecLocationValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"availability_zone": ovhtypes.TfStringType{},
			"region":            ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"availability_zone": v.AvailabilityZone,
			"region":            v.Region,
		})

	return objVal, diags
}

func (v TargetSpecLocationValue) Equal(o attr.Value) bool {
	other, ok := o.(TargetSpecLocationValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.AvailabilityZone.Equal(other.AvailabilityZone) {
		return false
	}

	if !v.Region.Equal(other.Region) {
		return false
	}

	return true
}

func (v TargetSpecLocationValue) Type(ctx context.Context) attr.Type {
	return TargetSpecLocationType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v TargetSpecLocationValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"availability_zone": ovhtypes.TfStringType{},
		"region":            ovhtypes.TfStringType{},
	}
}
