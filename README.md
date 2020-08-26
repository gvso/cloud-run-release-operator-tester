# Configurable Server Response

This is a simple webservice that have configurable parameters for 500 errors
and latency to allow manipulating the values needed to test things like metrics.

This was initially built to test metrics for
[GoogleCloudPlatform/cloud-run-release-manager](https://github.com/GoogleCloudPlatform/cloud-run-release-manager).

## Environment variables

The following environment variables are used to determine the webservice's
behavior:

* `PERCENT_500_RESPONSES`: Percentage of requests that would received an HTTP
  500 error
* `LATENCY_P99`: The minimum latency for the 99th percentile (in ms).
* `LATENCY_P95`: The minimum latency for the 95th percentile (in ms).
* `LATENCY_P50`: The minimum latency for the 50th percentile (in ms).

## Turning environment variables on/off

The above environment variables can be turned off at the beginning to test a
healthy service, and then turned on to start the behavior set with the variables
in the above section.

For this, this web service relies on the [Runtime
Configurator](https://cloud.google.com/deployment-manager/runtime-configurator):

1. [Create](https://cloud.google.com/deployment-manager/runtime-configurator/create-and-delete-runtimeconfig-resources#gcloud)
   Runtime Configurator resource:
  
    ```sh
    gcloud beta runtime-config configs create [CONFIG_NAME]
    ```

2. Set a variable on whether the behavior set the by the environment variables
   above should be followed:

    ```sh
    gcloud beta runtime-config configs variables set respect-variables \
    "no" --config-name [CONFIG_NAME] --is-text
    ```

3. Update this webservice's environment variable, so it knows which runtime
   configurator resource and variables to use.

    ```sh
    gcloud run services update [SERVICE] \
      --update-env-vars PROJECT=<YOUR_PROJECT>,CONFIG_NAME=<CONFIG_NAME>
    ```

4. Starts sending requests to your service (e.g. using [hey](https://github.com/rakyll/hey)).

5. Once you want your service to respect the environment variables set in the
   above section, update the Runtime Configurator resource's variable to `yes`:

    ```sh
      gcloud beta runtime-config configs variables set respect-variables \
        "yes" --config-name [CONFIG_NAME] --is-text
    ```
