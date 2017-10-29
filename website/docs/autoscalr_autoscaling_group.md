# AutoScalr Terraform Provider
# resource autoscalr_autoscaling_group

A resource that provides additional cluster management capabilities designed to reduce operational costs to existing Auto Scaling Group resource.

## Example Usage

```hcl
# Extend the MyASG resource with AutoScalr
resource "autoscalr_autoscaling_group" "asrEnableMyASG" {
  aws_region                  = "us-east-1"
  aws_autoscaling_group_name  = "${aws_autoscaling_group.MyASG.name}"
  instance_types              = ["c3.large","c3.xlarge"]
  display_name                = "myFirstAutoScalrApp"
  max_spot_percent_total      = 85
  max_spot_percent_one_market = 20
}
```

## Argument Reference

The following arguments are supported:

* `aws_autoscaling_group_name` - (Required) Name of AWS autoscaling group (ASG) to extend
* `display_name` - (Optional) Short name to be used in AutoScalr web UI display
* `aws_region` - (Required) AWS Region the autoscaling group to extend is in
* `instance_types` - (Required) List of instance types to use
* `scale_mode` - (Optional, Default: cpu) Options: cpu, queue, ecs
* `target_spare_cpu_percent` - (Optional, Default: 20) Target spare cpu percentage to scale to, e.g. 20% spare capacity = 80% cpu utilization
* `max_spot_percent_total` - (Optional, Default: 80) Maximum percentage of capacity to allow in Spot instances
* `max_spot_percent_one_market` - (Optional, Default: 20) Maximum percentage of capacity to allow in a single Spot market
* `detailed_monitoring_enabled` - (Optional, Default: true) Enables AWS per minute metrics which improves scaling decisions
* `os_family` - (Optional, Default: Linux/UNIX) Options: Linux/Unix, SUSE Linux, Windows
* `max_hours_instance_age` - (Optional, Default: off) When set, AutoScalr will schedule instance replacement if an instance's age exceeds this setting
* `autoscalr_enabled` - (Optional, Default: true) Flag to quickly allow disabling AutoScalr actions temporarily if desired

The following arguments are also supported when scale_mode is ecs:

* `cluster_name` - (Required) The name of the ECS cluster the target ASG is associated with
* `target_spare_memory_percent` - (Optional, Default: off) Target spare memory percentage to scale to in addition to target_spare_cpu_percent
* `aws_region` - (Required) AWS Region the autoscaling group to extend is in

The following arguments are also supported when scale_mode is queue:

* `queue_name` - (Required) SQS Queue name to use for scaling input
* `target_queue_size` - (Optional, Default: 1000) Target queue size to scale to. Should be non-zero for efficient scaling
* `max_minutes_to_target_queue_size` - (Optional, Default: 1000) Number of minutes to return to target queue size.  Lower values will trigger more aggressive scaling up to reach size faster.

The following arguments are also supported when scale_mode is fixed:

* `target_capacity` - (Required) The fixed number of vCpus AutoScalr should maintain

