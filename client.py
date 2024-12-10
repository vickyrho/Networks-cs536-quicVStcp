import subprocess
import time

def download_function(file_name):
    # Read fingerprints from the file
    with open("quic/certs/fingerprints.txt", "r") as f:
        fingerprints = f.read().strip()

    # Define the command as a list
    make_host = [
        "google-chrome",
        "--no-sandbox",
        "--headless",
        "--disable-gpu",
        "--user-data-dir=/tmp/chrome-profile",
        "--no-proxy-server",
        "--enable-quic",
        "--origin-to-force-quic-on=www.example.org:443",
        "--host-resolver-rules=MAP www.example.org:443 10.0.0.1:6121",
        "--ignore-certificate-errors-spki-list={}".format(fingerprints),
        "https://www.example.org/test.txt"
    ]

    # Run the subprocess and capture output/errors
    try:
        proc = subprocess.Popen(make_host, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout, stderr = proc.communicate()  # Get output and error
        if proc.returncode != 0:
            print "Error:", stderr
        else:
            print "Output:", stdout
    except Exception as e:
        print "Error running subprocess:", str(e)


total_time = 0
iteration = 1

for cnt in range(iteration):
    st = time.time()
    file_name = "index.html"
    download_function(file_name)
    en = time.time()
    total_time += ((en - st) * 1000)

avg_time = total_time / iteration
print "Total Time: ", total_time
print "Avg Time: ", avg_time
