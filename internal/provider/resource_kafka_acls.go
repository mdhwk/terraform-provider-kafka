package kafka

// import (
// 	"context"
// 	"strings"

// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/mdhwk/terraform-provider-kafka/kafka/client"
// 	"github.com/twmb/franz-go/pkg/kadm"
// 	"github.com/twmb/franz-go/pkg/kmsg"
// )

// func resourceACLs() *schema.Resource {
// 	return &schema.Resource{
// 		CreateContext: resourceACLCreate,
// 		ReadContext:   resourceACLRead,
// 		DeleteContext: resourceACLDelete,
// 		Schema: map[string]*schema.Schema{
// 			"resource_name": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 				ForceNew: true,
// 			},
// 			"resource_type": {
// 				Type:     schema.TypeString,
// 				Required: true,
// 				ForceNew: true,
// 				//ExactlyOneOf: []string{"Topic", "Group", "Cluster", "TransactionalID", "TransactionalID", "Any"},
// 			},
// 			"resource_pattern_type_filter": {
// 				Type:     schema.TypeString,
// 				Default:  "Literal",
// 				Optional: true,
// 				ForceNew: true,
// 				//ExactlyOneOf: []string{"Prefixed", "Any", "Match", "Literal"},
// 			},
// 			"acl_permission_type": {
// 				Type:     schema.TypeString,
// 				Default:  "Allow",
// 				Optional: true,
// 				ForceNew: true,
// 				//AtLeastOneOf: []string{"Allow", "Deny", "Any", "Unknown"},
// 			},
// 			"acl_principals": {
// 				Type: schema.TypeSet,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 					//ValidateFunc: verify.ValidARN,
// 				},
// 				Required: true,
// 				ForceNew: true,
// 			},
// 			"acl_hosts": {
// 				Type: schema.TypeSet,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 					//ValidateFunc: verify.ValidARN,
// 				},
// 				Optional: true,
// 				ForceNew: true,
// 			},
// 			"acl_operations": {
// 				Type: schema.TypeSet,
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 					//ValidateFunc: verify.ValidARN,
// 				},
// 				Optional: true,
// 				ForceNew: true,
// 				//AtLeastOneOf: []string{"Unknown", "Any", "All", "Read", "Write", "Create", "Delete", "Alter", "Describe", "ClusterAction", "DescribeConfigs", "AlterConfigs", "IdempotentWrite"},
// 			},
// 		},
// 	}
// }

// func resourceACLsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(*kadm.Client)
// 	var diags diag.Diagnostics

// 	aclDetails, diags := getACLDetails(ctx, d)
// 	if diags != nil {
// 		return diags
// 	}

// 	acls, err := client.CreateACLs(ctx, c, &aclDetails)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	d.SetId(makeACLIdentifier(&acls[0]))

// 	resourceACLRead(ctx, d, m)

// 	return diags
// }

// func resourceACLsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(*kadm.Client)

// 	// Warning or errors can be collected in a slice type
// 	var diags diag.Diagnostics

// 	aclDetails, diags := getACLDetails(ctx, d)
// 	if diags != nil {
// 		return diags
// 	}

// 	acls, err := client.ReadACLs(ctx, c, &aclDetails)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}

// 	if len(acls) < 1 {
// 		return nil
// 	}

// 	d.SetId(makeACLIdentifier(&acls[0]))

// 	if err := d.Set("resource_name", acls[0].Name); err != nil {
// 		diag.FromErr(err)
// 	}
// 	if err := d.Set("resource_type", acls[0].Type); err != nil {
// 		diag.FromErr(err)
// 	}
// 	if err := d.Set("resource_pattern_type_filter", acls[0].Pattern); err != nil {
// 		diag.FromErr(err)
// 	}
// 	if err := d.Set("acl_principals", acls[0].Principal); err != nil {
// 		diag.FromErr(err)
// 	}
// 	if err := d.Set("acl_hosts", acls[0].Host); err != nil {
// 		diag.FromErr(err)
// 	}
// 	if err := d.Set("acl_operations", acls[0].Operation); err != nil {
// 		diag.FromErr(err)
// 	}
// 	if err := d.Set("acl_permission_type", acls[0].Permission); err != nil {
// 		diag.FromErr(err)
// 	}

// 	return diags
// }

// func tfSet() {

// }

// func resourceACLsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	_ = m.(*kadm.Client)
// 	// Warning or errors can be collected in a slice type
// 	// var diags diag.Diagnostics

// 	// applicationID := d.Id()

// 	// err := c.DeleteACL(applicationID)
// 	// if err != nil {
// 	// 	return diag.FromErr(err)
// 	// }

// 	d.SetId("")

// 	return nil

// }

// func getACLsDetails(ctx context.Context, d *schema.ResourceData) (acl kadm.ACLBuilder, diags diag.Diagnostics) {
// 	var (
// 		aclHosts      []string
// 		aclOperations []string
// 		aclPrincipals []string
// 	)

// 	resourceName := d.Get("resource_name").(string)
// 	resourceType := d.Get("resource_type").(string)
// 	resourcePatternTypeFilter := d.Get("resource_pattern_type_filter").(string)
// 	aclPermissionType := d.Get("acl_permission_type").(string)
// 	if val, ok := d.GetOkExists("acl_hosts"); ok {
// 		aclHosts = stringValueSlice(val.(*schema.Set).List())
// 	}
// 	if val, ok := d.GetOkExists("acl_principals"); ok {
// 		aclPrincipals = stringValueSlice(val.(*schema.Set).List())
// 	}
// 	if val, ok := d.GetOkExists("acl_operations"); ok {
// 		aclOperations = stringValueSlice(val.(*schema.Set).List())
// 	}

// 	return acl, diags
// }

// func makeACLsIdentifier(res *client.ACLResult) string {
// 	return strings.Join([]string{*res.Principal, *res.Host, res.Operation.String(), res.Type.String(), *res.Name, res.Pattern.String(), res.Permission.String()}, "|")
// }

// func populateAclBuilder(acl *kadm.ACLBuilder) kadm.ACLBuilder {
// 	switch strings.ToLower(resourceType) {
// 	case "topic":
// 		acl.Topics(resourceName)
// 	case "group":
// 		acl.Groups(resourceName)
// 	// case "cluster":
// 	// 	acl.(resourceName)
// 	// case "transactionalid":
// 	// 	acl.TransactionalIDs(resourceName)
// 	case "unknown", "any":
// 		acl.AnyResource(resourceName)
// 	}

// 	for _, host := range aclHosts {
// 		switch aclPermissionType {
// 		case "Allow", "Any", "Unknown":
// 			acl.AllowHosts(host)
// 		case "Deny":
// 			acl.DenyHosts(host)
// 		}
// 	}

// 	switch aclPermissionType {
// 	case "Allow", "Any", "Unknown":
// 		acl.Allow(aclPrincipal)
// 	case "Deny":
// 		acl.Deny(aclPrincipal)
// 	}

// 	pt, err := kmsg.ParseACLResourcePatternType(resourcePatternTypeFilter)
// 	if err != nil {
// 		diags = append(diags, diag.FromErr(err)...)
// 	} else {
// 		acl.ResourcePatternType(pt)
// 	}

// 	var kop kmsg.ACLOperation
// 	for _, op := range aclOperations {
// 		var err error
// 		if kop, err = kmsg.ParseACLOperation(strings.ToLower(op)); err != nil {
// 			diags = append(diags, diag.FromErr(err)...)
// 			continue
// 		}
// 		acl.Operations(kop)
// 	}
// }
