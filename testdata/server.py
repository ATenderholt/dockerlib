from functools import partial
import http.server
import socketserver
import logging
import sys

LOGGER = logging.getLogger()
LOGGER.setLevel(logging.INFO)

PORT = 8000

if __name__ == "__main__":
    logging.basicConfig(stream=sys.stdout)

    Handler = partial(http.server.SimpleHTTPRequestHandler, directory=sys.argv[1])

    with socketserver.TCPServer(("", PORT), Handler) as httpd:
        LOGGER.info("Server started on port %d", PORT)
        httpd.serve_forever()
