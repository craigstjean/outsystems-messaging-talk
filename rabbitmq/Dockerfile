FROM rabbitmq:3.10-management
RUN rabbitmq-plugins enable --offline \
    rabbitmq_mqtt rabbitmq_web_mqtt \
    rabbitmq_stomp rabbitmq_web_stomp \
    rabbitmq_web_stomp_examples \
    rabbitmq_web_mqtt_examples
COPY ./rabbitmq.config /etc/rabbitmq/advanced.config

