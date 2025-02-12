# pushgateway-admin
[Prometheus Pushgateway](https://github.com/prometheus/pushgateway) admin tool.

## Features
```
Prometheus Pushgateway admin tool

Usage:
  pushgateway-admin [command]

Available Commands:
  clean       Clean stale metrics
  delete      Delete all metrics in group
  help        Help about any command
  metrics     Returns all metrics
  post        Send metrics to group via POST method
  put         Send metrics to group via PUT method
  status      Returns status
  wipe        Delete all metrics in all groups

Flags:
  -h, --help           help for pushgateway-admin
  -p, --pretty         Enable pretty print
  -t, --timeout uint   Pushgateway timeout (default 5)
  -u, --url string     Pushgateway URL
```

### Clean
```
Clean stale metrics

Usage:
  pushgateway-admin clean [flags]

Flags:
  -n, --dry-run         Enable dry-run mode
  -g, --group strings   Grouping key (label=value,...)
  -s, --sleep uint      Sleep before exit (sec)
  -l, --stale uint      Stale metrics threshold (sec) (default 300)
```

### Delete
```
Delete all metrics in group

Usage:
  pushgateway-admin delete [flags]

Flags:
  -g, --group strings   Grouping key (label=value,...)
```

### Metrics
```
Returns all metrics

Usage:
  pushgateway-admin metrics [flags]
```

### Post
```
POST works exactly like the PUT but only metrics with the same name as the newly pushed metrics are replaced.

Usage:
  pushgateway-admin post [flags]

Flags:
  -g, --group strings        Grouping key (label=value,...)
  -h, --help                 help for post
  -m, --metric string        Metric
  -f, --metric-file string   Read metrics from file
```

### Put
```
PUT is used to push a group of metrics. All metrics with the grouping key are replaced by the metrics pushed with put.

Usage:
  pushgateway-admin put [flags]

Flags:
  -g, --group strings        Grouping key (label=value,...)
  -h, --help                 help for put
  -m, --metric string        Metric
  -f, --metric-file string   Read metrics from file
```

### Status
```
Returns status

Usage:
  pushgateway-admin status [flags]
```

### Wipe
```
Delete all metrics in all groups

Usage:
  pushgateway-admin wipe [flags]
```

## Demo
```
$ ./make.sh start

$ docker run --rm -it --entrypoint=sh --network=pushgateway_default lystor/pushgateway-admin

$ pushgateway-admin status --url http://pushgateway:9091 --pretty

$ cat << EOF > /tmp/metrics.txt
# TYPE m1 counter
m1{l1="v1",l2="v2"} 1
# TYPE m2 gauge
# HELP m2 Metric m2 description
m2{l1="v1"} 2
m3 3
EOF

$ pushgateway-admin put --url http://pushgateway:9091 --group=job=j1,instance=i1 --metric-file /tmp/metrics.txt
Request:  PUT /metrics/job/j1/instance/i1
Response: status_code=202 duration=851.295µs

$ pushgateway-admin post --url http://pushgateway:9091 --group=job=j1,instance=i2 --metric="m4 4"
Request:  POST /metrics/job/j1/instance/i1
Response: status_code=202 duration=843.305µs

$ echo "m5 5" | pushgateway-admin post --url http://pushgateway:9091 --group=job=j1,instance=i2
Request:  POST /metrics/job/j1/instance/i1
Response: status_code=202 duration=793.901µs

$ cat << EOF | pushgateway-admin post --url http://pushgateway:9091 --group=job=j1,instance=i2
# TYPE m6 counter
m6{l1="v1"} 6
# TYPE m7 gauge
# HELP m7 Metric m7 description
m7 7
EOF

$ pushgateway-admin metrics --url http://pushgateway:9091 --pretty

$ pushgateway-admin delete --url http://pushgateway:9091 --group=job=j1,instance=i2
Request:  DELETE /metrics/job/j1/instance/i2
Response: status_code=202 duration=808.459µs

$ pushgateway-admin clean --url http://pushgateway:9091 --group=job=j1 --stale=60 --dry-run

$ pushgateway-admin wipe --url http://pushgateway:9091
Request:  PUT /api/v1/admin/wipe
Response: status_code=202 duration=1.397993ms

$ ./make.sh stop
```

## Build
```
$ ./make.sh build
```

## Author
* [Mykola Ulianytskyi](https://github.com/lystor)

## License
* Apache License 2.0