package autoscalr

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"fmt"
)

type Config struct {
	AccessKey string
	apiUrl    string
}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"autoscalr_autoscaling_group": resourceAutoScalrAutoscalingGroup(),
		},
		/*
			DataSourcesMap: map[string]*schema.Resource{
				"null_data_source": dataSource(),
			},
		*/
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	api_key := d.Get("api_key").(string)
	if api_key == "" {
		api_key = os.Getenv("AUTOSCALR_API_KEY")
	}
	if (len(api_key) < 1) {
		err := fmt.Errorf("AUTOSCALR_API_KEY environment variable must be set or api_key must be set in provider.")
		return nil, err
	}
	config := Config{
		apiUrl:    "https://app.autoscalr.com/api/autoScalrApp",
		AccessKey: api_key,
	}

	return &config, nil
}
