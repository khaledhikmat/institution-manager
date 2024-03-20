dapr stop --app-id campaign-manager-externalizer
(lsof -i:8083 | grep main) | awk '{print $2}' | xargs kill