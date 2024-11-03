for file in ./*; do
height=$(identify -format "%h" $file)
width1=$(echo "$height * 0.5625" | bc)
width2=$(printf "%.0f" $width1)

echo "height : $width2 width: $height"

#    magick "$file" -resize "${width2}x${height}^" -gravity center -crop "${width2}:${height}" +repage "wallpapers/$(basename "$file")"
    magick "$file" -resize "${width2}x${height}" -crop "${new_width}:${height}" "$(basename "$file")"

    echo "done $file"
#-gravity center -extent "${height}x${width2}"
done
