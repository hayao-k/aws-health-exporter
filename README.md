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

`--service`, `-s`: Filter events by service name (e.g., RDS).
`--status`, `-t`: Filter events by status. Possible values are open, closed, and upcoming.
`--echo`, `-e`: Echo CSV content to standard output.

### Example Commands

```bash
# Describe RDS events with open status and export to CSV
./health-exporter --service RDS --status open

# Describe upcoming EC2 events and echo the output to STDOUT
./health-exporter --service LAMBDA --status upcoming --echo
```

### Execution Example
```bash
$ health-exporter --service=RDS --status=upcoming
Use the arrow keys to navigate: ↓ ↑ → ← 
? Select an event: 
  ▸ RDS - AWS_RDS_PLANNED_LIFECYCLE_EVENT (ap-northeast-1, 2024-08-22 07:00:00)
    RDS - AWS_RDS_PLANNED_LIFECYCLE_EVENT (us-east-2, 2024-05-31 07:00:00)
    RDS - AWS_RDS_PLANNED_LIFECYCLE_EVENT (ap-northeast-1, 2024-05-31 07:00:00)
    RDS - AWS_RDS_PLANNED_LIFECYCLE_EVENT (us-west-1, 2024-08-22 07:00:00)
↓   RDS - AWS_RDS_PLANNED_LIFECYCLE_EVENT (us-east-1, 2024-08-22 07:00:00)

✔ RDS - AWS_RDS_PLANNED_LIFECYCLE_EVENT (ap-northeast-1, 2024-08-22 07:00:00)
Event details have been written to AWS_RDS_PLANNED_LIFECYCLE_EVENT_2024-08-22_07-00-00_ap-northeast-1.csv.
```

## Contributing
Contributions to Health Exporter are welcome! Please feel free to submit issues, pull requests, or enhancements to improve the tool.

## License
This project is licensed under the MIT License.

## Disclaimer
This tool is not officially supported by AWS. Please use it at your own risk.
