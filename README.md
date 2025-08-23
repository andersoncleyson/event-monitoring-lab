# Monitoring Analyze Test

This project implements a complete, virtualized monitoring environment for analyzing transaction data from CSV files. The stack uses Python to understand anomaly behavior, Prometheus for collecting and storing metrics, Grafana for visualization and dashboards, and a custom Go Exporter to read the data and expose it to Prometheus. The entire environment is orchestrated with Docker and Docker Compose, making it portable and easy to configure.

### Analyzing CSVs

To understand checkout behavior, I compared today's data with historical data (yesterday, the same day last week, last week's average, and last month's average). The chart below illustrates this comparison.

![Checkouts](data-analysis/img/checkouts.png)

**Checkout Comparison: Today vs. Historical Data**

![Checkout Comparison](data-analysis/img/checkout_comparison.png)

**Observations and Anomalies:**

Abnormal spikes at 8:00 AM-9:00 AM and 8:00 PM: Note that there is an extremely sharp increase in the number of checkouts "today" at 8:00 AM and 9:00 AM, and another spike at 8:00 PM. These values are much higher than all historical metrics, suggesting atypical behavior.

Anomalous drop at 3:00 PM: There was a drop to zero in the number of checkouts at 3:00 PM, which is unusual compared to the data from "yesterday" and "the same day last week." This may indicate a possible instability or technical issue.

**Transactions by Status Over Time**

This chart shows the transaction count for each status over time. Due to the large amount of data, it is easier to observe overall patterns and volatility.

**Transaction Count by Hour and Status**

When the data was aggregated hourly, a clearer view of daily patterns was obtained:

![Transactions by status](data-analysis/img/transactions_by_status.png)

**Observations:**

Predominance of "approved" transactions: The vast majority of transactions are approved, which is expected.

Hourly Patterns: Transaction activity appears to increase throughout the day, with peaks at certain times, which is a common business pattern.

**Conclusion and Recommendations**

Analysis of checkout data reveals clear and significant anomalies that should be investigated:

Checkout Spikes (8:00 AM-9:00 AM and 8:00 PM): Checkout spikes can be the result of:

Marketing actions: A successful campaign or promotion may have generated a sudden increase in traffic and conversions.

Fraudulent Activity: Such a drastic increase could also be a sign of bot activity or fraud attempts.

Specific Event: An external event may have led to this behavior.

Checkout Drop (3:00 PM): A drop to zero is a strong indication of:

Technical Issues: There may have been a failure in the checkout system, payment gateway, or some other part of the infrastructure.

### Explaining anomaly behavior

To further analyze the anomalous behavior, I simulated a SQL query on the transaction data to focus on the times we identified as problematic (8 AM, 9 AM, 3 PM, and 8 PM).

The "query" performed was essentially:

```
SELECT HOUR(timestamp), status, SUM(count)
FROM transactions
WHERE HOUR(timestamp) IN (8, 9, 10, 11, 12, 13, 14, 15, 16, 20)
GROUP BY HOUR(timestamp), status;
```

![SQL result](data-analysis/img/sql-result.png)

I created a bar chart from the output to visualize the distribution of transaction statuses at these specific times.

**Transaction Analysis Chart by Status and Time**

![hourly_status_counts_anamolous](data-analysis/img/hourly_status_counts_anomalous.png)

**Explanation of Anomalous Behavior**

By comparing the information from checkout_2.csv with this new view of transactions.csv, we can formulate much more precise hypotheses:

1. **The Drastic Drop in Checkouts at 3 PM**

**Checkout Anomaly:** The first log (checkout_2.csv) showed zero checkouts at 3 PM.

**Transaction Analysis:** The graph above shows that at 3 PM, there was a normal transaction volume, with 21,537 "approved" transactions and others within the expected range.

**Conclusion:** This is a crucial contradiction. The most likely explanation is that the checkout system was working correctly, but the system that records the checkout logs failed specifically at 3 PM. The problem wasn't a drop in sales, but rather a failure in monitoring or data collection.

2. **Checkout Spikes at 8:00 AM-9:00 AM**

**Checkout Anomaly:** The first log showed an extreme spike in checkouts at these times.

**Transaction Analysis:** The transaction graph reveals that, along with the high volume of approved transactions, 8:00 AM and 9:00 AM had the highest volume of "denied" transactions compared to other times.

**Conclusion:** The checkout spike was not entirely due to successful sales. The disproportionate increase in denied transactions may indicate:

**Fraud Attempts:** A "card testing" attack, where fraudsters mass-test a list of stolen credit cards. Most are denied, but it generates a high volume of checkout attempts.

**Misdirected Marketing Campaign:** A campaign that attracted an audience without sufficient funds, resulting in many denied purchase attempts.

3. **The 8 PM Checkout Spike**

**Checkout Anomaly:** There was also a checkout spike at 8 PM.

**Transaction Analysis:** Unlike 8 AM–9 AM, the distribution of transaction statuses at 8 PM appears normal, with no significant increase in declined or failed transactions.

**Conclusion:** This 8 PM spike appears to be a legitimate increase in customer traffic, perhaps due to it being a prime time for online shopping.

4. **New Anomaly Found at 4 PM**

**Transaction Analysis:** Our SQL query revealed something not visible in the checkout data: at 4 PM, transactions with "backend_reversed" and "failed" statuses appeared, which did not appear at the other times analyzed.

**Conclusion:** This indicates that a new type of technical issue may have started occurring after 4 PM, possibly related to communication failures with the payment system or backend issues.

In short, by combining the two data sources, we can go beyond simply identifying "spikes" and "dips" and begin to understand the potential causes behind each anomaly, separating likely technical failures from potential fraudulent activity or legitimate traffic spikes.

### Monitoring stack

**Architectural overview**

The data flow follows the Prometheus *pull* model:


1. **Transaction Exporter**: A custom service written in Go that reads the `transactions.csv` and `checkout_2.csv` files in real time. It processes this data and exposes it to an HTTP endpoint (`/metrics`) in a format Prometheus understands.
2. **Prometheus:** The Prometheus server is configured to scrape data from our exporter's `/metrics` endpoint at regular intervals. It stores these metrics in its time series database.
3. **Grafana:** Grafana connects to Prometheus as a data source. It queries Prometheus using PromQL and displays the data in interactive, visual dashboards, allowing analysis of transaction and checkout behavior.

