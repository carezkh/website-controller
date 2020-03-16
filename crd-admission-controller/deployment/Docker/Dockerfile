FROM debian
COPY webhook-controller /
COPY key/ /key/
ENTRYPOINT ["/webhook-controller","--ip","0.0.0.0","--port","443"]
