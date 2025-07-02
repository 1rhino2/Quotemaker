# Changelog

## [1.2.0] - 2025-07-02

### Added

- Robust refresh button (prevents spamming and API errors)
- Author avatars (Gravatar)
- More visible dark mode toggle
- 404 page for unknown routes
- `/api/quote` endpoint for programmatic access
- Favicon and meta tags

### Fixed

- No more redirect to GitHub
- Refresh button now works reliably even if clicked quickly

### Planned (future)

- Add local fallback quotes if all APIs fail
- Add quote categories (inspirational, funny, etc.)
- Allow users to submit their own quotes
- Add simple analytics (quote count, most popular author)
- Add deployment instructions for cloud providers
- Add Dockerfile for easy deployment
- Add basic tests for quote-fetching logic
