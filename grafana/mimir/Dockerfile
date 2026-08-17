FROM grafana/mimir:3.1.1

COPY mimir-config.yaml /etc/mimir/mimir.yml

CMD ["--config.file=/etc/mimir/mimir.yml", "--target=all"]
