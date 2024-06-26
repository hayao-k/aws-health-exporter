# AWS Health Exporter
Health Exporter is a command-line tool designed to describe AWS Health events for your organization. It allows you to filter events by service name and status, and export the details to a CSV file. Optionally, you can echo the CSV content to standard output.

## Features
* Event Filtering: Filter events by service name, event status, and other criteria to get precisely the data you need.
* Entity Filtering: Filter affected entities by status code (IMPAIRED, UNIMPAIRED, UNKNOWN, PENDING, or RESOLVED).
* AWS Organizations Support: Works seamlessly with AWS Organizations, allowing you to get a health overview of all accounts.
* CSV Export: Automatically formats and exports the data into a CSV file, making it simple to store, share, and analyze.

## Prerequisites
* AWS credentials with appropriate permissions to access AWS Health and AWS Organizations services
* You must have a Business, Enterprise On-Ramp, or Enterprise Support plan from AWS Support to use the AWS Health API. 

## Usage
To use AWS Health Exporter, run the binary with the desired flags. Below are the available flags:

* `--event-filter`, `--filter`, `-f`: Filter events by service name, event status, and other criteria.
* `--status-code`, `-c`: Filter entitiy by status code. Possible values are IMPAIRED, UNIMPAIRED, UNKNOWN, PENDING and RESOLVED
* `--echo`, `-e`: Echo CSV content to standard output.
* `--profile`, `-p`: Specify the AWS credential profile to use.
* `--account-id`, `-i`: Specify a single account ID to process (optional).
* `--output-file`, `--file-name`, `o`: Specify the output CSV file name.

### Details of the event filtering option
The `--event-filter` option allows you to specify complex filtering criteria. Below is a table of the available fields that can be included in the filter criteria:

| Field             | Description                         | Possible Values                                                   |
|-------------------|-------------------------------------|-------------------------------------------------------------------|
| `service`         | Filter events by AWS service name.  | e.g., `LAMBDA`, `RDS`, `EKS`                                      |
| `status`          | Filter events by status.            | `open`, `closed`, `upcoming`                                      |
| `category`        | Filter events by category.          | `issue`, `accountNotification`, `scheduledChange`, `investigation`|
| `region`          | Filter events by region.            | AWS region codes, e.g., `us-east-1`                               |
| `startTime`       | Filter events by start time.        | ISO 8601 date format                                              |
| `endTime`         | Filter events by end time.          | ISO 8601 date format                                              |
| `lastUpdatedTime` | Filter events by last updated time. | ISO 8601 date format                                              |

For `startTime`, `endTime` and `lastUpdatedTime`, you can specify a time range using `from` and `to` in ISO 8601 date format. Here is the structure for specifying the time range:

- `{from:YYYY-MM-DDTHH:MM:SSZ,to:YYYY-MM-DDTHH:MM:SSZ}`


### Example Commands

```bash
# Describe RDS events with open status and export to CSV
./health-exporter --event-filter service=RDS,status=open

# Describe upcoming LAMBDA events and echo the output to STDOUT
./health-exporter --event-filter service=LAMBDA,status=upcoming --echo

# Describe only events in the Tokyo region and specify their last updated time.
./health-exporter ----event-filter "lastUpdatedTime={from=2024-03-01T00:00:00Z,to=2024-05-02T23:59:59Z},region=ap-northeast-1"

# Get entities with pending status only and specify a custom file name
./health-exporter --status-code PENDING --output-file my_event_details.csv

# Get events using the specified profile
./health-exporter --profile my-profile

# Process only a single account
./health-exporter --account-id 123456789012
```

### Execution Example
```bash
$ health-exporter --event-filter service=LAMBDA,status=upcoming
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
