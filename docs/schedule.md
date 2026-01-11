# Sekisho プロジェクトスケジュール

## 前提条件 (Assumptions)
- **稼働想定**: 週5日稼働 + 必要に応じて週末バッファ使用 (1日8時間想定)
- **並行作業**: Phase 1の環境構築と設定実装は部分的に並行して進める
- **リスク管理**: Phase 2 (認証) には技術的不確実性が高いため、専用のバッファを設ける

```mermaid
gantt
    title Sekisho プロジェクト開発スケジュール
    dateFormat  YYYY-MM-DD
    axisFormat  %m/%d
    
    section Phase 1: 基盤構築
    開発環境セットアップ (Docker/TLS)      :p1_1, 2026-01-12, 2d
    Config実装 (並行作業)                 :p1_2, 2026-01-13, 2d
    HTTPサーバー実装 & テスト             :p1_3, after p1_2, 3d

    section Phase 2: OAuth2/OIDC
    Keycloak設定 (Realm/Client)       :p2_1, 2026-01-19, 1d
    OIDCクライアント実装 (PKCE含む)      :p2_2, 2026-01-20, 3d
    トークン処理 & テスト (Auth)         :p2_3, after p2_2, 3d
    Phase 2 予備バッファ                :p2_buf, after p2_3, 2d

    section Phase 3: セッション & セキュリティ
    セッション管理実装 (Encrypted Cookie) :p3_1, 2026-01-29, 3d
    セッション暗号化 & テスト            :p3_2, after p3_1, 2d
    セキュリティ対策 (CSRF/State)        :p3_3, 2026-02-02, 2d
    認証ミドルウェア実装 & テスト         :p3_4, after p3_3, 2d

    section Phase 4: 統合 & 仕上げ
    プロキシヘッダー転送処理             :p4_1, 2026-02-05, 1d
    統合テスト (E2E) & バグ修正          :p4_2, 2026-02-06, 3d
    ドキュメント整備 & リファクタリング    :p4_3, after p4_2, 1d

    section 最終バッファ
    プロジェクト予備日                  :buffer, 2026-02-10, 3d
```

## マイルストーン
- **2026-01-18**: 基盤機能完了（HTTPサーバー起動、プロキシ動作確認）
- **2026-01-28**: 認証フロー疎通（ログイン→トークン取得まで、バッファ込み）
- **2026-02-04**: セキュアなセッション管理完了（暗号化Cookie、CSRF対策）
- **2026-02-10**: MVP機能完成（全機能統合、テストパス）
