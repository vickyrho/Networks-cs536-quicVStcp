import subprocess
make_host = ["./quic/Quic/quic_server", "--quic_response_cache_dir=www.example.org", "--certificate_file=quic/certs/leaf_cert.pem",
 "--key_file=quic/certs/leaf_cert.pkcs8"]
subprocess.call(make_host)


# import subprocess
# import requests
# import time

# # Step 1: Start the QUIC server in a subprocess
# make_host = [
#     "./quic/Quic/quic_server",
#     "--quic_response_cache_dir=www.example.org",
#     "--certificate_file=quic/certs/leaf_cert.pem",
#     "--key_file=quic/certs/leaf_cert.pkcs8"
# ]

# # Start the server as a subprocess
# server_process = subprocess.Popen(make_host)

# try:
#     # Wait for the server to start (if needed)
#     time.sleep(2)  # Adjust the delay based on server start time

#     # Step 2: Make a request to the server
#     url = "https://www.example.org/test.txt"  # Update URL as needed
#     headers = {"User-Agent": "quic-client"}

#     # Use a custom library or HTTP/3-enabled client if needed
#     # Since `requests` does not support QUIC/HTTP3 natively, you need a compatible library.
#     # Example for HTTP/1.1 or HTTP/2:
#     response = requests.get(url, headers=headers, verify=False)  # Disable SSL verification for testing

#     # Step 3: Print the response details
#     print("Status Code:", response.status_code)
#     print("Headers:", response.headers)
#     print("Response Content:", response.text)

# except Exception as e:
#     print("Error during request or response handling:", e)

# finally:
#     # Step 4: Stop the QUIC server
#     server_process.terminate()
#     server_process.wait()
