# Sekisho プロジェクト WBS (Work Breakdown Structure)

## 概要
このWBSは、`project_plan.md`（スケジュール）、`requirements.md`（要件）、`architecture.md`（設計）に基づいて作成されています。

## 構造
- **Level 1**: フェーズ (Phase) - プロジェクト計画書の週次スケジュールに対応
- **Level 2**: 機能群 (Feature) - まとまった機能単位
- **Level 3**: タスク (Task) - 具体的な実装・作業項目

---

## Phase 1: 基盤構築 & プロキシ基本機能 (Week 1)
**完了条件 (KPI)**: 
- 開発環境がチーム全員のローカルで動作する
- `/health` エンドポイントが 200 OK を返す
- アップストリームへのプロキシ転送が成功する


### 1.1 開発環境セットアップ
- [ ] **プロジェクト初期化**
    - `go mod init`
    - ディレクトリ構造の作成 (`cmd`, `internal`, `docs` etc.)
    - Makefileの作成 (build, test, run)
- [ ] **Docker環境構築**
    - `docker-compose.yml` 作成
    - Keycloak設定 (import用JSON作成)
    - サンプルアップストリームサーバー (echo server) 用コンテナ設定
- [ ] **ローカルTLS環境構築**
    - `mkcert` インストール & 証明書発行
    - Docker/Goサーバーへの証明書組み込み

### 1.2 設定管理機能 (Config)
- [ ] **Configuration Manager実装** `internal/config`
    - 設定構造体 (`ServerConfig`, `OAuth2Config`, `SessionConfig`) 定義 (REQ: CONFIG-01, 02)
    - 環境変数からの読み込み処理実装 (`kelseyhightower/envconfig` 等)
    - 設定ファイル (YAML) 読み込み処理実装

### 1.3 HTTPサーバー & プロキシ基本実装
- [ ] **HTTP Server実装** `internal/server`
    - Graceful Shutdown対応 (REQ: OPS-02)
    - ヘルスチェックエンドポイント (`/health`) 実装 (REQ: BASIC-02)
- [ ] **リバースプロキシ実装** `internal/handler/proxy.go`
    - `httputil.ReverseProxy` の基本実装 (REQ: PROXY-01)
    - アップストリームへのリクエスト転送確認
    - **単体テスト作成** (`internal/config`, `internal/server`)


---

## Phase 2: OAuth2/OIDC 認証実装 (Week 2)
**完了条件 (KPI)**:
- Keycloak上のテストユーザーでログインフローが完遂する
- アクセストークンが正しく取得・検証できる
- 単体テストカバレッジ 80% 以上 (認証モジュール)


### 2.1 Keycloak設定
- [ ] **Realm & Client設定**
    - テスト用Realm `sekisho` 作成
    - Client `sekisho-proxy` 作成 (Confidential, Valid Redirect URIs設定)
    - テストユーザー作成

### 2.2 OIDCクライアント実装 (Auth)
- [ ] **OIDC基盤実装** `internal/oidc`
    - `OIDCClient` インターフェース定義 (ARCH: 3.2)
    - Discovery Endpoint (`.well-known/openid-configuration`) からの設定取得
- [ ] **PKCE実装** `internal/oidc/pkce.go`
    - `code_verifier` 生成 (REQ: AUTH-04)
    - `code_challenge` 生成 (S256)
- [ ] **認証フロー開始** `internal/handler/oauth2.go`
    - 認可URL生成 (`CreateAuthURL`)
    - `state`, `nonce` 生成

### 2.3 トークン制御 (Token)
- [ ] **トークン交換処理**
    - Authorization Code Grant 実装 (REQ: AUTH-01)
    - トークンレスポンス取得
- [ ] **IDトークン検証** `internal/oidc/token.go`
    - 署名検証 (JWKS) (REQ: AUTH-11)
    - `iss`, `aud`, `exp`, `iat` クレーム検証 (REQ: AUTH-06~10)
    - `nonce` 検証 (REQ: AUTH-05)
- [ ] **トークンリフレッシュ処理**
    - Refresh Tokenを用いたアクセストークン更新 (REQ: AUTH-07)
- [ ] **単体テスト作成** (`internal/oidc`)
    - トークン検証ロジック等のテスト

---

## Phase 3: セッション管理 & セキュリティ (Week 3)
**完了条件 (KPI)**:
- 暗号化されたセッションCookieがブラウザに保存される
- CSRFトークン検証が機能する
- 未認証リクエストが正しくリダイレクトされる


### 3.1 セッション管理 (Session)
- [ ] **セッション構造体定義** `internal/session`
    - Session Model定義 (User, Tokens)
    - Session Model定義 (User, Tokens)
    - Cookie操作用メソッド (`Get`, `Set`, `Clear`) 実装 (Store: Encrypted Cookie)
- [ ] **セッション暗号化**
    - AES-CFB 暗号化/復号化 実装 (REQ: SESSION-03)
    - HKDF-SHA256 による鍵導出 実装 (REQ: SESSION-04)

### 3.2 セキュリティ対策 (Security)
- [ ] **Cookieセキュリティ設定**
    - Secure, HttpOnly, SameSite属性の適用 (REQ: SESSION-01, SEC-02)
- [ ] **CSRF対策**
    - `state` パラメータの検証処理 (REQ: SEC-01)
- [ ] **Open Redirect対策**
    - コールバック後のリダイレクト先検証 (REQ: SEC-04)

### 3.3 認証ミドルウェア (Middleware)
- [ ] **Auth Middleware実装** `internal/middleware/auth.go`
    - セッション有無のチェック
    - 未認証時のリダイレクト処理 (REQ: PROXY-03)
    - 未認証時のリダイレクト処理 (REQ: PROXY-03)
    - トークン有効期限チェックと自動更新
- [ ] **単体テスト作成** (`internal/session`, `internal/middleware`)
    - 暗号化/復号化、ミドルウェアの挙動テスト

---

## Phase 4: 統合 & 仕上げ (Week 4)
**完了条件 (KPI)**:
- E2Eテストシナリオが全てパスする
- ユーザー情報ヘッダーがバックエンドに正しく渡っている


### 4.1 プロキシ拡張
- [ ] **ヘッダー転送処理**
    - `X-Forwarded-User`, `X-Forwarded-Email` 等の付与 (REQ: PROXY-02)
- [ ] **エラーハンドリング強化**
    - エラーレスポンスのJSON化 (REQ: ERROR-02)
    - ユーザーフレンドリーなエラーページ (REQ: ERROR-03)

### 4.2 テスト & 検証
- [ ] **統合テスト (E2E)**
    - ブラウザを用いたログインフロー確認
    - トークン更新フロー確認
    - ログアウト確認

### 4.3 ドキュメント & デプロイ準備
- [ ] **ドキュメント作成**
    - README更新 (セットアップ手順)
    - 学習記録 (ブログ記事下書き)
- [ ] **コードリファクタリング**
    - ログ出力の適正化 (構造化ログ) (REQ: OPS-01)

---

## コンティンジェンシープラン (非常時の対応検討)

### リスク緩和策
1. **Phase 2 (認証) での技術的ハマり**
    - **Trigger**: Week 2の水曜日時点でトークン取得ができていない場合
    - **Action**: 独自実装を諦め、`coreos/go-oidc` などの高レベルライブラリへの依存度を一時的に高める、またはPKCE等の高度な要件をPhase 4へ先送りする（Must要件の「認証そのもの」を優先）。

2. **スケジュール遅延時のスコープ調整 (Scope Cut)**
    - **Priority Low (削減候補)**:
        - 構造化ログの実装 (REQ: OPS-01)
        - ユーザーフレンドリーなエラーページ (REQ: ERROR-03)
        - Open Redirect対策の厳密なホワイトリスト管理 (一旦同一ドメインのみ許可とする等で簡略化)
    - **Priority High (死守)**:
        - OIDC認証フロー
        - セッション暗号化

