# sansay exporter


The sansay exporter allows sansay probing of endpoints over
HTTP, HTTPS.  Much of this code was taken from [snmp_exporter](https://github.com/prometheus/snmp_exporter).

## Running this software

### From binaries

Download the most suitable binary from [the releases tab](https://github.com/ringsq/sansay_exporter/releases)

Then:

    ./sansay_exporter <flags>


### Using the docker image

*Note: You may want to [enable ipv6 in your docker configuration](https://docs.docker.com/v17.09/engine/userguide/networking/default_network/ipv6/)*

    docker run --rm -d -p 9116:9116 --name sansay_exporter

### Checking the results

Visiting [http://localhost:9116/sansay?target=localhost:8888](http://localhost:9116/sansay?target=localhost:8888&username=user&password=password)
will return metrics against localhost:8888.

## Building the software

### Local Build

> **IMPORTANT NOTE:** For the SOAP interface to work correctly took a modification to the gowsdl library, which was made in the vendored version.  
> Therefore, to build and run successfully you *MUST* use the vendored version:
>

```
$ go build -mod=vendor -ldflags '-X main.Version=x.x'
```

### Building with Docker

After a successful local build:

    docker build  -t sansay_exporter .

## Configuration

sansay exporter is configured via command-line flags (such as what port to listen on, and the logging format and level).

To view all available command-line flags, run `./sansay_exporter -h`.

The timeout of each probe is automatically determined from the `scrape_timeout` in the [Prometheus config](https://prometheus.io/docs/operating/configuration/#configuration-file), slightly reduced to allow for network delays.
If not specified, it defaults to 10 seconds.

## Prometheus Configuration

The sansay exporter needs to be passed the target as a parameter, this can be
done with relabelling.

Example config:
```yml
scrape_configs:
  - job_name: 'sansay'
    metrics_path: /sansay
    params:
      username: ['user']
      password: ['password']
    static_configs:
      - targets:
        - localhost:8888    
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9116  # The sansay exporter's real hostname:port.
```
