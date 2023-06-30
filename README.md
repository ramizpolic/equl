# equl
equl is a dynamic and performant rule-based object equality comparator best suited for checks of mutating objects often
found in Kubernetes environments.

NOTE: Work in progress

## Usage
```go


```

### TODO
- [ ] Expose better interface for equality check (which also includes custom filtration function)
- [ ] Improve performances and add benchmarking stats
- [ ] Write more complete test scenarios
- [ ] Update docs

### Algorithm -- Selective object comparator
For non-transformative and selective approach to comparison, maps that contain fields with non-zero default values,
or for granular comparison.
Think: _selective, exact, close, or not really_.
Example: Structure that contains a lot of fields, but we are interested in specific ones to avoid dynamic injection checks.

Extends https://github.com/cisco-open/k8s-objectmatcher

- Let _expected_ and _actual_ be two objects which are marshall-able or reflect-able
- Recurse through the objects to create field map with names and types
- Compare to Options by depth and type, prioritize Exclusion by default, but offer options for filtering
- Expose following options for comparison:
    - `WithFields(...string)` => e.g. `WithFields(".Spec.[]Ports.*", ".Spec.Type")`
    - `WithoutFields(...string)` => e.g. `WithoutFields(".Spec.[]Ports.NodePort", ".Spec.IPFamilies")`
- Invoked as `Compare(expected, actual, ...Options)` and implements:
    - `Equal() bool` - returns if they are equal
    - `Error() error` - helps catch runtime errors, but pririotized through: base types, transforms, option misses and paradox cases.
