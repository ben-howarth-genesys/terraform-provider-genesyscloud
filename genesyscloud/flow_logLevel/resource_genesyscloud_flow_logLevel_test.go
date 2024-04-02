package flow_logLevel

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mypurecloud/platform-client-sdk-go/v125/platformclientv2"
	"terraform-provider-genesyscloud/genesyscloud/architect_flow"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"testing"
)

func TestAccResourceFlowLogLevel(t *testing.T) {
	var (
		flowResource         = "test_logLevel_flow1"
		resourceId           = "flow_log_level" + uuid.NewString()
		flowName             = "Terraform Test Flow log level " + uuid.NewString()
		flowLoglevelBase     = "Base"
		flowLoglevelAll      = "All"
		flowLogLevelDisabled = "Disabled"
		flowId               = "${genesyscloud_flow." + flowResource + ".id}"
		filePath             = "../../examples/resources/genesyscloud_flow/inboundcall_flow_example.yaml"
		inboundCallConfig    = fmt.Sprintf("inboundCall:\n  name: %s\n  defaultLanguage: en-us\n  startUpRef: ./menus/menu[mainMenu]\n  initialGreeting:\n    tts: Archy says hi!!!\n  menus:\n    - menu:\n        name: Main Menu\n        audio:\n          tts: You are at the Main Menu, press 9 to disconnect.\n        refId: mainMenu\n        choices:\n          - menuDisconnect:\n              name: Disconnect\n              dtmf: digit_9", flowName)
	)

	flowResourceConfig := architect_flow.GenerateFlowResource(
		flowResource,
		filePath,
		inboundCallConfig,
		true,
	)

	var ()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			util.TestAccPreCheck(t)
		},
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				// Create using only flow log level
				Config: flowResourceConfig + generateFlowLogLevelResource(
					flowId,
					flowLoglevelBase,
					resourceId,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_log_level", flowLoglevelBase),
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
			{
				// Create using only flow log level
				Config: flowResourceConfig + generateFlowLogLevelResource(
					flowId,
					flowLoglevelAll,
					resourceId,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_log_level", flowLoglevelAll),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.communications", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.event_error", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.event_other", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.event_warning", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.execution_input_outputs", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.execution_items", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.names", "true"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.variables", "true"),
				),
			},
			{
				// Create using only flow log level
				Config: flowResourceConfig + generateFlowLogLevelResource(
					flowId,
					flowLogLevelDisabled,
					resourceId,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_log_level", flowLogLevelDisabled),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.communications", "false"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.event_error", "false"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.event_other", "false"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.event_warning", "false"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.execution_input_outputs", "false"),
					resource.TestCheckResourceAttr("genesyscloud_flow_loglevel."+resourceId, "flow_characteristics.0.execution_items", "false"),
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