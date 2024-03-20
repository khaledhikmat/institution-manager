dapr run --app-id campaign-manager-membership \
         --app-protocol http \
         --app-port 8085 \
         --dapr-http-port 3505 \
         --log-level debug \
         --resources-path ./deploy/local/components \
         go run ./membership/main.go