package idp_salesforce

import (
	"github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/provider"
	resourceExporter "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_exporter"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	registrar "github.com/mypurecloud/terraform-provider-genesyscloud/genesyscloud/resource_register"
)

/*
resource_genesycloud_idp_salesforce_schema.go holds four functions within it:

1.  The registration code that registers the Datasource, Resource and Exporter for the package.
2.  The resource schema definitions for the idp_salesforce resource.
3.  The datasource schema definitions for the idp_salesforce datasource.
4.  The resource exporter configuration for the idp_salesforce exporter.
*/
const ResourceType = "genesyscloud_idp_salesforce"

// SetRegistrar registers all of the resources, datasources and exporters in the package
func SetRegistrar(regInstance registrar.Registrar) {
	regInstance.RegisterResource(ResourceType, ResourceIdpSalesforce())
	regInstance.RegisterExporter(ResourceType, IdpSalesforceExporter())
}

// ResourceIdpSalesforce registers the genesyscloud_idp_salesforce resource with Terraform
func ResourceIdpSalesforce() *schema.Resource {
	return &schema.Resource{
		Description: "Genesys Cloud Single Sign-on Salesforce Identity Provider. See this page for detailed configuration instructions: https://help.mypurecloud.com/articles/add-salesforce-as-a-single-sign-on-provider/",

		CreateContext: provider.CreateWithPooledClient(createIdpSalesforce),
		ReadContext:   provider.ReadWithPooledClient(readIdpSalesforce),
		UpdateContext: provider.UpdateWithPooledClient(updateIdpSalesforce),
		DeleteContext: provider.DeleteWithPooledClient(deleteIdpSalesforce),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(8 * time.Minute),
			Read:   schema.DefaultTimeout(8 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			`name`: {
				Description: `Name of the provider`,
				Optional:    true,
				Type:        schema.TypeString,
			},
			"certificates": {
				Description: "PEM or DER encoded public X.509 certificates for SAML signature validation.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"issuer_uri": {
				Description: "Issuer URI provided by Salesforce.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"target_uri": {
				Description: "Target URI provided by Salesforce.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			`slo_uri`: {
				Description: `Provided on app creation.`,
				Optional:    true,
				Type:        schema.TypeString,
			},
			`slo_binding`: {
				Optional: true,
				Type:     schema.TypeString,
			},
			`relying_party_identifier`: {
				Description: `String used to identify Genesys Cloud to Ping.`,
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
			},
			"disabled": {
				Description: "True if Salesforce is disabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func IdpSalesforceExporter() *resourceExporter.ResourceExporter {
	return &resourceExporter.ResourceExporter{
		GetResourcesFunc: provider.GetAllWithPooledClient(getAllIdpSalesforce),
		RefAttrs:         map[string]*resourceExporter.RefAttrSettings{}, // No references
	}
}
