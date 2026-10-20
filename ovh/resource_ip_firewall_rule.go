package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*ipFirewallRuleResource)(nil)
var _ resource.ResourceWithImportState = (*ipFirewallRuleResource)(nil)
var _ resource.ResourceWithValidateConfig = (*ipFirewallRuleResource)(nil)

func NewIpFirewallRuleResource() resource.Resource {
	return &ipFirewallRuleResource{}
}

type ipFirewallRuleResource struct {
	config *Config
}

func (r *ipFirewallRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip_firewall_rule"
}

func (d *ipFirewallRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

func (d *ipFirewallRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = IpFirewallRuleResourceSchema(ctx)
}

func (r *ipFirewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "|")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <ip>|<ip_on_firewall>|<sequence>")
		return
	}

	ip := splits[0]
	ipOnFirewall := splits[1]
	sequence, err := strconv.Atoi(splits[2])
	if err != nil {
		resp.Diagnostics.AddError("Given firewall sequence number must be an integer", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip"), ip)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip_on_firewall"), ipOnFirewall)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("sequence"), sequence)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), fmt.Sprintf("%s/%s/%d", ip, ipOnFirewall, sequence))...)
}

func (r *ipFirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var (
		data         IpFirewallRuleModel
		responseData IpFirewallRuleResponseModel
	)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("/ip/%s/firewall/%s/rule",
		url.PathEscape(data.Ip.ValueString()),
		url.PathEscape(data.IpOnFirewall.ValueString()),
	)
	if err := r.config.OVHClient.Post(endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	endpoint = fmt.Sprintf("/ip/%s/firewall/%s/rule/%d",
		url.PathEscape(data.Ip.ValueString()),
		url.PathEscape(data.IpOnFirewall.ValueString()),
		data.Sequence.ValueInt64(),
	)

	// Wait for state to be ok
	err := retry.RetryContext(ctx, 10*time.Minute, func() *retry.RetryError {
		readErr := r.config.OVHClient.Get(endpoint, &responseData)
		if readErr != nil {
			return retry.NonRetryableError(readErr)
		}

		if responseData.State.ValueString() == "ok" {
			return nil
		}

		return retry.RetryableError(errors.New("waiting for resource state to be ok"))
	})

	if err != nil {
		resp.Diagnostics.AddError("error waiting status to be ok", err.Error())
		return
	}

	responseData.MergeWith(&data)
	responseData.ID = ovhtypes.NewTfStringValue(responseData.Ip.ValueString() + "/" + responseData.IpOnFirewall.ValueString() + "/" + strconv.FormatInt(responseData.Sequence.ValueInt64(), 10))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *ipFirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		data         IpFirewallRuleModel
		responseData IpFirewallRuleResponseModel
	)

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fmt.Sprintf("/ip/%s/firewall/%s/rule/%d",
		url.PathEscape(data.Ip.ValueString()),
		url.PathEscape(data.IpOnFirewall.ValueString()),
		data.Sequence.ValueInt64(),
	)
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)

	// The API returns the ports only as desc strings, such as "eq 22" or "range 8000 8010". They are rebuilt
	// from the current descs, so that a rule changed outside of Terraform is not hidden by the values merged
	// from the state. A desc that cannot be parsed keeps the merged values.
	if port, portRange, err := parsePortDesc(responseData.SourcePortDesc.ValueString()); err != nil {
		tflog.Warn(ctx, fmt.Sprintf("failed to parse source port from desc: %s", err))
	} else {
		responseData.SourcePort, responseData.SourcePortRange = port, portRange
	}
	if port, portRange, err := parsePortDesc(responseData.DestinationPortDesc.ValueString()); err != nil {
		tflog.Warn(ctx, fmt.Sprintf("failed to parse destination port from desc: %s", err))
	} else {
		responseData.DestinationPort, responseData.DestinationPortRange = port, portRange
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *ipFirewallRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update on API side
}

func (r *ipFirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IpFirewallRuleModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := fmt.Sprintf("/ip/%s/firewall/%s/rule/%d",
		url.PathEscape(data.Ip.ValueString()),
		url.PathEscape(data.IpOnFirewall.ValueString()),
		data.Sequence.ValueInt64(),
	)
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for rule to be removed
	err := retry.RetryContext(ctx, 10*time.Minute, func() *retry.RetryError {
		readErr := r.config.OVHClient.Get(endpoint, nil)
		if readErr != nil {
			if errOvh, ok := readErr.(*ovh.APIError); ok && errOvh.Code == 404 {
				return nil
			} else {
				return retry.NonRetryableError(readErr)
			}
		}

		return retry.RetryableError(errors.New("waiting for resource to be removed"))
	})

	if err != nil {
		resp.Diagnostics.AddError("error verifying that resource was deleted", err.Error())
	}
}

// ValidateConfig ensures the end of a port range is greater than its start. A range of a single port is refused too,
// since the API silently converts it into a single port rule, causing a permanent replacement on every plan.
// It also ensures port ranges are only used with the tcp and udp protocols, and tcp_option and fragments only with
// the tcp protocol, as the API refuses them otherwise.
func (r *ipFirewallRuleResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// Ensure port options are set only for TCP or UDP rules
	var protocol ovhtypes.TfStringValue
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("protocol"), &protocol)...)
	// Used to ensure port rules are configured with TCP or UDP transport protocols
	protocolIsTCPOrUDP := protocol.ValueString() == "tcp" || protocol.ValueString() == "udp"

	// Ensure both range properties are not equal to each other
	for _, leafName := range []string{"destination_port_range", "source_port_range"} {
		var from, to ovhtypes.TfInt64Value
		p := path.Root(leafName)
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, p.AtName("from"), &from)...)
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, p.AtName("to"), &to)...)

		// Unknown or null, skip it
		if (from.IsUnknown() || from.IsNull()) || (to.IsUnknown() || to.IsNull()) {
			continue
		}
		if to.ValueInt64() <= from.ValueInt64() {
			resp.Diagnostics.AddAttributeError(p.AtName("to"), "Invalid port range values", "to must be greater than from. For a single port, use source_port or destination_port instead.")
		}
		// Ensure protocol is UDP or TCP if the option is set, and protocol is known
		if (!protocol.IsNull() && !protocol.IsUnknown()) && !protocolIsTCPOrUDP {
			resp.Diagnostics.AddAttributeError(p, "Invalid protocol for port range", "Port ranges can only be used with the tcp or udp protocol.")
		}
	}

	// The API only accepts tcp_option and fragments, sent together as the TCP option, with the tcp protocol
	if !protocol.IsNull() && !protocol.IsUnknown() && protocol.ValueString() != "tcp" {
		var tcpOption ovhtypes.TfStringValue
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("tcp_option"), &tcpOption)...)
		if !tcpOption.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("tcp_option"), "Invalid protocol for tcp_option", "tcp_option can only be used with the tcp protocol.")
		}

		var fragments ovhtypes.TfBoolValue
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("fragments"), &fragments)...)
		if !fragments.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root("fragments"), "Invalid protocol for fragments", "fragments can only be used with the tcp protocol.")
		}
	}
}

// parsePortDesc parses a port desc returned by the API: "eq 22" gives a single port, "range 8000 8010" a range,
// and an empty desc (no port filter) gives a null port and no range.
func parsePortDesc(desc string) (ovhtypes.TfInt64Value, *IPFirewallPortRange, error) {
	noPort := ovhtypes.TfInt64Value{Int64Value: basetypes.NewInt64Null()}
	fields := strings.Fields(desc)
	switch {
	case len(fields) == 0:
		return noPort, nil, nil
	case len(fields) == 2 && fields[0] == "eq":
		port, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return noPort, nil, fmt.Errorf("parsing port %q: %w", desc, err)
		}
		return ovhtypes.NewTfInt64Value(port), nil, nil
	case fields[0] == "range":
		portRange, err := parsePortRange(desc)
		if err != nil {
			return noPort, nil, err
		}
		return noPort, portRange, nil
	default:
		return noPort, nil, fmt.Errorf("unexpected port desc %q", desc)
	}
}

// parsePortRange parses a port desc such as "range 443 445".
// it will ensure the correct types are parsed and returned for the range objects
func parsePortRange(desc string) (*IPFirewallPortRange, error) {
	fields := strings.Fields(desc)
	if len(fields) != 3 || fields[0] != "range" {
		return nil, fmt.Errorf("unexpected port range %q", desc)
	}
	from, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing port range %q: %w", fields[1], err)
	}
	to, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing port range %q: %w", fields[2], err)
	}

	return &IPFirewallPortRange{
		From: ovhtypes.NewTfInt64Value(from),
		To:   ovhtypes.NewTfInt64Value(to),
	}, nil
}
