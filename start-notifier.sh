dapr run --app-id campaign-manager-notifier \
         --app-protocol http \
         --app-port 8081 \
         --dapr-http-port 3501 \
         --log-level debug \
         --resources-path ./deploy/local/components \
         go run ./notifier/main.go