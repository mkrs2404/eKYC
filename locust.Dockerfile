FROM python:3.6

RUN apt-get update && \ 
    pip3 install locust && \
    locust -V

RUN cd /opt/ && \
    mkdir load-testing && \
    cd load-testing/

EXPOSE 8089

CMD ["locust", "-f", "/opt/load-testing/locustfile.py"]
