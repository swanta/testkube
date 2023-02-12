---
sidebar_position: 1
sidebar_label: Creating
---

# Testkube Tests

Tests are single executor oriented objects. Test can have different types, which depends on which executors are installed in your cluster.

Testkube includes `postman/collection`, `cypress/project`, `curl/test`, `k6/script` and `soapui/xml` test types which are auto registered during the Testkube install by default.

As Testkube was designed with flexibility in mind, you can add your own executors to handle additional test types.

## **Test Source**

Tests can be currently created from multiple sources:

1. A simple `file` with the test content, For example, with Postman collections, we're exporting the collection as a JSON file. For cURL executors, we're passing a JSON file with the configured cURL command.
2. String - we can also define the content of the test as a string
3. Git directory - we can pass `repository`, `path` and `branch` where our tests are stored. This is used in Cypress executor as Cypress tests are more like npm-based projects which can have a lot of files. We are handling sparse checkouts which are fast even in the case of huge mono-repos.
4. Git file - similarly to Git directories, we can use files located on Git by specifying `git-uri` and `branch`.

Note: not all executors support all input types. Please refer to the individual executors' documentation to see which options are available.

## **Create a Test**

### **Create Your First Test from a File (Postman Collection Test)**

To create your first Postman collection in Testkube, export your collection into a file.

Right click on your collection name:

![create postman collection step 1](../../img/test-create-1.png)

Click the **Export** button:

![create postman collection step 2](../../img/test-create-2.png)

Save in a convenient location. In this example, we are using `~/Downloads/TODO.postman_collection.json` path.

![create postman collection step 3](../../img/test-create-3.png)

Create a Testkube test using the exported JSON and give it a unique and fitting name. For simplicity's sake we used `test` in this example.

```bash
kubectl testkube create test --file ~/Downloads/TODO.postman_collection.json --name test
```

Output:

```bash
Detected test type postman/collection
Test created test 🥇
```

Test created! Now we have a reusable test.

### **Updating Tests**

If you need to update your test after change in Postman, re-export it to a file and run the update command:

```bash
kubectl testkube update test --file ~/Downloads/TODO.postman_collection.json --name test
```

To check if the test was created correctly, look at the `Test Custom Resource` in your Kubernetes cluster:

Output:

```bash
Detected test test type postman/collection
Test updated test 🥇
```

Testkube will override all test settings and content with the `update` method.

### **Checking Test Content**

Let's see what has been created:

```bash
kubectl get tests -ntestkube
```

Output:

```bash
NAME   AGE
test   32s
```

Get the details of a test:

```bash
kubectl get tests -ntestkube test-example -oyaml
```bash
$ kubectl testkube get test test

name: test
type_: postman/collection
content: |-
    {
        "info": {
                "_postman_id": "b40de9fe-9201-4b03-8ca2-3064d9027dd6",
                "name": "TODO",
                "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
        },
        "item": [
                {
                        "name": "Create TODO",
                        "event": [
                                {
                                        "listen": "test",
                                        "script": {
                                                "exec": [
                                                        "pm.test(\"Status code is 201 CREATED\", function () {",
                                                        "    pm.response.to.have.status(201);",
                                                        "});",
                                                        "",
                                                        "",
                                                        "pm.test(\"Check if todo item craeted successfully\", function() {",
                                                        "    var json = pm.response.json();",
                                                        "    pm.environment.set(\"item\", json.url);",
                                                        "    pm.sendRequest(json.url, function (err, response) {",
                                                        "        var json = pm.response.json();",
                                                        "        pm.expect(json.title).to.eq(\"Create video for conference\");",
                                                        "",
                                                        "    });",
                                                        "    console.log(\"creating\", pm.environment.get(\"item\"))",
                                                        "})",
                                                        "",
                                                        ""
                                                ],
                                                "type": "text/javascript"
                                        }
                                },
                                {
                                        "listen": "prerequest",
                                        "script": {
                                                "exec": [
                                                        ""
                                                ],
                                                "type": "text/javascript"
                                        }
                                }
                        ],
                        "protocolProfileBehavior": {
                                "disabledSystemHeaders": {}
                        },
                        "request": {
                                "method": "POST",
                                "header": [
                                        {
                                                "key": "Content-Type",
                                                "value": "application/json",
                                                "type": "text"
                                        }
                                ],
                                "body": {
                                        "mode": "raw",
                                        "raw": "{\"title\":\"Create video for conference\",\"order\":1,\"completed\":false}"
                                },
                                "url": {
                                        "raw": "{{uri}}",
                                        "host": [
                                                "{{uri}}"
                                        ]
                                }
                        },
                        "response": []
                },
                {
                        "name": "Complete TODO item",
                        "event": [
                                {
                                        "listen": "prerequest",
                                        "script": {
                                                "exec": [
                                                        "console.log(\"completing\", pm.environment.get(\"item\"))"
                                                ],
                                                "type": "text/javascript"
                                        }
                                }
                        ],
                        "request": {
                                "method": "PATCH",
                                "header": [
                                        {
                                                "key": "Content-Type",
                                                "value": "application/json",
                                                "type": "text"
                                        }
                                ],
                                "body": {
                                        "mode": "raw",
                                        "raw": "{\"completed\": true}"
                                },
                                "url": {
                                        "raw": "{{item}}",
                                        "host": [
                                                "{{item}}"
                                        ]
                                }
                        },
                        "response": []
                },
                {
                        "name": "Delete TODO item",
                        "event": [
                                {
                                        "listen": "prerequest",
                                        "script": {
                                                "exec": [
                                                        "console.log(\"deleting\", pm.environment.get(\"item\"))"
                                                ],
                                                "type": "text/javatest"
                                        }
                                },
                                {
                                        "listen": "test",
                                        "script": {
                                                "exec": [
                                                        "pm.test(\"Status code is 204 no content\", function () {",
                                                        "    pm.response.to.have.status(204);",
                                                        "});"
                                                ],
                                                "type": "text/javascript"
                                        }
                                }
                        ],
                        "request": {
                                "method": "DELETE",
                                "header": [],
                                "url": {
                                        "raw": "{{item}}",
                                        "host": [
                                                "{{item}}"
                                        ]
                                }
                        },
                        "response": []
                }
        ],
        "event": [
                {
                        "listen": "prerequest",
                        "script": {
                                "type": "text/javascript",
                                "exec": [
                                        ""
                                ]
                        }
                },
                {
                        "listen": "test",
                        "script": {
                                "type": "text/javascript",
                                "exec": [
                                        ""
                                ]
                        }
                }
        ],
        "variable": [
                {
                        "key": "uri",
                        "value": "http://34.74.127.60:8080/todos"
                },
                {
                        "key": "item",
                        "value": null
                }
        ]
    }

```

We can see that the test resource was created with Postman collection JSON content.

You can also check tests with the standard `kubectl` command which will list the tests Custom Resource.

```bash
kubectl get tests -n testkube test -oyaml
```

### **Create a Test from Git**

Some executors can handle files and some can handle only git resources. You'll need to follow the particular executor **readme** file to be aware which test types the executor handles.

Let's assume that a Cypress project is created in a git repository - <https://github.com/kubeshop/testkube-executor-cypress/tree/main/examples> - where **examples** is a test directory in the `https://github.com/kubeshop/testkube-executor-cypress.git` repository.

Now we can create our Cypress-based test as shown below. In git based tests, we need to pass the test type.

```bash
kubectl testkube create test --uri https://github.com/kubeshop/testkube-executor-cypress.git --git-branch main --git-path examples --name kubeshop-cypress --type cypress/project
```

Output:

```bash
Test created kubeshop-cypress 🥇
```

Let's check how the test created by Testkube is defined in the cluster:

```bash
$ kubectl get tests -n testkube kubeshop-cypress -o yaml
apiVersion: tests.testkube.io/v1
kind: Test
metadata:
  creationTimestamp: "2021-11-17T12:29:32Z"
  generation: 1
  name: kubeshop-cypress
  namespace: testkube
  resourceVersion: "225162"
  uid: f0d856aa-04fc-4238-bb4c-156ff82b4741
spec:
  repository:
    branch: main
    path: examples
    type: git
    uri: https://github.com/kubeshop/testkube-executor-cypress.git
  type: cypress/project
```

As we can see, this test has `spec.repository` with git repository data. This data can now be used by the executor to download test data.
#### **Providing Certificates**

If the git repository is using a self-signed certificate, you can provide the certificate using Kubernetes secrets and passing the secret name to the `--git-certificate-secret` flag.

In order to create a secret, use the following command:

```bash
kubectl create secret generic git-cert --from-file=git-cert.crt --from-file=git-cert.key
```
Then you can pass the secret name to the `--git-certificate-secret` flag and, during the test execution, the certificate will be mounted to the test container and added to the trusted authorities.

### **Mapping local files**

Local files can be added into a Testkube Test. This can be set on Test level passing the file in the format `source_path:destination_path` using the flag `--copy-files`. The file will be copied upon execution from the machine running `kubectl`. The files will be then available in the `/data/uploads` folder inside the test container.

```bash
kubectl testkube create test --name maven-example-file-test --git-uri https://github.com/kubeshop/testkube-executor-maven.git --git-path examples/hello-maven-settings --type maven/test  --git-branch main --copy-files "/Users/local_user/local_maven_settings.xml:settings.xml"
```

Output:

```bash
Test created maven-example-file-test 🥇
```

To run this test, refer to `settings.xml` from the `/data/uploads` folder:

```bash
testkube run test maven-example-file-test --args "--settings" --args "/data/uploads/settings.xml" -v "TESTKUBE_MAVEN=true" --args "-e" --args "-X" --env "DEBUG=\"true\""
```

### **Changing the default job template used for test execution**

You can always create your own custom executor with its own job template definition used for test execution. But sometimes you just need to adjust an existing job template of a standard Testkube executor with a few parameters. In this case you can use additional parameter `--job-template` when you create or run the test:

```bash
kubectl testkube create test --git-branch main --git-uri https://github.com/kubeshop/testkube-example-cypress-project.git --git-path "cypress" --name template-test --type cypress/project --job-template job.yaml
```

Where `job.yaml` file contains adjusted job template parts for merging with default job template:

```yaml
apiVersion: batch/v1
kind: Job
spec:
  template:
    spec:
      containers:
      - name: {{ .Name }}
        image: {{ .Image }}
        imagePullPolicy: Always
        command:
          - "/bin/runner"
          - '{{ .Jsn }}'
        volumeMounts:
        - name: data-volume
          mountPath: /data
        resources:
          limits:
            memory: 128Mi
```

When you run such a test you will face a memory limit for the test executor pod, when the default job template doesn't have any resource constraints.

### **Executing pre run script**
If you need to provide additional configuration for your executor environment, you can submit prerun script to be executed before test started. For example, we have such a simple shell script stored in `script.sh` file:

```sh
!/bin/sh

echo "Storing ssl certificate in file from env secret env"
echo "$SSL_CERT" > /data/ssl.crt
```

Then just provide it when you create or run the test using `--prerun-script` parameter:

```bash
kubectl testkube create test --file test/postman/LocalHealth.postman_collection.json --name script-test --type postman/collection --prerun-script script.sh --secret-env SSL_CERT=your-k8s-secret
```

### **Changing the default scraper job template used for container executor tests**

When you use container executor tests generating artifacts for scraping, then we launch 2 sequential kubernetes jobs, one is for test execution and other one is for scraping test results. Sometimes you need to adjust an existing scraper job template of a standard Testkube scraper with a few parameters. In this case you can use additional parameter `--scraper-template` when you create or run the test:

```bash
kubectl testkube create test --name scraper-test --type scraper/test --artifact-storage-class-name standard --artifact-volume-mount-path /share --artifact-dir test/files --scraper-template scraper.yaml
```

Where `scraper.yaml` file contains adjusted scraper job template parts for merging with default scraper job template:

```yaml
apiVersion: batch/v1
kind: Job
spec:
  template:
    spec:
      containers:
      - name: {{ .Name }}-scraper
        image: {{ .ScraperImage }}
        imagePullPolicy: Always
        command:
          - "/bin/runner"
          - '{{ .Jsn }}'
        {{- if .ArtifactRequest }}
        volumeMounts:
          - name: artifact-volume
            mountPath: {{ .ArtifactRequest.VolumeMountPath }}
        {{- end }}
        resources:
          limits:
            memory: 512Mi
```

When you run such a test you will face a memory limit for the scraper pod, when the default scraper job template doesn't have any resource constraints.

## **Summary**

Tests are the main smallest abstractions over test suites in Testkube, they can be created with different sources and used by executors to run on top of a particular test framework.
