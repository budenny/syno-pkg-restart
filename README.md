# syno-pkg-restart

## A webhook-based container for restarting packages/services on Synology NAS(host for this container).  


**Important: Requires to be privileged container.**


# How to toggle webhook:
restarting "Tailscale" on Synology NAS:
```bash
curl -X GET -H "Authorization: Bearer <token>" http://<nas_ip>:5480/\?svc\=Tailscale
```

# How to build docker container:
```bash
docker build -t syno-pkg-restart .
```

# How to run docker container:
```bash
  docker run -d \
  --name=syno-pkg-restart \
  --privileged \
  -p 5480:80 \
  -v /:/host \
  -e TOKEN="<token>" \
  -e HOST_ROOT="/host" \
  --restart unless-stopped \
  syno-pkg-restart
```

# Token:
Optional. Arbitrary secret value between client and container. E.g. generate uuid set it in container, then use during toggling webhook.
