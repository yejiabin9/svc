FROM alpine
ADD svc /svc
ENTRYPOINT [ "/svc" ]