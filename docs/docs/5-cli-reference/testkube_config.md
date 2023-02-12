## testkube config

Set feature configuration value

```
testkube config <feature> <value> [flags]
```

### Options

```
  -h, --help   help for config
```

### Options inherited from parent commands

```
  -a, --api-uri string     api uri, default value read from config if set (default "http://localhost:8088")
  -c, --client string      client used for connecting to Testkube API one of proxy|direct (default "proxy")
      --namespace string   Kubernetes namespace, default value read from config if set (default "testkube")
      --oauth-enabled      enable oauth (default true)
      --verbose            show additional debug messages
```

### SEE ALSO

* [testkube](testkube.md)	 - Testkube entrypoint for kubectl plugin
* [testkube config api-uri](testkube_config_api-uri.md)	 - Set api uri for testkube client
* [testkube config namespace](testkube_config_namespace.md)	 - Set namespace for testkube client
* [testkube config oauth](testkube_config_oauth.md)	 - Set oauth credentials for api uri in testkube client

