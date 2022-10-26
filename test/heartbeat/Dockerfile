FROM ubuntu

RUN apt -y update && apt install -y ca-certificates curl bash iputils-ping 

COPY heartbeat.sh .

USER 1001

CMD ["./heartbeat.sh"]
