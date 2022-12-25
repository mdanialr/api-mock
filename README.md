# API Mocker
Simple API mock that will help to serve JSON response that's been defined in config (_yaml_) file

# How to Use
1. Download the binary from [GitHub Releases](https://github.com/mdanialr/api-mock/releases)
2. Extract to dir `api-mock/`
    ```bash
    mkdir api-mock
    tar -C api-mock -xzf api-mock_version_os_arch.tar.gz
    ```
3. Create app config file.
    ```bash
    cp app.yml.example app.yml
    ```
   Above app config will read json response in current directory with `response.yml` as the config filename.
4. Create yaml file to define the endpoint along with their respective response
    ```bash
    cp response.yml.example response.yml
    ```
5. Run
    ```bash
    ./api-mock
    ```
6. Hit endpoint `/api/v1/test` will return status code `200` and json response like
    ```json
    {
      "messages": "data was found",
      "data": {
        "status": "okay",
        "detail": true
      },
      "errors": null,
      "code": 200
    }
    ```
7. You can add another endpoint along with their json response as many as you want, and of course you does not need to restart the app since it support hot reload thanks to [Viper](https://github.com/spf13/viper) 

# Known Limitation
- Only support `JSON` response
- Does not support auth. so anyone can hit the endpoint and get the JSON response
- Does not support any request validation (future development)
- Does not support giving multiple responses in single endpoint. say you want give response either `200` or `400` according to request in single endpoint