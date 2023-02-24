def application(env, start_response):
    body = "Hello, 世界".encode()
    start_response(
        "200 OK",
        [("Content-Type", "text/plain"), ("Content-Length", str(len(body)))],
    )
    return [body]
