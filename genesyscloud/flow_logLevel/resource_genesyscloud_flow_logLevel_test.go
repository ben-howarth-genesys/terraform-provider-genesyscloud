package flow_logLevel

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mypurecloud/platform-client-sdk-go/v125/platformclientv2"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"testing"
)

// these operations need to be performed on a given flow.
func getTestFlow() {

}

func TestAccResourceFlowLogLevel(t *testing.T) {
	var (
		resourceId            = "flow_log_level" + uuid.NewString()
		communications        = false
		eventError            = true
		eventOther            = false
		eventWarning          = true
		executionInputOutputs = false
		executionItems        = true
		flowLoglevel          = "Base"
		names                 = false
		variables             = false
		flowId                = "e3aebe90-5a65-409e-9775-43d547b66e07"
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			util.TestAccPreCheck(t)
		},
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				// Create using only flow log level
				Config: generateFlowLogLevelWithoutCharacteristicsResource(
					flowId,
					flowLoglevel,
					resourceId,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_log_level", flowLoglevel),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_id", flowId),
				),
			},
			{
				// Update using all fields
				Config: generateFlowLogLevelResource(
					communications,
					eventError,
					eventOther,
					eventWarning,
					executionInputOutputs,
					executionItems,
					flowId,
					flowLoglevel,
					names,
					resourceId,
					variables,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_log_level", flowLoglevel),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_id", flowId),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.communications", "false"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.event_error", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.event_other", "false"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.event_warning", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.execution_input_outputs", "false"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.execution_items", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.names", "false"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.variables", "false"),
				),
			},
		},
		CheckDestroy: testVerifyFlowLogLevelDestroyed,
	})
}

func testVerifyFlowLogLevelDestroyed(state *terraform.State) error {
	architectAPI := platformclientv2.NewArchitectApi()
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "genesyscloud_flow_loglevel" {
			continue
		}
		expandArray := []string{"logLevelCharacteristics.characteristics"}
		flowLogLevel, resp, err := architectAPI.GetFlowInstancesSettingsLoglevels(rs.Primary.ID, expandArray)
		if flowLogLevel != nil {
			return fmt.Errorf("flowLogLevel for flowId (%s) still exists", rs.Primary.ID)
		} else if util.IsStatus404(resp) {
			// Language not found as expected
			continue
		} else {
			// Unexpected error
			return fmt.Errorf("Unexpected error: %s", err)
		}
	}
	// Success. All grammar languages deleted
	return nil
}
