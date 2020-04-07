# Peanut
![](https://github.com/orelb/peanut/workflows/Build/badge.svg)

**Peanut** helps you compose your static site from multiple sources with a static-site generator of your choice.

It is meant to be used as a complementary tool to a static-site generator such as Jekyll or Hugo.
You use peanut to fetch documents from configured sources and then run your site build command. 

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

## Configuration file
The configuration for Peanut is quite simple. It allows you to describe one or more "sources" to fetch documents from.
In addition to that, you can selectively choose which files to include from the source and where to place them locally.

### Placement
Create a `peanut.yml` file in your site's root.

 ### File mapping
 The process of choosing which files to include from a source is called *file mapping* and it takes place in the `files` array of a source.
 
 Each entry in the `files` list is a mapping which consists of two paths separated by a colon. The left-hand side describes the match pattern in the source; and the right-hand side describes the local destination path for the matched files.

 ```yaml
sources:
    - name: Cool Product
      # ...
      files:
        # Glob matching, match multiple files at once
        - /docs/**/*.md:/docs/cool-product
          
        # Match single file using explicit filename
        - /artwork/welcome.png:/assets/cool-product-welcome.png
```
Do note that some properties are omitted for brevity.

**Important**: In the case of glob matching, the destination path must be a directory. The case is different for single-file matching, where the destination path is the destination filename (giving you the option to rename the file if you wish).

