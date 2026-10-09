# Fork 说明 / About this fork

This is a fork of [OLCUBO/cubox-cli](https://github.com/OLCUBO/cubox-cli) that **unlocks the full web/app API group** on top of the official CLI: folder write operations, batch card operations, recycle bin, marks (highlights) management, reading lists, import/export tasks, account info and AI features.

这是官方 cubox-cli 的 fork：在官方能力之外，解锁了 **web/app 组的全部接口**——收藏夹写操作、批量卡片操作、回收站、标注管理、阅读清单、导入/导出任务、账号信息与 AI 功能。

```
# Folder management (web/app group)
cubox-cli folder delete --id ID[,ID2,...] [--dry-run] [--hard]
cubox-cli folder new --name NAME [--parent PARENT_ID]
cubox-cli folder rename --id ID --new-name NAME
cubox-cli folder archive --id ID[,ID2] --on|--off
cubox-cli folder move --id ID --parent PARENT_ID [--index N]   # N: 0-based, default last

# Batch card operations (web/app group)
cubox-cli card star --id ID[,ID2] --on|--off
cubox-cli card read --id ID[,ID2,...]            # batch mark as read
cubox-cli card move --id ID[,ID2] --folder-id FID  # batch move to folder
cubox-cli card tag --id ID[,ID2] --add name,name [--delete-tag-ids ID]

# Recycle bin
cubox-cli recycle list [--page N]
cubox-cli recycle recover --id ID[,ID2,...]
cubox-cli recycle clean --id ID[,ID2,...]        # permanent, destructive

# Tags (web/app group)
cubox-cli tag new --name NAME [--parent TAG_ID]
cubox-cli tag sort --json '{...}'                # experimental drag-order payload

# Marks (highlights)
cubox-cli mark list [--page N] [--keyword K]
cubox-cli mark count
cubox-cli mark delete --id ID[,ID2,...]
cubox-cli mark color --id ID[,ID2] --color 1..5  # 1=Yellow 2=Green 3=Blue 4=Pink 5=Purple
cubox-cli mark export [--card ID] [--mark ID]
cubox-cli mark export-text [--card ID] [--mark ID]

# Reading lists
cubox-cli lists list
cubox-cli lists new --title TITLE [--intro INTRO]
cubox-cli lists delete --list ID
cubox-cli lists update --list ID [--title T] [--intro I]
cubox-cli lists publish --list ID --on|--off
cubox-cli lists cards --list ID [--page N] [--size M]
cubox-cli lists add-item --list ID --card ID [--highlight] [--note]
cubox-cli lists remove-item --list ID --card ID
cubox-cli lists collect --list ID --card ID
cubox-cli lists collect-all --list ID
cubox-cli lists sort --list ID --card ID1,ID2    # full ordered list
cubox-cli lists check --card ID[,ID2]            # which lists contain these cards

# Export
cubox-cli export status        # today's export count + import readiness
cubox-cli export bookmarks     # request full-library export (v3 endpoint)
cubox-cli export mail --email ADDR --card ID[,ID2,...]   # async export via email (verified)

# Import (official task pipeline)
cubox-cli import bookmarks --file bookmarks.html   # upload (HTML; folder tree preserved)
cubox-cli import progress      # poll the import task

# Account
cubox-cli account apikey       # API extension key(s)
cubox-cli account settings     # reading settings
cubox-cli account settings-update --json '{...}'   # experimental: full settings object
cubox-cli account sync         # Notion / Flowus / Readwise sync status
cubox-cli account insight --card CARD_ID   # a card's AI insight
cubox-cli account insight-state --card CARD_ID
cubox-cli account more-qas --card CARD_ID  # follow-up Q&As of the insight
cubox-cli account mail         # Mail Drop (email-in) settings

# AI
cubox-cli ai generate --card CARD_ID   # stream AI insight generation (SSE; consumes AI quota)
cubox-cli ai ask --question Q [--context TEXT]   # ⚠ experimental: endpoint verified
                                       # (POST, form-urlencoded, OpenAI-delta SSE)
                                       # but the required field set is not fully
                                       # mapped yet — needs one live request capture
```

- `folder delete` defaults to **safe mode**: remaining cards are first moved to
  Uncategorized, then the folder is deleted (mirroring the Cubox app's own
  behavior — cards are never lost). `--hard` skips the move step.
- `--dry-run` previews the operation without touching anything.

### Request-encoding notes (reverse-engineered)

The web/app group uses **several different body encodings**, and the server
rejects the wrong one with `-5000 Invalid input` / `must not be blank`:

| Endpoints | Encoding |
|---|---|
| `group/new`, `group/update`, `group/move/another`, `norm/cards/export/async`, `settings/read/update`, `marks/*`, `recycle/*` | **form-urlencoded** |
| `group/delete/{id}`, `group/moveSearchEnginesOut/{id}` | no body (path param) |
| `reading-list*`, `lists/items/*`, `v2/bookmark/import` | multipart/form-data |
| `ai/ask`, `card/insight/generate/stream` | form-urlencoded request, **SSE** response (OpenAI delta format) |

(The web app's own HTTP client serializes every POST body with a
query-string encoder — including bodies that look like JSON in its source,
e.g. the recycle endpoints' `searchEngines` field is a JSON-*stringified*
array inside a form body.)

`folder move` may take a while server-side (tree reorder); the client uses a
120s timeout for form posts, and a timed-out request may still have been
applied — re-check with `folder list`.

### Not covered / 未覆盖

- `ai/ask` free-form Q&A: command shipped as **experimental** — the wire
  format is mapped but one required form field is still unidentified
  (needs one live request capture from the web/app).
- Reading-list per-item pin/update.
- Everything else in the web/app group is covered.

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
