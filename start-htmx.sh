dapr run --app-id campaign-manager-htmx \
         --app-protocol http \
         --app-port 3000 \
         --dapr-http-port 3510 \
         --log-level debug \
         --resources-path ./deploy/local/components \
         go run ./htmx/main.go