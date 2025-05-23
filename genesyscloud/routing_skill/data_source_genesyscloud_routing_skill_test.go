package routing_skill

import (
	"fmt"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/util"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRoutingSkill(t *testing.T) {
	var (
		skillResourceLabel   = "routing-skill"
		skillDataSourceLabel = "routing-skill-data"
		skillName            = "Terraform Skill-" + uuid.NewString()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { util.TestAccPreCheck(t) },
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				Config: GenerateRoutingSkillResource(
					skillResourceLabel,
					skillName,
				),
				Check: resource.ComposeTestCheckFunc(
					waitSeconds(3 * time.Second),
				),
			},
			{
				Config: GenerateRoutingSkillResource(
					skillResourceLabel,
					skillName,
				) + generateRoutingSkillDataSource(skillDataSourceLabel, "genesyscloud_routing_skill."+skillResourceLabel+".name", "genesyscloud_routing_skill."+skillResourceLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.genesyscloud_routing_skill."+skillDataSourceLabel, "id", "genesyscloud_routing_skill."+skillResourceLabel, "id"),
				),
			},
		},
	})
}

func waitSeconds(duration time.Duration) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("Sleeping for %v", duration)
		time.Sleep(duration)
		return nil
	}
}

func generateRoutingSkillDataSource(
	resourceLabel string,
	name string,
	// Must explicitly use depends_on in terraform v0.13 when a data source references a resource
	// Fixed in v0.14 https://github.com/hashicorp/terraform/pull/26284
	dependsOnResource string) string {
	return fmt.Sprintf(`data "genesyscloud_routing_skill" "%s" {
		name = %s
        depends_on=[%s]
	}
	`, resourceLabel, name, dependsOnResource)
}
