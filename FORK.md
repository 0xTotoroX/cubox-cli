# Fork 说明 / About this fork

This is a fork of [OLCUBO/cubox-cli](https://github.com/OLCUBO/cubox-cli) that adds **folder write operations** (create / rename / delete), which the official CLI does not expose (its `folder` command only supports `list`).

这是官方 cubox-cli 的 fork，新增了官方 CLI 没有的**收藏夹写操作**（新建 / 重命名 / 删除）。

## What was added / 新增内容

```
# Folder management (web/app group)
cubox-cli folder delete --id ID[,ID2,...] [--dry-run] [--hard]
cubox-cli folder new --name NAME [--parent PARENT_ID]
cubox-cli folder rename --id ID --new-name NAME
cubox-cli folder archive --id ID[,ID2] --on|--off

# Recycle bin
cubox-cli recycle list [--page N]
cubox-cli recycle recover --id ID[,ID2,...]
cubox-cli recycle clean --id ID[,ID2,...]        # permanent, destructive

# Marks (highlights)
cubox-cli mark list [--page N] [--keyword K]
cubox-cli mark delete --id ID[,ID2,...]
cubox-cli mark color --id ID[,ID2] --color 1..5  # 1=Yellow 2=Green 3=Blue 4=Pink 5=Purple

# Reading lists
cubox-cli lists list
cubox-cli lists add-item --list ID --card ID [--highlight] [--note]
cubox-cli lists remove-item --list ID --card ID

# Export
cubox-cli export status        # today's export count + import readiness
cubox-cli export bookmarks     # request full-library export (v3 endpoint)

# Account
cubox-cli account apikey       # API extension key(s)
cubox-cli account settings     # reading settings
cubox-cli account sync         # Notion / Flowus / Readwise sync status
cubox-cli account insight --card CARD_ID   # a card's AI insight
```

- `folder delete` defaults to **safe mode**: remaining cards are first moved to
  Uncategorized, then the folder is deleted (mirroring the Cubox app's own
  behavior — cards are never lost). `--hard` skips the move step.
- `--dry-run` previews the operation without touching anything.

### Not covered / 未覆盖

- `AI ask` and streaming insight generation (SSE streams; reading an existing
  insight is covered by `account insight`)
- Drag-reorder of folders (`/group/move/another` payload is drag-session shaped)
- Reading-list create/delete and per-item sort/pin/update
- Writing settings back (`/settings/read/update`), per-card export formats
  (`/norm/cards/export/*` need format-specific payloads)
- These are straightforward to add — the auth plumbing is all in place.

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
