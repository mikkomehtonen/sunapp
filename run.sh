docker run -d -p 3104:8080 \
  --restart unless-stopped \
  --name sunapp \
  -e LOGO_LINK_URL \
   sunapp
