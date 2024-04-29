# AWS Health Exporter
Health Exporterは、組織の AWS Health イベント情報を取得するためのコマンドラインツールです。サービス名とステータスでイベントをフィルタリングし、対象アカウントやリソースなどの詳細をCSVファイルにエクスポートできます。オプションで、CSV内容を標準出力にエコーすることもできます。

## 機能
* イベントフィルタリング: サービス名とステータスでイベントをフィルタリングし、必要なデータのみを取得できます。
* AWS Organizations 対応: AWS Health の組織 View から情報を取得します。スタンドアロンアカウントでは使用できません。
* CSVエクスポート: データを自動的にCSV形式で整形してエクスポートするため、保存、共有、分析が簡単です。

## 前提条件
* AWS Health APIとAWS Organizations APIにアクセスするための適切な権限を持つAWS認証情報
* AWS HealthAPIを使用するには、AWSサポートからBusiness、Enterprise On-Ramp、またはEnterpriseサポートプランが必要です

## 使い方
AWS Health Exporterを使用するには、必要なフラグを付けてコマンドを実行します。利用可能なフラグは以下の通りです。

--service, -s: サービス名でイベントをフィルタリングします (例: RDS)。
--status, -t: ステータスでイベントをフィルタリングします。指定可能な値は open、closed、upcoming です。
--echo, -e: CSV内容を標準出力にエコーします。

### コマンドの例
```bash
# ステータスがopenのRDSイベントをCSVにエクスポート
./health-exporter --service RDS --status open

# 今後の LAMBDA イベントを標準出力にエコーし、CSVにエクスポート
./health-exporter --service LAMBDA --status upcoming --echo
```

### 実行例
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
AWS Health Exporterへの貢献を歓迎します！問題の提出、プルリクエストの送信、ツールの改善のための拡張機能の提案など、お気軽にどうぞ。

## License
このプロジェクトはMITライセンスの下で公開されています。

## 免責事項
自己責任で使用してください。
