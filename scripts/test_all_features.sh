#!/bin/bash

echo "🧪 STORMFINDER COMPREHENSIVE FEATURE TEST"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run test
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_exit_code="${3:-0}"
    
    echo -e "\n${BLUE}Testing: $test_name${NC}"
    echo "Command: $command"
    
    if eval "$command" > /dev/null 2>&1; then
        local exit_code=$?
        if [ $exit_code -eq $expected_exit_code ]; then
            echo -e "${GREEN}✅ PASSED${NC}"
            ((TESTS_PASSED++))
        else
            echo -e "${RED}❌ FAILED (Exit code: $exit_code, Expected: $expected_exit_code)${NC}"
            ((TESTS_FAILED++))
        fi
    else
        echo -e "${RED}❌ FAILED (Command execution error)${NC}"
        ((TESTS_FAILED++))
    fi
}

# Build the tool first
echo -e "${YELLOW}🔨 Building Stormfinder...${NC}"
if go build ./cmd/stormfinder; then
    echo -e "${GREEN}✅ Build successful${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Test 1: Help command
run_test "Help Command" "./stormfinder -h"

# Test 2: Version command  
run_test "Version Command" "./stormfinder -version"

# Test 3: List sources
run_test "List Sources" "./stormfinder -ls"

# Test 4: Basic enumeration
run_test "Basic Enumeration" "./stormfinder -d example.com -silent"

# Test 5: Enhanced enumeration (brute force + permutations)
run_test "Enhanced Enumeration" "./stormfinder -d example.com -b -p -silent"

# Test 6: Caching feature
run_test "Caching Feature" "./stormfinder -d example.com --cache -silent"

# Test 7: AI feature (placeholder)
run_test "AI Feature" "./stormfinder -d example.com --ai -silent"

# Test 8: Advanced CT mining (placeholder)
run_test "Advanced CT Mining" "./stormfinder -d example.com --advanced-ct -silent"

# Test 9: Social mining (placeholder)
run_test "Social Mining" "./stormfinder -d example.com --social -silent"

# Test 10: Relationship mapping (placeholder)
run_test "Relationship Mapping" "./stormfinder -d example.com --map -silent"

# Test 11: JSON output
run_test "JSON Output" "./stormfinder -d example.com -oJ -o test_output.json -silent"

# Test 12: Verbose mode
run_test "Verbose Mode" "./stormfinder -d example.com -v"

# Test 13: Performance optimization
run_test "Speed Optimization" "./stormfinder -d example.com --optimize-speed -silent"

# Test 14: Memory optimization  
run_test "Memory Optimization" "./stormfinder -d example.com --optimize-memory -silent"

# Test 15: Custom wordlist (if exists)
if [ -f "wordlist.txt" ]; then
    run_test "Custom Wordlist" "./stormfinder -d example.com -b -w wordlist.txt -silent"
else
    echo -e "\n${YELLOW}⚠️  Skipping Custom Wordlist test (wordlist.txt not found)${NC}"
fi

# Test 16: Multiple domains from file
echo "example.com" > test_domains.txt
echo "google.com" >> test_domains.txt
run_test "Multiple Domains" "./stormfinder -dL test_domains.txt -silent"

# Clean up test files
rm -f test_output.json test_domains.txt

# Final results
echo -e "\n${BLUE}========================================${NC}"
echo -e "${BLUE}🏁 TEST RESULTS SUMMARY${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ Tests Passed: $TESTS_PASSED${NC}"
echo -e "${RED}❌ Tests Failed: $TESTS_FAILED${NC}"

TOTAL_TESTS=$((TESTS_PASSED + TESTS_FAILED))
SUCCESS_RATE=$((TESTS_PASSED * 100 / TOTAL_TESTS))

echo -e "${BLUE}📊 Success Rate: $SUCCESS_RATE%${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 ALL TESTS PASSED! Stormfinder is ready for GitHub release! 🚀${NC}"
    exit 0
else
    echo -e "\n${YELLOW}⚠️  Some tests failed. Please review and fix issues before release.${NC}"
    exit 1
fi
