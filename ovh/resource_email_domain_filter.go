package ovh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
)

var (
	_ resource.ResourceWithConfigure      = (*emailDomainFilterResource)(nil)
	_ resource.ResourceWithImportState    = (*emailDomainFilterResource)(nil)
	_ resource.ResourceWithValidateConfig = (*emailDomainFilterResource)(nil)
)

// Every write on a filter returns a task rather than the result, so each step
// waits for its tasks to leave the queue and then for the change to be visible.
const (
	emailFilterPollInterval = 3 * time.Second
	emailFilterPollTimeout  = 10 * time.Minute
)

// emailDomainFilterAPI mirrors email.domain.Filter. OVHcloud stores the name
// lowercased and matches it case-insensitively.
type emailDomainFilterAPI struct {
	Name        string `json:"name"`
	Priority    int64  `json:"priority"`
	Active      bool   `json:"active"`
	Action      string `json:"action"`
	ActionParam string `json:"actionParam"`
}

// emailDomainFilterRuleAPI mirrors email.domain.Rule.
type emailDomainFilterRuleAPI struct {
	ID      int64  `json:"id"`
	Header  string `json:"header"`
	Operand string `json:"operand"`
	Value   string `json:"value"`
}

// emailDomainFilterTaskAPI is the part of email.domain.TaskFilter this needs.
type emailDomainFilterTaskAPI struct {
	ID int64 `json:"id"`
}

// emailDomainFilterCreate is the creation payload, which carries the first rule.
type emailDomainFilterCreate struct {
	Name        string `json:"name"`
	Priority    int64  `json:"priority"`
	Active      bool   `json:"active"`
	Action      string `json:"action"`
	ActionParam string `json:"actionParam,omitempty"`
	Header      string `json:"header"`
	Operand     string `json:"operand"`
	Value       string `json:"value"`
}

type emailDomainFilterModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	AccountName types.String `tfsdk:"account_name"`
	Name        types.String `tfsdk:"name"`
	Priority    types.Int64  `tfsdk:"priority"`
	Active      types.Bool   `tfsdk:"active"`
	Action      types.String `tfsdk:"action"`
	ActionParam types.String `tfsdk:"action_param"`
	Rules       types.Set    `tfsdk:"rules"`
}

type emailDomainFilterRuleModel struct {
	Header  types.String `tfsdk:"header"`
	Operand types.String `tfsdk:"operand"`
	Value   types.String `tfsdk:"value"`
}

var emailDomainFilterRuleAttrTypes = map[string]attr.Type{
	"header":  types.StringType,
	"operand": types.StringType,
	"value":   types.StringType,
}

// emailFilterRule is a rule reduced to what identifies it. Rules have no order: every
// one of them must match for the filter to apply.
type emailFilterRule struct {
	Header, Operand, Value string
}

// emailFilterActiveDefault plans active as true when the configuration leaves
// it out. It stands in for booldefault, which this repository does not vendor,
// and still lets a filter switched off outside Terraform show up as a diff.
type emailFilterActiveDefault struct{}

func (emailFilterActiveDefault) Description(context.Context) string {
	return "Defaults to true when not configured."
}

func (m emailFilterActiveDefault) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (emailFilterActiveDefault) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.BoolValue(true)
	}
}

func NewEmailDomainFilterResource() resource.Resource {
	return &emailDomainFilterResource{}
}

type emailDomainFilterResource struct {
	config *Config
}

func (r *emailDomainFilterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_domain_filter"
}

func (r *emailDomainFilterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *emailDomainFilterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a filter on a mailbox of an OVHcloud email domain.",
		MarkdownDescription: "Manages a filter on a mailbox of an OVHcloud email domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Identifier of the filter, <domain>/<account_name>/<name> with the name lowercased as OVHcloud stores it",
				MarkdownDescription: "Identifier of the filter, `<domain>/<account_name>/<name>` with the name lowercased as OVHcloud stores it",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Required:            true,
				Description:         "Name of the email domain",
				MarkdownDescription: "Name of the email domain",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_name": schema.StringAttribute{
				Required:            true,
				Description:         "Mailbox the filter applies to, as the local part of its address",
				MarkdownDescription: "Mailbox the filter applies to, as the local part of its address",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				Description: "Name of the filter, unique per mailbox. OVHcloud stores it lowercased and matches it " +
					"case-insensitively, so a change of case alone is applied in place",
				MarkdownDescription: "Name of the filter, unique per mailbox. OVHcloud stores it lowercased and matches it " +
					"case-insensitively, so a change of case alone is applied in place",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(
						func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
							resp.RequiresReplace = !strings.EqualFold(req.StateValue.ValueString(), req.PlanValue.ValueString())
						},
						"Renaming a filter forces a new one, unless only the case changes.",
						"Renaming a filter forces a new one, unless only the case changes.",
					),
				},
			},
			"priority": schema.Int64Attribute{
				Required:            true,
				Description:         "Order in which the mailbox's filters apply, lowest first",
				MarkdownDescription: "Order in which the mailbox's filters apply, lowest first",
			},
			"active": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Whether the filter applies. Defaults to true",
				MarkdownDescription: "Whether the filter applies. Defaults to `true`",
				PlanModifiers: []planmodifier.Bool{
					emailFilterActiveDefault{},
				},
			},
			"action": schema.StringAttribute{
				Required: true,
				Description: "What to do with matching mail: accept, account (move it to another mailbox), " +
					"delete or redirect. There is no endpoint to change it, so a change forces a new filter",
				MarkdownDescription: "What to do with matching mail: `accept`, `account` (move it to another mailbox), " +
					"`delete` or `redirect`. There is no endpoint to change it, so a change forces a new filter",
				Validators: []validator.String{
					stringvalidator.OneOf("accept", "account", "delete", "redirect"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"action_param": schema.StringAttribute{
				Optional: true,
				Description: "Full address the mail goes to, required by the account and redirect actions and " +
					"refused by the others. A change forces a new filter",
				MarkdownDescription: "Full address the mail goes to, required by the `account` and `redirect` actions and " +
					"refused by the others. A change forces a new filter",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules": schema.SetNestedAttribute{
				Required: true,
				Description: "Conditions on the mail's headers. All of them must match for the filter to apply, " +
					"so they have no order",
				MarkdownDescription: "Conditions on the mail's headers. **All** of them must match for the filter to apply, " +
					"so they have no order",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"header": schema.StringAttribute{
							Required:            true,
							Description:         "Header to test, such as From, Subject or Authentication-Results",
							MarkdownDescription: "Header to test, such as `From`, `Subject` or `Authentication-Results`",
						},
						"operand": schema.StringAttribute{
							Required:            true,
							Description:         "How to test it: checkspf, contains or noContains",
							MarkdownDescription: "How to test it: `checkspf`, `contains` or `noContains`",
							Validators: []validator.String{
								stringvalidator.OneOf("checkspf", "contains", "noContains"),
							},
						},
						"value": schema.StringAttribute{
							Required:            true,
							Description:         "Value to test for",
							MarkdownDescription: "Value to test for",
						},
					},
				},
			},
		},
	}
}

func (r *emailDomainFilterResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data emailDomainFilterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || data.Action.IsUnknown() || data.ActionParam.IsUnknown() {
		return
	}

	action := data.Action.ValueString()
	param := data.ActionParam.ValueString()

	switch action {
	case "account", "redirect":
		if param == "" {
			resp.Diagnostics.AddAttributeError(path.Root("action_param"), "Missing action_param",
				fmt.Sprintf("action %q needs action_param, the address the mail goes to.", action))
		} else if !strings.Contains(param, "@") {
			// The API accepts a bare local part and resolves it against OVHcloud's
			// own MX host rather than the domain, so "postmaster" silently means
			// postmaster@mx1.ovh.net and the filter hands matching mail to OVHcloud.
			resp.Diagnostics.AddAttributeError(path.Root("action_param"), "action_param must be a full address",
				fmt.Sprintf("%q is not a full address. OVHcloud resolves a bare local part against its own MX host, "+
					"not against this domain, so the mail would leave it. Write %s@<domain>.", param, param))
		}
	case "accept", "delete":
		if param != "" {
			resp.Diagnostics.AddAttributeError(path.Root("action_param"), "Unexpected action_param",
				fmt.Sprintf("action %q takes no action_param.", action))
		}
	}
}

func emailFilterBase(domain, account string) string {
	return "/email/domain/" + url.PathEscape(domain) + "/account/" + url.PathEscape(account) + "/filter"
}

func emailFilterPath(domain, account, name string) string {
	return emailFilterBase(domain, account) + "/" + url.PathEscape(name)
}

func emailFilterNotFound(err error) bool {
	var apiErr *ovh.APIError
	return errors.As(err, &apiErr) && apiErr.Code == http.StatusNotFound
}

func (r *emailDomainFilterResource) getFilter(domain, account, name string) (*emailDomainFilterAPI, error) {
	var res emailDomainFilterAPI
	if err := r.config.OVHClient.Get(emailFilterPath(domain, account, name), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *emailDomainFilterResource) getRules(domain, account, name string) ([]emailDomainFilterRuleAPI, error) {
	endpoint := emailFilterPath(domain, account, name) + "/rule"

	var ids []int64
	if err := r.config.OVHClient.Get(endpoint, &ids); err != nil {
		return nil, fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	rules := make([]emailDomainFilterRuleAPI, 0, len(ids))
	for _, id := range ids {
		var rule emailDomainFilterRuleAPI
		ruleEndpoint := fmt.Sprintf("%s/%d", endpoint, id)
		if err := r.config.OVHClient.Get(ruleEndpoint, &rule); err != nil {
			return nil, fmt.Errorf("error calling Get %s: %w", ruleEndpoint, err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func emailFilterRuleKeys(rules []emailDomainFilterRuleAPI) map[emailFilterRule]int64 {
	keys := make(map[emailFilterRule]int64, len(rules))
	for _, rule := range rules {
		keys[emailFilterRule{rule.Header, rule.Operand, rule.Value}] = rule.ID
	}
	return keys
}

func emailFilterSameRules(live []emailDomainFilterRuleAPI, want []emailFilterRule) bool {
	if len(live) != len(want) {
		return false
	}
	keys := emailFilterRuleKeys(live)
	for _, rule := range want {
		if _, ok := keys[rule]; !ok {
			return false
		}
	}
	return true
}

// emailFilterRulesFromSet returns the configured rules in a stable order, so the rule sent
// with the creation call does not vary from one run to the next.
func emailFilterRulesFromSet(ctx context.Context, set types.Set) ([]emailFilterRule, diag.Diagnostics) {
	var models []emailDomainFilterRuleModel
	diags := set.ElementsAs(ctx, &models, false)

	rules := make([]emailFilterRule, 0, len(models))
	for _, m := range models {
		rules = append(rules, emailFilterRule{m.Header.ValueString(), m.Operand.ValueString(), m.Value.ValueString()})
	}
	sort.Slice(rules, func(i, j int) bool {
		a, b := rules[i], rules[j]
		if a.Header != b.Header {
			return a.Header < b.Header
		}
		if a.Operand != b.Operand {
			return a.Operand < b.Operand
		}
		return a.Value < b.Value
	})

	return rules, diags
}

func emailFilterRulesToSet(ctx context.Context, live []emailDomainFilterRuleAPI) (types.Set, diag.Diagnostics) {
	models := make([]emailDomainFilterRuleModel, 0, len(live))
	for _, rule := range live {
		models = append(models, emailDomainFilterRuleModel{
			Header:  types.StringValue(rule.Header),
			Operand: types.StringValue(rule.Operand),
			Value:   types.StringValue(rule.Value),
		})
	}
	return types.SetValueFrom(ctx, types.ObjectType{AttrTypes: emailDomainFilterRuleAttrTypes}, models)
}

// emailFilterPoll retries check until it reports done, or gives up after the timeout.
func emailFilterPoll(ctx context.Context, what string, check func() (bool, error)) error {
	deadline := time.Now().Add(emailFilterPollTimeout)
	for {
		done, err := check()
		if err != nil {
			return err
		}
		if done {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("%s: still not done after %s", what, emailFilterPollTimeout)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(emailFilterPollInterval):
		}
	}
}

// waitForTasks waits until none of the given tasks is queued any more. It reads
// the domain-wide queue and looks for the ids themselves, rather than trusting
// the queue's account filter to match however the account is spelled. A task
// whose id could not be read makes it wait for the whole queue to drain.
func (r *emailDomainFilterResource) waitForTasks(ctx context.Context, domain string, tasks []int64) error {
	if len(tasks) == 0 {
		return nil
	}

	pending := make(map[int64]bool, len(tasks))
	drainAll := false
	for _, id := range tasks {
		if id == 0 {
			drainAll = true
		}
		pending[id] = true
	}

	endpoint := "/email/domain/" + url.PathEscape(domain) + "/task/filter"
	return emailFilterPoll(ctx, "waiting for the filter tasks to complete", func() (bool, error) {
		var queued []int64
		if err := r.config.OVHClient.Get(endpoint, &queued); err != nil {
			return false, fmt.Errorf("error calling Get %s: %w", endpoint, err)
		}
		if drainAll {
			return len(queued) == 0, nil
		}
		for _, id := range queued {
			if pending[id] {
				return false, nil
			}
		}
		return true, nil
	})
}

// post sends a write and returns the id of the task it queued.
func (r *emailDomainFilterResource) post(endpoint string, body any) (int64, error) {
	var task emailDomainFilterTaskAPI
	if err := r.config.OVHClient.Post(endpoint, body, &task); err != nil {
		return 0, fmt.Errorf("error calling Post %s: %w", endpoint, err)
	}
	return task.ID, nil
}

// del sends a delete and returns the ids of the tasks it queued.
func (r *emailDomainFilterResource) del(endpoint string) ([]int64, error) {
	var raw json.RawMessage
	if err := r.config.OVHClient.Delete(endpoint, &raw); err != nil {
		return nil, err
	}
	return emailFilterTaskIDs(raw), nil
}

// emailFilterTaskIDs reads the tasks a delete queued. The API schema documents
// an array of email.domain.TaskFilter, but the API answers with a single
// object, so both are accepted. Anything else yields the id 0, which makes
// waitForTasks wait for the whole queue to drain instead.
func emailFilterTaskIDs(raw json.RawMessage) []int64 {
	var many []emailDomainFilterTaskAPI
	if err := json.Unmarshal(raw, &many); err == nil && len(many) > 0 {
		ids := make([]int64, 0, len(many))
		for _, t := range many {
			ids = append(ids, t.ID)
		}
		return ids
	}

	var one emailDomainFilterTaskAPI
	if err := json.Unmarshal(raw, &one); err == nil && one.ID != 0 {
		return []int64{one.ID}
	}

	return []int64{0}
}

func (r *emailDomainFilterResource) waitForRules(ctx context.Context, domain, account, name string, want []emailFilterRule) error {
	return emailFilterPoll(ctx, fmt.Sprintf("waiting for the rules of filter %q to match", name), func() (bool, error) {
		live, err := r.getRules(domain, account, name)
		if err != nil {
			return false, err
		}
		return emailFilterSameRules(live, want), nil
	})
}

func (r *emailDomainFilterResource) waitForFilter(ctx context.Context, domain, account, name, what string, done func(*emailDomainFilterAPI) bool) error {
	return emailFilterPoll(ctx, what, func() (bool, error) {
		f, err := r.getFilter(domain, account, name)
		if err != nil {
			if emailFilterNotFound(err) {
				return false, nil
			}
			return false, err
		}
		return done(f), nil
	})
}

func (r *emailDomainFilterResource) setActivity(ctx context.Context, domain, account, name string, active bool) error {
	task, err := r.post(emailFilterPath(domain, account, name)+"/changeActivity", map[string]bool{"activity": active})
	if err != nil {
		return err
	}
	if err := r.waitForTasks(ctx, domain, []int64{task}); err != nil {
		return err
	}
	return r.waitForFilter(ctx, domain, account, name, fmt.Sprintf("waiting for filter %q to become active=%t", name, active),
		func(f *emailDomainFilterAPI) bool { return f.Active == active })
}

// readInto refreshes data from the API. It keeps the configured spelling of the
// name, since OVHcloud returns it lowercased, and reports whether it still exists.
func (r *emailDomainFilterResource) readInto(ctx context.Context, data *emailDomainFilterModel) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	domain := data.Domain.ValueString()
	account := data.AccountName.ValueString()
	name := data.Name.ValueString()

	f, err := r.getFilter(domain, account, name)
	if err != nil {
		if emailFilterNotFound(err) {
			return false, diags
		}
		diags.AddError("Error reading the filter", err.Error())
		return true, diags
	}

	live, err := r.getRules(domain, account, name)
	if err != nil {
		diags.AddError("Error reading the filter's rules", err.Error())
		return true, diags
	}

	if !strings.EqualFold(f.Name, name) {
		data.Name = types.StringValue(f.Name)
	}
	data.ID = types.StringValue(domain + "/" + account + "/" + strings.ToLower(data.Name.ValueString()))
	data.Priority = types.Int64Value(f.Priority)
	data.Active = types.BoolValue(f.Active)
	data.Action = types.StringValue(f.Action)
	// The API returns an empty string for the actions that take no parameter.
	if f.ActionParam == "" {
		data.ActionParam = types.StringNull()
	} else {
		data.ActionParam = types.StringValue(f.ActionParam)
	}

	rules, d := emailFilterRulesToSet(ctx, live)
	diags.Append(d...)
	data.Rules = rules

	return true, diags
}

func (r *emailDomainFilterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data emailDomainFilterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	account := data.AccountName.ValueString()
	name := data.Name.ValueString()

	rules, diags := emailFilterRulesFromSet(ctx, data.Rules)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	wantActive := data.Active.ValueBool()

	// The creation call carries only one rule. Rules are ANDed, so a filter
	// holding a subset of them matches more mail than intended, and an active
	// one would act on all of it until the rest arrived. So it is created
	// inactive, completed and checked, and only then activated.
	first := rules[0]
	task, err := r.post(emailFilterBase(domain, account), emailDomainFilterCreate{
		Name:        name,
		Priority:    data.Priority.ValueInt64(),
		Active:      false,
		Action:      data.Action.ValueString(),
		ActionParam: data.ActionParam.ValueString(),
		Header:      first.Header,
		Operand:     first.Operand,
		Value:       first.Value,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating the filter", err.Error())
		return
	}

	// The creation is queued, and the filter's identity is already known. Save
	// it before waiting, so a failure below leaves a tainted resource that the
	// next apply replaces, rather than a filter Terraform does not know about
	// whose name then collides with the retry. The next refresh corrects the
	// other attributes, or drops the resource if the filter never appeared.
	data.ID = types.StringValue(domain + "/" + account + "/" + strings.ToLower(name))
	data.Active = types.BoolValue(false)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.waitForTasks(ctx, domain, []int64{task}); err != nil {
		resp.Diagnostics.AddError("Error waiting for the filter to be created", err.Error())
		return
	}
	if err := r.waitForFilter(ctx, domain, account, name, fmt.Sprintf("waiting for filter %q to appear", name),
		func(*emailDomainFilterAPI) bool { return true }); err != nil {
		resp.Diagnostics.AddError("Error waiting for the filter to be created", err.Error())
		return
	}

	// From here the filter exists, so a failure refreshes the state saved above
	// from the API, keeping the tainted resource as close to the truth as it can.
	fail := func(summary string, err error) {
		resp.Diagnostics.AddError(summary, err.Error())
		if _, d := r.readInto(ctx, &data); !d.HasError() {
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}
	}

	var tasks []int64
	for _, rule := range rules[1:] {
		task, err := r.post(emailFilterPath(domain, account, name)+"/rule", rule.body())
		if err != nil {
			fail("Error adding a rule to the filter", err)
			return
		}
		tasks = append(tasks, task)
	}
	if err := r.waitForTasks(ctx, domain, tasks); err != nil {
		fail("Error waiting for the filter's rules to be added", err)
		return
	}
	if err := r.waitForRules(ctx, domain, account, name, rules); err != nil {
		fail("The filter's rules do not match the configuration, so it was left inactive", err)
		return
	}

	if wantActive {
		if err := r.setActivity(ctx, domain, account, name, true); err != nil {
			fail("Error activating the filter", err)
			return
		}
	}

	if _, d := r.readInto(ctx, &data); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (f emailFilterRule) body() map[string]string {
	return map[string]string{"header": f.Header, "operand": f.Operand, "value": f.Value}
}

func (r *emailDomainFilterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data emailDomainFilterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	exists, diags := r.readInto(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *emailDomainFilterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state emailDomainFilterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	account := state.AccountName.ValueString()
	name := state.Name.ValueString()

	// Switch off before touching the rules, and back on only once they are
	// final, so an active filter never runs on a half-edited rule set.
	if state.Active.ValueBool() && !plan.Active.ValueBool() {
		if err := r.setActivity(ctx, domain, account, name, false); err != nil {
			resp.Diagnostics.AddError("Error deactivating the filter", err.Error())
			return
		}
	}

	if !plan.Rules.Equal(state.Rules) {
		want, diags := emailFilterRulesFromSet(ctx, plan.Rules)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := r.replaceRules(ctx, domain, account, name, want); err != nil {
			resp.Diagnostics.AddError("Error updating the filter's rules", err.Error())
			return
		}
	}

	if !plan.Priority.Equal(state.Priority) {
		task, err := r.post(emailFilterPath(domain, account, name)+"/changePriority",
			map[string]int64{"priority": plan.Priority.ValueInt64()})
		if err == nil {
			err = r.waitForTasks(ctx, domain, []int64{task})
		}
		if err == nil {
			want := plan.Priority.ValueInt64()
			err = r.waitForFilter(ctx, domain, account, name, fmt.Sprintf("waiting for filter %q to take priority %d", name, want),
				func(f *emailDomainFilterAPI) bool { return f.Priority == want })
		}
		if err != nil {
			resp.Diagnostics.AddError("Error changing the filter's priority", err.Error())
			return
		}
	}

	if !state.Active.ValueBool() && plan.Active.ValueBool() {
		if err := r.setActivity(ctx, domain, account, name, true); err != nil {
			resp.Diagnostics.AddError("Error activating the filter", err.Error())
			return
		}
	}

	// A change of case in the name needs no call: OVHcloud matches it
	// case-insensitively, so the new spelling is simply what state records.
	if _, d := r.readInto(ctx, &plan); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// replaceRules moves the filter's rules to want. It adds the new rules before
// deleting the old ones: rules are ANDed, so in between the filter matches only
// what both versions match, never more than either. Deleting first could leave
// it matching far more, or with no rule at all.
func (r *emailDomainFilterResource) replaceRules(ctx context.Context, domain, account, name string, want []emailFilterRule) error {
	live, err := r.getRules(domain, account, name)
	if err != nil {
		return err
	}
	have := emailFilterRuleKeys(live)

	var tasks []int64
	for _, rule := range want {
		if _, ok := have[rule]; ok {
			continue
		}
		task, err := r.post(emailFilterPath(domain, account, name)+"/rule", rule.body())
		if err != nil {
			return err
		}
		tasks = append(tasks, task)
	}
	if err := r.waitForTasks(ctx, domain, tasks); err != nil {
		return err
	}

	keep := make(map[emailFilterRule]bool, len(want))
	for _, rule := range want {
		keep[rule] = true
	}

	tasks = nil
	for rule, id := range have {
		if keep[rule] {
			continue
		}
		ids, err := r.del(fmt.Sprintf("%s/rule/%d", emailFilterPath(domain, account, name), id))
		if err != nil && !emailFilterNotFound(err) {
			return fmt.Errorf("error deleting rule %d: %w", id, err)
		}
		tasks = append(tasks, ids...)
	}
	if err := r.waitForTasks(ctx, domain, tasks); err != nil {
		return err
	}

	return r.waitForRules(ctx, domain, account, name, want)
}

func (r *emailDomainFilterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data emailDomainFilterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	account := data.AccountName.ValueString()
	name := data.Name.ValueString()

	tasks, err := r.del(emailFilterPath(domain, account, name))
	if err != nil {
		if emailFilterNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting the filter", err.Error())
		return
	}

	if err := r.waitForTasks(ctx, domain, tasks); err != nil {
		resp.Diagnostics.AddError("Error waiting for the filter to be deleted", err.Error())
		return
	}

	// Wait until it is really gone, so that a replacement with the same name
	// does not collide with it.
	err = emailFilterPoll(ctx, fmt.Sprintf("waiting for filter %q to disappear", name), func() (bool, error) {
		_, err := r.getFilter(domain, account, name)
		if emailFilterNotFound(err) {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the filter to be deleted", err.Error())
	}
}

func (r *emailDomainFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The name goes last and may itself contain a slash.
	parts := strings.SplitN(req.ID, "/", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			fmt.Sprintf("Expected an import ID of the form <domain>/<account_name>/<name>, got %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("account_name"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[2])...)
}
