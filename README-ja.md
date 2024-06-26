# AWS Health Exporter
Health Exporterは、組織の AWS Health イベント情報を取得するためのコマンドラインツールです。サービス名とステータスでイベントをフィルタリングし、対象アカウントやリソースなどの詳細をCSVファイルにエクスポートできます。オプションで、CSV内容を標準出力にエコーすることもできます。

## 機能
* イベントフィルタリング: サービス名やステータスなどの条件でイベントをフィルタリングし、必要なデータのみを取得できます。
* リソースフィルタリング: 特定のステータスコード (IMPAIRED, UNIMPAIRED, UNKNOWN, PENDING, or RESOLVED) にマッチするリソースのみを取得できます。
* AWS Organizations 対応: AWS Health の組織 View から情報を取得します。スタンドアロンアカウントでは使用できません。
* CSVエクスポート: データを自動的にCSV形式で整形してエクスポートするため、保存、共有、分析が簡単です。

## 前提条件
* AWS Health APIとAWS Organizations APIにアクセスするための適切な権限を持つAWS認証情報
* AWS HealthAPIを使用するには、AWSサポートからBusiness、Enterprise On-Ramp、またはEnterpriseサポートプランが必要です

## 使い方
AWS Health Exporterを使用するには、必要なフラグを付けてコマンドを実行します。利用可能なフラグは以下の通りです。

* `--event-filter`, `--filter`, `-f`: サービス、ステータスなどの条件でイベントをフィルタリングします。
* `--status-code`, `-c`: ステータスコードでエンティティをフィルタリングします。指定可能な値は IMPAIRED, UNIMPAIRED, UNKNOWN, PENDING, RESOLVED です。
* `--echo`, `-e`: CSV内容を標準出力にエコーします。
* `--profile`, `-p`: AWS クレデンシャルプロファイルを指定します。
* `--account-id`, `-a`: 単一のアカウント ID を指定して、そのアカウントのみを処理します (指定しない場合は全アカウントを出力)。
* `--output-file`, `--file-name`, `o`: 任意のファイル名を指定して、CSVを出力できます (指定しない場合はツールデフォルトのファイル名を使用)。

### イベントのフィルタリング
`--event-filter`オプションを使用すると、複雑なフィルタリング条件を指定できます。以下の表に、フィルタ条件に含めることができる利用可能なフィールドを示します。

| フィールド        | 説明                                | 指定可能な値                                              　　    |
|-------------------|-------------------------------------|-------------------------------------------------------------------|
| `service`         | Filter events by AWS service name.  | e.g., `LAMBDA`, `RDS`, `EKS`                                      |
| `status`          | Filter events by status.            | `open`, `closed`, `upcoming`                                      |
| `category`        | Filter events by category.          | `issue`, `accountNotification`, `scheduledChange`, `investigation`|
| `region`          | Filter events by region.            | AWS region codes, e.g., `us-east-1`                               |
| `startTime`       | Filter events by start time.        | ISO 8601 date format                                              |
| `endTime`         | Filter events by end time.          | ISO 8601 date format                                              |
| `lastUpdatedTime` | Filter events by last updated time. | ISO 8601 date format                                              |

`startTime`, `endTime`, `lastUpdatedTime` では、ISO 8601日付形式を使用して `from` と `to` で時間範囲を指定できます。

- `{from:YYYY-MM-DDTHH:MM:SSZ,to:YYYY-MM-DDTHH:MM:SSZ}`


### コマンドの例
```bash
# ステータスがopenのRDSイベントをフィルタリング
./health-exporter ----event-filter service=RDS,status=open

# ステータスがupcomingのLAMBDA イベントをフィルタリングし、結果は標準出力にもエコーする
./health-exporter --event-filter service=LAMBDA,status=upcoming --echo

# 東京リージョンのイベントを最終更新日を指定してフィルタリング
./health-exporter ----event-filter "lastUpdatedTime={from=2024-03-01T00:00:00Z,to=2024-05-02T23:59:59Z},region=ap-northeast-1"

# 未対応のリソースのみを指定したファイル名でCSVにエクスポート
./health-exporter --status-code PENDING --output-file my_event_details.csv

# 任意の AWS Profile を指定
./health-exporter --profile my-profile

# 特定の 1 アカウントのみをフィルタリングして出力
./health-exporter --account-id 123456789012
```

### 実行例
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

### 出力例
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
AWS Health Exporterへの貢献を歓迎します！問題の提出、プルリクエストの送信、ツールの改善のための拡張機能の提案など、お気軽にどうぞ。

## License
このプロジェクトはMITライセンスの下で公開されています。

## 免責事項
自己責任で使用してください。
