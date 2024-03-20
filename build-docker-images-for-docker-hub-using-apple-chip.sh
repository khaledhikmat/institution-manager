docker buildx build --platform linux/amd64 -t khaledhikmat/campaignmanager-campaign:latest . -f Dockerfile-campaign
docker buildx build --platform linux/amd64 -t khaledhikmat/campaignmanager-htmx:latest . -f Dockerfile-htmx
