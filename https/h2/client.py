import os
import time
import csv

if __name__ == '__main__':
    # Remove existing 'files' directory if it exists and create a new one
    if os.path.exists("./files"):
        os.system("rm -rf ./files")
    os.makedirs("./files", exist_ok=True)

    # Create or overwrite the result.csv file with the headers
    with open("result.csv", "w", newline="") as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["each file", "file number", "time"])

    # Define sizes and corresponding file numbers
    sizes = ["5KB", "10KB", "100KB", "200KB", "500KB", "1MB", "10MB"]
    file_numbers = {
        "5KB": [1, 200],
        "10KB": [1, 100],
        "100KB": [1, 10],
        "200KB": [1, 5],
        "500KB": [1, 2],
        "1MB": [1],
        "10MB": [1]
    }

    # Download files and record times
    for total_size in sizes:
        for number in file_numbers[total_size]:
            start_time = time.time()  # Record start time
            for i in range(number):
                os.system(f"wget --no-check-certificate -P ./files https://10.0.0.1:8000/files/test_{total_size}_{number}_{i}.txt")
            end_time = time.time()  # Record end time

            elapsed_time_ms = (end_time - start_time) * 1000  # Convert time to milliseconds
            print(f"size: {total_size}, number: {number}, download time: {elapsed_time_ms:.2f} ms")

            # Append results to the CSV file
            with open("result.csv", "a", newline="") as csvfile:
                writer = csv.writer(csvfile)
                writer.writerow([total_size, number, elapsed_time_ms])
