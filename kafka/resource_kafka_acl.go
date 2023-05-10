package kafka

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mdhwk/terraform-provider-kafka/kafka/client"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kmsg"
)

func resourceACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceACLCreate,
		ReadContext:   resourceACLRead,
		DeleteContext: resourceACLDelete,
		Schema: map[string]*schema.Schema{
			"resource_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"Topic", "Group", "Cluster", "TransactionalID", "TransactionalID", "Any"},
			},
			"resource_pattern_type_filter": {
				Type:     schema.TypeString,
				Default:  "Literal",
				Optional: true,
				ForceNew: true,
				//ValidateFunc: validateResourcePatterns,
				//ExactlyOneOf: []string{"Prefixed", "Any", "Match", "Literal"},
			},
			"acl_principal": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"acl_host": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"acl_operation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				//ValidateFunc: validateAclOperations,
				//AtLeastOneOf: []string{"Unknown", "Any", "All", "Read", "Write", "Create", "Delete", "Alter", "Describe", "ClusterAction", "DescribeConfigs", "AlterConfigs", "IdempotentWrite"},
			},
			"acl_permission_type": {
				Type:     schema.TypeString,
				Default:  "Allow",
				Optional: true,
				ForceNew: true,
				//AtLeastOneOf: []string{"Allow", "Deny", "Any", "Unknown"},
			},
		},
	}
}

// func validateAclOperations(i interface{}, s string) (ws []string, err []error) {
// 	ops := []string{"Unknown", "Any", "All", "Read", "Write", "Create", "Delete", "Alter", "Describe", "ClusterAction", "DescribeConfigs", "AlterConfigs", "IdempotentWrite"}
// 	set := make(map[string]struct{}, len(ops))

// 	for _, o := range ops {
// 		set[strings.ToLower(o)] = struct{}{}
// 	}

// 	if _, ok := set[s]; !ok {
// 		err = append(err, fmt.Errorf("%s is not a valid operation"))
// 	}

// 	return
// }

// func validateResourcePatterns(i interface{}, s string) ([]string, []error) {

// }

func resourceACLCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*kadm.Client)
	var diags diag.Diagnostics

	aclDetails, diags := getACLDetails(ctx, d)
	if diags != nil {
		return diags
	}

	acls, err := client.CreateACLs(ctx, c, &aclDetails)
	if err != nil {
		return diag.FromErr(err)
	}

	if acls[0].Err != nil {
		return diag.Errorf("ACL: %v, err: %v", acls[0], err)
	}

	d.SetId(makeACLIdentifier(&acls[0]))

	resourceACLRead(ctx, d, m)

	return diags
}

func resourceACLRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*kadm.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	aclDetails, diags := getACLDetails(ctx, d)
	if diags != nil {
		return diags
	}

	acls, err := client.ReadACLs(ctx, c, &aclDetails)
	if err != nil {
		return diag.FromErr(err)
	}

	if acls[0].Err != nil {
		return diag.Errorf("ACL: %v, err: %v", acls[0], err)
	}

	if len(acls) < 1 {
		return nil
	}

	d.SetId(makeACLIdentifier(&acls[0]))

	if err := d.Set("resource_name", acls[0].Name); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("resource_type", acls[0].Type); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("resource_pattern_type_filter", acls[0].Pattern); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("acl_principal", acls[0].Principal); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("acl_host", acls[0].Host); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("acl_operation", acls[0].Operation); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("acl_permission_type", acls[0].Permission); err != nil {
		diag.FromErr(err)
	}

	return diags
}

func resourceACLDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*kadm.Client)
	var diags diag.Diagnostics

	aclDetails, diags := getACLDetails(ctx, d)
	if diags != nil {
		return diags
	}

	acls, err := client.DeleteACLs(ctx, c, &aclDetails)
	if err != nil {
		return diag.FromErr(err)
	}

	if acls[0].Err != nil {
		return diag.Errorf("ACL: %v, err: %v", acls[0], err)
	}

	d.SetId("")

	return diags
}

func getACLDetails(ctx context.Context, d *schema.ResourceData) (acl kadm.ACLBuilder, diags diag.Diagnostics) {
	resourceName := d.Get("resource_name").(string)
	resourceType := d.Get("resource_type").(string)
	resourcePatternTypeFilter := d.Get("resource_pattern_type_filter").(string)
	aclPrincipal := d.Get("acl_principal").(string)
	aclPermissionType := d.Get("acl_permission_type").(string)
	aclHost := d.Get("acl_host").(string)
	aclOperation := d.Get("acl_operation").(string)

	switch strings.ToLower(resourceType) {
	case "topic":
		acl.Topics(resourceName)
	case "group":
		acl.Groups(resourceName)
	// case "cluster":
	// 	acl.(resourceName)
	// case "transactionalid":
	// 	acl.TransactionalIDs(resourceName)
	case "unknown", "any":
		acl.AnyResource(resourceName)
	}

	switch aclPermissionType {
	case "Allow", "Any", "Unknown":
		acl.AllowHosts(aclHost)
	case "Deny":
		acl.DenyHosts(aclHost)
	}

	switch aclPermissionType {
	case "Allow", "Any", "Unknown":
		acl.Allow(aclPrincipal)
	case "Deny":
		acl.Deny(aclPrincipal)
	}

	pt, err := kmsg.ParseACLResourcePatternType(resourcePatternTypeFilter)
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		acl.ResourcePatternType(pt)
	}

	kop, err := kmsg.ParseACLOperation(strings.ToLower(aclOperation))
	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
	} else {
		acl.Operations(kop)
	}

	return acl, diags
}

func makeACLIdentifier(res *client.ACLResult) string {
	return strings.ToLower(strings.Join([]string{*res.Principal, *res.Host, res.Operation.String(), res.Type.String(), *res.Name, res.Pattern.String(), res.Permission.String()}, "|"))
}
