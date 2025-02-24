#!/bin/bash
set -e

INPUT_FILE="$1"
OUTPUT_DIR="./output"

mkdir -p "$OUTPUT_DIR"

ffmpeg -i "$INPUT_FILE" \
    -preset fast \
    -g 48 \
    -sc_threshold 0 \
    -map 0:v:0 \
    -map 0:a:0 \
    -b:v:0 2500k \
    -maxrate 2675k \
    -bufsize 3750k \
    -hls_time 4 \
    -hls_playlist_type vod \
    -hls_segment_filename "$OUTPUT_DIR/output_%03d.ts" \
    "$OUTPUT_DIR/output.m3u8"

echo "Conversion complete. Files are in $OUTPUT_DIR"
