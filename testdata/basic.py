import logging
import time
import sys

LOGGER = logging.getLogger()
LOGGER.setLevel(logging.INFO)


if __name__ == "__main__":
    logging.basicConfig(stream=sys.stdout)

    LOGGER.info("Hello!")
    time.sleep(3)

