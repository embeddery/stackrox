Tag for [build #{{.Build.Number}}]({{.Build.URL}}) is `{{.Env.TAG}}`.

💻 For deploying this image using the dev scripts, run the following first:

```sh
export MAIN_IMAGE_TAG='{{.Env.TAG}}'
```

🕹️ A `roxctl` binary artifact can be [downloaded from CircleCI](https://circleci.com/gh/stackrox/stackrox/{{.Build.Number}}/artifacts).
