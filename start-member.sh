dapr run --app-id campaign-manager-member \
         --app-protocol http \
         --app-port 8084 \
         --dapr-http-port 3504 \
         --log-level debug \
         --resources-path ./deploy/local/components \
         go run ./member/main.go