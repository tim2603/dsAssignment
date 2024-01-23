import random
import names

# Function to generate random data for bank customers
def generate_bank_customer_data(num_entries):
    data = ""
    for _ in range(num_entries):
        unique_identifier = f"{random.choice(['ah', 'bh', 'ch', 'dh', 'eh'])}{random.randint(1000, 9999)}{random.choice(['shfn', 'thfn', 'uhfn', 'vhfn', 'whfn'])}{random.randint(1000, 9999)}{random.choice(['ks', 'ls', 'ms', 'ns', 'os'])}"
        full_name = names.get_full_name(gender=random.choice(["male", "female"]))
        amount = round(random.uniform(100.0, 100000.0), 2)
        data_line = f"{unique_identifier} {full_name} {amount:.2f}\n"
        data += data_line
    return data

# Set the number of customer data entries
num_entries = 10000

# Create 10 text files with customer data
for file_num in range(1, 11):
    file_name = f"{file_num}_bank_customers.txt"
    customer_data = generate_bank_customer_data(num_entries)
    
    with open(file_name, "w") as file:
        file.write(customer_data)

    print(f"File '{file_name}' created.")