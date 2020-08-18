# Homework with Golang

This time the homework will require you to pick up Go (or Golang) and learn the basics in order to complete the homework assignment.

The task seems simple on the surface: write a service with an HTTP API for measuring HTTP response latency for various hosts. i.e. how long does it take to get a response from https://15min.lt after
sending HTTP GET request to this host. 

## Requirements
Below are the requirements for this task

### General

* You need to fork this repository. :rotating_light: Remember the number of forks before clicking a button, you will need that number later
* The service (http server) needs to be written in Go
* Any package or framework can be used, but I suggest to just use Standard Go packages like `net/http`, `time`, `json`

### HTTP API

* Your service should have an HTTP API with at least one endpoint: `/measure`
* `/measure` endpoint query parameters (`/measure?host=reddit.com&protocol=https&samples=3`)
   * `host` - specifies which host should be targeted, i.e. centric.eu, reddit.com, amazon.com
   * `protocol` - allowed values are: _http_ and _https_
   * `samples` - how many samples of response time to measure, i.e. _5_ means that 5 requests should be sent to the provided host and every response latency measured
* response should be a valid JSON object, with host, protocol and results keys on root object, example:
  ```json
  {
      "host": "reddit.com",
      "protocol": "https",
      "results": {
          "measurements": [
              "1099ms",
              "1052ms",
              "1303ms"
          ],
          "averageLatency": "1151ms"
      }
  }
  ```

### Things to consider

* error handling
  * provided host is not reachable
  * request timeouts for unreachable or unresponsive hosts
* the deployment (where it will run?)
  * your own laptop (perfectly ok)
  * Azure
  * AWS
  * Kubernetes
* concurency - perhaps it will take less time to collect all samples at the same time and aggregate the results afterwards?
* automatic time units. It might make sense to use automatic time units, i.e. show milliseconds when the latency is under 5 seconds and show seconds when the latency is higher


### Bonus points

* CI/CD pipeline in gitlab
* Service is deployed in either:
  * Azure (VM)
  * AWS (EC2)
  * AWS Lambda (the whole service is written as lambda function)
  * Kubernetes
* Automated tests with at least 50% code coverage (`go test -cover ./...`)

### :rotating_light: Deployment on provided Azure VM :rotating_light:
When deploying your application to the provided Azure VM, please use the following formula to calculate the TCP port for your API:

`PORT = 8080 + your github fork number`

### Solution

The homework is deployed on the ainize.ai platform and is accessible [here](https://ainize.ai/Ben10k/go-homework)

The link to measure the latency to https://15min.lt: https://master-go-homework-ben10k.endpoint.ainize.ai/measure?host=15min.lt&protocol=https&samples=10
