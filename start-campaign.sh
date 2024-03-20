dapr run --app-id campaign-manager-campaign \
         --app-protocol http \
         --app-port 8080 \
         --dapr-http-port 3500 \
         --log-level debug \
         --resources-path ./deploy/local/components \
         go run ./campaign/main.go