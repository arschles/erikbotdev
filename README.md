# ErikBotDev

**For all things k8s: https://k8syaml.com/**

>h/t rockerBoo

TODO: Add description and usage here later.

Modular support. Pick which modules you like.

### Modules

- Hue (lights)
- Elgato (fill lights)
- OBS (scenes)
- Browser Source (use bot as a web endpoint to update a browser source)
- Voice Effects (Through VST in OBS)

## Configuration

Many of the commands in this bot interact with OBS through the [OBS websocket plugin](https://obsproject.com/forum/resources/obs-websocket-remote-control-obs-studio-from-websockets.466/). By default, this bot _requires_ that it can connect to the websocket in order to even run. If it can't it marks itself as offline and won't respond to any commands.

### If you want to force this into "streaming mode"

If you do not have the OBS websocket plugin running, you have two options for forcing the bot to configure itself into streaming mode:

- Pass `-s` or `--streaming-on`
- Mark an individual command `"offline": true` to enable just that command

## Sounds

Any sound you reference in the config file ([sample](./erikbotdev.json)) needs to be a WAV file in the media directory.

# ErikBotServer

Hosted Twitch bot server. Accessible at http://51.105.203.101 (no HTTPS right now)

To install with Helm:

```shell
helm install \
    erikbotserver \
    ./charts/erikbotserver \
    -n erikbotserver \
    --create-namespace \
    --set server.clientID=${TWITCH_CLIENT_ID} \
    --set server.clientSecret=${TWITCH_CLIENT_SECRET} \
    --set server.oauthToken=${TWITCH_OAUTH_TOKEN}
```

Update:

```shell
helm upgrade \
    erikbotserver \
    ./charts/erikbotserver \
    -n erikbotserver \
    --set server.clientID=${TWITCH_CLIENT_ID} \
    --set server.clientSecret=${TWITCH_CLIENT_SECRET} \
    --set server.oauthToken=${TWITCH_OAUTH_TOKEN}
```

Delete:

```shell
helm uninstall -n erikbotserver erikbotserver
```

# Ingress Controller

Using nginx-ingress and cert-manager

https://cert-manager.io/docs/tutorials/acme/ingress/

Use nginx-ingress and cert-manager

Installing cert-manager:

```shell
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.0.1 \
  --set installCRDs=true \
  --create-namespace
```