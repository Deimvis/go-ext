# xrtpoint

* xrtpoint (Runtime Point) — package that helps with injecting points changing runtime flow.
* Runtime Points do nothing when built without `debug` tag.
* The package is similar to [failpoint](https://github.com/pingcap/failpoint) and [gofail](https://github.com/etcd-io/gofail?tab=readme-ov-file), but behaves diffirently.
  We do not inject code on build time allowing to fully customize points in runtime.
  And we offer control exclusively over runtime behaviour, forbidding any other changes as much as possible.
* xrtpoint imposes small overhead for release builds, so it's not zero-cost, but overhead is negligible (push/pop latency on call stack)
* You may want to use this package if you want to:
  * Inject panic (crash) in arbitrary line of code (e.g. in order to test consistency guarantees)
  * Inject wait in arbitrary line of code (e.g. in order to reorder operations in the desired manner)
  * Test that particular line of code is reachable or unreachable
  * TODO: Inject scheduler yield
* You should not use this package if you need to:
  * Modify function/method return values — use mocking
  * Modify state (local or global variables) in arbitrary line of code
