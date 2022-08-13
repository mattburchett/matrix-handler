# Matrix Handler

### Usage
> POST /generic/{matrixRoom}/{matrixUser}/{matrixPassword}
This is a generic handler that will accept any form of data and pass it through to Matrix using the appropriate JSON format.

### Usage
> POST /slack/{matrixRoom}/{matrixUser}/{matrixPassword}
This is a Slack-friendly handler that will accept any form of Slack data and pass through just the text to Matrix using the appropriate JSON format.


### Docker Usage

#### Prebuilt Image
There is a prebuilt image already on `ghcr.io/mattburchett/matrix-handler:latest`

#### Building Your Own

You can build the image by running 

```bash
docker build -t myrepo/matrix-handler:latest .
```

#### Docker Compose
You can use these settings for deploying via Docker Compose.

```yaml
version: '3.5'
services:
    matrix-handler:
      image: ghcr.io/mattburchett/matrix-handler:latest
      restart: always
      container_name: matrix-handler
      ports:
      - 3000:3000
      volumes:
        - ./config.json:/config/config.json:ro
```

