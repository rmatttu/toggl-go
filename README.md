# toggl-go

Toggl Track data to file.

## Usage

```bash
./toggl-go -email "your-mail-address@example.com" -token "___YOUR_TOGGL_API_TOKEN___" -since 2021 | tee output.jsonl
```

```bash
go run main.go -email "your-mail-address@example.com" -token "___YOUR_TOGGL_API_TOKEN___" -since 2021 | tee output.jsonl
```

### jq sample

```bash
cat output.jsonl | jq -r '.data[] | [.id, .pid, .tid, .uid, .description, .start, .end, .project, .project_hex_color, dur, (.tags | join("|"))] | @tsv' | sort | uniq | less
# id pid tid uid description start end project tags
```

## Requirements

## Installation

## License

## Author

## References

### toggl

* [toggl_api_docs/detailed.md at master · toggl/toggl_api_docs](https://github.com/toggl/toggl_api_docs/blob/master/reports/detailed.md)
* [PythonでTogglのAPIをたたく - Qiita](https://qiita.com/stu345/items/abebeb2da6d97382a7b1)
