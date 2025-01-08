#!/bin/bash
set -e

source /opt/environments/python/comfyui/bin/activate

echo "=========================="
echo "=== Pinning Dependencies ==="
echo "=========================="

pip install "numpy<2"

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
if [ -f "requirements.txt" ]; then
    echo "=== installing dependencies from requirements.txt"
    pip install -r requirements.txt 
fi
if [ -f "install.py" ]; then
    echo "=== executing install.py"
    python install.py
fi
