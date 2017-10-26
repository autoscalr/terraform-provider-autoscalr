package autoscalr

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAutoScalr_autoscaling_group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAsrAsgDestroyed,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAsrConfig,
				Check: resource.ComposeTestCheckFunc(
					testAsrAsgExists("autoscalr_autoscaling_group.myApp"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "aws_autoscaling_group_name", "testASG"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "aws_region", "us-east-1"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "scale_mode", "cpu"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "autoscalr_enabled", "false"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "target_spare_cpu_percent", "20"),
				),
			},
			resource.TestStep{
				Config: testAsrConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAsrAsgExists("autoscalr_autoscaling_group.myApp"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "aws_autoscaling_group_name", "testASG"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "aws_region", "us-east-1"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "scale_mode", "cpu"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "autoscalr_enabled", "false"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "target_spare_cpu_percent", "30"),
					resource.TestCheckResourceAttr("autoscalr_autoscaling_group.myApp", "max_spot_percent_one_market", "35"),
				),
			},
		},
	})
}

const testAsrConfig = `
resource "autoscalr_autoscaling_group" "myApp" {
  aws_region = "us-east-1"
  aws_autoscaling_group_name = "testASG"
  display_name = "testASG"
  instance_types = ["t1.micro", "m1.medium"]
  scale_mode = "cpu"
  autoscalr_enabled = false
  target_spare_cpu_percent = 20
}
`
const testAsrConfigUpdated = `
resource "autoscalr_autoscaling_group" "myApp" {
  aws_region = "us-east-1"
  aws_autoscaling_group_name = "testASG"
  display_name = "testASG"
  instance_types = ["t1.micro", "m1.medium"]
  scale_mode = "cpu"
  autoscalr_enabled = false
  target_spare_cpu_percent = 30
  max_spot_percent_one_market = 35
}
`

func testAsrAsgDestroyed(s *terraform.State) error {
	//for _, r := range s.RootModule().Resources {
	//	fmt.Println(r.Primary.ID)
	//	return fmt.Errorf("Unexpected resource found in terraform state: %s", r.Primary.ID)
	//}
	return nil
}

func testAsrAsgExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found in terraform state: %s", n)
		}
		return nil
	}
}
