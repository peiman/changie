#!/bin/bash

# This script adds a 'v' prefix to all non-prefixed version tags and optionally removes the old tags

# Exit on any error
set -e

echo "Updating Git tags to add 'v' prefix..."

# Get all numeric version tags
tags=$(git tag -l | grep "^[0-9]" | sort -V)

if [ -z "$tags" ]; then
    echo "No numeric version tags found. Nothing to do."
    exit 0
fi

echo "Found the following version tags to update: $tags"
echo "This script will:"
echo "1. Create new tags with 'v' prefix for each version"
echo "2. Optionally remove the old tags (if --remove-old is specified)"
echo "3. Push the new tags to remote (if --push is specified)"
echo ""
echo "Press ENTER to continue or CTRL+C to abort..."
read -r

# Process each tag
for tag in $tags; do
    new_tag="v$tag"
    
    # Check if the v-prefixed tag already exists
    if git rev-parse "$new_tag" >/dev/null 2>&1; then
        echo "Tag $new_tag already exists, skipping..."
        continue
    fi
    
    # Get the commit hash for the original tag
    commit=$(git rev-parse "$tag")
    
    # Create the new tag
    echo "Creating new tag $new_tag at commit $commit..."
    git tag "$new_tag" "$commit"
    
    echo "Created tag $new_tag"
done

echo "All tags updated successfully!"

# Remove old tags if requested
if [[ "$1" == "--remove-old" || "$2" == "--remove-old" || "$3" == "--remove-old" ]]; then
    echo "Removing old tags without 'v' prefix..."
    for tag in $tags; do
        echo "Removing tag $tag..."
        git tag -d "$tag"
    done
    echo "Old tags removed successfully!"
fi

# Push tags if requested
if [[ "$1" == "--push" || "$2" == "--push" || "$3" == "--push" ]]; then
    echo "Pushing tags to remote..."
    git push --tags
    
    # If we removed old tags and want to push, we need to remove them from the remote too
    if [[ "$1" == "--remove-old" || "$2" == "--remove-old" || "$3" == "--remove-old" ]]; then
        echo "Removing old tags from remote..."
        for tag in $tags; do
            echo "Removing remote tag $tag..."
            git push origin ":refs/tags/$tag"
        done
        echo "Old remote tags removed successfully!"
    fi
    
    echo "Tags pushed successfully!"
fi

echo "Done!" 