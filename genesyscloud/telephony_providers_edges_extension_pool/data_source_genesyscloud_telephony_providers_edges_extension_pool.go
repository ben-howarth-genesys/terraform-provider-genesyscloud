package telephony_providers_edges_extension_pool

import (
	"context"
	"fmt"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceExtensionPoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	sdkConfig := m.(*provider.ProviderMeta).ClientConfig
	extensionPoolProxy := getExtensionPoolProxy(sdkConfig)

	extensionPoolStartPhoneNumber := d.Get("start_number").(string)
	extensionPoolEndPhoneNumber := d.Get("end_number").(string)

	return util.WithRetries(ctx, 15*time.Second, func() *retry.RetryError {

		extensionPools, resp, getErr := extensionPoolProxy.getAllExtensionPools(ctx)

		if getErr != nil {
			return retry.NonRetryableError(util.BuildWithRetriesApiDiagnosticError(ResourceType, fmt.Sprintf("error requesting list of extension pools: %s", getErr), resp))
		}

		if extensionPools == nil || len(*extensionPools) == 0 {
			return retry.RetryableError(util.BuildWithRetriesApiDiagnosticError(ResourceType, fmt.Sprintf("no extension pools found with start phone number: %s and end phone number: %s", extensionPoolStartPhoneNumber, extensionPoolEndPhoneNumber), resp))
		}

		for _, extensionPool := range *extensionPools {
			if extensionPool.StartNumber != nil && *extensionPool.StartNumber == extensionPoolStartPhoneNumber &&
				extensionPool.EndNumber != nil && *extensionPool.EndNumber == extensionPoolEndPhoneNumber &&
				extensionPool.State != nil && *extensionPool.State != "deleted" {
				d.SetId(*extensionPool.Id)
			}
		}
		return nil
	})

}
