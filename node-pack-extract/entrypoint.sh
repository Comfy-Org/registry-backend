#!/bin/bash
set -e

# run the original endpoint
init.sh &

# loop until we can extract the node information
TIMEOUT=${TIMEOUT:-600}
OUTPUTFILE=${1:-"/tmp/output.json"}
echo -n >"$OUTPUTFILE"
until cat "$OUTPUTFILE" | grep ''; do
    sleep 1

    if ((SECONDS >= TIMEOUT)); then
        jq -n '{success: false, reason: "timeout"}' | tee "$OUTPUTFILE"
        break
    fi

    curl -sf localhost:8188/object_info |
        jq -c '
            to_entries |
            map(  
                select(.value.python_module == "custom_nodes.'$CUSTOM_NODE_NAME'")
                | .value |= {
                    category : .category,
                    description : .description,
                    deprecated : .deprecated,
                    experimental : .experimental,
                    input_types : (.input | tojson),
                    return_names : (.output_name | tojson),
                    return_types : (.output | tojson),
                    output_is_list : .output_is_list,
                }
            ) |
            if length > 0 then 
                {success: true, nodes: from_entries} 
            else 
                {success: false, reason: "node cannot be loaded into comfy ui"} 
            end' |
        tee "$OUTPUTFILE"
done

jq -e '.success' "$OUTPUTFILE"
