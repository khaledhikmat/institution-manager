dapr run --app-id campaign-manager-institution \
         --app-protocol http \
         --app-port 8082 \
         --dapr-http-port 3502 \
         --log-level debug \
         --resources-path ./deploy/local/components \
         go run ./institution/main.go