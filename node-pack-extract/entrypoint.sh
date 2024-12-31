#!/bin/sh
set -e

# run the original endpoint 
init.sh &

# loop until we can extract the node information
OUTPUTFILE=${1:-"/tmp/output.json"}
echo -n > "$OUTPUTFILE"
until cat "$OUTPUTFILE" | grep ''; do
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
                    return_names : .output_name,
                    return_types : .output,
                    output_is_list : .output_is_list,
                }
            ) |
            if length > 0 then {nodes: from_entries} else "" end' | 
        tee "$OUTPUTFILE" 

    sleep 1
done

# make sure its json or we fail
grep '{' "$OUTPUTFILE"
