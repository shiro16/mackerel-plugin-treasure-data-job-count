# mackerel-plugin-treasure-data-job-count

Treasure Data custom metrics plugin for mackerel.io agent.

## Synopsis

```
mackerel-plugin-treasure-data-job-count [-treasure-data-api-key=<API KEY>] [-from=<FROM>] [-to=<TO>] [-status=<STATUS>]
```

## Example Of mackerel-agent.conf

```
[plugin.metrics.treasure-data-job-count]
command = "/path/to/mackerel-plugin-treasure-data-job-count -treasure-data-api-key=API_KEY -from=0 -to=99 -status=running"
```
