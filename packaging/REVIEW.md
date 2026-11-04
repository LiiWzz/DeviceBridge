# DB-011 independent review

Two fresh read-only Orca/Codex reviewers independently examined implementation
`474922df55e2c62ee96173b96a87ff96d4ea7b12` against fixed point
`c8fae1481fe657f385d7747878780d1f4485e986`. Standards and spec remained separate.
Coordinator owned fixes and final evidence; reviewers did not edit, commit, push,
approve or merge. This is source/spec review, not a production audit.

## Standards

> Completed the read-only Standards review of commit 474922d against resolved fixed point c8fae148, inspected the nonempty three-dot diff and commit list, read README.md, AGENTS.md, docs/M01_EXECUTION.md and ADR-001, and examined the supplied host/Linux/BusyBox/ARM preview logs without rerunning broad tests or modifying files. Found no documented standards violations and no actionable baseline smell findings: the build-host tool streams artifacts, the installer preserves offline/unprivileged operation with integrity and destination checks, and the documentation states compatibility, crash-cleanup and real-hardware limits without claiming production transport. No Standards fixes remain from this review; the coordinator retains ownership of final evidence/artifact additions, PR workflow and actual ARM hardware validation.

## Spec

> I completed the read-only Spec review of commit 474922d against c8fae148, fetched issue #16 as the primary specification, inspected the required repository documents and supplied evidence, and left the checkout unchanged. One P2 finding at packaging/install.sh:205-206 (also :212-213): issue #16 requires “atomic activation with rollback”, but updating previous before activating the new binary destroys the rollback target when activation fails; on unprivileged Darwin ARM64 I installed versions 1 and 2, made bin non-writable with chmod 500, attempted version 3, observed exit 1 with Permission denied, restored bin permissions, and found active=2/previous=2 with --rollback failing “no different previous release available”. The coordinator should preserve or restore the old previous link on activation failure and add a focused failed-activation regression test; no other actionable Spec findings were identified, broad tests were not rerun, and final evidence/PR work remains coordinator-owned.

### Fix verification

Coordinator commit `f9784d44a36bed1feb932115fb77de0b496adc1d` activates the
executable before history and restores it if history replacement fails. Focused
tests reproduce an unwritable binary directory and injected history-rename
failure, verify both old links survive, then perform a successful rollback. Both
passed on unprivileged macOS and Debian dash. BusyBox passes the permission case;
PATH-shim injection is unavailable with its built-in mv and explicitly skipped.

The Spec reviewer examined the fix in a follow-up dispatch:

> I reviewed fix commit f9784d4 against 474922d, inspected the current temporary source and TestFailedActivationPreservesRollback, and made no edits or test runs. The original P2 is fixed by inspection: packaging/install.sh:159 leaves history untouched when executable activation fails, lines 161-163 restore the original executable when history replacement fails, and both upgrade and rollback use this helper; the regression exercises unwritable bin and injected history-rename failure, checks both original links, and verifies rollback remains usable. I found no remaining concrete normal-failure defect in this change; final package regeneration, full test evidence and PR ownership remain with the coordinator, with power-loss durability outside the claimed scope.

Coordinator commit `c9fcf1a1b2531e3e81323503a11a65bcc06bd6f8` also refreshes
the installed manager on upgrade, so fixes reach existing installations. A
regression simulates an obsolete exit-47 manager and checks working status after
upgrade. The Spec reviewer separately verified that delta:

> I reviewed the f9784d4..c9fcf1a installer and test delta through shared Git objects, including surrounding validation and activation order, without editing files or running tests. No concrete remaining issue was found: refresh_manager runs after payload and installed-state validation under the installation lock, replaces only a differing regular manager through a verified staged copy and rename, and the obsolete-manager regression confirms installed status works after upgrade; the PATH-shim skip is limited to unavailable failure injection. Final full-suite evidence, regenerated preview packages and PR work remain coordinator-owned, and the explicitly allowed manager refresh on failed binary activation does not reopen the rollback-link finding.

Final host/Linux suites, race checks, ShellCheck and actual archive smokes ran
again after these fixes; [EVIDENCE.md](EVIDENCE.md) links actual output. Final
`0.1.0-preview.3` archives supersede internal preview.1/preview.2 builds. The
coordinator added final documentation/evidence after review, without changing
production installer/builder code.

Standards: 0 findings (worst: none). Spec: 1 P2 finding, fixed and rechecked; 0 unresolved findings.
