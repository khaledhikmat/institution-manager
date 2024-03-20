dapr stop --app-id campaign-manager-institution
(lsof -i:8082 | grep main) | awk '{print $2}' | xargs kill