package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
)

var (
	_ resource.ResourceWithConfigure   = (*emailDomainDkimResource)(nil)
	_ resource.ResourceWithImportState = (*emailDomainDkimResource)(nil)
)

// Values of email.domain.DKIMStatusEnum.
const (
	dkimStatusDisabled    = "disabled"
	dkimStatusEnabled     = "enabled"
	dkimStatusError       = "error"
	dkimStatusModifying   = "modifying"
	dkimStatusToConfigure = "toConfigure"
)

// The API applies the change asynchronously.
const (
	dkimActivationPollInterval = 5 * time.Second
	dkimActivationPollTimeout  = 2 * time.Minute
)

// emailDomainDkimSelectorAPI mirrors email.domain.DKIMSelector.
type emailDomainDkimSelectorAPI struct {
	SelectorName string `json:"selectorName"`
	Cname        string `json:"cname"`
	Status       string `json:"status"`
}

// emailDomainDkimAPI mirrors the payload of GET /email/domain/{domain}/dkim.
type emailDomainDkimAPI struct {
	Status         string                       `json:"status"`
	Autoconfig     bool                         `json:"autoconfig"`
	ActiveSelector *string                      `json:"activeSelector"`
	Selectors      []emailDomainDkimSelectorAPI `json:"selectors"`
}

type emailDomainDkimModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	Status         types.String `tfsdk:"status"`
	Autoconfig     types.Bool   `tfsdk:"autoconfig"`
	ActiveSelector types.String `tfsdk:"active_selector"`
	Selectors      types.List   `tfsdk:"selectors"`
}

var emailDomainDkimSelectorAttrTypes = map[string]attr.Type{
	"selector_name": types.StringType,
	"cname":         types.StringType,
	"status":        types.StringType,
}

func NewEmailDomainDkimResource() resource.Resource {
	return &emailDomainDkimResource{}
}

type emailDomainDkimResource struct {
	config *Config
}

func (r *emailDomainDkimResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_domain_dkim"
}

func (r *emailDomainDkimResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.config = config
}

func (r *emailDomainDkimResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Enables DKIM signing on an OVHcloud email domain (MX Plan / Zimbra). " +
			"Destroying the resource disables DKIM again.",
		MarkdownDescription: "Enables DKIM signing on an OVHcloud email domain (MX Plan / Zimbra). " +
			"Destroying the resource disables DKIM again.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Unique identifier for the resource, the domain name",
				MarkdownDescription: "Unique identifier for the resource, the domain name",
			},
			"domain": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the email domain",
				MarkdownDescription: "Name of the email domain",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "DKIM status: disabled, enabled, error, modifying or toConfigure",
				MarkdownDescription: "DKIM status: `disabled`, `enabled`, `error`, `modifying` or `toConfigure`",
			},
			"autoconfig": schema.BoolAttribute{
				Computed: true,
				Description: "Whether OVHcloud publishes the selector records itself. " +
					"False when the DNS zone is hosted elsewhere, in which case the selector CNAMEs must be created manually",
				MarkdownDescription: "Whether OVHcloud publishes the selector records itself. " +
					"`false` when the DNS zone is hosted elsewhere, in which case the selector CNAMEs must be created manually",
			},
			"active_selector": schema.StringAttribute{
				Computed:            true,
				Description:         "Name of the selector currently used to sign outgoing messages, null while none is active",
				MarkdownDescription: "Name of the selector currently used to sign outgoing messages, null while none is active",
			},
			"selectors": schema.ListNestedAttribute{
				Computed: true,
				Description: "Selectors allocated for this domain. When autoconfig is false, each cname must be " +
					"published in the DNS zone before DKIM can leave the toConfigure status",
				MarkdownDescription: "Selectors allocated for this domain. When `autoconfig` is `false`, each `cname` must be " +
					"published in the DNS zone before DKIM can leave the `toConfigure` status",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"selector_name": schema.StringAttribute{
							Computed:            true,
							Description:         "Name of the selector",
							MarkdownDescription: "Name of the selector",
						},
						"cname": schema.StringAttribute{
							Computed:            true,
							Description:         "The CNAME record to publish, as a full zone-file line",
							MarkdownDescription: "The CNAME record to publish, as a full zone-file line",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							Description:         "Selector status: set, toFix or toSet",
							MarkdownDescription: "Selector status: `set`, `toFix` or `toSet`",
						},
					},
				},
			},
		},
	}
}

// get reads the current DKIM state of a domain.
func (r *emailDomainDkimResource) get(domain string) (*emailDomainDkimAPI, error) {
	endpoint := "/email/domain/" + url.PathEscape(domain) + "/dkim"

	var res emailDomainDkimAPI
	if err := r.config.OVHClient.Get(endpoint, &res); err != nil {
		return nil, fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	return &res, nil
}

// waitForActivation polls until the domain reaches a settled status after an
// enable call, so that the state saved by Terraform is not a transient one.
//
// Both "modifying" and "disabled" are transient here. The API acknowledges the
// enable, reports "modifying" for a while, and can then still answer "disabled"
// for several seconds before it flips to "enabled". Treating "disabled" as
// terminal saves a state that says DKIM is off moments before it comes on,
// which then only self-corrects on the next refresh.
func (r *emailDomainDkimResource) waitForActivation(ctx context.Context, domain string, res *emailDomainDkimAPI) (*emailDomainDkimAPI, error) {
	deadline := time.Now().Add(dkimActivationPollTimeout)

	for res.Status == dkimStatusModifying || res.Status == dkimStatusDisabled {
		if time.Now().After(deadline) {
			return res, fmt.Errorf("DKIM for domain %q is still %q after %s", domain, res.Status, dkimActivationPollTimeout)
		}

		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case <-time.After(dkimActivationPollInterval):
		}

		var err error
		if res, err = r.get(domain); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// toModel maps an API payload onto the Terraform model.
func (data *emailDomainDkimModel) toModel(domain string, res *emailDomainDkimAPI) error {
	data.ID = types.StringValue(domain)
	data.Domain = types.StringValue(domain)
	data.Status = types.StringValue(res.Status)
	data.Autoconfig = types.BoolValue(res.Autoconfig)

	if res.ActiveSelector == nil {
		data.ActiveSelector = types.StringNull()
	} else {
		data.ActiveSelector = types.StringValue(*res.ActiveSelector)
	}

	// The API does not guarantee an order, and has been observed returning the
	// same two selectors swapped between consecutive calls. A Terraform list is
	// order-sensitive, so sort to keep the attribute stable across refreshes.
	ordered := slices.Clone(res.Selectors)
	slices.SortFunc(ordered, func(a, b emailDomainDkimSelectorAPI) int {
		return strings.Compare(a.SelectorName, b.SelectorName)
	})

	selectors := make([]attr.Value, 0, len(ordered))
	for _, s := range ordered {
		obj, diags := types.ObjectValue(emailDomainDkimSelectorAttrTypes, map[string]attr.Value{
			"selector_name": types.StringValue(s.SelectorName),
			"cname":         types.StringValue(s.Cname),
			"status":        types.StringValue(s.Status),
		})
		if diags.HasError() {
			return fmt.Errorf("error building selector %q", s.SelectorName)
		}
		selectors = append(selectors, obj)
	}

	list, diags := types.ListValue(types.ObjectType{AttrTypes: emailDomainDkimSelectorAttrTypes}, selectors)
	if diags.HasError() {
		return fmt.Errorf("error building the selectors list")
	}
	data.Selectors = list

	return nil
}

// warnIfNotEnabled surfaces the states where the API accepted the call but DKIM
// is not actually signing yet, which is otherwise invisible from a successful apply.
func warnIfNotEnabled(diags interface {
	AddWarning(summary, detail string)
}, domain string, res *emailDomainDkimAPI) {
	switch res.Status {
	case dkimStatusEnabled:
		// Signing is active, nothing to report.
	case dkimStatusToConfigure:
		detail := fmt.Sprintf(
			"DKIM for domain %q is enabled but still waiting on DNS, so outgoing messages are not signed yet.",
			domain,
		)
		if !res.Autoconfig {
			detail += " The DNS zone is not hosted by OVHcloud (autoconfig is false), so the CNAME of every selector" +
				" listed in the selectors attribute has to be published manually."
		}
		diags.AddWarning("DKIM is not active yet", detail)
	case dkimStatusError:
		diags.AddWarning(
			"DKIM is in error state",
			fmt.Sprintf("The API reports status %q for domain %q. Check the selector records in the DNS zone.",
				dkimStatusError, domain),
		)
	}
}

func (r *emailDomainDkimResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data emailDomainDkimModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()

	endpoint := "/email/domain/" + url.PathEscape(domain) + "/dkim/enable"
	if err := r.config.OVHClient.Put(endpoint, nil, nil); err != nil {
		// DKIM is already on, so the object exists but Terraform does not know
		// about it. Point at import rather than leaving the raw 409.
		var apiErr *ovh.APIError
		if errors.As(err, &apiErr) && apiErr.Code == http.StatusConflict {
			resp.Diagnostics.AddError(
				"DKIM is already enabled for this domain",
				fmt.Sprintf("DKIM is already enabled for %q, so it cannot be created. Import the existing "+
					"configuration instead:\n\n    terraform import ovh_email_domain_dkim.<name> %s\n\nAPI said: %s",
					domain, domain, err.Error()),
			)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
		return
	}

	res, err := r.get(domain)
	if err != nil {
		resp.Diagnostics.AddError("Error reading DKIM configuration", err.Error())
		return
	}

	if res, err = r.waitForActivation(ctx, domain, res); err != nil {
		resp.Diagnostics.AddError("Error waiting for DKIM activation", err.Error())
		return
	}

	if err := data.toModel(domain, res); err != nil {
		resp.Diagnostics.AddError("Error mapping DKIM configuration", err.Error())
		return
	}

	warnIfNotEnabled(&resp.Diagnostics, domain, res)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *emailDomainDkimResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data emailDomainDkimModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()

	res, err := r.get(domain)
	if err != nil {
		resp.Diagnostics.AddError("Error reading DKIM configuration", err.Error())
		return
	}

	// DKIM was turned off outside Terraform: the resource no longer exists.
	if res.Status == dkimStatusDisabled {
		resp.State.RemoveResource(ctx)
		return
	}

	if err := data.toModel(domain, res); err != nil {
		resp.Diagnostics.AddError("Error mapping DKIM configuration", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update only ever runs for changes that do not require replacement. The domain
// is the sole configurable attribute and it forces a new resource, so this
// refreshes the computed attributes rather than calling the API.
func (r *emailDomainDkimResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data emailDomainDkimModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()

	res, err := r.get(domain)
	if err != nil {
		resp.Diagnostics.AddError("Error reading DKIM configuration", err.Error())
		return
	}

	if err := data.toModel(domain, res); err != nil {
		resp.Diagnostics.AddError("Error mapping DKIM configuration", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *emailDomainDkimResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data emailDomainDkimModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()

	endpoint := "/email/domain/" + url.PathEscape(domain) + "/dkim/disable"
	if err := r.config.OVHClient.Put(endpoint, nil, nil); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Put %s", endpoint), err.Error())
	}
}

func (r *emailDomainDkimResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}
