FROM rabbitmq:3.10-management
RUN rabbitmq-plugins enable --offline \
    rabbitmq_mqtt rabbitmq_web_mqtt \
    rabbitmq_stomp rabbitmq_web_stomp \
    rabbitmq_web_stomp_examples \
    rabbitmq_web_mqtt_examples
COPY ./advanced.config /etc/rabbitmq/advanced.config

COPY ./cert/ /etc/rabbitmq/cert
RUN chown -R rabbitmq:rabbitmq /etc/rabbitmq/cert
COPY ./rabbitmq.conf /etc/rabbitmq/rabbitmq.conf

