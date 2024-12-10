import os

def add_headers(directory, output_dir):
    for root, _, files in os.walk(directory):
        for file in files:
            file_path = os.path.join(root, file)
            output_path = os.path.join(output_dir, os.path.relpath(file_path, directory))
            os.makedirs(os.path.dirname(output_path), exist_ok=True)

            # Determine Content-Type
            if file.endswith('.mp4'):
                content_type = "video/mp4"
            elif file.endswith('.mpd'):
                content_type = "application/dash+xml"
            elif file.endswith('.js'):
                content_type = "application/javascript"
            elif file.endswith('.html'):
                content_type = "text/html"
            else:
                content_type = "application/octet-stream"

            # Get File Size
            file_size = os.path.getsize(file_path)

            # Prepare Headers
            headers = f"""HTTP/1.1 200 OK
Content-Type: {content_type}
Content-Length: {file_size}
X-Original-Url: /{os.path.relpath(file_path, directory).replace(os.sep, '/')}
"""

            # Write to Output File
            with open(file_path, 'rb') as f_in, open(output_path, 'wb') as f_out:
                f_out.write(headers.encode('utf-8') + b'\r\n')
                f_out.write(f_in.read())

input_dir = "./vids/test"
output_dir = "./www.example.org/cache_with_headers"
add_headers(input_dir, output_dir)

