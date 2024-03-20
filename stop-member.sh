dapr stop --app-id campaign-manager-member
(lsof -i:8084 | grep main) | awk '{print $2}' | xargs kill