import pandas as pd

# Load the checkout data
checkout_df = pd.read_csv('checkout_2.csv')

# Display the first few rows and info of the checkout DataFrame
print("Checkout Data Head: ")
print(checkout_df.head())
print("\nCheckout Data Info:")
checkout_df.info()

#Load the transactions data
transactions_df = pd.read_csv('transactions.csv')

#Display the first few rows and info of the transactions Dataframe
print("\nTransactions Data Head:")
print(transactions_df.head())
print("\nTransactions Data Info:")
transactions_df.info()