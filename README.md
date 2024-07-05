# registry-backend

The backend API server for [Comfy Registry](https://comfyregistry.org) and [Comfy CI/CD](https://comfyci.org).

Join us at our discord: [https://discord.gg/comfycontrib](https://discord.gg/comfycontrib)

Registry React Frontend [Github](https://github.com/Comfy-Org/registry-web)
Registry CLI [Github](https://github.com/yoland68/comfy-cli)

## Local Development

### Golang

Install Golang:

<https://go.dev/doc/install>

Install go packages

`go get`

### Supabase

Install [Supabase Cli](https://supabase.com/docs/guides/cli/getting-started)

`brew install supabase/tap/supabase`

`supabase start`

Open [Supabase Studio](http://127.0.0.1:54323/project/default) locally.

### Start API Server

`docker compose up`

This commands starts the server with Air that listens to changes. It connects to the Supabase running locally.

### Set up local ADC credentials

These are needed for authenticating Firebase JWT token auth + calling other GCP APIs.

When testing login with registry, use this:
`gcloud config set project dreamboothy-dev`

`gcloud auth application-default login`

If you are testing creating a node, you need to impersonate a service account because it requires signing cloud storage urls.

`gcloud auth application-default login --impersonate-service-account 357148958219-compute@developer.gserviceaccount.com`

TODO(robinhuang): Create a service account suitable for dev.

# Code Generation

Make sure you install the golang packages locally.

`go get`

## Schema Change

Update the files in `ent/schema`.

### Regenerate code

This should search all directories and run go generate. This will run all the commands in the `generate.go` files in the repository.

`go generate ./...`

Or manually run:

`go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --feature sql/lock --feature sql/modifier ./ent/schema`

### Generate Migration Files

Run this command to generate migration files needed for staging/prod database schema changes:

```shell
atlas migrate diff migration \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public"
```

## API Spec Change (openapi.yml)

### Regenerate code

This should search all directories and run go generate. This will run all the commands in the `generate.go` files in the repository.

`go generate ./...`

Or manually run:

`export PATH="$PATH:$HOME/bin:$HOME/go/bin"`

<https://github.com/deepmap/oapi-codegen/issues/795>

`oapi-codegen --config drip/codegen.yaml openapi.yml`

## TroubleShooting / Common Errors

Here are some common errors and how to resolve them.

### Security Scan

If you are calling the `security-scan` endpoint, you need to add the endpoint url to `docker-compose.yml` and then make sure you have the correct permissions to call that function.

Check the `security-scan` Cloud Function repo for instructions on how to do that with `gcloud`.

For non Comfy-Org contributors, you can use your own hosted function or just avoid touching this part. We keep the security scan code private to avoid exploiters taking advantage of it.

### Firebase Token Errors

Usually in localdev, we use dreamboothy-dev Firebase project for authentication. This conflicts with our machine creation logic because all of those machine images are in dreamboothy. TODO(robinhuang): Figure out a solution for this. Either we replicate things in dreamboothy-dev, or we pass project information separately when creating machine images.

### Creating VM instance error

**Example:**

```
{
    "severity": "ERROR",
    "error": "error creating instance: Post \"https://compute.googleapis.com/compute/v1/projects/dreamboothy/zones/us-central1-a/instances\": oauth2: \"invalid_grant\" \"reauth related error (invalid_rapt)\" \"https://support.google.com/a/answer/9368756\"",
    "time": "2024-02-26T01:32:27Z",
    "message": "Error creating instance:"
}

{
    "severity": "ERROR",
    "error": "failed to get session using author id 'nz0vAxfqWLSrqPcUhspyuOEp03z2': error creating instance: Post \"https://compute.googleapis.com/compute/v1/projects/dreamboothy/zones/us-central1-a/instances\": oauth2: \"invalid_grant\" \"reauth related error (invalid_rapt)\" \"https://support.google.com/a/answer/9368756\"",
    "time": "2024-02-26T01:32:27Z",
    "message": "Error occurred Path: /workflows/:id, Method: GET\n"
}
```

**Resolution:**

You would likely need to run `gcloud auth application-default login` again and
restart your docker containers/services to pick up the new credentials.

### Calling CreateSession endpoint

Use the postman collection to call the CreateSession endpoint. You should be able to import changes with `openapi.yml`
file.
You should use this as a request body since there are list of supported GPU type.

```json
{
  "gpu-type": "nvidia-tesla-t4"
}
```

### Bypass Authentication Error

In order to bypass authentication error, you can add make the following changes in `firebase_auth.go` file.

```go
package drip_middleware

func FirebaseMiddleware(entClient *ent.Client) echo.MiddlewareFunc {
 return func(next echo.HandlerFunc) echo.HandlerFunc {
  return func(ctx echo.Context) error {
   userDetails := &UserDetails{
    ID:    "test-james-token-id",
    Email: "test-james-email@gmail.com",
    Name:  "James",
   }

   authdCtx := context.WithValue(ctx.Request().Context(), UserContextKey, userDetails)
   ctx.SetRequest(ctx.Request().WithContext(authdCtx))
   newUserError := db.UpsertUser(ctx.Request().Context(), entClient, userDetails.ID, userDetails.Email, userDetails.Name)
   if newUserError != nil {
    log.Ctx(ctx).Info().Ctx(ctx.Request().Context()).Err(newUserError).Msg("error User upserted successfully.")
   }
   return next(ctx)
  }
 }
}

```

### Machine Image Not Found

We use a custom machine image to create VM instances. That machine image is specified in `docker-compose.yml` file.

```yaml
MACHINE_IMAGE: "comfy-cloud-template-3"
```

If you are getting an error that the machine image is not found, you can create a new machine image by following the
steps below:

**TODO**: explore steps to create machine image with comfy setup.

For the purpose of just testing endpoints, you don't really need to worry about Comfy specific machine image.
You can simply create a new VM on the GCP console and use that VM's image to create a new machine image.
And then update the `docker-compose.yml` file with the new machine image name.

## Clean Up Resources

You can use this script to cleanup resources for specific user.

```shell
`docker compose -f scripts/cleanup/docker-compose.cleanup.yml run --rm cleanup -u <user id>`
```
