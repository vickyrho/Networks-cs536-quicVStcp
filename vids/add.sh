for file in ./video/*.mp4; do
    size=$(stat -c%s "$file")
    {
        echo -e "HTTP/1.1 200 OK"
        echo -e "Content-Type: video/mp4"
        echo -e "Content-Length: $size"
        echo
        cat "$file"
    } > "$file.http"
done

