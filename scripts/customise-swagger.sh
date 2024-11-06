#!/usr/bin/env bash

# Converting yaml to json
docker run --rm -v "${PWD}":/workdir mikefarah/yq -p yaml -o json ./docs/static/swagger.yaml > ./docs/static/swagger.json

echo "Customizing swagger files..."

swagger_data=$(jq '.' ./docs/static/swagger.json)

# Adding license
swagger_data=$(echo "${swagger_data}" | jq '.info += {"license":{"name":"MIT License","url":"https://github.com/arkeonetwork/arkeo/blob/master/LICENSE"}}')

# Adding external docs
swagger_data=$(echo "${swagger_data}" | jq '. += {"externalDocs":{"description":"Find out more about Arkeo","url":"https://docs.arkeo.network"}}')


# Removing specific tags
tags_to_remove='["Query", "Service"]'
swagger_data=$(echo "$swagger_data" | jq --argjson tags_to_remove "$tags_to_remove" '
    # Check if "tags" exists before filtering top-level tags
    if .tags then
        .tags |= map(select(.name as $tag | $tags_to_remove | index($tag) | not))
    else . end |

    # Check if "paths" exists and iterate to remove tags from each operation
    if .paths then
        .paths |= with_entries(
            .value |= with_entries(
                if .value.tags then
                    .value.tags |= map(select(. as $tag | $tags_to_remove | index($tag) | not))
                else .
                end
            )
        )
    else . end
')

# Save the modified Swagger data back to JSON file
echo "$swagger_data" > ./docs/static/swagger.json

# Minifying json
jq -c . < ./docs/static/swagger.json > ./docs/static/swagger.min.json

# Cleanup
rm ./docs/static/swagger.yaml
rm ./docs/static/swagger.json