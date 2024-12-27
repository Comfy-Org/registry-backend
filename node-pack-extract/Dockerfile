FROM ghcr.io/ai-dock/comfyui:v2-cpu-22.04-v0.2.7

COPY ./provisioning.sh /opt/ai-dock/bin/provisioning.sh
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV WEB_ENABLE_AUTH false
ENTRYPOINT [ "/entrypoint.sh" ]