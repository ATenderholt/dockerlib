import logging

LOGGER = logging.getLogger()
LOGGER.setLevel(logging.INFO)


def handler(event, _):
    LOGGER.info("Got event: %s", event)
