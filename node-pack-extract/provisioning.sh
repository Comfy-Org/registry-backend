#!/bin/bash
set -e

echo "==========================================="
echo "=== Downloading the custom node archive ==="
echo "==========================================="
cd /opt/ComfyUI/custom_nodes/
rm -rf "$CUSTOM_NODE_NAME"
mkdir -p "$CUSTOM_NODE_NAME"
wget -O "$CUSTOM_NODE_NAME.zip" "$CUSTOM_NODE_URL"

echo "=================================="
echo "=== Installing the custom node ==="
echo "=================================="
unzip "${CUSTOM_NODE_NAME}.zip" -d "$CUSTOM_NODE_NAME"
cd "$CUSTOM_NODE_NAME"
source /opt/environments/python/comfyui/bin/activate
pip install -r requirements.txt 
if [ -f "install.py" ]; then
    python install.py
fi
