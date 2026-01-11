# Sekisho API仕様書

## 1. エンドポイント一覧

| エンドポイント | メソッド | 認証 | 説明 |
|----------------|----------|------|------|
| `/oauth2/start` | GET | 不要 | OAuth2認証フロー開始 |
| `/oauth2/callback` | GET | 不要 | IdPからのコールバック |
| `/oauth2/sign_out` | GET/POST | 必要 | ログアウト |
| `/health` | GET | 不要 | ヘルスチェック |
| `/*` | ANY | 必要 | リバースプロキシ |

---

## 2. エンドポイント詳細

### 2.1 GET `/oauth2/start`

OAuth2認証フローを開始し、IdPへリダイレクトする。

#### リクエスト

| パラメータ | 位置 | 必須 | 説明 |
|------------|------|------|------|
| `rd` | Query | × | 認証後のリダイレクト先（デフォルト: `/`） |

```http
GET /oauth2/start?rd=/dashboard HTTP/1.1
Host: localhost:4180
```

#### レスポンス

**成功時: 302 Found**

```http
HTTP/1.1 302 Found
Location: https://keycloak.example.com/realms/sekisho/protocol/openid-connect/auth
  ?client_id=sekisho
  &redirect_uri=http://localhost:4180/oauth2/callback
  &response_type=code
  &scope=openid+email+profile
  &state=<random_state>
  &nonce=<random_nonce>
  &code_challenge=<pkce_challenge>
  &code_challenge_method=S256
Set-Cookie: _sekisho_csrf=<encrypted_state_data>; Path=/; HttpOnly; Secure; SameSite=Lax
```

---

### 2.2 GET `/oauth2/callback`

IdPからの認可コードを受け取り、トークン交換を行う。

#### リクエスト

| パラメータ | 位置 | 必須 | 説明 |
|------------|------|------|------|
| `code` | Query | ○ | 認可コード |
| `state` | Query | ○ | CSRF検証用state |

```http
GET /oauth2/callback?code=abc123&state=xyz789 HTTP/1.1
Host: localhost:4180
Cookie: _sekisho_csrf=<encrypted_state_data>
```

#### レスポンス

**成功時: 302 Found**

```http
HTTP/1.1 302 Found
Location: /dashboard
Set-Cookie: _sekisho=<encrypted_session>; Path=/; HttpOnly; Secure; SameSite=Lax; Max-Age=86400
Set-Cookie: _sekisho_csrf=; Path=/; Max-Age=0
```

**エラー時: 400 Bad Request**

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": "invalid_state",
  "error_description": "State parameter mismatch"
}
```

---

### 2.3 GET/POST `/oauth2/sign_out`

セッションを破棄し、ログアウトする。

#### リクエスト

| パラメータ | 位置 | 必須 | 説明 |
|------------|------|------|------|
| `rd` | Query | × | ログアウト後のリダイレクト先 |

```http
POST /oauth2/sign_out HTTP/1.1
Host: localhost:4180
Cookie: _sekisho=<encrypted_session>
```

#### レスポンス

**成功時: 302 Found**

```http
HTTP/1.1 302 Found
Location: /
Set-Cookie: _sekisho=; Path=/; Max-Age=0
```

---

### 2.4 GET `/health`

アプリケーションのヘルス状態を返す。

#### リクエスト

```http
GET /health HTTP/1.1
Host: localhost:4180
```

#### レスポンス

**成功時: 200 OK**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "ok",
  "version": "1.0.0"
}
```

---

### 2.5 ANY `/*`（リバースプロキシ）

認証済みリクエストを上流サーバーに転送する。

#### 認証済みリクエスト

```http
GET /api/users HTTP/1.1
Host: localhost:4180
Cookie: _sekisho=<encrypted_session>
```

#### 上流への転送リクエスト

```http
GET /api/users HTTP/1.1
Host: upstream-app:8080
X-Forwarded-User: user@example.com
X-Forwarded-Email: user@example.com
X-Forwarded-Preferred-Username: johndoe
X-Forwarded-Groups: admin,users
X-Forwarded-Access-Token: eyJhbGciOiJSUzI1NiIs...
X-Forwarded-For: 192.168.1.100
X-Forwarded-Proto: https
X-Forwarded-Host: localhost:4180
```

#### 未認証リクエスト

```http
HTTP/1.1 302 Found
Location: /oauth2/start?rd=%2Fapi%2Fusers
```

---

## 3. 転送ヘッダー仕様

### 3.1 ユーザー情報ヘッダー

| ヘッダー | 説明 | 取得元 |
|----------|------|--------|
| `X-Forwarded-User` | ユーザー識別子 | IDToken `sub` |
| `X-Forwarded-Email` | メールアドレス | IDToken `email` |
| `X-Forwarded-Preferred-Username` | ユーザー名 | IDToken `preferred_username` |
| `X-Forwarded-Groups` | グループ（カンマ区切り） | IDToken `groups` |
| `X-Forwarded-Access-Token` | アクセストークン | TokenResponse |

### 3.2 プロキシ情報ヘッダー

| ヘッダー | 説明 |
|----------|------|
| `X-Forwarded-For` | クライアントIP |
| `X-Forwarded-Proto` | オリジナルプロトコル（http/https） |
| `X-Forwarded-Host` | オリジナルHost |
| `X-Real-IP` | クライアントIP（単一） |

---

## 4. Cookie仕様

### 4.1 セッションCookie

| 属性 | 値 | 説明 |
|------|-----|------|
| Name | `_sekisho` | Cookie名（設定可能） |
| Value | 暗号化されたセッションデータ | AES-CFBで暗号化 |
| Path | `/` | 全パスで有効 |
| HttpOnly | `true` | JavaScript からアクセス不可 |
| Secure | `true`（本番） | HTTPS必須 |
| SameSite | `Lax` | CSRF対策 |
| Max-Age | 86400（24時間） | 有効期限（設定可能） |

### 4.2 CSRF Cookie

| 属性 | 値 |
|------|-----|
| Name | `_sekisho_csrf` |
| Value | 暗号化された `{state, nonce, verifier, rd}` |
| Path | `/` |
| HttpOnly | `true` |
| Secure | `true`（本番） |
| SameSite | `Lax` |
| Max-Age | 300（5分） |

---

## 5. セッションデータ構造

Cookie内に暗号化して格納されるデータ：

```json
{
  "session_id": "uuid-v4",
  "access_token": "eyJhbGciOi...",
  "refresh_token": "eyJhbGciOi...",
  "id_token": "eyJhbGciOi...",
  "expires_at": "2026-01-12T15:00:00Z",
  "user": {
    "sub": "12345678-abcd-efgh",
    "email": "user@example.com",
    "name": "John Doe",
    "preferred_username": "johndoe",
    "groups": ["admin", "users"]
  },
  "created_at": "2026-01-11T15:00:00Z"
}
```

---

## 6. エラーコード一覧

### 6.1 OAuth2エラー

| コード | HTTPステータス | 説明 |
|--------|---------------|------|
| `invalid_state` | 400 | stateパラメータ不一致 |
| `missing_code` | 400 | 認可コードがない |
| `token_exchange_failed` | 500 | トークン交換失敗 |
| `invalid_id_token` | 401 | IDトークン検証失敗 |
| `invalid_nonce` | 401 | nonce不一致 |
| `invalid_audience` | 401 | aud不一致 |
| `session_expired` | 401 | セッション期限切れ |
| `refresh_failed` | 401 | トークンリフレッシュ失敗 |

### 6.2 プロキシエラー

| コード | HTTPステータス | 説明 |
|--------|---------------|------|
| `upstream_unreachable` | 502 | 上流サーバー接続失敗 |
| `upstream_timeout` | 504 | 上流サーバータイムアウト |

### 6.3 エラーレスポンス形式

```json
{
  "error": "error_code",
  "error_description": "Human readable description",
  "request_id": "req-uuid-for-debugging"
}
```

---

## 7. ログ形式

構造化ログ（JSON）：

```json
{
  "timestamp": "2026-01-11T15:54:25+09:00",
  "level": "info",
  "message": "request completed",
  "request_id": "req-12345",
  "method": "GET",
  "path": "/api/users",
  "status": 200,
  "duration_ms": 45,
  "user": "user@example.com",
  "remote_addr": "192.168.1.100"
}
```

---

## 8. 設定可能項目

| 設定項目 | 環境変数 | デフォルト | 説明 |
|----------|----------|------------|------|
| リスンアドレス | `LISTEN_ADDRESS` | `:4180` | サーバーのリスンアドレス |
| 上流URL | `UPSTREAM_URL` | - | 転送先URL（必須） |
| Issuer URL | `OAUTH2_ISSUER_URL` | - | OIDC Issuer URL（必須） |
| Client ID | `OAUTH2_CLIENT_ID` | - | OAuth2 Client ID（必須） |
| Client Secret | `OAUTH2_CLIENT_SECRET` | - | OAuth2 Client Secret（必須） |
| Redirect URL | `OAUTH2_REDIRECT_URL` | - | コールバックURL（必須） |
| Cookie Name | `COOKIE_NAME` | `_sekisho` | セッションCookie名 |
| Cookie Secret | `COOKIE_SECRET` | - | Cookie暗号化キー（必須） |
| Cookie Expire | `COOKIE_EXPIRE` | `24h` | Cookie有効期限 |
| Cookie Secure | `COOKIE_SECURE` | `true` | Secure属性 |
