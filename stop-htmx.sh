dapr stop --app-id campaign-manager-htmx
(lsof -i:3000 | grep main) | awk '{print $2}' | xargs kill