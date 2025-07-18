#!/bin/bash

# Escaping Eden Server Utilities
# A menu-driven interface for server maintenance tools

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_color() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print header
print_header() {
    echo
    print_color $CYAN "=================================="
    print_color $CYAN "  Escaping Eden Server Utilities"
    print_color $CYAN "=================================="
    echo
}

# Function to check if we're in the right directory
check_directory() {
    if [[ ! -f "go.mod" ]] || [[ ! -d "tools" ]]; then
        print_color $RED "Error: Please run this script from the Escaping Eden project root directory"
        print_color $YELLOW "Expected files: go.mod, tools/ directory"
        exit 1
    fi
}

# Function to get tool description
get_tool_description() {
    local tool_name=$1
    case "$tool_name" in
        "sync_items.go")
            echo "Sync item properties from JSON to database (preserves player inventories)"
            ;;
        "fix_active_character_items.go")
            echo "Fix equippable field and other properties for existing character items"
            ;;
        "debug_database.go")
            echo "Debug tool to examine database contents and find characters"
            ;;
        "test_equipment.go")
            echo "Test the character equipment system functionality"
            ;;
        *)
            echo "No description available"
            ;;
    esac
}

# Function to display available tools
show_tools_menu() {
    print_color $GREEN "Available Tools:"
    echo
    
    local counter=1
    local tool_files=()
    
    # Scan tools directory for .go files
    for file in tools/*.go; do
        if [[ -f "$file" ]]; then
            local basename=$(basename "$file")
            tool_files+=("$basename")
            local description=$(get_tool_description "$basename")
            print_color $BLUE "  $counter) $basename"
            print_color $YELLOW "     $description"
            echo
            ((counter++))
        fi
    done
    
    if [[ ${#tool_files[@]} -eq 0 ]]; then
        print_color $RED "No tools found in the tools/ directory"
        return 1
    fi
    
    print_color $PURPLE "  0) Exit"
    echo
    
    # Get user selection
    while true; do
        read -p "Select a tool to run (0-$((${#tool_files[@]})): " choice
        
        if [[ "$choice" == "0" ]]; then
            print_color $GREEN "Goodbye!"
            exit 0
        elif [[ "$choice" =~ ^[1-9][0-9]*$ ]] && [[ "$choice" -le "${#tool_files[@]}" ]]; then
            local selected_tool="${tool_files[$((choice-1))]}"
            run_tool "$selected_tool"
            break
        else
            print_color $RED "Invalid selection. Please choose a number between 0 and ${#tool_files[@]}"
        fi
    done
}

# Function to run a selected tool
run_tool() {
    local tool_name=$1
    local tool_path="tools/$tool_name"
    
    echo
    print_color $GREEN "Running: $tool_name"
    print_color $YELLOW "Path: $tool_path"
    echo
    
    # Confirm before running
    read -p "Are you sure you want to run this tool? (y/N): " confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        print_color $YELLOW "Operation cancelled"
        return
    fi
    
    echo
    print_color $CYAN "Executing tool..."
    echo "----------------------------------------"
    
    # Change to tools directory and run the tool
    cd tools
    if go run "$tool_name"; then
        echo "----------------------------------------"
        print_color $GREEN "Tool completed successfully!"
    else
        echo "----------------------------------------"
        print_color $RED "Tool failed with exit code $?"
    fi
    cd ..
    
    echo
    read -p "Press Enter to continue..."
}

# Function to show help
show_help() {
    print_header
    echo "This script provides a menu interface for Escaping Eden server utilities."
    echo
    print_color $GREEN "Usage:"
    echo "  ./server-utils.sh          - Show interactive menu"
    echo "  ./server-utils.sh --help   - Show this help message"
    echo
    print_color $GREEN "Adding New Tools:"
    echo "1. Place your .go file in the tools/ directory"
    echo "2. Edit this script to add a description in the 'tools' array"
    echo "3. The tool will automatically appear in the menu"
    echo
    print_color $GREEN "Requirements:"
    echo "- Must be run from the project root directory"
    echo "- Go must be installed and available in PATH"
    echo "- Tools directory must exist with .go files"
    echo
}

# Main function
main() {
    # Check for help flag
    if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
        show_help
        exit 0
    fi
    
    # Check directory and prerequisites
    check_directory
    
    # Show header and menu
    print_header
    show_tools_menu
}

# Run main function with all arguments
main "$@"