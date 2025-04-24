#  (WIP) Snowflake ID

[![Go Report Card](https://goreportcard.com/badge/github.com/bymerk/snowflake)](https://goreportcard.com/report/github.com/bymerk/snowflake)
![Go Version](https://img.shields.io/github/go-mod/go-version/bymerk/snowflake)
![License](https://img.shields.io/github/license/bymerk/snowflake)

| 1 bit | 41 bits          | 6 bits      | 4 bits     | 12 bits       |
|-------|------------------|-------------|------------|----------------|
| unused (sign) | timestamp (ms since custom epoch) | cluster ID | node ID | sequence number |

---

## 🧩 Field Descriptions

### • 1 bit — Unused (Sign Bit)
- Always `0`.
- Ensures the ID remains positive when stored as a signed `int64`.

### • 41 bits — Timestamp
- Represents the number of **milliseconds since a custom epoch**.
- **Default epoch is set to `2025-01-01 00:00:00 UTC`**, but it can be configured.
- Provides chronological sorting.
- Supports ~69 years of IDs.

### • 6 bits — Cluster ID
- Identifies the **data center** or **Kubernetes cluster region**.
- Supports up to `2^6 = 64` clusters.

### • 4 bits — Node ID
- Identifies a specific **machine, pod, or instance** within a cluster.
- Supports up to `2^4 = 16` nodes per cluster.

### • 12 bits — Sequence Number
- Used to ensure uniqueness when multiple IDs are generated in the **same millisecond**.
- Resets every millisecond.
- Supports up to `2^12 = 4096` IDs per node per millisecond.

---

## 💡 Example

If we generate an ID at time `t = 1735678800000` (in ms),
with `clusterID = 5`, `nodeID = 3`, and `sequence = 17`,
and using the default epoch `2025-01-01`,
the Snowflake ID is constructed as:

```go
id := ((t - epoch) << 22) | (clusterID << 16) | (nodeID << 12) | sequence
```

## 📦 Docker Image

The Docker image is available on Docker Hub:

```bash
docker pull bymerk/snowflake

