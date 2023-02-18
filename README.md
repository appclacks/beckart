# Beckart

Beckart is a tool allowing you to define and execute HTTP test scenarios. You can use it for example to test HTTP APIs.

In Beckart, each step can reuse information gathered from previous ones. For exampla, extract values from body responses and HTTP headers and reuse them later in another step.

Beckart also supports`transformers` which are functions applied on HTTP requests to modify them before being sent.
Only one transformer is supported for now as a proof of concept but new ones will soon be released.

## Example

We will in this example write a test scenario that will:

- Send a GET request to `https://jsonplaceholder.typicode.com/posts`. This API returns a list of posts (you can try it by [following this link](jsonplaceholder.typicode.com/posts).
We will then extract the first post title and set it in a variable named `post-title`, and extract the first post id in a variable named `post-id`.
- The next step will log the post title
- The next step will send a PATCH request to `https://jsonplaceholder.typicode.com/posts/<post-id>` with a JSON body, to simulate an update. The title will be set the value associated to the `new-title` variable defined in the configuration file.
- The next step will send a DELETE request to `https://jsonplaceholder.typicode.com/posts/<post-id>`, by reusing the variable set before.
We will extract the full response body and set it in the `delete-body` variable
- The next step will log the previous response body.

Beckart uses the [Golang template](https://pkg.go.dev/text/template) syntax. Requests body, headers and URLs support templating.

This is the configuration file for our example:

```yaml
variables:
  new-title: "new post title"

actions:
  - name: get-post
    http:
      valid-status:
        - 200
      timeout: "5s"
      target: jsonplaceholder.typicode.com
      method: GET
      port: 443
      protocol: https
      path: "/posts"
      extractors:
        body-json:
          post-title:
            - 0
            - "title"
          post-id:
            - 0
            - "id"

  - name: log-title
    log:
      message: 'title = {{ index .Variables "post-title" }}'

  - name: patch-post
    http:
      valid-status:
        - 200
      timeout: "5s"
      target: jsonplaceholder.typicode.com
      method: PATCH
      headers:
        content-type: "application/json; charset=UTF-8"
      body: |
        {"title": "{{ index .Variables "new-title"}}"}

      port: 443
      protocol: https
      path: '/posts/{{ index .Variables "post-id"}}'

  - name: delete-post
    http:
      valid-status:
        - 200
      timeout: "5s"
      target: jsonplaceholder.typicode.com
      method: DELETE
      port: 443
      protocol: https
      path: '/posts/{{ index .Variables "post-id"}}'
      extractors:
        body: "delete-body"

  - name: log-delete-body
    log:
      message: 'body = {{ index .Variables "delete-body" }}'
```

Run `beckart run --config config.yaml` and you should see beckart being executed successfully:

```json
{"level":"info","ts":1676062923.6442242,"caller":"runner/run.go:14","msg":"starting test scenario"}
{"level":"info","ts":1676062923.6442842,"caller":"runner/run.go:16","msg":"start executing action","action":"get-post"}
{"level":"info","ts":1676062923.7279494,"caller":"runner/run.go:30","msg":"successfully executed action","action":"get-post"}
{"level":"info","ts":1676062923.7280793,"caller":"runner/run.go:16","msg":"start executing action","action":"log-title"}
{"level":"info","ts":1676062923.728326,"caller":"runner/logger.go:15","msg":"title = sunt aut facere repellat provident occaecati excepturi optio reprehenderit"}
{"level":"info","ts":1676062923.7283976,"caller":"runner/run.go:30","msg":"successfully executed action","action":"log-title"}
{"level":"info","ts":1676062923.7284405,"caller":"runner/run.go:16","msg":"start executing action","action":"patch-post"}
{"level":"info","ts":1676062924.2569144,"caller":"runner/run.go:30","msg":"successfully executed action","action":"patch-post"}
{"level":"info","ts":1676062924.257055,"caller":"runner/run.go:16","msg":"start executing action","action":"delete-post"}
{"level":"info","ts":1676062924.5637982,"caller":"runner/run.go:30","msg":"successfully executed action","action":"delete-post"}
{"level":"info","ts":1676062924.5638957,"caller":"runner/run.go:16","msg":"start executing action","action":"log-delete-body"}
{"level":"info","ts":1676062924.564099,"caller":"runner/logger.go:15","msg":"body = {}"}
{"level":"info","ts":1676062924.5641642,"caller":"runner/run.go:30","msg":"successfully executed action","action":"log-delete-body"}
{"level":"info","ts":1676062924.5642056,"caller":"runner/run.go:32","msg":"test scenario finished successfully"}
```

Beckart will fail if the HTTP call timeouts, or if the HTTP endpoint returns an expected status code (not present in `valid-status`).

## Transformers

Transfoerms are functions that can be applied to HTTP requests before sending them.
They're still a work in progress, only one transformer exists for now as a proof of concept

### Exoscale transformer

This transformer will sign the HTTP request as described in the Exoscale [API documentation)(https://openapi-v2.exoscale.com/):

```yaml
transformers:
  ## transformers definitions
  exo:
    exoscale:
      api-key: <exoscale API key>
      api-secret: <exoscale API secret>
actions:
  - name: a1
    http:
      valid-status:
        - 200
      timeout: "5s"
      target: api-ch-gva-2.exoscale.com
      method: GET
      port: 443
      protocol: https
      path: "/v2/event"
    ## list of transformers to apply
    transformers:
      - exo
```
