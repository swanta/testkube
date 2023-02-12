## testkube delete

Delete resources

```
testkube delete <resourceName> [flags]
```

### Options

```
  -c, --client string   Client used for connecting to testkube API one of proxy|direct (default "proxy")
  -h, --help            help for delete
      --verbose         should I show additional debug messages
```

### Options inherited from parent commands

```
  -a, --api-uri string     api uri, default value read from config if set (default "http://localhost:8088")
      --namespace string   Kubernetes namespace, default value read from config if set (default "testkube")
      --oauth-enabled      enable oauth (default true)
```

### SEE ALSO

* [testkube](testkube.md)	 - Testkube entrypoint for kubectl plugin
* [testkube delete executor](testkube_delete_executor.md)	 - Delete Executor
* [testkube delete test](testkube_delete_test.md)	 - Delete Test
* [testkube delete testsource](testkube_delete_testsource.md)	 - Delete test source
* [testkube delete testsuite](testkube_delete_testsuite.md)	 - Delete test suite
* [testkube delete webhook](testkube_delete_webhook.md)	 - Delete webhook

