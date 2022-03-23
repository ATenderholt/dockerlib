import logging
import sys
import urllib.request

LOGGER = logging.getLogger()
LOGGER.setLevel(logging.INFO)

if __name__ == '__main__':
    logging.basicConfig(stream=sys.stdout)

    with urllib.request.urlopen(sys.argv[1] + "/hello.txt") as r:
        LOGGER.info("Status: %d", r.status)
