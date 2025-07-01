# Monitoring & Alerting for Key-Server

To keep our **key-server** reliable and performant, we will use proven monitoring and alerting methods — **USE**, **RED**, and **Golden Signals**. Here’s what we can do:

---

## 1. Monitoring Approach

| Method             | Metrics / Focus                            | What to Monitor & Why                                                                                                       |
| ------------------ | ------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------- |
| **USE**            | **Utilization, Saturation, Errors**        | CPU, memory usage, network saturation, and system-level errors. Helps identify resource bottlenecks and failures.           |
| **RED**            | **Rate, Errors, Duration**                 | Request per second, error rate (4xx/5xx), and latency (p50, p90, p99). Helps understand API behavior and client experience. |
| **Golden Signals** | **Availability, Latency, Traffic, Errors** | Key reliability signals: success ratio, request volume, latency, and error counts.                                          |

---

## 2. Grafana Dashboard Setup

| Panel Category        | Metrics / Panels                                                                                   | Purpose                                                                               |
| --------------------- | -------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| **USE Metrics**       | - CPU Utilization per pod<br>- Memory Utilization<br>- Network I/O Saturation<br>- 5xx Error count | Monitor pod resource usage and saturation to detect bottlenecks or failures.          |
| **RED Metrics**       | - Request Rate (RPS)<br>- Error Rate (%)<br>- Request Duration (p50, p90, p99 latency)             | Understand traffic patterns, error trends, and latency to ensure smooth API behavior. |
| **Golden Signals**    | - Availability (% of 200 OK)<br>- Total Traffic (RPS)<br>- Latency P99<br>- Error Count (5xx)      | Track overall service health and user experience.                                     |
| **Additional Panels** | - Active Connections Gauge<br>- Key Length Request Histogram                                       | Detect connection pressure and possible abuse patterns (large key requests).          |

---

## 3. Alerting Strategy

| Alert Name                  | Condition / Expression                            | Severity | Purpose                                           |
| --------------------------- | ------------------------------------------------- | -------- | ------------------------------------------------- |
| **High Error Rate**         | 5xx error rate > 5% over 5 minutes                | Critical | Detects server errors affecting availability.     |
| **High Latency**            | p99 request duration > 1 second                   | Warning  | Flags slow requests degrading user experience.    |
| **Low Availability**        | Successful requests < 95% over 5 minutes          | Critical | Indicates service health issues or downtime.      |
| **High CPU/Memory Usage**   | CPU or Memory > 90% per pod                       | Warning  | Prevents resource exhaustion and pod failures.    |
| **High Request Rate (RPS)** | Total RPS > 1000 sustained for 2 minutes          | Critical | Detects potential abuse or DoS attacks.           |
| **Key Length Abuse**        | Spike in large key requests near max allowed size | Warning  | Flags possible resource abuse from expensive ops. |

---

## 4. Configs

- Dashboard JSON file under `grafana/`
- Alert rules YAML under `prometheus/`
