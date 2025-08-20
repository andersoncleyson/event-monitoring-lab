import pandas as pd
import matplotlib.pyplot as plt

# Load the checkout data
checkout_df = pd.read_csv('checkout_2.csv')

# Plotting the checkout data
plt.figure(figsize=(12, 6))
plt.plot(checkout_df['time'], checkout_df['today'], label='Today', marker='o')
plt.plot(checkout_df['time'], checkout_df['yesterday'], label='Yesterday', marker='x')
plt.plot(checkout_df['time'], checkout_df['same_day_last_week'], label='Same Day Last Week', marker='s')
plt.plot(checkout_df['time'], checkout_df['avg_last_week'], label='Avg. Last Week', linestyle='--')
plt.plot(checkout_df['time'], checkout_df['avg_last_month'], label='Avg. Last Month', linestyle=':')
plt.xlabel("Time of Day")
plt.ylabel("Checkout Count")
plt.title("Checkout Comparison: Today vs. Historical Data")
plt.xticks(rotation=45)
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig('checkout_comparison.png')

# Load the transactions data
transactions_df = pd.read_csv('transactions.csv')
transactions_df['timestamp'] = pd.to_datetime(transactions_df['timestamp'])

# Separate data by status
approved_df = transactions_df[transactions_df['status'] == 'approved'].set_index('timestamp')
refunded_df = transactions_df[transactions_df['status'] == 'refunded'].set_index('timestamp')
denied_df = transactions_df[transactions_df['status'] == 'denied'].set_index('timestamp')
failed_df = transactions_df[transactions_df['status'] == 'failed'].set_index('timestamp')
reversed_df = transactions_df[transactions_df['status'] == 'reversed'].set_index('timestamp')

# Plotting transactions data
plt.figure(figsize=(15, 10))

# Plot for approved transactions
plt.subplot(3, 2, 1)
plt.plot(approved_df.index, approved_df['count'], label='Approved')
plt.title('Approved Transactions Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Count')
plt.xticks(rotation=45)
plt.legend()
plt.grid(True)

# Plot for refunded transactions
plt.subplot(3, 2, 2)
plt.plot(refunded_df.index, refunded_df['count'], label='Refunded', color='orange')
plt.title('Refunded Transactions Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Count')
plt.xticks(rotation=45)
plt.legend()
plt.grid(True)

# Plot for denied transactions
plt.subplot(3, 2, 3)
plt.plot(denied_df.index, denied_df['count'], label='Denied', color='red')
plt.title('Denied Transactions Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Count')
plt.xticks(rotation=45)
plt.legend()
plt.grid(True)

# Plot for failed transactions
plt.subplot(3, 2, 4)
plt.plot(failed_df.index, failed_df['count'], label='Failed', color='purple')
plt.title('Failed Transactions Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Count')
plt.xticks(rotation=45)
plt.legend()
plt.grid(True)

# Plot for reversed transactions
plt.subplot(3, 2, 5)
plt.plot(reversed_df.index, reversed_df['count'], label='Reversed', color='brown')
plt.title('Reversed Transactions Over Time')
plt.xlabel('Timestamp')
plt.ylabel('Count')
plt.xticks(rotation=45)
plt.legend()
plt.grid(True)

plt.tight_layout()
plt.savefig('transactions_by_status.png')

# Agregando por hora para a análise de anomalias
transactions_df['hour'] = transactions_df['timestamp'].dt.hour
hourly_transactions = transactions_df.groupby(['hour', 'status'])['count'].sum().unstack().fillna(0)

# Plotando transações por hora
hourly_transactions.plot(kind='bar', stacked=True, figsize=(15, 7))
plt.title('Hourly Transaction Counts by Status')
plt.xlabel('Hour of Day')
plt.ylabel('Total Transaction Count')
plt.xticks(rotation=0)
plt.grid(axis='y')
plt.savefig('hourly_transactions.png')