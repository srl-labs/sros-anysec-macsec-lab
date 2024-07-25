# Automation Panel Frontend

Build image using: `docker build -t ghcr.io/srl-labs/sros-anysec-macsec-lab/frontend .`

In a proxy environment, use build arguments:

```bash
docker build --build-arg https_proxy=$(HTTP_PROXY) -t ghcr.io/srl-labs/sros-anysec-macsec-lab/frontend .
```
