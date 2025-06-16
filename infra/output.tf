output "asg_name" {
  value = aws_autoscaling_group.web_asg.name
}

output "asg_instance_ids" {
  description = "The EC2 instance IDs currently part of the web ASG"
  value       = data.aws_instances.web_asg.ids
}
