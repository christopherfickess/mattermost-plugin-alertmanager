<!-- PR template. Keep it tight — reviewers will skip long preambles.
     Delete sections that don't apply rather than writing "N/A". -->

## Summary

<!-- 1-3 bullets. What changed and why. -->

## Test plan

<!-- Concrete checklist of what you verified. Reviewers will look for:
     - lint + tests green
     - security pipeline green (SBOM/CodeQL)
     - if behavior changed, what scenario was exercised
     Replace these bullets with what you actually did. -->

- [ ] `make check-style` clean
- [ ] `make test` green
- [ ] `make sbom-audit` clean
- [ ] Manually exercised the changed behavior (describe how)

## Risk + rollback

<!-- One sentence: what's the worst-case if this lands broken, and
     how do we roll back. Skip if change is pure refactor or docs. -->

## Related issues

<!-- Closes #N, refs #M -->
