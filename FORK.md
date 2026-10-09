# Fork 说明 / About this fork

This is a fork of [OLCUBO/cubox-cli](https://github.com/OLCUBO/cubox-cli) that adds **folder write operations** (create / rename / delete), which the official CLI does not expose (its `folder` command only supports `list`).

这是官方 cubox-cli 的 fork，新增了官方 CLI 没有的**收藏夹写操作**（新建 / 重命名 / 删除）。

## What was added / 新增内容

```
cubox-cli folder delete --id ID[,ID2,...] [--dry-run] [--hard]
cubox-cli folder new --name NAME [--parent PARENT_ID]
cubox-cli folder rename --id ID --new-name NAME
```

- `folder delete` defaults to **safe mode**: remaining cards are first moved to
  Uncategorized, then the folder is deleted (mirroring the Cubox app's own
  behavior — cards are never lost). `--hard` skips the move step.
- `--dry-run` previews the operation without touching anything.

## How it works / 工作原理

Cubox exposes two separate API groups with **mutually isolated credentials**:

| API group | Endpoints | Auth style | Token |
|---|---|---|---|
| CLI group | `/c/api/cli/*` | `Authorization: Bearer <token>` | API Extension token (`auth login`) |
| Web/App group | `/c/api/*` | `Authorization: <token>` (**bare**, no Bearer prefix) | Cubox.app login token |

All pre-existing commands keep using the CLI group (unchanged). The new folder
write commands use the Web/App group — the same endpoints the Cubox web app and
the native Cubox.app call:

- `POST /c/api/group/moveSearchEnginesOut/{id}` — move cards to Uncategorized + delete folder
- `POST /c/api/group/delete/{id}` — direct delete
- `POST /c/api/group/new` — create (expects `application/x-www-form-urlencoded`)
- `POST /c/api/group/update` — rename (expects `application/x-www-form-urlencoded`)
- `GET /c/api/v2/group/my` — list folders (camelCase payload)

Note the request format details discovered while implementing: `group/new` and
`group/update` reject JSON bodies with `-5000 Invalid input` — they require
form-urlencoded encoding.

## App token resolution / App token 来源

The web/app group token (`Ctoken`) is resolved in this order:

1. `CUBOX_CTOKEN` environment variable
2. macOS: read automatically from the logged-in Cubox.app
   (`~/Library/Group Containers/group.com.guaiqi.cubox/...plist`)

## Caveats / 注意事项

- These endpoints are **undocumented**; Cubox may change them at any time.
- The app token is tied to the Cubox.app login session (logging out
  invalidates it).
- Use responsibly: your own account, moderate request rates.

## License

MIT (inherited from the upstream project).
