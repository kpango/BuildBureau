#!/bin/bash

# BuildBureau Bootstrap Mode
# Enables BuildBureau to build and improve itself

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔══════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   BuildBureau Bootstrap Mode        ║${NC}"
echo -e "${BLUE}║   Self-Hosting/Self-Improvement     ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════╝${NC}"
echo

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "bootstrap" ]; then
    echo -e "${RED}Error: Must run from BuildBureau root directory${NC}"
    exit 1
fi

# Check for required API keys
API_KEYS_FOUND=0
if [ ! -z "$GEMINI_API_KEY" ]; then
    echo -e "${GREEN}✓${NC} GEMINI_API_KEY found"
    API_KEYS_FOUND=1
fi
if [ ! -z "$OPENAI_API_KEY" ]; then
    echo -e "${GREEN}✓${NC} OPENAI_API_KEY found"
    API_KEYS_FOUND=1
fi
if [ ! -z "$CLAUDE_API_KEY" ]; then
    echo -e "${GREEN}✓${NC} CLAUDE_API_KEY found"
    API_KEYS_FOUND=1
fi

if [ $API_KEYS_FOUND -eq 0 ]; then
    echo -e "${RED}Error: No LLM API keys found${NC}"
    echo "Set at least one of: GEMINI_API_KEY, OPENAI_API_KEY, CLAUDE_API_KEY"
    exit 1
fi

# Ensure build directory exists
echo -e "\n${BLUE}Building BuildBureau...${NC}"
make build || {
    echo -e "${RED}Build failed${NC}"
    exit 1
}

# Ensure bootstrap database directory exists
mkdir -p data

# Set bootstrap configuration
export BUILDBUREAU_CONFIG="bootstrap/config.yaml"

echo -e "\n${GREEN}✓${NC} Bootstrap environment ready"
echo -e "${BLUE}Configuration:${NC} $BUILDBUREAU_CONFIG"
echo -e "${BLUE}Database:${NC} ./data/bootstrap.db"
echo -e "${BLUE}Agents:${NC} bootstrap/agents/"
echo

echo -e "${YELLOW}═══════════════════════════════════════${NC}"
echo -e "${YELLOW}  BOOTSTRAP MODE: SELF-IMPROVEMENT${NC}"
echo -e "${YELLOW}═══════════════════════════════════════${NC}"
echo
echo -e "BuildBureau will now run in ${GREEN}self-hosting mode${NC}."
echo -e "Agents have deep knowledge of BuildBureau's codebase."
echo
echo -e "${BLUE}What you can do:${NC}"
echo -e "  • Add new features to BuildBureau"
echo -e "  • Refactor existing code"
echo -e "  • Optimize performance"
echo -e "  • Add tests"
echo -e "  • Fix bugs"
echo
echo -e "${YELLOW}⚠  Review all generated code before applying!${NC}"
echo -e "${YELLOW}⚠  Test changes thoroughly!${NC}"
echo
echo -e "Press Enter to start..."
read

# Run BuildBureau in bootstrap mode
echo -e "\n${GREEN}Starting BuildBureau in bootstrap mode...${NC}\n"
./build/buildbureau

echo -e "\n${GREEN}Bootstrap session completed${NC}"
echo -e "Review changes with: ${BLUE}git diff${NC}"
