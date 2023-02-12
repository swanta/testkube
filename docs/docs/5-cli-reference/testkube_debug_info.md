## testkube debug info

Show debug info

### Synopsis

Get all the necessary information to debug an issue in Testkube

```
testkube debug info [flags]
```

### Options

```
  -h, --help   help for info
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

* [testkube debug](testkube_debug.md)	 - Print environment information for debugging

