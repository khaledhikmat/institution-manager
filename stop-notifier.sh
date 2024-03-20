dapr stop --app-id campaign-manager-notifier
(lsof -i:8081 | grep main) | awk '{print $2}' | xargs kill