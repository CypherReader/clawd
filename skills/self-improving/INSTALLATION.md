# Self-Improving Skill - Installation Guide

Enable your Clawdbot instance to learn and improve from interactions.

## Quick Install

```bash
cd /root/clawd

# 1. Create memory structure
mkdir -p memory/learnings memory/patterns

# 2. Initialize files
touch memory/self-improvement.md
touch memory/patterns/coding-patterns.md
touch memory/patterns/communication-patterns.md
touch memory/patterns/problem-solving-patterns.md

# 3. Install skill (if not already in skills/)
# Skill should be at: /root/clawd/skills/self-improving/SKILL.md

# 4. Test the skill
skills/self-improving/scripts/capture_learning.sh
```

## Verification

The skill should be automatically available. Test by asking:

```
"What have you learned recently?"
"Show me your patterns for coding"
"Learn from this interaction"
```

## How It Works

### Automatic Learning

The agent will automatically:
1. **Capture** significant interactions in daily logs
2. **Analyze** patterns of success and failure
3. **Update** the knowledge base with learnings
4. **Apply** learnings to future similar tasks

### Manual Learning

You can explicitly trigger learning:
- **"Learn from this"** - Capture current interaction
- **"What patterns do you have for X?"** - Query patterns
- **"Improve your approach to X"** - Focus improvement

### Memory Structure

```
memory/
├── self-improvement.md          # Curated learnings (main file)
├── learnings/
│   └── YYYY-MM-DD.md           # Daily raw learning data
└── patterns/
    ├── coding-patterns.md       # Software development
    ├── communication-patterns.md # User interaction
    ├── problem-solving-patterns.md # Debug & troubleshoot
    └── [domain]-patterns.md     # Domain-specific patterns
```

## Integration with Existing Memory

Works alongside your current memory system:
- `MEMORY.md` - Long-term important facts
- `memory/YYYY-MM-DD.md` - Daily logs
- `memory/self-improvement.md` - **Learnings (NEW)**
- `memory/patterns/` - **Pattern library (NEW)**
- `memory/learnings/` - **Learning data (NEW)**

No conflicts - complementary systems.

## Configuration

### Adjust Learning Sensitivity

Edit `skills/self-improving/SKILL.md` frontmatter:

```yaml
description: Enable continuous learning... [customize when it triggers]
```

### Privacy Settings

By default, the skill:
- ✅ Does NOT store sensitive data
- ✅ Uses placeholders for keys/passwords
- ✅ Generalizes learnings (not user-specific)
- ✅ Can forget on request: "Forget learnings about X"

## Usage Examples

### After Completing a Task

```
User: "Learn from building the Oracle system"
Agent: [Analyzes, captures patterns, updates knowledge base]
```

### Before Similar Task

```
User: "I need to build another trading system"
Agent: [Searches memory, reviews Oracle learnings, applies patterns]
```

### Reviewing Progress

```
User: "What have you learned this week?"
Agent: [Shows curated learnings from memory/self-improvement.md]
```

### Querying Patterns

```
User: "Show me your coding patterns"
Agent: [Displays memory/patterns/coding-patterns.md]
```

## Metrics Tracking

The skill tracks:
- Task completion rate
- User satisfaction signals
- Efficiency improvements
- Learning velocity
- Pattern application success

View metrics in `memory/self-improvement.md` under "Metrics"

## Maintenance

### Weekly Review (Recommended)

```bash
# Review recent learnings
cat memory/learnings/$(date +%Y-%m-%d).md

# Consolidate into main file if significant
# Edit memory/self-improvement.md

# Update patterns if new ones discovered
# Edit memory/patterns/*.md
```

### Monthly Cleanup

1. Review old learning entries
2. Archive outdated patterns
3. Refine confidence levels
4. Remove superseded learnings

## Troubleshooting

### "Skill not triggering"

Check skill description is clear:
```bash
head -n 10 skills/self-improving/SKILL.md
```

Should clearly describe when to use the skill.

### "Memory files not found"

Ensure structure exists:
```bash
ls -la memory/learnings/
ls -la memory/patterns/
```

### "Can't capture learnings"

Check script permissions:
```bash
chmod +x skills/self-improving/scripts/capture_learning.sh
```

## Advanced Usage

### Custom Pattern Categories

Create new pattern files for your domains:

```bash
touch memory/patterns/oracle-trading-patterns.md
touch memory/patterns/data-analysis-patterns.md
touch memory/patterns/project-management-patterns.md
```

### Automated Learning Triggers

Set up cron jobs to review learnings:

```bash
# Review and consolidate weekly
cron --add "0 0 * * 0" "Review weekly learnings"
```

### Integration with Other Skills

The self-improving skill works with:
- **GitHub skill**: Learn from code reviews
- **Notion skill**: Document learnings in Notion
- **Slack skill**: Share learnings with team

## Success Indicators

You'll know it's working when:
- ✅ Fewer corrections needed
- ✅ Faster task completion
- ✅ Better anticipation of needs
- ✅ Consistent improvement trends
- ✅ More successful first attempts

## Example: Oracle Project Learning

The skill has already captured learnings from the Oracle trading system:

```bash
# View Oracle learnings
cat memory/self-improvement.md

# View coding patterns learned
cat memory/patterns/coding-patterns.md
```

These patterns will be applied to future similar projects automatically!

## Getting Help

If the skill isn't working as expected:

1. Check skill file exists: `ls skills/self-improving/SKILL.md`
2. Verify memory structure: `ls -R memory/`
3. Test script manually: `skills/self-improving/scripts/capture_learning.sh`
4. Ask agent: "Why isn't self-improvement working?"

## Uninstallation

To remove (not recommended):

```bash
rm -rf skills/self-improving
rm -rf memory/learnings
rm -rf memory/patterns
rm memory/self-improvement.md
```

This will not affect other memory files (MEMORY.md, daily logs).

---

**Remember**: The goal is continuous improvement, not perfection. Start small, build your pattern library over time, and watch performance improve!
