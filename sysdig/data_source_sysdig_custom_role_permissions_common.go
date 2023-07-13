package sysdig

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func dataSourceSysdigCustomRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"requested_permissions": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"permissions": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func dataSourceSysdigCustomRoleMonitorPermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()

	if err != nil {
		return diag.FromErr(err)
	}
	rp := d.Get("requested_permissions").([]interface{})

	rps := readPermissions(rp)
	permissions, err := client.GetMonitorPermissionsWithDependencies(ctx, rps)

	if err != nil {
		return diag.FromErr(err)
	}
	ps := make([]string, len(permissions))
	for i, permission := range permissions {
		ps[i] = permission.Authority
	}

	cdefChecksum := sha256.Sum256([]byte(strings.Join(rps, ",")))
	d.SetId(fmt.Sprintf("%x", cdefChecksum))
	_ = d.Set("permissions", ps)

	return nil
}

// TODO create a common function
func dataSourceSysdigCustomRoleSecurePermissionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()

	if err != nil {
		return diag.FromErr(err)
	}
	rp := d.Get("requested_permissions").([]interface{})

	rps := readPermissions(rp)
	permissions, err := client.GetSecurePermissionsWithDependencies(ctx, rps)

	if err != nil {
		return diag.FromErr(err)
	}
	ps := make([]string, len(permissions))
	for i, permission := range permissions {
		ps[i] = permission.Authority
	}

	cdefChecksum := sha256.Sum256([]byte(strings.Join(rps, ",")))
	d.SetId(fmt.Sprintf("%x", cdefChecksum))
	_ = d.Set("permissions", ps)

	return nil
}

func readPermissions(rp []interface{}) []string {
	permissions := make([]string, len(rp))
	for i, permission := range rp {
		permissions[i] = permission.(string)
	}
	return permissions
}
