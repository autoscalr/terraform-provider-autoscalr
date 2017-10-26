provider "aws" {
  region = "us-east-1"
}
provider "autoscalr" {
  // You either need to specify the api_key here or via the AUTOSCLAR_API_KEY enviroment variable
  //api_key = "yourKey"
}

resource "aws_launch_configuration" "test_lc" {
  name_prefix   = "test-lc-"
  image_id      = "ami-8c1be5f6" // Base Amazon Linux AMI in us-east-1
  instance_type = "t1.micro"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "myAppASG" {
  availability_zones          = ["us-east-1a", "us-east-1b","us-east-1c"]
  name                        = "myAppASG"
  max_size                    = 2
  min_size                    = 0
  desired_capacity            = 0
  health_check_grace_period   = 300
  health_check_type           = "EC2"
  force_delete                = true
  launch_configuration        = "${aws_launch_configuration.test_lc.name}"
  lifecycle {
    create_before_destroy     = true
  }
  suspended_processes         = ["AZRebalance"] // Recommended to keep ASG from fighting AutoScalr AZ Rebalancing
}

resource "autoscalr_autoscaling_group" "asr4myAppASG" {
  aws_region                  = "us-east-1"
  aws_autoscaling_group_name  = "${aws_autoscaling_group.myAppASG.name}"
  instance_types              = ["c3.large","c3.xlarge"]
  display_name                = "myFirstAutoScalrApp"
  max_spot_percent_total      = 85
  max_spot_percent_one_market = 25
}
