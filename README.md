# slack-memobot

A memo bot written in go language using slack socket mode

See k8s directory if running on kubernetes.

The dictionary files are stored in the `data` directory of the execution directory. Dictionary files are created for each slack channel.

The slack bot token and api token are read from environment variables. If a `.env` file exists, it will be read from the .env file.
