Configurable Server Response
============================

This is a simple webservice that have configurable parameters for 500 errors
and latency to allow manipulating the values needed to test things like metrics.

This was initially built to test metrics for [GoogleCloudPlatform/cloud-run-release-operator](https://github.com/GoogleCloudPlatform/cloud-run-release-operator).

Environment variables
--------------------

The following environment variables are used to determine the webservice's 
behavior:

 * `PERCENT_500_RESPONSES`: Percentage of requests that would received an HTTP
   500 error
 * `LATENCY_P99`: The minimum latency for the 99th percentile (in ms).
 * `LATENCY_P95`: The minimum latency for the 95th percentile (in ms).
 * `LATENCY_P50`: The minimum latency for the 50th percentile (in ms).
