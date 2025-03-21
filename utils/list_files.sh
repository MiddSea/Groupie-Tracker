#!/bin/sh
# Name: list_files.sh
# Purpose: List all tracked files in the repository according to .gitignore in a tree format
# Authors: Seán David Middleton & Björn Österman
# Date: March 19, 2025
#
# Usage: ./list_files.sh [options] [directory]
#
# This script uses the 'tree' command to display files in a hierarchical structure,
# while respecting the .gitignore patterns to exclude files that shouldn't be tracked.
# It's designed to be portable and reusable across different projects.
#
# Options:
#   -f, -F, --filenames   Output only filenames, one per line
#   -s, --showhidden      Show hidden directories (.notes, .util, .vscode, etc.)
#   -h, -?, --help        Display this help message
#
# Examples:
#   ./list_files.sh                    # Standard tree view, hiding hidden directories
#   ./list_files.sh -f                 # List only filenames, hiding hidden directories
#   ./list_files.sh -s                 # Standard tree view, showing hidden directories
#   ./list_files.sh -f -s              # List only filenames, showing hidden directories

# Function to show usage
show_usage() {
    echo "Usage: $0 [options] [directory]"
    echo ""
    echo "Options:"
    echo "  -f, -F, --filenames   Output only filenames, one per line"
    echo "  -s, --showhidden      Show hidden directories (.notes, .util, .vscode, etc.)"
    echo "  -h, -?, --help        Display this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Standard tree view, hiding hidden directories"
    echo "  $0 -f                 # List only filenames, hiding hidden directories"
    echo "  $0 -s                 # Standard tree view, showing hidden directories"
    echo "  $0 -f -s              # List only filenames, showing hidden directories"
    exit 0
}

# Initialize variables
FILENAMES_ONLY=0
SHOW_HIDDEN=0
TARGET_DIR="."

# Parse command line arguments
while [ "$#" -gt 0 ]; do
    case "$1" in
        -f|-F|--filenames)
            FILENAMES_ONLY=1
            shift
            ;;
        -s|--showhidden)
            SHOW_HIDDEN=1
            shift
            ;;
        -h|-\?|--help)
            show_usage
            ;;
        -*)
            echo "Unknown option: $1"
            show_usage
            ;;
        *)
            TARGET_DIR="$1"
            shift
            ;;
    esac
done

# Check if tree command is available
if ! command -v tree > /dev/null 2>&1; then
    echo "Error: 'tree' command is not installed. Please install it first."
    echo "macOS: brew install tree"
    echo "Ubuntu/Debian: apt-get install tree"
    echo "CentOS/RHEL: yum install tree"
    exit 1
fi

# Get ignore patterns from .gitignore if it exists
IGNORE_PATTERNS=""
if [ -f "${TARGET_DIR}/.gitignore" ]; then
    # Process .gitignore file to create tree exclude patterns
    # Skip comment lines and empty lines
    while IFS= read -r line; do
        # Skip empty lines and comments
        if [ -n "$line" ] && ! echo "$line" | grep -q "^#"; then
            # Remove leading slash if present
            pattern=$(echo "$line" | sed 's/^\///')
            # Add to ignore patterns
            IGNORE_PATTERNS="${IGNORE_PATTERNS} -I \"${pattern}\""
        fi
    done < "${TARGET_DIR}/.gitignore"
fi

# Add common patterns to always ignore unless showing hidden is enabled
if [ "$SHOW_HIDDEN" -eq 0 ]; then
    IGNORE_PATTERNS="${IGNORE_PATTERNS} -I \".git\" -I \".notes\" -I \".util\" -I \".vscode\""
else
    # Always ignore .git even when showing hidden directories
    IGNORE_PATTERNS="${IGNORE_PATTERNS} -I \".git\""
fi

# Create the command with all ignore patterns
if [ "$FILENAMES_ONLY" -eq 1 ]; then
    # Use -if option for filenames only, with full paths
    CMD="tree -a -if ${TARGET_DIR} ${IGNORE_PATTERNS} | grep -v '^${TARGET_DIR}\$'"
else
    # Regular tree display
    CMD="tree -a ${TARGET_DIR} ${IGNORE_PATTERNS}"
    echo "Displaying tracked files in ${TARGET_DIR}:"
    echo "-------------------------------------------"
fi

# Execute the command using eval to handle the complex string with quotes
eval "${CMD}"

if [ "$FILENAMES_ONLY" -eq 0 ]; then
    echo "-------------------------------------------"
    if [ "$SHOW_HIDDEN" -eq 1 ]; then
        echo "Note: Hidden directories are shown (except .git)"
    else
        echo "Note: Hidden directories and files matching patterns in .gitignore are excluded"
    fi
fi