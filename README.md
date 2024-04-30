# AWS Health Exporter
Health Exporter is a command-line tool designed to describe AWS Health events for your organization. It allows you to filter events by service name and status, and export the details to a CSV file. Optionally, you can echo the CSV content to standard output.

## Features
* Event Filtering: Filter events by service name and status, to get precisely the data you need.
* AWS Organizations Support: Works seamlessly with AWS Organizations, allowing you to get a health overview of all accounts.
* CSV Export: Automatically formats and exports the data into a CSV file, making it simple to store, share, and analyze.

## Prerequisites
* AWS credentials with appropriate permissions to access AWS Health and AWS Organizations services
* You must have a Business, Enterprise On-Ramp, or Enterprise Support plan from AWS Support to use the AWS Health API. 

## Usage
To use AWS Health Exporter, run the binary with the desired flags. Below are the available flags:

* `--service`, `-s`: Filter events by service name (e.g., RDS).
* `--status`, `-t`: Filter events by status. Possible values are open, closed, and upcoming.
* `--echo`, `-e`: Echo CSV content to standard output.

### Example Commands

```bash
# Describe RDS events with open status and export to CSV
./health-exporter --service RDS --status open

# Describe upcoming LAMBDA events and echo the output to STDOUT
./health-exporter --service LAMBDA --status upcoming --echo
```

### Execution Example
```bash
$ health-exporter --service=RDS --status=upcoming
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select an event: 
  ▸ LAMBDA - AWS_LAMBDA_PLANNED_LIFECYCLE_EVENT (us-east-1, 2024-10-14 07:00:00)
    LAMBDA - AWS_LAMBDA_PLANNED_LIFECYCLE_EVENT (ap-northeast-1, 2024-10-14 07:00:00)
    LAMBDA - AWS_LAMBDA_PLANNED_LIFECYCLE_EVENT (ap-northeast-1, 2024-06-12 07:00:00)
    LAMBDA - AWS_LAMBDA_PLANNED_LIFECYCLE_EVENT (ap-southeast-2, 2024-10-14 07:00:00)
↓   LAMBDA - AWS_LAMBDA_PLANNED_LIFECYCLE_EVENT (us-east-1, 2024-06-12 07:00:00)

✔ LAMBDA - AWS_LAMBDA_PLANNED_LIFECYCLE_EVENT (us-east-1, 2024-10-14 07:00:00)
Event details have been written to AWS_LAMBDA_PLANNED_LIFECYCLE_EVENT_2024-10-14_07-00-00_us-east-1.csv.
```

### Example of output
```csv
Account ID,Account Name,Region,Identifier,Status,Last Updated
000000000000,account-0000,us-east-1,arn:aws:lambda:us-east-1:000000000000:function:Old_Runtime_Lambda_Function-1PBKPZPFSJ058,PENDING,2024-04-21 20:11:29
111111111111,account-1111,us-east-1,arn:aws:lambda:us-east-1:111111111111:function:Old_Runtime_Lambda_Function-uuTi2u7DbooD,PENDING,2024-04-21 20:11:29
111111111111,account-1111,us-east-1,arn:aws:lambda:us-east-1:111111111111:function:Old_Runtime_Lambda_Function-omdieC8Umobo,PENDING,2024-04-21 20:11:29
222222222222,account-2222,us-east-1,arn:aws:lambda:us-east-1:222222222222:function:Old_Runtime_Lambda_Function-ULZ27BYSQ0MN,PENDING,2024-04-21 20:11:29
222222222222,account-2222,us-east-1,arn:aws:lambda:us-east-1:222222222222:function:Old_Runtime_Lambda_Function-10YNGBMU46VP9,PENDING,2024-04-21 20:11:29
222222222222,account-2222,us-east-1,arn:aws:lambda:us-east-1:222222222222:function:Old_Runtime_Lambda_Function-CEgHAu41udFy,PENDING,2024-04-21 20:11:29
333333333333,account-3333,us-east-1,arn:aws:lambda:us-east-1:333333333333:function:Old_Runtime_Lambda_Function-zNKRpLWP0pXB,PENDING,2024-04-21 20:11:29
333333333333,account-3333,us-east-1,arn:aws:lambda:us-east-1:333333333333:function:Old_Runtime_Lambda_Function-24ES8MRQJ9R6,PENDING,2024-04-21 20:11:29
444444444444,account-4444,us-east-1,arn:aws:lambda:us-east-1:444444444444:function:Old_Runtime_Lambda_Function-134QIS8IYF84K,PENDING,2024-04-21 20:11:29
444444444444,account-4444,us-east-1,arn:aws:lambda:us-east-1:444444444444:function:Old_Runtime_Lambda_Function-B97VeyrZNXIy,PENDING,2024-04-21 20:11:29
```

## Contributing
Contributions to Health Exporter are welcome! Please feel free to submit issues, pull requests, or enhancements to improve the tool.

## License
This project is licensed under the MIT License.

## Disclaimer
This tool is not officially supported by AWS. Please use it at your own risk.
