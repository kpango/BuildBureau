# Example: Self-Hosting Bootstrap Task

## Scenario

Using BuildBureau to add a new feature to itself.

## Task

"Add a caching layer to the LLM manager to avoid redundant API calls for
identical prompts."

## Expected Flow

### 1. President Analysis

- Understands the request is about optimizing LLM API usage
- Clarifies caching requirements (TTL, size limits, invalidation)
- Defines objectives:
  - Reduce API costs
  - Improve response time for repeated queries
  - Maintain accuracy (cache invalidation)

### 2. Secretary Coordination

- Records the optimization task
- Identifies affected component: `internal/llm/manager.go`
- Delegates to Director for architectural planning

### 3. Director Decomposition

- Researches caching strategies for LLM responses
- Proposes design:
  - In-memory cache with LRU eviction
  - Cache key: hash of (provider, model, prompt, options)
  - Configurable TTL and max size
- Delegates to Manager for detailed design

### 4. Manager Design

- Creates specification:

  ```go
  type LLMCache struct {
      cache    *lru.Cache
      ttl      time.Duration
      maxSize  int
      mutex    sync.RWMutex
  }

  func (m *Manager) Generate(ctx, provider, prompt, opts) (string, error) {
      // Check cache
      if cached := m.cache.Get(cacheKey); cached != nil {
          return cached.response, nil
      }

      // Generate if not cached
      response, err := m.generateUncached(...)

      // Store in cache
      m.cache.Set(cacheKey, response, m.ttl)

      return response, err
  }
  ```

- Specifies tests needed
- Delegates to Engineer for implementation

### 5. Engineer Implementation

- Implements the caching layer:
  - Adds `golang-lru` dependency
  - Implements cache structure
  - Updates `Manager.Generate()` method
  - Adds cache statistics methods
- Writes tests:
  - Test cache hits
  - Test cache misses
  - Test TTL expiration
  - Test max size eviction
- Updates configuration:
  ```yaml
  llms:
    cache:
      enabled: true
      ttl: 300s # 5 minutes
      max_size: 1000
  ```
- Updates documentation in `docs/MULTI_PROVIDER.md`

## Generated Changes

### Files Modified

- `internal/llm/manager.go` - Add caching logic
- `internal/llm/cache.go` - New cache implementation
- `internal/llm/cache_test.go` - Tests for cache
- `pkg/types/config.go` - Add cache configuration
- `internal/config/loader.go` - Load cache config
- `docs/MULTI_PROVIDER.md` - Document caching

### Test Output

```
=== RUN   TestLLMCache
=== RUN   TestLLMCache/CacheHit
=== RUN   TestLLMCache/CacheMiss
=== RUN   TestLLMCache/TTLExpiration
=== RUN   TestLLMCache/LRUEviction
--- PASS: TestLLMCache (0.03s)
    --- PASS: TestLLMCache/CacheHit (0.01s)
    --- PASS: TestLLMCache/CacheMiss (0.01s)
    --- PASS: TestLLMCache/TTLExpiration (0.01s)
    --- PASS: TestLLMCache/LRUEviction (0.00s)
PASS
```

## Human Review

Review the generated changes:

```bash
git diff
```

Run tests:

```bash
make test
```

Try it out:

```bash
make build
./build/buildbureau
# Make the same query twice - second should be instant
```

## Approval & Merge

If satisfied:

```bash
git add .
git commit -m "Add LLM response caching (self-implemented by BuildBureau)"
git push
```

## Learning Captured

BuildBureau's memory system stores:

- The caching pattern used
- Performance improvements observed
- Design decisions and rationale
- Test patterns that worked

Next time a similar optimization is needed, agents will reference this
implementation.

## Metrics

- **API Call Reduction**: 60% reduction in duplicate queries
- **Response Time**: 95% faster for cached responses
- **Implementation Time**: ~10 minutes (vs hours manually)
- **Code Quality**: Passes all linters, follows patterns
- **Test Coverage**: 95% of new code

## Meta-Learning

This bootstrap task itself demonstrates:

- BuildBureau understanding its own architecture
- Self-aware code generation
- Pattern recognition and reuse
- Recursive improvement capability

---

**Note**: This is an example flow. Actual implementation details may vary based
on agent decisions and LLM responses.
