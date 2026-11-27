---
title: Publishing Documentation
layout: default
---

# Publishing documentation

GitHub Actions builds `docs/` from `main` and publishes the generated site to
[GitHub Pages](https://open-and-sustainable.github.io/jacquard/).
The source lives alongside the Go library; no separate documentation branch is used.

The workflow is `.github/workflows/pages.yml`:

1. Jekyll builds the Markdown pages with the GreenDoc theme and local layout.
2. The build uploads the generated site as a Pages artifact.
3. The deploy job publishes that artifact to the `github-pages` environment when
   the repository's Pages source is set to GitHub Actions.

Pushes that change `docs/` or the workflow trigger a build and deployment.
Documentation pull requests run the build without deploying. The workflow also
supports manual runs from the repository's Actions tab.

In repository Settings → Pages, select **GitHub Actions** as the publishing source,
then run the Documentation Pages workflow from the Actions tab. The build produces
a site artifact even with publishing disabled; deployment is skipped until that
setting is enabled. Generated HTML belongs in the workflow artifact, not a Git branch.

## Branding and assets

`docs/_config.yml` sets the site title, description, navigation, indigo colors,
and the logo and favicon paths. The logo is `docs/assets/jacquard_logo.png` and the
favicon is `docs/assets/favicon.ico`.
The local layout omits an image tag if its asset file is absent.
