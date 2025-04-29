---
subcategory : "Logs Data Platform"
---

# ovh_dbaas_logs_input

Creates a dbaas logs input.

## Example Usage

```terraform
data "ovh_dbaas_logs_input_engine" "logstash" {
  name          = "logstash"
  version       = "9.x"
}

resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
  service_name = "...."
  title        = "my stream"
  description  = "my graylog stream"
}

resource "ovh_dbaas_logs_input" "input" {
  service_name     = ovh_dbaas_logs_output_graylog_stream.stream.service_name
  description      = ovh_dbaas_logs_output_graylog_stream.stream.description
  title            = ovh_dbaas_logs_output_graylog_stream.stream.title
  engine_id        = data.ovh_dbaas_logs_input_engine.logstash.id
  stream_id        = ovh_dbaas_logs_output_graylog_stream.stream.id

  allowed_networks = ["10.0.0.0/16"]
  exposed_port     = "6154"
  nb_instance      = 2

  configuration {
    logstash {
      input_section = <<EOF
  beats {
    port => 6514
    ssl => true
    ssl_certificate => "/etc/ssl/private/server.crt"
    ssl_key => "/etc/ssl/private/server.key"
  }
  EOF

    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `allowed_networks` - List of IP blocks
* `configuration` - (Required) Input configuration
  * `flowgger` - (Optional) Flowgger configuration
    * `log_format` - Type of format to decode. One of "RFC5424", "LTSV", "GELF", "CAPNP"
    * `log_framing` - Indicates how messages are delimited. One of "LINE", "NUL", "SYSLEN", "CAPNP"
  * `logstash` - (Optional) Logstash configuration
    * `filter_section` - (Optional) The filter section of logstash.conf
    * `input_section` - (Required) The filter section of logstash.conf
    * `pattern_section` - (Optional) The list of customs Grok patterns
* `description` - (Required) Input description
* `engine_id` - (Required) Input engine ID
* `exposed_port` - Port
* `nb_instance` - Number of instance running (input, mutually exclusive with parameter `autoscale`)
* `autoscale` - Whether the workload is auto-scaled (mutually exclusive with parameter `nb_instance`)
* `min_scale_instance` - Minimum number of instances in auto-scaled mode
* `max_scale_instance` - Maximum number of instances in auto-scaled mode
* `service_name` - (Required) service name
* `stream_id` - (Required) Associated Graylog stream
* `title` - (Required) Input title

## Attributes Reference

Id is set to the input Id. In addition, the following attributes are exported:

* `created_at` - Input creation
* `hostname` - Hostname
* `input_id` - Input ID
* `is_restart_required` - Indicate if input need to be restarted
* `public_address` - Input IP address
* `ssl_certificate` - Input SSL certificate
* `status` - init: configuration required, pending: ready to start, running: available
* `updated_at` - Input last update
* `current_nb_instance` - Number of instance running (returned by the API)
