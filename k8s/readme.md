# Usage

Prepare an .env file. In it, write the app token and bot token as follows

```
SLACK_APP_TOKEN=xapp-xxx
SLACK_BOT_TOKEN=xoxb-xxx
```

* run `k8s_create_secret.sh` for creating secret from dotenv.
* run `k8s_create_app.sh` for creating pv,pvc,deployment. 