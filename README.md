# Peanut

**Peanut** helps you compose your static site from multiple sources with a static-site generator of your choice.

It is especially useful for large sites who want to separate their content across multiple git repositories. 
(i.e: single documentation website for multiple products)

## Getting started
After installing Peanut, create a file called `peanut.yml` in your site's root:
```yaml
version: 1

sources:
  - name: Webhook Documenation
    type: github
    repository_url: https://github.com/adnanh/webhook.git
    files:
      # This allows you to replace the path of a file or directory in the destination
      - /docs/**/*.md:/tools/webhook
      - /docs/logo/logo-256x256.png:/assets/logos/webhook.png
```
Running `peanut` in site's root should pull all the files and organize them in the paths we chose. 