dapr run --app-id campaign-manager-externalizer \
         --app-protocol http \
         --app-port 8083 \
         --dapr-http-port 3503 \
         --log-level debug \
         --resources-path ./deploy/local/components \
         go run ./externalizer/main.go