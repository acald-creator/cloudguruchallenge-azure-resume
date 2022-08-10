### Cloud Resume Challenge in Azure
> Antonette Caldwell | https://resume.acaldwell.dev

* Used `Pulumi` to deploy:
    * Azure Resource Group
    * Azure Storage Account
    * Upload the contents needed for the static website
    * Azure CDN Profile
    * Azure CDN Endpoint
    * Azure CDN Custom Domain
    * Azure Function (custom handler in Go)(future)

**NOTE**: Due to `Pulumi` needing to delete a resource and recreate them, I used `az` cli tool to deploy Azure CosmosDB.

I also used `CUELang` as the `cue` support within `Pulumi-Yaml` template, which is new for me because I love testing out new technologies.

My website is fetching a Resume `JSON` file which makes updating information easier to manager.

**GitHub Actions**

Workflow - Create deployment credentials

`az ad sp create-for-rbac --name <ROLE_NAME> --role contributor --scopes /subscriptions/<SUBSCRIPTION_ID>/resourcegroups/<RESOURCE_GROUP> --sdk-auth`

After setting the secrets within GitHub Actions, the current workflow is setup to upload new website contents when the files have been changed, and the necessity to purge the CDN endpoint in order to propagate the new content.