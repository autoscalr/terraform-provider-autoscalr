provider "aws" {
  region = "us-east-1"
}
provider "autoscalr" {
  // You either need to specify the api_key here or via the AUTOSCLAR_API_KEY enviroment variable
  //api_key = "yourKey"
}

resource "aws_launch_configuration" "test_lc" {
  name_prefix   = "test-lc-"
  image_id      = "ami-9700e3fc"
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
  suspended_processes         = ["AZRebalance"]
}

resource "autoscalr_autoscaling_group" "asr4myAppASG" {
  aws_region                        = "us-east-1"                               // Required
  aws_autoscaling_group_name        = "${aws_autoscaling_group.myAppASG.name}"  // Required
  instance_types                    = ["c3.large","c3.xlarge"]                  // Required
  display_name                      = "myAppName"                               // Short name displayed in AutoScalr web UI
  scale_mode                        = "cpu"                                     // Values: cpu, queue, ecs. Default cpu
  target_spare_cpu_percent          = 20                                        // Target spare cpu percentage to scale to
  max_spot_percent_total            = 80                                        // Maximum percentage of capacity to allow in Spot instances
  max_spot_percent_one_market       = 20                                        // Maximum percentage of capacity to allow in a single Spot market
  detailed_monitoring_enabled       = true                                      // Enables AWS per minute metrics which improves scaling decisions
  os_family                         = "Linux/UNIX"                              // Values: Linux/Unix, SUSE Linux, Windows
  max_hours_instance_age            = 24                                        // When set, AutoScalr will schedule instance replacement if age exceeds this setting
  autoscalr_enabled                 = true                                      // Flag to quickly allow disabling AutoScalr actions temporarily if desired

  //For scale_mode = ecs
  cluster_name                      = "anEcsCluster"                            // If ASG is supporting an ECS cluster, putting the cluster name here turns on additional optimizations
                                                                                // Should use substitution reference to ECS resource to establish dependency if also created by Terraform
  target_spare_memory_percent       = 20                                        // Target spare memory percentage to scale to.  Only applicable with ecs scale_mode.

  // For scale_mode = queue
  queue_name                        = "anSqsQueue"                              // SQS Queue name to use for scaling input. Only applicable with queue scale_mode.
                                                                                // Should use substitution reference to SQS resource to establish dependency if also created by Terraform
  target_queue_size                 = 1000                                      // Target queue size to scale to. Should be non-zero for efficient scaling.
  max_minutes_to_target_queue_size  = 60                                        // Number of minutes to return to target queue size.
                                                                                // Lower values will trigger more aggressive scaling up to reach target_queue_size faster.
}
