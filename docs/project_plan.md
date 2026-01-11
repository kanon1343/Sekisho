# OAuth2Proxy 自作プロジェクト計画書

## 1. プロジェクト概要

### 1.1 プロジェクト名
**Sekisho（関所）** - OAuth2Proxy 自作プロジェクト

### 1.2 プロジェクトの目的
OAuth2 Proxyを自作することで、以下の技術・概念を深く理解する：
- OAuth 2.0 / OpenID Connect（OIDC）プロトコル
- リバースプロキシの仕組み
- セッション管理とCookie処理
- セキュリティベストプラクティス（CSRF、トークン管理等）
- セキュリティベストプラクティス（CSRF、トークン管理等）

### 1.3 前提条件と制約
- 開発者: 1名（パートタイム）
- 期間: 4週間 + バッファ1週間
- 予算: 0円（ローカル環境のみ）

### 1.4 成功基準（完了条件）
以下のすべてを満たした時点をプロジェクトの完了とする：
1. **機能**：MVP機能要件をすべて満たし、正常に動作すること。
2. **品質**：統合テストがすべてパスすること。
3. **成果**：理解した内容を「学習記録」としてドキュメント化し、Qiita/Zenn等の記事の下書きレベルまで仕上げること。
「関所（Sekisho）」は、日本の歴史において通行人の身元を確認し、許可されたものだけを通過させる施設。OAuth2 Proxyの役割（認証・認可のゲートウェイ）と一致するため、この名前を採用。

---

## 2. プロジェクトスコープ

### 2.1 対象範囲（MVP）
| 項目 | 説明 |
|------|------|
| 認証プロトコル | OAuth 2.0 Authorization Code Flow (w/ PKCE) + OIDC |
| 対応IdP | **Keycloak**（Docker環境） |
| 動作モード | リバースプロキシモード |
| セッション管理 | Cookieのみ |

### 2.2 対象外
- 複数IdP対応（Google, GitHub等）
- Redisセッションストア
- メトリクス/監視機能
- 本番環境での運用

---

## 3. 技術スタック

### 3.1 Go言語
| 技術 | 用途 |
|------|------|
| **Go 1.21+** | メインの実装言語 |
| `net/http` | HTTPサーバー・リバースプロキシ |
| `net/http/httputil` | ReverseProxy実装 |
| `golang.org/x/oauth2` | OAuth2クライアント実装 |
| `github.com/go-chi/chi` | ルーティング（軽量） |


### 3.2 開発環境
| ツール | 用途 |
|--------|------|
| Docker Compose | Keycloak + サンプルアプリ起動 |
| Keycloak | ローカルIdP |
| Makefile | ビルド・テスト自動化 |
| mkcert | ローカル開発用TLS証明書作成 |

---

## 4. 開発フェーズ（4週間）

### Week 1: 基盤構築
- [ ] プロジェクト構造のセットアップ
- [ ] Docker Compose環境構築（Keycloak）
- [ ] ローカルHTTPS環境構築 (mkcert)
- [ ] 基本的なHTTPサーバー実装
- [ ] リバースプロキシ機能の実装

### Week 2: OAuth2認証実装
- [ ] Keycloak設定（Realm, Client作成）
- [ ] Authorization Code Flowの実装 (PKCE対応)
- [ ] トークン取得・検証処理
- [ ] ID Token検証（OIDC nonce, aud check）
- [ ] Token Refresh処理（期限切れ時の自動更新）

### Week 3: セッション管理 & セキュリティ
- [ ] Cookieベースセッション管理
- [ ] CSRF対策（stateパラメータ）
- [ ] セキュアCookie設定
- [ ] リダイレクトURL検証 (Open Redirect対策)
- [ ] ログイン/ログアウト機能

### Week 4: 統合 & 仕上げ
- [ ] ユーザー情報のヘッダー転送
- [ ] エラーハンドリング実装
- [ ] 統合テスト
- [ ] ドキュメント整備

### Week 5: バッファ & リファクタリング
- [ ] 予備日（予期せぬ遅延への対応）
- [ ] コードリファクタリング
- [ ] 最終動作確認

---

## 5. 成果物

| 成果物 | 説明 |
|--------|------|
| ソースコード | OAuth2 Proxy実装（Go） |
| Docker Compose | Keycloak + サンプルアプリ環境 |
| README | セットアップ・使用方法 |
| 学習記録 | 実装中の気づき・学び |

---

## 6. 参考資料

- [OAuth 2.0 Authorization Framework (RFC 6749)](https://datatracker.ietf.org/doc/html/rfc6749)
- [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)
- [oauth2-proxy/oauth2-proxy](https://github.com/oauth2-proxy/oauth2-proxy)
- [Keycloak Documentation](https://www.keycloak.org/documentation)
