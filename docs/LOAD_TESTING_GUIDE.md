# Load Testing Guide - E-commerce System

This guide explains how to perform load testing on the e-commerce recommendation system using Locust.

## Table of Contents
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Test Scenarios](#test-scenarios)
- [Understanding Results](#understanding-results)
- [Performance Metrics](#performance-metrics)
- [Troubleshooting](#troubleshooting)

## Overview

Load testing simulates realistic user behavior to evaluate system performance under various load conditions. Our Locust implementation tests:

- User authentication (registration, login)
- Product browsing and searching
- Category navigation
- Product interactions (views, likes, purchases)
- Recommendation system
- Profile management

## Prerequisites

### 1. Install Dependencies

```bash
pip install -r requirements-loadtest.txt
```

This installs:
- `locust==2.20.0` - Load testing framework
- `flask==3.0.0` - Web server for Locust UI
- `gevent==23.9.1` - Async I/O library

### 2. Start the Server

```bash
make run
```

Ensure the server is running on `http://localhost:8080`

## Quick Start

### Method 1: Using the Run Script (Recommended)

```bash
./run_load_test.sh
```

Select from 6 test scenarios:
1. Web UI Mode (Interactive)
2. Light Load (10 users)
3. Medium Load (50 users)
4. Heavy Load (100 users)
5. Stress Test (200 users)
6. Custom Configuration

### Method 2: Manual Locust Commands

#### Web UI Mode
```bash
locust --host=http://localhost:8080
```
Open browser at `http://localhost:8089`

#### Headless Mode
```bash
locust --host=http://localhost:8080 \
       --headless \
       --users 50 \
       --spawn-rate 5 \
       --run-time 3m \
       --html report.html \
       --csv results
```

## Test Scenarios

### 1. Web UI Mode (Interactive)

**When to use:** Development, debugging, real-time monitoring

**Features:**
- Live statistics dashboard
- Dynamic user control
- Real-time charts
- Manual start/stop

**Usage:**
```bash
locust --host=http://localhost:8080
```

Access at: `http://localhost:8089`

### 2. Light Load Test

**Configuration:**
- Users: 10
- Spawn Rate: 2 users/sec
- Duration: 2 minutes

**Purpose:** Baseline performance testing

**Command:**
```bash
locust --host=http://localhost:8080 \
       --headless \
       --users 10 \
       --spawn-rate 2 \
       --run-time 2m \
       --html load_test_light.html \
       --csv load_test_light
```

### 3. Medium Load Test

**Configuration:**
- Users: 50
- Spawn Rate: 5 users/sec
- Duration: 3 minutes

**Purpose:** Typical traffic simulation

**Expected Performance:**
- Response time: < 500ms (p95)
- Error rate: < 1%
- Throughput: 50-100 req/sec

### 4. Heavy Load Test

**Configuration:**
- Users: 100
- Spawn Rate: 10 users/sec
- Duration: 5 minutes

**Purpose:** Peak traffic simulation

**Expected Performance:**
- Response time: < 1000ms (p95)
- Error rate: < 5%
- Throughput: 100-200 req/sec

### 5. Stress Test

**Configuration:**
- Users: 200
- Spawn Rate: 20 users/sec
- Duration: 5 minutes

**Purpose:** System limits discovery

**Goals:**
- Find breaking point
- Identify bottlenecks
- Test error handling

## Understanding Results

### HTML Report

The HTML report includes:

1. **Request Statistics**
   - Request count
   - Failure rate
   - Response times (median, avg, min, max)
   - Requests per second
   - Percentiles (50th, 66th, 75th, 80th, 90th, 95th, 98th, 99th)

2. **Response Time Charts**
   - Response time over time
   - Response time percentiles
   - Users over time

3. **Failure Statistics**
   - Error types
   - Failure count
   - Failure rate

### CSV Files

Generated CSV files:

- `{prefix}_stats.csv` - Request statistics
- `{prefix}_stats_history.csv` - Time-series data
- `{prefix}_failures.csv` - Detailed failure information

### Reading Statistics

```
Type   Name                 # reqs  # fails  Median  Average  Min   Max   RPS
GET    /api/v1/products     1500    10      45      52       12    450   25.0
POST   /api/v1/auth/login   100     0       120     135      85    280   1.67
```

**Key Metrics:**
- **# reqs**: Total requests made
- **# fails**: Number of failed requests
- **Median**: 50th percentile response time
- **Average**: Mean response time
- **RPS**: Requests per second

## Performance Metrics

### Response Time Targets

| Endpoint Type | p50 | p95 | p99 |
|--------------|-----|-----|-----|
| Read (GET) | < 100ms | < 300ms | < 500ms |
| Write (POST/PUT) | < 200ms | < 500ms | < 1000ms |
| Recommendations | < 300ms | < 800ms | < 1500ms |

### Throughput Targets

| Load Level | Concurrent Users | Target RPS |
|-----------|------------------|------------|
| Light | 10 | 20-30 |
| Medium | 50 | 80-120 |
| Heavy | 100 | 150-250 |
| Stress | 200 | 250-400 |

### Error Rate Targets

- **Acceptable:** < 0.1%
- **Warning:** 0.1% - 1%
- **Critical:** > 1%

## User Behavior Simulation

The Locust test simulates realistic e-commerce user behavior:

### Task Flow
1. **Registration/Login** (once on start)
2. **Browse Categories** (weight: 2)
3. **Browse Products** (weight: 3)
4. **View Product Details** (weight: 4)
5. **Record Product View** (weight: 3)
6. **Search Products** (weight: 2)
7. **Filter by Category** (weight: 2)
8. **Filter by Price Range** (weight: 2)
9. **Like Product** (weight: 2)
10. **Purchase Product** (weight: 1)
11. **Get Recommendations** (weight: 2)
12. **Get Similar Users** (weight: 1)
13. **Update Profile** (weight: 1)
14. **Get Interaction History** (weight: 1)
15. **View Profile** (weight: 1)

### Weight Distribution

Tasks are weighted to reflect realistic user behavior:
- **High frequency:** Browsing (weight 3-4)
- **Medium frequency:** Interactions (weight 2-3)
- **Low frequency:** Purchases (weight 1)

### Think Time

Wait time between tasks: 1-3 seconds (simulates reading/thinking time)

## Advanced Configuration

### Custom Locust Configuration

Create/edit `locust.conf`:

```conf
# Web interface
web-host = 0.0.0.0
web-port = 8089

# Target
host = http://localhost:8080

# Logging
loglevel = INFO
logfile = locust.log

# Headless mode options
# headless = true
# users = 50
# spawn-rate = 5
# run-time = 5m

# Output
# html = report.html
# csv = results
```

### Environment Variables

```bash
export HOST=http://localhost:8080
export WEB_PORT=8089
./run_load_test.sh
```

### Custom Test Scenarios

Modify `locustfile.py` to:
- Add new tasks
- Adjust task weights
- Change wait times
- Customize user behavior

Example:
```python
@task(5)  # Higher weight = more frequent
def my_custom_task(self):
    with self.client.get("/api/v1/custom", 
                         catch_response=True) as response:
        if response.status_code == 200:
            response.success()
        else:
            response.failure(f"Failed: {response.status_code}")
```

## Monitoring During Tests

### Real-time Metrics

When using Web UI mode, monitor:

1. **Current RPS** - Should match expected load
2. **Response Times** - Watch for spikes
3. **Failure Rate** - Should stay near 0%
4. **Active Users** - Ramping up correctly

### System Monitoring

Monitor server resources:

```bash
# CPU and Memory
top

# Database connections
docker stats

# Network traffic
netstat -an | grep 8080

# Application logs
tail -f app.log
```

## Interpreting Results

### Good Performance Indicators

✅ Response times within targets  
✅ Error rate < 0.1%  
✅ Stable throughput  
✅ Linear scaling with users  

### Performance Issues

⚠️ Increasing response times  
⚠️ High error rates  
⚠️ Throughput plateau  
⚠️ Memory leaks  

### Common Bottlenecks

1. **Database queries** - Slow recommendations
2. **MongoDB connections** - Connection pool exhaustion
3. **Memory** - Large result sets
4. **CPU** - Complex calculations
5. **Network** - High latency

## Troubleshooting

### Issue: Connection Refused

**Symptoms:** All requests fail immediately

**Solutions:**
```bash
# Check if server is running
curl http://localhost:8080/ping

# Start server
make run

# Check port
lsof -i :8080
```

### Issue: High Failure Rate

**Symptoms:** > 1% error rate

**Solutions:**
1. Check server logs for errors
2. Reduce concurrent users
3. Increase spawn rate time
4. Verify database health

### Issue: Slow Response Times

**Symptoms:** p95 > 1000ms

**Solutions:**
1. Check database query performance
2. Review MongoDB indexes
3. Monitor system resources
4. Reduce data set size

### Issue: Locust Not Installing

**Symptoms:** pip install fails

**Solutions:**
```bash
# Upgrade pip
pip install --upgrade pip

# Use Python 3.8+
python3 --version

# Install with specific version
pip install locust==2.20.0
```

### Issue: CSV Files Not Generated

**Symptoms:** Only HTML report created

**Solutions:**
```bash
# Specify CSV prefix explicitly
locust --csv=results

# Check write permissions
ls -la *.csv
```

## Best Practices

### 1. Gradual Load Increase

Don't spawn all users immediately:
- ✅ Spawn rate: 1-10 users/sec
- ❌ Spawn rate: 100 users/sec

### 2. Realistic Wait Times

Include think time between requests:
```python
wait_time = between(1, 3)  # seconds
```

### 3. Multiple Test Runs

Run tests multiple times:
- Warm-up run (discard results)
- 3-5 actual test runs
- Average the results

### 4. Clean State Between Runs

Reset database state:
```bash
make reset-db
make seed-db
```

### 5. Monitor Everything

Track:
- Application metrics
- Database performance
- System resources
- Network usage

### 6. Document Findings

Record:
- Test configuration
- Results summary
- Issues found
- Improvement recommendations

## Example Test Session

```bash
# 1. Install dependencies
pip install -r requirements-loadtest.txt

# 2. Start server
make run

# 3. Run baseline test
./run_load_test.sh
# Select option 2 (Light Load)

# 4. Review results
open load_test_light.html

# 5. Run medium load test
./run_load_test.sh
# Select option 3 (Medium Load)

# 6. Compare results
open load_test_medium.html

# 7. Run stress test
./run_load_test.sh
# Select option 5 (Stress Test)

# 8. Analyze bottlenecks
cat load_test_stress_stats.csv
```

## Performance Optimization Tips

Based on test results, consider:

1. **Database Optimization**
   - Add indexes on frequently queried fields
   - Optimize aggregation pipelines
   - Use projection to limit returned fields

2. **Caching**
   - Cache product lists
   - Cache category data
   - Cache user recommendations

3. **Connection Pooling**
   - Increase MongoDB connection pool size
   - Use persistent connections

4. **Code Optimization**
   - Profile slow endpoints
   - Reduce N+1 queries
   - Optimize recommendation algorithm

5. **Infrastructure**
   - Scale horizontally
   - Use load balancer
   - Deploy in production environment

## References

- [Locust Documentation](https://docs.locust.io/)
- [Performance Testing Best Practices](https://docs.locust.io/en/stable/writing-a-locustfile.html)
- [Interpreting Results](https://docs.locust.io/en/stable/analyzing-results.html)

## Support

For issues or questions:
1. Check server logs: `docker logs <container>`
2. Review Locust logs: `locust.log`
3. Verify configuration: `locust.conf`
4. Test manually: `curl http://localhost:8080/api/v1/products`
