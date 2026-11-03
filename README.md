# translate-kit

Shared Go document translation extracted from SkillView (MIT). No UI framework, Keychain, settings file, database or document cache.

Reuse one Client. Pass Config and APIKey for each request. Translate preserves Markdown code and link destinations and returns a complete result or an error; progress reports completed chunks, not model tokens. Applications own consent, file revisions, presentation, credentials and caching. Plain text currently follows the existing Markdown segmentation behavior.

Providers: Google web, Youdao web, LibreTranslate and OpenAI-compatible AI. Legacy Google Cloud remains for compatibility. Web endpoints are unofficial and may change. No automatic provider switching. Limits: 1 MiB document, 4096 chunks, 2 MiB response, 60-second requests; bounded rate-limit retries only. Anonymous Youdao state stays in memory.

Interface localization is NOT shared. This package never reads environment API keys or a Keychain.

This is an unpublished local development module. SkillView and CatView currently use a sibling checkout through `replace github.com/sogud/translate-kit => ../translate-kit`. No remote version or tag is available yet. The source is MIT licensed; dependency notices must be completed before distribution.


## Use

```go
client := translatekit.New(translatekit.Options{})
result, err := client.Translate(ctx, translatekit.Request{
    Config: translatekit.Config{
        Provider: "ai", TargetLanguage: "zh-CN",
        BaseURL: baseURL, Model: model,
    },
    APIKey: key, Text: document, Revision: revision,
}, func(progress translatekit.Progress) {
    // Update the app's progress display on its UI thread.
})
```

Call only after the user requests translation and the app identifies the selected provider. Never log Request or document contents. Keep one client to retain pacing and cooldown; use context cancellation when the reading task is canceled. The application must reject stale results when the file revision or active tab changes. TestConnection(ctx, config, key) translates only the fixed sample Hello. and must also be explicitly requested.

Existing focused checks are split between this module's internal provider checks and SkillView's adapter/integration checks:
`go test .` here; `go test ./internal/translation ./internal/app` in SkillView. No real provider requests are needed. Tests moved from SkillView retain its MIT attribution.
