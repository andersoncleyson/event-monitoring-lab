import pandas as pd
import matplotlib.pyplot as plt

# Load the transactions  data and convert timestamp
transactions_df = pd.read_csv('transactions.csv')
transactions_df['timestamp'] = pd.to_datetime(transactions_df['timestamp'])
transactions_df['hour'] = transactions_df['timestamp'].dt.hour

# "SQL Query" equivalent using pandas
# We will focuus on the anamalous hours (8, 9, 15, 20) and some normal hours for comparison
select_hours = [8, 9, 10, 11, 12, 13, 14, 15, 16, 20]
filtered_transactions = transactions_df[transactions_df['hour'].isin(select_hours)]

# Group by hour and status, and sum the counts
hourly_status_counts = filtered_transactions.groupby(['hour', 'status'])['count'].sum().unstack().fillna(0)

# Create the plot
hourly_status_counts.plot(kind='bar', figsize=(15, 8), width=0.8)

plt.title('Transaction Count by Status at Selected Times')
plt.xlabel('Time of day')
plt.ylabel('Total Transaction Count')
plt.xticks(rotation=0)
plt.legend(title='Status')
plt.grid(axis='y', linestyle='--')
plt.tight_layout()
plt.savefig('hourly_status_counts_anomalous.png')

# Print the "query" result to be used in the explanation
print("'SQL Query' Result (Transaction Count by Time and Status):")
print(hourly_status_counts)
