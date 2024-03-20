dapr stop --app-id campaign-manager-membership
(lsof -i:8085 | grep main) | awk '{print $2}' | xargs kill