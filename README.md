# Server Analyzer & Notification Dispatcher When It Crashed Hard

I just want to send a friend of mine an image to run to check from outside my network if something goes offline.

## Getting Started

```bash
podman run \
   --env="SANDWICH_CHECK_URL=https://my.fvtt.wtf/" \
   --env="SANDWICH_WEBHOOK=https://discordapp.com/api/webhooks/1234567890123456789/xxxxxxxxx_xxxxxxxxxxxxxxxxxxxxxxxxx_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
   --env="SANDWICH_TARGET_DOWN_IDS=123456789012345678,8765432109876543210" 
   \ ghcr.io/marvinlwenzel/sandwich:1.0.0
```

## How the sausage is made

SSL is tricky in low level images, so https://stackoverflow.com/a/45397221 helped.

I did not yet manage to copy certs manually, but I would like to make it work