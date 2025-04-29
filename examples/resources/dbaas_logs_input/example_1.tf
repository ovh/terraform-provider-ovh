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
