#!/usr/bin/env bash

input=$(cat)

# Colors
RESET='\033[0m'
BOLD='\033[1m'
DIM='\033[2m'
CYAN='\033[36m'
YELLOW='\033[33m'
GREEN='\033[32m'
RED='\033[31m'
MAGENTA='\033[35m'
WHITE='\033[37m'

# ============================================================================
# Context Window Usage Bar
# ============================================================================

render_usage_bar() {
  local pct=$1
  local bar_width=10
  local filled=$((pct * bar_width / 100))
  local empty=$((bar_width - filled))
  local bar=""

  [ "$filled" -gt 0 ] && printf -v fill "%${filled}s" && bar="${fill// /▓}"
  [ "$empty" -gt 0 ] && printf -v pad "%${empty}s" && bar="${bar}${pad// /░}"

  echo "$bar"
}

get_bar_color() {
  local pct=$1
  if [ "$pct" -ge 80 ]; then
    echo "$RED"
  elif [ "$pct" -ge 50 ]; then
    echo "$YELLOW"
  else
    echo "$GREEN"
  fi
}

# ============================================================================
# Cost Calculations
# ============================================================================

format_cost() {
  local usd=$1
  local jpy=$(echo "$usd * 150" | bc)
  printf '¥%.0f' "$jpy"
}

# ============================================================================
# Rate Limits
# ============================================================================

format_rate_limits() {
  local five_h=$1
  local week=$2
  local limits="undetermined"

  [ -n "$five_h" ] && limits="5h: $(printf '%.0f' "$five_h")%"
  [ -n "$week" ] && limits="${limits:+$limits }7d: $(printf '%.0f' "$week")%"

  echo "$limits"
}

# ============================================================================
# Git Status
# ============================================================================

get_git_status() {
  local dir=$1
  local git_info="none"

  local branch=$(git -C "$dir" symbolic-ref --short HEAD 2>/dev/null)
  if [ -n "$branch" ]; then
    local dirty=""
    if ! git -C "$dir" diff --quiet 2>/dev/null || ! git -C "$dir" diff --cached --quiet 2>/dev/null; then
      dirty=" ${RED}✗ dirty${RESET}"
    else
      dirty=" ${GREEN}✓ clean${RESET}"
    fi
    git_info="${BOLD}${CYAN}${branch}${RESET}${dirty}"
  fi

  echo "$git_info"
}

# ============================================================================
# Render Output
# ============================================================================

MODEL=$(echo "$input" | jq -r '.model.display_name')
DIR=$(echo "$input" | jq -r '.workspace.current_dir')
PCT=$(echo "$input" | jq -r '.context_window.used_percentage // 0' | cut -d. -f1)

BAR=$(render_usage_bar "$PCT")
BAR_COLOR=$(get_bar_color "$PCT")

COST=$(echo "$input" | jq -r '.cost.total_cost_usd // 0')
COST_FMT=$(format_cost "$COST")

FIVE_H=$(echo "$input" | jq -r '.rate_limits.five_hour.used_percentage // empty')
WEEK=$(echo "$input" | jq -r '.rate_limits.seven_day.used_percentage // empty')
LIMITS=$(format_rate_limits "$FIVE_H" "$WEEK")

GIT_INFO=$(get_git_status "$DIR")

echo -e "${DIM}Using${RESET} ${BOLD}${CYAN}${MODEL}${RESET} ${DIM}in${RESET} ${YELLOW}${DIR}${RESET}"
echo -e "${DIM}Usage:${RESET} ${BAR_COLOR}${BAR}${RESET} ${BOLD}${PCT}%${RESET} ${DIM}|${RESET} ${GREEN}~${COST_FMT} equiv${RESET} ${DIM}|${RESET} ${MAGENTA}${LIMITS}${RESET}"
echo -e "${DIM}Git:${RESET} ${GIT_INFO}"
