package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
)

var (
	_ resource.ResourceWithConfigure   = (*emailDomainRedirectionResource)(nil)
	_ resource.ResourceWithImportState = (*emailDomainRedirectionResource)(nil)
)

// Creating and updating a redirection return a task rather than the object, so
// the result has to be polled for.
const (
	emailRedirectionPollInterval = 2 * time.Second
	emailRedirectionPollTimeout  = 2 * time.Minute
)

// emailDomainRedirectionAPI mirrors email.domain.RedirectionGlobal. Note that
// localCopy is absent: the API accepts it on creation but never returns it.
type emailDomainRedirectionAPI struct {
	ID   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
}

// emailDomainRedirectionCreate mirrors email.domain.RedirectionCreation.
type emailDomainRedirectionCreate struct {
	From      string `json:"from"`
	To        string `json:"to"`
	LocalCopy bool   `json:"localCopy"`
}

type emailDomainRedirectionModel struct {
	ID        types.String `tfsdk:"id"`
	Domain    types.String `tfsdk:"domain"`
	From      types.String `tfsdk:"from"`
	To        types.String `tfsdk:"to"`
	LocalCopy types.Bool   `tfsdk:"local_copy"`
}

func NewEmailDomainRedirectionResource() resource.Resource {
	return &emailDomainRedirectionResource{}
}

type emailDomainRedirectionResource struct {
	config *Config
}

func (r *emailDomainRedirectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_domain_redirection"
}

func (r *emailDomainRedirectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *emailDomainRedirectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages an email address redirection on an OVHcloud email domain.",
		MarkdownDescription: "Manages an email address redirection on an OVHcloud email domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				// Deliberately no UseStateForUnknown: changing "to" goes through
				// changeRedirection, which assigns a new id, so promising the old
				// one would survive the apply produces "inconsistent result after
				// apply". Planning it as unknown on every update is the honest
				// description of what the API does.
				Description:         "Identifier of the redirection, assigned by OVHcloud. Changing the target reassigns it",
				MarkdownDescription: "Identifier of the redirection, assigned by OVHcloud. Changing the target reassigns it",
			},
			"domain": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the email domain",
				MarkdownDescription: "Name of the email domain",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"from": schema.StringAttribute{
				Required:            true,
				Description:         "Address to redirect, which must belong to the domain",
				MarkdownDescription: "Address to redirect, which must belong to the domain",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"to": schema.StringAttribute{
				Required:            true,
				Description:         "Address to redirect to, which may be outside the domain",
				MarkdownDescription: "Address to redirect to, which may be outside the domain",
			},
			"local_copy": schema.BoolAttribute{
				Optional: true,
				Description: "Whether to keep a copy in the source mailbox, which requires that mailbox to exist. " +
					"Defaults to false. " +
					"Set on creation only: the API never returns it, so its value is remembered from the configuration " +
					"and a change forces a new redirection",
				MarkdownDescription: "Whether to keep a copy in the source mailbox, which requires that mailbox to exist. " +
					"Defaults to `false`. " +
					"Set on creation only: the API never returns it, so its value is remembered from the configuration " +
					"and a change forces a new redirection",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// get reads a single redirection.
func (r *emailDomainRedirectionResource) get(domain, id string) (*emailDomainRedirectionAPI, error) {
	endpoint := "/email/domain/" + url.PathEscape(domain) + "/redirection/" + url.PathEscape(id)

	var res emailDomainRedirectionAPI
	if err := r.config.OVHClient.Get(endpoint, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// findID looks a redirection up by its endpoints. Creating one returns a task
// rather than the object, so this is how the assigned id is recovered.
func (r *emailDomainRedirectionResource) findID(domain, from, to string) (string, error) {
	endpoint := fmt.Sprintf("/email/domain/%s/redirection?from=%s&to=%s",
		url.PathEscape(domain), url.QueryEscape(from), url.QueryEscape(to))

	var ids []string
	if err := r.config.OVHClient.Get(endpoint, &ids); err != nil {
		return "", fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	if len(ids) == 0 {
		return "", nil
	}

	return ids[0], nil
}

// waitForRedirection polls until a redirection from -> to exists, since the
// create and update calls only hand back a task.
func (r *emailDomainRedirectionResource) waitForRedirection(ctx context.Context, domain, from, to string) (string, error) {
	deadline := time.Now().Add(emailRedirectionPollTimeout)

	for {
		id, err := r.findID(domain, from, to)
		if err != nil {
			return "", err
		}
		if id != "" {
			return id, nil
		}

		if time.Now().After(deadline) {
			return "", fmt.Errorf("redirection %s -> %s on domain %q did not appear within %s",
				from, to, domain, emailRedirectionPollTimeout)
		}

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(emailRedirectionPollInterval):
		}
	}
}

func (r *emailDomainRedirectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data emailDomainRedirectionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	from := data.From.ValueString()
	to := data.To.ValueString()

	endpoint := "/email/domain/" + url.PathEscape(domain) + "/redirection"
	body := emailDomainRedirectionCreate{
		From:      from,
		To:        to,
		LocalCopy: data.LocalCopy.ValueBool(),
	}
	if err := r.config.OVHClient.Post(endpoint, body, nil); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	id, err := r.waitForRedirection(ctx, domain, from, to)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the redirection to be created", err.Error())
		return
	}

	data.ID = types.StringValue(id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *emailDomainRedirectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data emailDomainRedirectionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()

	res, err := r.get(domain, data.ID.ValueString())
	if err != nil {
		var apiErr *ovh.APIError
		if errors.As(err, &apiErr) && apiErr.Code == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading the redirection", err.Error())
		return
	}

	data.ID = types.StringValue(res.ID)
	data.From = types.StringValue(res.From)
	data.To = types.StringValue(res.To)
	// local_copy is deliberately left as-is: the API does not return it, so the
	// configured value is the only record of what was asked for.

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *emailDomainRedirectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state emailDomainRedirectionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := plan.Domain.ValueString()
	id := state.ID.ValueString()

	// Only "to" can change without replacement, and it has its own endpoint.
	endpoint := "/email/domain/" + url.PathEscape(domain) + "/redirection/" + url.PathEscape(id) + "/changeRedirection"
	body := struct {
		To string `json:"to"`
	}{To: plan.To.ValueString()}

	if err := r.config.OVHClient.Post(endpoint, body, nil); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Changing the target replaces the redirection, so the id moves. Wait for
	// the new pairing to exist rather than assuming the old id still resolves.
	newID, err := r.waitForRedirection(ctx, domain, plan.From.ValueString(), plan.To.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the redirection to be updated", err.Error())
		return
	}

	plan.ID = types.StringValue(newID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *emailDomainRedirectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data emailDomainRedirectionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/email/domain/" + url.PathEscape(data.Domain.ValueString()) +
		"/redirection/" + url.PathEscape(data.ID.ValueString())

	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		var apiErr *ovh.APIError
		if errors.As(err, &apiErr) && apiErr.Code == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
	}
}

func (r *emailDomainRedirectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	domain, id, found := strings.Cut(req.ID, "/")
	if !found || domain == "" || id == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected an import ID of the form <domain>/<id>, got %q. "+
				"List the ids for a domain with: ovhcloud email-domain redirection list <domain>", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), domain)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
