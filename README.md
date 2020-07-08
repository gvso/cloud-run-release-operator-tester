Cloud Run Release Operator Testing App
======================================

This is a simple webservice that have configurable parameters for 500 errors
and latency to allow manipulating the values needed to test the metrics
component of [GoogleCloudPlatform/cloud-run-release-operator](https://github.com/GoogleCloudPlatform/cloud-run-release-operator).

Environment variables
--------------------

The following environment variables are used to determine the webservice's 
behavior:

 * `PERCENT_500_RESPONSES`: Percentage of requests that would received an HTTP
   500 error
 * `LATENCY_TRESHOLD`: The minimum latency that is considered above the
   expected limits of the webservice.
 * `PERCENT_OVER_LATENCY_TRESHOLD`: Percentage of requests that would receive a
   latency bigger than the treshold.