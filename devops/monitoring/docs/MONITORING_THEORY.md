# Monitoring & Observability Theory Guide

## 1. Observability Fundamentals

### 1.1 Observability vs Monitoring
```
┌─────────────────┐    ┌─────────────────┐
│   Monitoring    │    │  Observability  │
├─────────────────┤    ├─────────────────┤
│ • Known issues  │    │ • Unknown issues│
│ • Predefined    │    │ • Exploratory   │
│   dashboards    │    │   analysis      │
│ • Alerts        │    │ • Root cause    │
│ • Health checks │    │   analysis      │
└─────────────────┘    └─────────────────┘
```

### 1.2 The Three Pillars

#### **Metrics**
- Numerical measurements over time
- Aggregated data points
- Efficient storage and querying
- Examples: CPU usage, request rate, error count

#### **Logs**
- Discrete events with timestamps
- Detailed context information
- Searchable and filterable
- Examples: Application logs, error logs, audit logs

#### **Traces**
- Request flow through distributed systems
- Spans represent operations
- Shows dependencies and latency
- Examples: HTTP request trace, database query trace

## 2. Prometheus Deep Dive

### 2.1 Architecture Components
```
┌─────────────────────────────────────────┐
│              Prometheus                 │
├─────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────────┐ │
│  │   Server    │    │   Alertmanager  │ │
│  │ • Scraper   │    │ • Notifications │ │
│  │ • TSDB      │    │ • Grouping      │ │
│  │ • Query     │    │ • Routing       │ │
│  └─────────────┘    └─────────────────┘ │
└─────────────────────────────────────────┘
           ↑                    ↑
    ┌─────────────┐      ┌─────────────┐
    │  Exporters  │      │   Grafana   │
    │ • Node      │      │ • Dashboards│
    │ • App       │      │ • Alerts    │
    │ • Custom    │      │ • Users     │
    └─────────────┘      └─────────────┘
```

### 2.2 Data Model
```
<metric_name>{<label_name>=<label_value>, ...} value timestamp

# Example
http_requests_total{method="GET", status="200", instance="api-1"} 1234 1609459200
```

### 2.3 Metric Types

#### **Counter**
- Only increases (or resets to zero)
- Use for: requests, errors, tasks completed
```promql
# Rate of HTTP requests per second
rate(http_requests_total[5m])
```

#### **Gauge**
- Can go up or down
- Use for: memory usage, active connections
```promql
# Current memory usage
process_resident_memory_bytes
```

#### **Histogram**
- Samples observations into buckets
- Use for: request durations, response sizes
```promql
# 95th percentile response time
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

#### **Summary**
- Similar to histogram but calculates quantiles
- Use for: request durations with predefined quantiles
```promql
# Median response time
http_request_duration_seconds{quantile="0.5"}
```

### 2.4 PromQL Fundamentals

#### **Selectors**
```promql
# Instant vector
http_requests_total{job="api-server"}

# Range vector
http_requests_total[5m]

# Regex matching
http_requests_total{status=~"5.."}
```

#### **Functions**
```promql
# Rate - per-second rate
rate(http_requests_total[5m])

# Increase - total increase
increase(http_requests_total[1h])

# Avg over time
avg_over_time(cpu_usage[10m])
```

#### **Operators**
```promql
# Arithmetic
cpu_usage * 100

# Comparison
memory_usage > 0.8

# Logical
up == 1 and rate(http_requests_total[5m]) > 10
```

## 3. Logging Architecture

### 3.1 Centralized Logging Flow
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│ Application │───►│   Fluentd   │───►│Elasticsearch│───►│   Kibana    │
│    Logs     │    │ (collect)   │    │  (store)    │    │ (visualize) │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

### 3.2 Log Levels
```
TRACE → DEBUG → INFO → WARN → ERROR → FATAL
  ↑       ↑       ↑      ↑       ↑        ↑
Development    Production    Critical   System
```

### 3.3 Structured Logging
```json
{
  "timestamp": "2024-01-01T10:00:00Z",
  "level": "INFO",
  "service": "user-service",
  "trace_id": "abc123",
  "span_id": "def456",
  "message": "User created successfully",
  "user_id": "12345",
  "duration_ms": 150
}
```

### 3.4 Log Aggregation Patterns
- **Push Model**: Applications push logs to central system
- **Pull Model**: Central system pulls logs from applications
- **Agent-based**: Local agents collect and forward logs

## 4. Distributed Tracing

### 4.1 Trace Components
```
Trace
├── Span A (Root)
│   ├── Span B (Child)
│   │   └── Span D (Child)
│   └── Span C (Child)
│       └── Span E (Child)
└── Span F (Sibling)
```

### 4.2 Span Attributes
```json
{
  "trace_id": "trace-123",
  "span_id": "span-456", 
  "parent_span_id": "span-789",
  "operation_name": "http_get_user",
  "start_time": "2024-01-01T10:00:00Z",
  "end_time": "2024-01-01T10:00:01Z",
  "duration_ms": 1000,
  "tags": {
    "http.method": "GET",
    "http.url": "/users/123",
    "http.status_code": 200,
    "user.id": "123"
  }
}
```

### 4.3 Sampling Strategies
- **Probabilistic**: Sample X% of traces
- **Rate Limiting**: Sample max N traces per second
- **Adaptive**: Adjust sampling based on system load

## 5. Application Monitoring

### 5.1 Golden Signals (SRE)
```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Latency   │  │   Traffic   │  │   Errors    │  │ Saturation  │
├─────────────┤  ├─────────────┤  ├─────────────┤  ├─────────────┤
│ • Response  │  │ • Request   │  │ • Error     │  │ • CPU/Memory│
│   time      │  │   rate      │  │   rate      │  │   usage     │
│ • P50, P95, │  │ • RPS/QPS   │  │ • HTTP 5xx  │  │ • Queue     │
│   P99       │  │ • Bandwidth │  │ • Timeouts  │  │   length    │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
```

### 5.2 RED Method (Request Rate, Errors, Duration)
```promql
# Request Rate
sum(rate(http_requests_total[5m])) by (service)

# Error Rate
sum(rate(http_requests_total{status=~"5.."}[5m])) by (service) /
sum(rate(http_requests_total[5m])) by (service)

# Duration (P95)
histogram_quantile(0.95, 
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le, service)
)
```

### 5.3 USE Method (Utilization, Saturation, Errors)
```promql
# CPU Utilization
100 - (avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

# Memory Saturation
(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / 
node_memory_MemTotal_bytes

# Disk Errors
rate(node_disk_io_errors_total[5m])
```

## 6. Infrastructure Monitoring

### 6.1 Node Monitoring
```promql
# CPU Usage
100 - (avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

# Memory Usage
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100

# Disk Usage
(1 - (node_filesystem_avail_bytes / node_filesystem_size_bytes)) * 100

# Network Traffic
rate(node_network_receive_bytes_total[5m])
```

### 6.2 Kubernetes Monitoring
```promql
# Pod CPU Usage
sum(rate(container_cpu_usage_seconds_total[5m])) by (pod)

# Pod Memory Usage
sum(container_memory_working_set_bytes) by (pod)

# Pod Restart Count
increase(kube_pod_container_status_restarts_total[1h])

# Node Status
up{job="kubernetes-nodes"}
```

### 6.3 Application Monitoring
```promql
# HTTP Request Rate
sum(rate(http_requests_total[5m])) by (method, endpoint)

# Error Rate
sum(rate(http_requests_total{status=~"5.."}[5m])) /
sum(rate(http_requests_total[5m]))

# Response Time
histogram_quantile(0.95, 
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
)

# Database Connection Pool
database_connections_active / database_connections_max
```

## 7. Alerting Strategy

### 7.1 Alert Severity Levels
```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  CRITICAL   │  │   WARNING   │  │    INFO     │  │    DEBUG    │
├─────────────┤  ├─────────────┤  ├─────────────┤  ├─────────────┤
│ • Page      │  │ • Email     │  │ • Slack     │  │ • Log only  │
│   immediately│  │ • Slack     │  │ • Dashboard │  │ • Metrics   │
│ • Service   │  │ • Ticket    │  │ • No action │  │   only      │
│   down      │  │ • Trend     │  │   required  │  │             │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
```

### 7.2 Alert Rules Best Practices
```yaml
# Good Alert Rule
- alert: HighErrorRate
  expr: |
    sum(rate(http_requests_total{status=~"5.."}[5m])) /
    sum(rate(http_requests_total[5m])) > 0.1
  for: 5m  # Wait 5 minutes before firing
  annotations:
    summary: "High error rate detected"
    description: "Error rate is {{ $value | humanizePercentage }}"
```

### 7.3 Alert Fatigue Prevention
- **Meaningful Alerts**: Only alert on actionable issues
- **Proper Thresholds**: Avoid too sensitive alerts
- **Grouping**: Group related alerts
- **Escalation**: Implement escalation policies

## 8. Dashboard Design

### 8.1 Dashboard Hierarchy
```
┌─────────────────────────────────────┐
│         Executive Dashboard         │
│ • Business KPIs                     │
│ • High-level health                 │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│        Operational Dashboard        │
│ • Service health                    │
│ • Golden signals                    │
└─────────────────────────────────────┘
                  ↓
┌─────────────────────────────────────┐
│        Technical Dashboard          │
│ • Detailed metrics                  │
│ • Troubleshooting data             │
└─────────────────────────────────────┘
```

### 8.2 Effective Visualization
- **USE principle**: Understand, Select, Execute
- **Progressive disclosure**: Start high-level, drill down
- **Consistent color coding**: Red=bad, Green=good
- **Appropriate chart types**: Time series, gauges, tables

## 9. SLI/SLO/SLA Framework

### 9.1 Definitions
```
SLI (Service Level Indicator)
├── Metrics that matter to users
├── Examples: latency, availability, throughput
└── Should be measurable and meaningful

SLO (Service Level Objective)  
├── Target value for SLI
├── Examples: 99.9% availability, <100ms latency
└── Internal goal for service quality

SLA (Service Level Agreement)
├── Contract with consequences
├── Examples: 99.5% uptime or credit
└── External commitment to customers
```

### 9.2 Error Budget
```
Error Budget = 1 - SLO

If SLO = 99.9% availability
Then Error Budget = 0.1% = 43.8 minutes/month

Error Budget Policy:
• > 50% remaining: Ship new features
• < 50% remaining: Focus on reliability
• 0% remaining: Feature freeze
```

## 10. Performance Monitoring

### 10.1 Application Performance
```promql
# Apdex Score (Application Performance Index)
(
  sum(rate(http_request_duration_seconds_bucket{le="0.1"}[5m])) +
  sum(rate(http_request_duration_seconds_bucket{le="0.4"}[5m])) * 0.5
) / sum(rate(http_request_duration_seconds_count[5m]))
```

### 10.2 Database Performance
```promql
# Query Duration
histogram_quantile(0.95, 
  sum(rate(mysql_query_duration_seconds_bucket[5m])) by (le)
)

# Connection Pool Usage
mysql_connections_active / mysql_connections_max

# Slow Queries
rate(mysql_slow_queries_total[5m])
```

### 10.3 Cache Performance
```promql
# Cache Hit Rate
redis_keyspace_hits_total / 
(redis_keyspace_hits_total + redis_keyspace_misses_total)

# Cache Memory Usage
redis_memory_used_bytes / redis_memory_max_bytes
```

## 11. Capacity Planning

### 11.1 Growth Prediction
```promql
# Linear prediction for next 4 hours
predict_linear(http_requests_total[1h], 4*3600)

# Capacity based on current growth
(
  predict_linear(cpu_usage[1h], 24*3600) > 0.8
) * on() group_left(instance) (cpu_usage)
```

### 11.2 Resource Utilization Trends
- **CPU trends**: Identify peak hours, growth patterns
- **Memory trends**: Detect memory leaks, allocation patterns
- **Storage trends**: Plan for disk space expansion
- **Network trends**: Bandwidth planning

## 12. Troubleshooting Methodology

### 12.1 Incident Response Process
```
1. Detection
   ├── Monitoring alerts
   ├── User reports
   └── Proactive monitoring

2. Triage
   ├── Assess impact
   ├── Assign severity
   └── Initial response

3. Investigation
   ├── Gather data
   ├── Form hypotheses
   └── Test theories

4. Resolution
   ├── Apply fix
   ├── Verify resolution
   └── Monitor stability

5. Post-mortem
   ├── Root cause analysis
   ├── Process improvements
   └── Prevention measures
```

### 12.2 Debugging with Metrics
```promql
# Service dependency mapping
up{job=~".*service.*"}

# Resource constraints
(rate(container_cpu_cfs_throttled_periods_total[5m]) / 
 rate(container_cpu_cfs_periods_total[5m])) > 0.25

# Memory pressure
container_memory_working_set_bytes / 
container_spec_memory_limit_bytes > 0.8
```

## Next Steps

Với kiến thức này, bạn có thể:
1. Thiết kế monitoring strategy phù hợp
2. Implement metrics, logging, và tracing
3. Tạo meaningful dashboards và alerts
4. Troubleshoot issues hiệu quả
5. Optimize application performance

Bạn muốn tôi giải thích chi tiết hơn về phần nào không? 