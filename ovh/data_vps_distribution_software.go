package ovh

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
)

const vpsDistributionSoftwareFanoutWorkers = 4

var (
	_ datasource.DataSource              = (*vpsDistributionSoftwareDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*vpsDistributionSoftwareDataSource)(nil)
)

func NewVpsDistributionSoftwareDataSource() datasource.DataSource {
	return &vpsDistributionSoftwareDataSource{}
}

type vpsDistributionSoftwareDataSource struct {
	config *Config
}

type vpsDistributionSoftwareModel struct {
	ServiceName  types.String `tfsdk:"service_name"`
	TypeFilter   types.String `tfsdk:"type_filter"`
	StatusFilter types.String `tfsdk:"status_filter"`
	SoftwareIDs  types.List   `tfsdk:"software_ids"`
	Software     types.List   `tfsdk:"software"`
}

func vpsDistributionSoftwareObjectAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":     types.Int64Type,
		"name":   types.StringType,
		"type":   types.StringType,
		"status": types.StringType,
	}
}

func (d *vpsDistributionSoftwareDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps_distribution_software"
}

func (d *vpsDistributionSoftwareDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.config = config
}

func (d *vpsDistributionSoftwareDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists the software currently installed on a VPS. " +
			"Detail for each software id is fetched in parallel and returned in the software block.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				Required:    true,
				Description: "The internal name of the VPS.",
			},
			"type_filter": schema.StringAttribute{
				Optional:    true,
				Description: "Restrict results to a software type: database, environment or webserver.",
				Validators: []validator.String{
					stringvalidator.OneOf("database", "environment", "webserver"),
				},
			},
			"status_filter": schema.StringAttribute{
				Optional:    true,
				Description: "Restrict results to a software status: stable, testing or deprecated.",
				Validators: []validator.String{
					stringvalidator.OneOf("stable", "testing", "deprecated"),
				},
			},
			"software_ids": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "Sorted list of installed software ids (after filtering).",
			},
			"software": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Detail for each installed software id (after filtering).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "Software id.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Software name.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Software type.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Software status.",
						},
					},
				},
			},
		},
	}
}

func (d *vpsDistributionSoftwareDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vpsDistributionSoftwareModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	listEndpoint := fmt.Sprintf("/vps/%s/distribution/software", url.PathEscape(serviceName))

	var ids []int64
	if err := d.config.OVHClient.GetWithContext(ctx, listEndpoint, &ids); err != nil {
		if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == 404 {
			msg := apiErr.Message
			switch {
			case strings.Contains(msg, "Got an invalid (or empty) URL"):
				resp.Diagnostics.AddError(
					"VPS API endpoint not available",
					fmt.Sprintf(
						"the OVHcloud API endpoint %s is not available on this VPS lineup. "+
							"This data source may only work on legacy VPS plans, or the endpoint "+
							"may have been deprecated. See the data source's documentation for "+
							"supported VPS generations.",
						listEndpoint),
				)
				return
			case strings.Contains(msg, "does not exist"):
				resp.Diagnostics.AddError(
					"VPS resource not found",
					fmt.Sprintf(
						"the requested resource at %s does not exist (the VPS may not have "+
							"the required option subscribed, or the resource ID is wrong)",
						listEndpoint),
				)
				return
			}
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling GET %s", listEndpoint),
			err.Error(),
		)
		return
	}

	// Fan out detail lookups, capped at vpsDistributionSoftwareFanoutWorkers in flight.
	details := make([]VPSDistributionSoftware, len(ids))
	errs := make([]error, len(ids))

	var wg sync.WaitGroup
	sem := make(chan struct{}, vpsDistributionSoftwareFanoutWorkers)
	for i, id := range ids {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, sid int64) {
			defer wg.Done()
			defer func() { <-sem }()
			itemEndpoint := fmt.Sprintf("/vps/%s/distribution/software/%d",
				url.PathEscape(serviceName), sid)
			var sw VPSDistributionSoftware
			if err := d.config.OVHClient.GetWithContext(ctx, itemEndpoint, &sw); err != nil {
				errs[idx] = fmt.Errorf("GET %s: %w", itemEndpoint, err)
				return
			}
			details[idx] = sw
		}(i, id)
	}
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			resp.Diagnostics.AddError("Error fetching installed software detail", err.Error())
			return
		}
	}

	typeFilter := data.TypeFilter.ValueString()
	statusFilter := data.StatusFilter.ValueString()

	filtered := make([]VPSDistributionSoftware, 0, len(details))
	for _, sw := range details {
		if typeFilter != "" && sw.Type != typeFilter {
			continue
		}
		if statusFilter != "" && sw.Status != statusFilter {
			continue
		}
		filtered = append(filtered, sw)
	}

	sort.Slice(filtered, func(i, j int) bool { return filtered[i].ID < filtered[j].ID })

	idVals := make([]attr.Value, 0, len(filtered))
	objVals := make([]attr.Value, 0, len(filtered))
	objType := types.ObjectType{AttrTypes: vpsDistributionSoftwareObjectAttrTypes()}
	for _, sw := range filtered {
		idVals = append(idVals, types.Int64Value(int64(sw.ID)))
		obj, diags := types.ObjectValue(vpsDistributionSoftwareObjectAttrTypes(), map[string]attr.Value{
			"id":     types.Int64Value(int64(sw.ID)),
			"name":   types.StringValue(sw.Name),
			"type":   types.StringValue(sw.Type),
			"status": types.StringValue(sw.Status),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		objVals = append(objVals, obj)
	}

	idList, diags := types.ListValue(types.Int64Type, idVals)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	swList, diags := types.ListValue(objType, objVals)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.SoftwareIDs = idList
	data.Software = swList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
