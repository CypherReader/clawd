# Coding Patterns

Successful software development approaches learned from experience.

## ✅ Success Patterns

### Framework-First Pattern

**When**: Building systems with multiple similar components (strategies, plugins, modules)

**Approach**:
1. Define interface/contract first
2. Write framework/manager for the interface
3. Implement first component fully (with tests)
4. Add remaining components incrementally

**Benefits**:
- New components are trivial to add
- Consistent behavior across components
- Easy to test and validate
- Clear contracts reduce integration bugs

**Example**: Oracle strategy system
- Strategy interface defined first
- Manager handles all strategies uniformly
- Each strategy implements same contract
- Adding strategy #6 would take <1 hour

**Confidence**: HIGH

---

### Test-As-Documentation Pattern

**When**: Building libraries, frameworks, or complex systems

**Approach**:
1. Write tests that show how to use the code
2. Test names describe behavior clearly
3. Test cases cover common use cases
4. Tests serve as executable examples

**Benefits**:
- Tests don't just validate, they explain
- New developers learn by reading tests
- Refactoring safe (tests catch breaks)
- Documentation always up-to-date

**Example**: Oracle execution tests
```go
func TestRiskManager_ValidateTrade()
func TestOrchestrator_ExecuteArbitrage_Success()
func TestOrchestrator_ExecuteArbitrage_Leg2Failure()
```

Each test name tells a story, code shows how to use it.

**Confidence**: HIGH

---

### Safety-First Pattern

**When**: Building systems that handle money, data, or critical operations

**Approach**:
1. Build safety controls BEFORE execution logic
2. Risk management is not optional
3. Emergency stops are mandatory
4. Validate before acting, rollback on errors

**Benefits**:
- User confidence from day one
- Prevents costly mistakes
- Easier to get approval for deployment
- Less fear, more trust

**Example**: Oracle risk manager built before orchestrator
- Position limits enforced
- Daily loss caps active
- Emergency stop works
- Can't execute without validation

**Confidence**: CRITICAL - Never skip this

---

### Incremental Delivery Pattern

**When**: Long projects that could benefit from early feedback

**Approach**:
1. Define minimal viable functionality
2. Build and deliver that first
3. Get feedback, iterate
4. Add features incrementally

**Benefits**:
- User sees progress quickly
- Can change direction based on feedback
- Each increment adds value
- Can stop at any point with working system

**Example**: Oracle strategies
- Could have stopped after strategy 1 (had working system)
- Each strategy added more value
- User could have tested early
- No big-bang risk

**Confidence**: HIGH

---

### Document-As-You-Build Pattern

**When**: Building anything that will be deployed or used by others

**Approach**:
1. Write README early (describe what it will do)
2. Update docs as you build features
3. Write deployment guide alongside code
4. Don't leave docs for "later"

**Benefits**:
- Documentation is better (context fresh)
- Forces clarity of thought
- No rush at the end
- Users can deploy immediately when done

**Example**: Oracle
- README updated throughout
- Docker guide written with Docker setup
- Each strategy has inline docs
- Deployment ready when code done

**Confidence**: HIGH

---

## ⚠️ Anti-Patterns (Avoid These)

### Big-Bang Integration

**Problem**: Build everything separately, integrate at end
**Why Bad**: Integration issues found too late
**Instead**: Integrate continuously, test often

### Docs-As-Afterthought

**Problem**: "We'll document it later"
**Why Bad**: Never happens, or rushed and poor quality
**Instead**: Document as you build

### Skip-The-Tests

**Problem**: "We'll add tests later" or "It's simple, doesn't need tests"
**Why Bad**: Regressions, refactoring fear, no examples
**Instead**: Test first or test immediately after

### Feature-Before-Safety

**Problem**: Build cool features before safety controls
**Why Bad**: Can't deploy safely, risk of harm
**Instead**: Safety first, features second

---

## Context-Specific Patterns

### For Financial Systems
1. Paper trading mode mandatory
2. Risk limits before execution
3. Emergency stop always available
4. Audit trail for all trades
5. Start with small limits, scale up

### For Multi-Strategy Systems
1. Framework before strategies
2. Each strategy independent and testable
3. Manager handles orchestration
4. Easy to enable/disable strategies
5. Per-strategy performance tracking

### For Long Development Sessions
1. Monitor budget/context usage
2. Suggest breaks before limits
3. Plan multi-session splits
4. Save progress frequently
5. Document handoff points

---

Last Updated: 2026-03-04
