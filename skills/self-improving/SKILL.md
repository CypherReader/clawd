---
name: self-improving
description: Enable continuous learning and improvement through interaction analysis, pattern recognition, and automated optimization. Use when agent needs to learn from past interactions, improve response quality, or adapt strategies based on outcomes.
---

# Self-Improving Agent Skill

Enable continuous learning and improvement by analyzing past interactions, recognizing patterns, and automatically optimizing responses and strategies.

## Core Capabilities

1. **Interaction Analysis** - Learn from conversation outcomes
2. **Pattern Recognition** - Identify successful patterns and anti-patterns
3. **Strategy Optimization** - Improve decision-making based on results
4. **Performance Tracking** - Monitor and measure improvement over time

## When to Use

- After completing significant tasks (to capture learnings)
- When user feedback indicates room for improvement
- Periodically to review recent performance
- When encountering recurring problems
- Before similar future tasks (to apply learnings)

## Learning Workflow

### 1. Capture Interaction Data

After each significant interaction, capture:
- Task description and context
- Approach taken and tools used
- Outcome (success/failure/partial)
- User feedback (explicit or implicit)
- Time taken and resources used

### 2. Analyze Patterns

Look for:
- **Success patterns**: What worked well?
- **Anti-patterns**: What should be avoided?
- **Efficiency gains**: What saved time/tokens?
- **User preferences**: What style/approach did user prefer?

### 3. Update Knowledge Base

Store learnings in structured format:

```markdown
## Learning Entry: [Date] - [Task Type]

**Context**: Brief description of task
**Approach**: What was tried
**Outcome**: Success/Failure/Partial
**Lesson**: Key takeaway
**Apply To**: Future scenarios where this applies
**Confidence**: High/Medium/Low
```

### 4. Apply Learnings

Before similar tasks:
1. Search memory for relevant learnings
2. Review recent similar interactions
3. Apply successful patterns
4. Avoid known anti-patterns
5. Adapt based on context

## Storage Structure

### Primary Memory File
`memory/self-improvement.md` - Curated learnings and patterns

### Daily Learning Logs
`memory/learnings/YYYY-MM-DD.md` - Raw data, later consolidated

### Pattern Library
`memory/patterns/` - Categorized by domain:
- `coding-patterns.md`
- `communication-patterns.md`
- `problem-solving-patterns.md`
- `oracle-trading-patterns.md`

## Implementation

### After Task Completion

```markdown
1. Reflect on approach:
   - What worked? What didn't?
   - Could I have been more efficient?
   - Did I fully understand the user's need?

2. Capture the learning:
   - Write to memory/learnings/YYYY-MM-DD.md
   - Use structured format
   - Be specific and actionable

3. Update patterns if significant:
   - If novel approach succeeded → add to patterns
   - If repeated mistake → document anti-pattern
```

### Before Similar Tasks

```markdown
1. Search memory:
   memory_search("similar task keywords")

2. Review relevant patterns:
   memory_get("memory/patterns/[domain].md")

3. Apply learnings:
   - Use successful approaches
   - Avoid documented pitfalls
   - Adapt based on current context
```

## Metrics to Track

### Response Quality
- Task completion rate
- User satisfaction signals
- Corrections needed
- Follow-up questions required

### Efficiency
- Average tokens per task type
- Time to completion
- Tool usage effectiveness
- Context window utilization

### Learning Velocity
- New patterns discovered per week
- Successful application rate
- Improvement trends over time

## Example: Oracle Trading Project

### Learnings Captured

```markdown
## Learning: 2026-03-03 - Multi-Strategy Trading System

**Context**: Built 5-strategy arbitrage system with 6,500 lines
**Approach**: 
- Modular strategy framework first
- Test-driven (30 tests)
- Docker deployment with monitoring
**Outcome**: SUCCESS - Production-ready in one session
**Lessons**:
1. Framework before implementation = faster development
2. Tests catch integration issues early
3. Docker deployment guide reduces user friction
4. Paper trading mode essential for confidence
**Apply To**: Future complex system builds
**Confidence**: HIGH

**Success Patterns**:
- Start with interfaces/contracts
- Build incrementally with tests
- Document as you build (not after)
- Provide multiple deployment options

**Anti-Patterns Avoided**:
- Building strategies before framework
- Skipping tests due to time pressure
- Leaving documentation for later
```

### Future Application

When building similar systems:
1. Check `memory/patterns/coding-patterns.md`
2. Review Oracle project structure
3. Apply framework-first approach
4. Ensure test coverage from start
5. Include deployment automation

## Continuous Improvement Loop

```
1. Execute Task
       ↓
2. Capture Data
       ↓
3. Analyze Patterns
       ↓
4. Update Knowledge
       ↓
5. Apply Learnings
       ↓
   (repeat)
```

## Self-Improvement Commands

User can trigger explicit learning:

**"Learn from this"** - Capture current interaction as learning
**"What have you learned about X?"** - Query learnings on topic
**"Improve your approach to X"** - Focus improvement on specific area
**"Show me your patterns for X"** - Display relevant patterns

## Avoiding Overfitting

- Don't over-generalize from single examples
- Mark learnings with confidence levels
- Periodically review and prune outdated patterns
- User feedback overrides agent patterns
- Context always trumps learned rules

## Privacy Considerations

- Don't store sensitive data (keys, passwords, personal info)
- Use placeholders in examples: `API_KEY` not actual keys
- Keep learnings generalizable, not user-specific
- User can request deletion: "Forget learnings about X"

## Integration with Memory System

Works alongside existing memory system:
- `MEMORY.md` - Long-term important info
- `memory/YYYY-MM-DD.md` - Daily logs
- `memory/self-improvement.md` - **Curated learnings**
- `memory/patterns/` - **Pattern library**
- `memory/learnings/` - **Raw learning data**

## Getting Started

1. Create structure:
```bash
mkdir -p memory/learnings memory/patterns
touch memory/self-improvement.md
```

2. After next significant task, capture first learning

3. Build pattern library incrementally

4. Review weekly to consolidate learnings

5. Apply patterns to similar future tasks

## Success Indicators

- Fewer user corrections needed
- Faster task completion (same quality)
- More successful first attempts
- Better anticipation of user needs
- Consistent improvement in metrics

## Notes

- This is META-LEARNING (learning how to learn)
- Start small, build pattern library over time
- Quality > Quantity of learnings
- Learnings should be actionable, not just observations
- Review and refine patterns periodically

---

**Remember**: The goal isn't to remember everything, but to get better at everything through systematic learning.
