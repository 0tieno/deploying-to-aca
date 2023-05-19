---
published: false
type: workshop
title: Deploying to Azure Container Apps
short_title: Deploying to ACA
description: Learn how to deploy containers to Azure Container Apps.
level: beginner                         
authors:                                
  - Josh Duffney
contacts:
  - joshduffney
duration_minutes: 30
tags: Azure, Azure Container Apps, Azure Container Registry, Go, Golang, Containers, Docker
---

# Deploying to Azure Container Apps

In this workshop, you'll learn how deploy a containerized application to Azure Container Apps. Azure Container Apps allows you to deploy containerized applications without having to manage the underlying infrastructure, leaving you to focus on your application.

## Objectives

You'll learn how to:
- Create an Azure Container Apps environment
- Deploy a simple Go web application to Azure Container Apps
- Allow access to the web application with an external ingress
- Deploy revisions of the web application

## Prerequisites

| | |
|----------------------|------------------------------------------------------|
| Azure account        | [Get a free Azure account](https://azure.microsoft.com/free) |
| Azure CLI            | [Install Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) |
| Docker               | [Install Docker](https://docs.docker.com/get-docker/) |

---

## Build and Push the Container Image

In this section, you'll build and push a container image to Azure Container Registry (ACR). That image will be used in the next section to deploy a container app to Azure Container Apps.

### Populate variables

Before you begin, set the following variables:

```powershell
$resource_group="<resource_group>"
$acr_name="<acr_name>"
$tag="<tag>"
$container_app_environment_name="welcome-to-build-aca-env"
$container_app_name="welcome-to-build-aca"
$image_name="welcometobuild"
```

Replace the values in the angle brackets with the following:
- `<resource_group>` - The name of the resource group to deploy the Azure Container Registry instance to.
- `<acr_name>` - The name of the Azure Container Registry instance.
- `<tag>` - A tag to use for the container image.

<div class="info" data-title="note">

> Use your lab machine's login username as the name of your ACR instance and tag name. For example, if your lab machine's login username is `labuser123`, then your ACR instance should be named `labuser123acr` and your image tag would be `labuser123`. This will ensure that your ACR instance name is unique.

</div>

### Deploy an Azure Container Registry instance

Next, run the following command to deploy an Azure Container Registry instance:

```powershell
az acr create `
  --name $acr_name `
  --resource-group $resource_group `
  --sku Basic `
  --admin-enabled true `
  --location eastus
```

### Log in to the Azure CLI and Azure Container Registry

Run the following commands to log in to the Azure CLI and ACR:

```powershell
az login;
az acr login --name $acr_name
```

### Create a container image and push it to ACR

Clone the [Deploying to Azure Container Apps](https://github.com/Duffney/deploying-to-azure-container-apps) repository to your local machine using the following command:

```powershell
git clone https://github.com/Duffney/deploying-to-azure-container-apps.git
```

Next, change into the root of the repository and create a container image using the following command:

```powershell
cd deploying-to-azure-container-apps;

az acr build --registry $acr_name --image "$acr_name.azurecr.io/$image_name:$tag" .
```

<details>
<summary>Example</summary>

```output
az acr build --registry acaworkshopdemo -t acaworkshopdemo.azurecr.io/welcometobuild:labuser123 .
```

</details>

---

## Deploy a Container image to Azure Container Apps

In this section, you'll deploy a containerized Go web application to Azure Container Apps. The application will be accessible via an external ingress and will use environment variables and Azure Container Registry secrets to modify the application's behavior.

### Add the containerapp extension to the Azure CLI

Run the following command to add the containerapp extension to the Azure CLI:

```powershell
az extension add --name containerapp
```

### Create an Azure Container Apps environment

An Azure Container Apps environment is a logical grouping of resources that are used to deploy containerized applications. Within an environment, you can deploy one or more container apps and share resources such as a container registry and secrets.

Run the following command to create an Azure Container Apps environment:

```powershell
az containerapp env create `
  --name $container_app_environment_name `
  --resource_group $resource_group `
  --location eastus
```

<details>

<summary>Example</summary>

```powershell
az containerapp env create `
  --name welcometobuild-aca-env `
  --resource_group aca-demo-rg `
  --location eastus
```

</details>

### Create the container app

Container apps define the container image to deploy, the environment variables to set, and the secrets and or volumes to mount. You can pull imags from Azure Container Registry or Docker Hub and set environment variables and secrets from Azure Key Vault. Container apps can also be deployed with an external ingress, which allows you to access the application from outside the environment. Internal ingress is also available, which allows you to access the application from within the environment.

Run the following commands to create a container app:

```powershell
$token=az acr login --name $acr_name --expose-token --output tsv --query accessToken;

$loginServer=az acr show --name $acr_name --query loginServer --output tsv;
```

```powershell
az containerapp create `
    --name $container_app_name `
    --resource_group $resource_group `
    --environment $container_app_environment_name  `
    --image "$loginServer/$image:$tag" `
    --target-port 8080 `
    --ingress 'external' `
    --registry-server $loginServer `
    --registry-username 00000000-0000-0000-0000-000000000000 `
    --registry-password $token `
    --query properties.configuration.ingress.fqdn
```

<details>

<summary>Example</summary>

```output
az containerapp create \
    --name welcome-to-build-aca \
    --resource_group aca-demo-rg \
    --environment welcometobuild-aca-env  \
    --image $loginServer/welcometobuild:labuser123 \
    --target-port 8080 \
    --ingress 'external' \
    --registry-server $loginServer \
    --registry-username 00000000-0000-0000-0000-000000000000 \
    --registry-password $token \
    --query properties.configuration.ingress.fqdn
```

</details>

Browse to the URL returned by the command to view the application. You should see the following:

![Welcome To Microsoft Build](../images/welcome_to_build.gif)

---

## Deploy a Revision

In this section, you'll deploy a revision of the container app. 

Revisions allow you to deploy new versions of the container app without having to create a new container app. Revisions can be deployed with a new container image, environment variables, secrets, and volumes. 

You'll trigger a new deployment by updating updating the container app's environment variables using a contanier app secret.

### Create a secret

In the Azure Portal, navigate to your Azure Container App. Next, follow the steps below to create a secret:

1. Select **Secrets** from the left-hand menu under **Settings**.
2. Select **+ Add**.
3. Enter `welcome-secret` as the secret's **Key**.
4. Leave `Container Apps Secret` selected.
5. Enter `<YourName>` for the **Value**.
6. Click **Add**.

Replace `<YourName>` with your first and last name.

![Create a secret](../images/create_secret.png)

### Edit the container app

Next, you need to update the container app to use the new secret as an environment variable to change the configuration of the web app. Once the seed is updated, a new revision will be deployed.

Follow the steps below to update the container app:

1. Select **Containers** from the left-hand menu under **Application**.
2. Click **Edit and Deploy**.
3. Check the box next to your container app, and then click **Edit**.

![Add the secret to the container app as an environment variable](../images/add_secret_to_container_app.png)

### Add the secret to the container app as an environment variable

Once the container app is open for editing, follow the steps below to add the secret as an environment variable:

1. Under **Environment Variables**, click **+ Add**.
2. Enter `WELCOME_NAME` for the **Name**.
3. Select `Reference a secret` as the source.
4. Then, select `welcome-secret` as the **Value**.
5. Click **Save**.
6. Click **Create** to deploy the new revision.

![Edit the container app](../images/edit_container_app.png)

### View the new revision

Once the container app is updated, a new revision will be deployed. Follow the steps below to view the new revision:

1. Select **Revision Management** from the left-hand menu under **Application**.
2. Click the revision with the latest **Created** date.
3. Click the link next to the **Revision URL** to view the application.

![View the new revision](../images/view_new_revision.png)